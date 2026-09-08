package deps

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// ExternalDependencies define the external app dependencies needed
	// to fullfill all automation tasks
	ExternalDependencies = []*Dependency{
		Golint,
	}
)

// Installable describes something that can be installed
type Installable interface {
	Install() error
}

// GithubRelease describes how to fetch a binary from a GitHub Releases asset.
// AssetPattern is a Go text/template rendered with [ReleaseAssetData].
// BinPath is the relative path inside the extracted archive to the executable.
// InstallDir is optional; defaults to the directory returned by [defaultInstallDir].
// Version is optional; when empty the latest published release is resolved via the GitHub API.
type GithubRelease struct {
	// Repo is "owner/repo", e.g. "sass/dart-sass"
	Repo string
	// Version is the release tag, e.g. "1.103.1". Empty = resolve latest via API.
	Version string
	// AssetPattern is a Go text/template applied to [ReleaseAssetData] to produce
	// the asset filename. Example: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.tar.gz"
	AssetPattern string
	// BinPath is the relative path inside the extracted archive to the binary.
	// Example: "dart-sass/sass"
	BinPath string
	// InstallDir is the directory where the binary is placed.
	// Defaults to $GOPATH/bin, then $HOME/.local/bin.
	InstallDir string
}

// ReleaseAssetData is the template context available in [GithubRelease.AssetPattern].
type ReleaseAssetData struct {
	// Version is the resolved release tag, e.g. "1.103.1"
	Version string
	// OS is the GOOS-mapped value used in release asset names: "linux", "macos", "windows"
	OS string
	// Arch is the GOARCH-mapped value used in release asset names: "x64", "arm64", "arm"
	Arch string
}

// Dependency encapsulates attributes of a depending application
type Dependency struct {
	Bin string
	// GoInstall holds the arguments passed to `go install` (build flags + package path).
	// Example: []string{"-tags", "extended", "github.com/gohugoio/hugo@latest"}
	GoInstall []string
	// Env holds optional environment variables set when running go install.
	// Example: map[string]string{"CGO_ENABLED": "1"}
	Env map[string]string
	// GithubRelease, when set, installs the binary by downloading a GitHub release asset.
	// Mutually exclusive with GoInstall.
	GithubRelease *GithubRelease
}

// Install will install the dependency
func (d *Dependency) Install(ctx context.Context) (result bool, err error) {
	if err = CheckDependencies(ctx, d); err != nil {
		if os.IsNotExist(err) {
			return InstallDependencies(ctx, d)
		}
		return true, nil
	}
	return false, nil
}

// CheckDependencies determines if provided dependencies are available.
// If binary is not available [os.ErrNotExist] is returned
func CheckDependencies(ctx context.Context, dependencies ...*Dependency) error {
	if dependencies == nil || len(dependencies) == 0 {
		dependencies = ExternalDependencies
	}
	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep.Bin); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				return &os.PathError{
					Path: dep.Bin,
					Err:  os.ErrNotExist,
				}
			}
			return err
		}
	}
	return nil
}

// InstallDependencies makes sure to install provided dependencies
// and return if all installing all dependencies were successful
func InstallDependencies(ctx context.Context, dependencies ...*Dependency) (bool, error) {
	var warnings error
	if dependencies == nil || len(dependencies) == 0 {
		dependencies = ExternalDependencies
	}
	for _, dep := range dependencies {
		if len(dep.GoInstall) > 0 {
			args := append([]string{"install"}, dep.GoInstall...)
			if err := sh.RunWithV(dep.Env, mg.GoCmd(), args...); err != nil {
				return false, fmt.Errorf("Dependency cannot be installed: %w", err)
			}
			continue
		}

		if dep.GithubRelease != nil {
			if err := installFromGithubRelease(ctx, dep); err != nil {
				return false, fmt.Errorf("dependency %q cannot be installed: %w", dep.Bin, err)
			}
			continue
		}

		warnings = errors.Join(warnings, fmt.Errorf("Installation of '%s' not supported, yet. "+
			"Thus installation needs to be handled, externally", dep.Bin))
	}
	return true, warnings
}

// installFromGithubRelease downloads, extracts, and installs a binary from a GitHub release asset.
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

	assetName, err := renderTemplate(gr.AssetPattern, assetData)
	if err != nil {
		return fmt.Errorf("render asset pattern: %w", err)
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", gr.Repo, version, assetName)

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

	srcBin := filepath.Join(extractDir, filepath.FromSlash(gr.BinPath))

	installDir := gr.InstallDir
	if installDir == "" {
		installDir = defaultInstallDir()
	}
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("create install dir %s: %w", installDir, err)
	}

	destBin := filepath.Join(installDir, dep.Bin)
	if err := copyFile(srcBin, destBin, 0755); err != nil {
		return fmt.Errorf("install binary to %s: %w", destBin, err)
	}

	return nil
}

// resolveLatestRelease queries the GitHub API for the latest release tag of the given repo.
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

// renderTemplate renders a Go text/template string with the given data.
func renderTemplate(pattern string, data ReleaseAssetData) (string, error) {
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

// mapOS maps runtime.GOOS to the OS identifier used in Dart Sass (and common) release assets.
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

// mapArch maps runtime.GOARCH to the architecture identifier used in Dart Sass release assets.
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

// defaultInstallDir returns $GOPATH/bin if GOPATH is set, otherwise $HOME/.local/bin.
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

// downloadFile fetches url and writes the body to dest.
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

// extractTarGz extracts a .tar.gz archive into destDir.
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
		// Guard against path traversal
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

// extractZip extracts a .zip archive into destDir.
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

// copyFile copies src to dest with the given permission bits.
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
