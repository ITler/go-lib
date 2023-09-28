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
		// Golint,
	}
)

// Dependency encapsulates attributes of a depending application
type Dependency struct {
	Bin            string
	InstallMethods []InstallMethod
	InstallCmdF    []InstallCmdF
}

func (o *Dependency) Check() error {
	return checkDependency(o)
}

func (o *Dependency) Install(ctx context.Context) error {
	return installDependency(ctx, o)
}

func (o *Dependency) Via(meths ...InstallMethod) Installable {
	if meths != nil {
		o.InstallMethods = append(o.InstallMethods, meths...)
	}

	return o
}

// Check determines if provided dependencies are available.
// If binary is not available [os.ErrNotExist] is returned
func Check(ctx context.Context, dependencies ...*Dependency) error {
	for _, dep := range ensureDependencies(dependencies...) {
		if err := dep.Check(); err != nil {
			return err
		}
	}
	return nil
}

func checkDependency(dep *Dependency) error {
	if _, err := exec.LookPath(dep.Bin); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return &os.PathError{
				Path: dep.Bin,
				Err:  os.ErrNotExist,
			}
		}
		return err
	}
	return nil
}

// InstallIfMissing makes sure to install provided dependencies
// and return if all installing all dependencies were successful
func InstallIfMissing(ctx context.Context, dependencies ...*Dependency) (bool, error) {
	for _, dep := range ensureDependencies(dependencies...) {
		if err := dep.Check(); err != nil {
			if os.IsNotExist(err) {
				if err = dep.Install(ctx); err != nil {
					return false, err
				}
			} else {
				return false, err
			}
		}
	}
	return true, nil
}

func Install(ctx context.Context, dependencies ...*Dependency) (bool, error) {
	var warnings error

	for _, dep := range ensureDependencies(dependencies...) {
		if len(dep.InstallMethods) == 0 {
			warnings = errors.Join(warnings, fmt.Errorf(
				"Installation of '%s' not supported, yet. "+
					"Thus installation needs to be handled, externally", dep.Bin))
			continue
		}
		if err := dep.Install(ctx); err != nil {
			return false, err
		}
	}
	return true, warnings
}

func installDependency(ctx context.Context, dep *Dependency) error {
	if dep.GoInstall != "" {
		if err := sh.RunV(mg.GoCmd(), "install", dep.GoInstall); err != nil {
			return fmt.Errorf("Dependency install: %w", err)
		}
	}

	return nil
}

func ensureDependencies(deps ...*Dependency) []*Dependency {
	if deps == nil || len(deps) == 0 {
		deps = ExternalDependencies
	}
	return deps
}
