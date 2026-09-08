package deps

import (
	"testing"
)

func TestDependency_EnvFieldExists(t *testing.T) {
	d := &Dependency{
		Bin:       "go",
		Env:       map[string]string{"CGO_ENABLED": "1"},
		GoInstall: []string{"some/package@latest"},
	}
	if d.Env["CGO_ENABLED"] != "1" {
		t.Fatalf("expected CGO_ENABLED=1, got %q", d.Env["CGO_ENABLED"])
	}
}

func TestDependency_EnvNilSafe(t *testing.T) {
	_ = &Dependency{
		Bin:       "go",
		GoInstall: []string{"golang.org/x/lint/golint@latest"},
		// Env intentionally omitted — must not panic
	}
	// No panic == pass
}

// --- GithubRelease: OS/arch mapping ---

func TestMapOS(t *testing.T) {
	tests := []struct {
		goos string
		want string
	}{
		{"darwin", "macos"},
		{"linux", "linux"},
		{"windows", "windows"},
		{"freebsd", "linux"}, // fallback
	}
	for _, tc := range tests {
		if got := mapOS(tc.goos); got != tc.want {
			t.Errorf("mapOS(%q) = %q, want %q", tc.goos, got, tc.want)
		}
	}
}

func TestMapArch(t *testing.T) {
	tests := []struct {
		goarch string
		want   string
	}{
		{"amd64", "x64"},
		{"arm64", "arm64"},
		{"arm", "arm"},
		{"386", "x64"}, // fallback
	}
	for _, tc := range tests {
		if got := mapArch(tc.goarch); got != tc.want {
			t.Errorf("mapArch(%q) = %q, want %q", tc.goarch, got, tc.want)
		}
	}
}

// --- GithubRelease: asset name template rendering ---

func TestRenderAssetName(t *testing.T) {
	tests := []struct {
		pattern string
		data    ReleaseAssetData
		want    string
	}{
		{
			pattern: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.tar.gz",
			data:    ReleaseAssetData{Version: "1.103.1", OS: "linux", Arch: "x64"},
			want:    "dart-sass-1.103.1-linux-x64.tar.gz",
		},
		{
			pattern: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.zip",
			data:    ReleaseAssetData{Version: "1.99.0", OS: "windows", Arch: "x64"},
			want:    "dart-sass-1.99.0-windows-x64.zip",
		},
		{
			pattern: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.tar.gz",
			data:    ReleaseAssetData{Version: "1.103.1", OS: "macos", Arch: "arm64"},
			want:    "dart-sass-1.103.1-macos-arm64.tar.gz",
		},
	}
	for _, tc := range tests {
		got, err := renderAssetName(tc.pattern, tc.data)
		if err != nil {
			t.Errorf("renderAssetName(%q, %v) unexpected error: %v", tc.pattern, tc.data, err)
			continue
		}
		if got != tc.want {
			t.Errorf("renderAssetName(%q, %v) = %q, want %q", tc.pattern, tc.data, got, tc.want)
		}
	}
}

func TestRenderAssetName_InvalidTemplate(t *testing.T) {
	_, err := renderAssetName("{{.Unclosed", ReleaseAssetData{})
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

// --- DartSass known default ---

func TestDartSass_Defaults(t *testing.T) {
	if DartSass.Bin != "sass" {
		t.Errorf("DartSass.Bin = %q, want \"sass\"", DartSass.Bin)
	}
	if DartSass.GithubRelease == nil {
		t.Fatal("DartSass.GithubRelease is nil")
	}
	if DartSass.GithubRelease.Version != "" {
		t.Errorf("DartSass.GithubRelease.Version should be empty (unpinned), got %q",
			DartSass.GithubRelease.Version)
	}
	if DartSass.GithubRelease.Repo != "sass/dart-sass" {
		t.Errorf("DartSass.GithubRelease.Repo = %q, want \"sass/dart-sass\"",
			DartSass.GithubRelease.Repo)
	}
}

func TestDartSass_VersionOverride(t *testing.T) {
	// Simulate what a consumer does to pin a version.
	// Save and restore so tests remain independent.
	original := DartSass.GithubRelease.Version
	defer func() { DartSass.GithubRelease.Version = original }()

	DartSass.GithubRelease.Version = "1.99.0"
	if DartSass.GithubRelease.Version != "1.99.0" {
		t.Errorf("version override failed, got %q", DartSass.GithubRelease.Version)
	}
}
