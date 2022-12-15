package deps

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/rs/zerolog/log"
)

// Installable describes something that can be installed
type Installable interface {
	Install() error
}

// Dependency encapsulates attributes of a depending application
type Dependency struct {
	Bin       string
	GoInstall string
	NoInstall bool
}

// Install will install the dependency
func (d *Dependency) Install(ctx context.Context) (result bool, err error) {
	if err = CheckDependencies(ctx, d); err != nil {
		if os.IsNotExist(err) {
			if err = InstallDependencies(ctx, d); err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
}

// CheckDependencies determines if provided dependencies are available.
// If binary is not available [os.ErrNotExist] is returned
func CheckDependencies(ctx context.Context, dependencies ...*Dependency) error {
	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep.Bin); err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				return &os.PathError{
					Op:   "",
					Path: dep.Bin,
					Err:  os.ErrNotExist,
				}
			}
		}
	}
	return nil
}

// InstallDependencies makes sure to install provided dependencies
func InstallDependencies(ctx context.Context, dependencies ...*Dependency) error {
	for _, dep := range dependencies {
		if dep.GoInstall != "" {
			if err := sh.RunV(mg.GoCmd(), strings.Split("install "+dep.GoInstall, " ")...); err != nil {
				return fmt.Errorf("Dependency cannot be installed - %w", err)
			}
			continue
		}

		if !dep.NoInstall {
			log.Info().Msgf("Installation of '%s' will not be managed.", dep.Bin)
		}
	}
	return nil
}
