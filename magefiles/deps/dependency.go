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

// Dependency encapsulates attributes of a depending application
type Dependency struct {
	Bin       string
	GoInstall string
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
		if dep.GoInstall != "" {
			if err := sh.RunV(mg.GoCmd(), "install", dep.GoInstall); err != nil {
				return false, fmt.Errorf("Dependency cannot be installed: %w", err)
			}
			continue
		}

		warnings = errors.Join(warnings, fmt.Errorf("Installation of '%s' not supported, yet. "+
			"Thus installation needs to be handled, externally", dep.Bin))
	}
	return true, warnings
}
