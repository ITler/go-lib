package deps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

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

// Dependency encapsulates attributes of a depending application.
// Exactly one install strategy field should be set (GoInstall or GithubRelease).
// If none is set, installation is a no-op with a warning; the binary is expected
// to be provided externally.
type Dependency struct {
	// Bin is the executable name used for PATH lookup and as the installed binary name.
	Bin string
	// GoInstall holds the arguments passed to `go install` (build flags + package path).
	// Example: []string{"-tags", "extended", "github.com/gohugoio/hugo@latest"}
	GoInstall []string
	// Env holds optional environment variables set when running go install.
	// Example: map[string]string{"CGO_ENABLED": "1"}
	Env map[string]string
	// GithubRelease, when set, installs the binary by downloading a GitHub Releases asset.
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
