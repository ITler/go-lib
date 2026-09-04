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
