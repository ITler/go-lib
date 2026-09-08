package deps

// This file implements the GithubRelease install strategy for [Dependency].
//
// It handles the full lifecycle of fetching a binary from a GitHub Releases page:
//   - resolving the latest release tag via the GitHub API (when no version is pinned)
//   - constructing the platform-specific asset filename from a Go template
//   - downloading and extracting .tar.gz or .zip archives
//   - installing the result into a bin directory on PATH
//
// Two installation modes are supported, selected by whether [GithubRelease.BundleDir] is set:
//
//  1. Single-binary install (BundleDir is empty):
//     The executable at BinPath is copied directly into InstallDir as Dependency.Bin.
//
//  2. Bundle install (BundleDir is set):
//     Some tools ship as a directory bundle where the entry-point script relies on
//     sibling files at a fixed relative path (e.g. Dart Sass, whose wrapper script
//     exec's src/dart and src/sass.snapshot from the same directory). In this mode
//     the entire BundleDir subtree is copied into InstallDir, and a symlink is created
//     at InstallDir/Dependency.Bin → InstallDir/BinPath so the binary is on PATH while
//     the bundle's internal structure remains intact.
//
// In both modes, BinPath is always relative to the archive root, which keeps its
// meaning consistent regardless of whether BundleDir is also set.

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// GithubRelease describes how to install a [Dependency] from a GitHub Releases asset.
//
// Minimal example — single relocatable binary, always latest:
//
//	&GithubRelease{
//	    Repo:         "cli/cli",
//	    AssetPattern: "gh_{{.Version}}_{{.OS}}_{{.Arch}}.tar.gz",
//	    BinPath:      "gh_{{.Version}}_{{.OS}}_{{.Arch}}/bin/gh",
//	}
//
// Bundle example — wrapper script with runtime siblings (e.g. Dart Sass):
//
//	&GithubRelease{
//	    Repo:         "sass/dart-sass",
//	    AssetPattern: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.tar.gz",
//	    BundleDir:    "dart-sass",
//	    BinPath:      "dart-sass/sass",
//	}
type GithubRelease struct {
	// Repo is the GitHub repository in "owner/repo" form.
	// Example: "sass/dart-sass"
	Repo string

	// Version is the release tag to install, e.g. "1.103.1".
	// Leave empty to resolve the latest published release automatically via the GitHub API.
	// Pin this in consumer code when reproducibility matters:
	//   deps.DartSass.GithubRelease.Version = "1.103.1"
	Version string

	// AssetPattern is a Go text/template that is rendered with [ReleaseAssetData] to produce
	// the release asset filename to download.
	// Example: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.tar.gz"
	AssetPattern string

	// BinPath is the path to the executable inside the archive, always relative to the
	// archive root — regardless of whether BundleDir is set.
	//
	// Single-binary install (BundleDir empty): the file at BinPath is copied to
	// InstallDir/Dependency.Bin.
	//
	// Bundle install (BundleDir set): BinPath locates the entry-point within the installed
	// bundle; a symlink InstallDir/Dependency.Bin → InstallDir/BinPath is created so the
	// binary appears on PATH while the bundle's internal layout stays intact.
	//
	// Example: "dart-sass/sass"
	BinPath string

	// BundleDir is the top-level directory inside the archive that must be installed as a
	// whole unit. Set this when the entry-point binary is a wrapper script (or otherwise
	// not a self-contained executable) that depends on sibling files at a fixed relative
	// path within the same directory.
	//
	// When set, the entire subtree is copied to InstallDir/BundleDir and a symlink is
	// created at InstallDir/Dependency.Bin pointing to InstallDir/BinPath.
	//
	// Leave empty for self-contained binaries that can be freely relocated.
	//
	// Example: "dart-sass"
	BundleDir string

	// InstallDir is the directory into which the binary (or bundle) is placed.
	// Defaults to $GOPATH/bin, falling back to $HOME/.local/bin.
	// Override when a specific location is required.
	InstallDir string
}

// ReleaseAssetData is the template context passed to [GithubRelease.AssetPattern].
// All fields are populated automatically from the resolved version and the current
// platform at install time.
type ReleaseAssetData struct {
	// Version is the resolved release tag, e.g. "1.103.1".
	Version string
	// OS is the GOOS-derived platform token used in release asset names.
	// Mapping: darwin → "macos", windows → "windows", anything else → "linux".
	OS string
	// Arch is the GOARCH-derived architecture token used in release asset names.
	// Mapping: arm64 → "arm64", arm → "arm", anything else → "x64".
	Arch string
}

// installFromGithubRelease is the entry point called by [InstallDependencies].
// It resolves the version, builds the download URL, fetches and extracts the archive,
// then delegates to the appropriate install mode based on [GithubRelease.BundleDir].
func installFromGithubRelease(ctx context.Context, dep *Dependency) error {
	gr := dep.GithubRelease

	version := gr.Version
	if version == "" {
		var err error
		version, err = resolveLatestRelease(ctx, gr.Repo)
		if err != nil {
			return fmt.Errorf("resolve latest release for %s: %w", gr.Repo, err)
		}
	}

	assetData := ReleaseAssetData{
		Version: version,
		OS:      mapOS(runtime.GOOS),
		Arch:    mapArch(runtime.GOARCH),
	}

	assetName, err := renderAssetName(gr.AssetPattern, assetData)
	if err != nil {
		return fmt.Errorf("render asset pattern: %w", err)
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/%s",
		gr.Repo, version, assetName,
	)

	tmpDir, err := os.MkdirTemp("", "gh-release-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, assetName)
	if err := downloadFile(ctx, downloadURL, archivePath); err != nil {
		return fmt.Errorf("download %s: %w", downloadURL, err)
	}

	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("create extract dir: %w", err)
	}

	switch {
	case strings.HasSuffix(assetName, ".tar.gz"):
		if err := extractTarGz(archivePath, extractDir); err != nil {
			return fmt.Errorf("extract tar.gz: %w", err)
		}
	case strings.HasSuffix(assetName, ".zip"):
		if err := extractZip(archivePath, extractDir); err != nil {
			return fmt.Errorf("extract zip: %w", err)
		}
	default:
		return fmt.Errorf("unsupported archive format: %s", assetName)
	}

	installDir := gr.InstallDir
	if installDir == "" {
		installDir = defaultInstallDir()
	}
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("create install dir %s: %w", installDir, err)
	}

	if gr.BundleDir != "" {
		return installBundle(dep, gr, extractDir, installDir)
	}
	return installBinary(dep, gr, extractDir, installDir)
}

// installBinary copies the single executable at BinPath (archive-relative) into installDir
// as Dependency.Bin. Used when the binary is self-contained and freely relocatable.
func installBinary(dep *Dependency, gr *GithubRelease, extractDir, installDir string) error {
	src := filepath.Join(extractDir, filepath.FromSlash(gr.BinPath))
	dest := filepath.Join(installDir, dep.Bin)
	if err := copyFile(src, dest, 0755); err != nil {
		return fmt.Errorf("install binary to %s: %w", dest, err)
	}
	return nil
}

// installBundle copies the entire BundleDir subtree into installDir, then creates a
// symlink installDir/Dependency.Bin → installDir/BinPath so the entry-point is on PATH
// while the bundle's internal layout (and any relative paths within it) stays intact.
func installBundle(dep *Dependency, gr *GithubRelease, extractDir, installDir string) error {
	srcBundle := filepath.Join(extractDir, filepath.FromSlash(gr.BundleDir))
	destBundle := filepath.Join(installDir, filepath.FromSlash(gr.BundleDir))
	if err := copyTree(srcBundle, destBundle); err != nil {
		return fmt.Errorf("install bundle to %s: %w", destBundle, err)
	}

	// Symlink installDir/Bin → installDir/BinPath so the entry-point is on PATH.
	symlinkPath := filepath.Join(installDir, dep.Bin)
	symlinkTarget := filepath.Join(installDir, filepath.FromSlash(gr.BinPath))
	_ = os.Remove(symlinkPath) // remove stale symlink if present
	if err := os.Symlink(symlinkTarget, symlinkPath); err != nil {
		return fmt.Errorf("create symlink %s → %s: %w", symlinkPath, symlinkTarget, err)
	}
	return nil
}

// resolveLatestRelease queries the GitHub API for the latest published release tag
// of the given "owner/repo". Used when [GithubRelease.Version] is not pinned.
func resolveLatestRelease(ctx context.Context, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if release.TagName == "" {
		return "", fmt.Errorf("empty tag_name in response")
	}
	return release.TagName, nil
}

// renderAssetName renders the AssetPattern template with the given [ReleaseAssetData]
// to produce the filename of the release asset to download.
func renderAssetName(pattern string, data ReleaseAssetData) (string, error) {
	tmpl, err := template.New("asset").Parse(pattern)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// mapOS maps runtime.GOOS to the OS token used in GitHub release asset filenames.
// darwin → "macos", windows → "windows", everything else → "linux".
func mapOS(goos string) string {
	switch goos {
	case "darwin":
		return "macos"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

// mapArch maps runtime.GOARCH to the architecture token used in GitHub release asset filenames.
// arm64 → "arm64", arm → "arm", everything else → "x64".
func mapArch(goarch string) string {
	switch goarch {
	case "arm64":
		return "arm64"
	case "arm":
		return "arm"
	default:
		return "x64"
	}
}

// defaultInstallDir returns the default directory for installed binaries.
// Prefers $GOPATH/bin (always on PATH when actions/setup-go is used in CI),
// falling back to $HOME/.local/bin.
func defaultInstallDir() string {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return filepath.Join(gopath, "bin")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/usr", "local", "bin")
	}
	return filepath.Join(home, ".local", "bin")
}

// downloadFile fetches the given URL and writes the response body to dest.
func downloadFile(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %s downloading %s", resp.Status, url)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// extractTarGz extracts a .tar.gz archive into destDir, preserving file modes.
// Rejects entries whose resolved path escapes destDir (path-traversal guard).
func extractTarGz(src, destDir string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, filepath.FromSlash(hdr.Name))
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal path in archive: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
	return nil
}

// extractZip extracts a .zip archive into destDir, preserving file modes.
// Rejects entries whose resolved path escapes destDir (path-traversal guard).
func extractZip(src, destDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(destDir, filepath.FromSlash(f.Name))
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal path in archive: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		rc.Close()
		out.Close()
		if copyErr != nil {
			return copyErr
		}
	}
	return nil
}

// copyFile copies the file at src to dest with the given permission bits.
func copyFile(src, dest string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// copyTree recursively copies the directory tree rooted at src into dest,
// preserving file permission bits.
func copyTree(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target, info.Mode())
	})
}
