//go:build mage

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ITler/go-lib/magefiles/deps"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	externalDependencies = []*deps.Dependency{
		deps.Golint,
	}
)

type Code mg.Namespace

type Deps mg.Namespace

// Test run all tests
func Test(ctx context.Context) error {
	log.Info().Msg("Run unit tests...")
	mg.Deps(Code.Test)
	log.Info().Msg("Run lint and vet...")
	mg.Deps(Code.Lint)
	log.Info().Msg("...done\n")

	return nil
}

// Lint validates static site configuration
func (Code) Lint(ctx context.Context) error {
	err := sh.RunV(deps.Golint.Bin, strings.Split("-set_exit_status ./...", " ")...)
	if err != nil {
		return err
	}
	err = sh.RunV(mg.GoCmd(), strings.Split("vet ./...", " ")...)
	if err != nil {
		return err
	}
	return nil
}

// Test validates static site configuration
func (Code) Test(ctx context.Context) error {
	err := sh.RunV(mg.GoCmd(), strings.Split("test ./... -short -v -race -coverprofile=coverage.out -covermode=atomic -tags=\"\"", " ")...)
	if err != nil {
		return err
	}
	return nil
}

// Ci installs dependencies in a quick way, suitable for temporary pipeline runners
func (Deps) Ci(ctx context.Context) error {
	for _, dep := range externalDependencies {
		installed, err := dep.Install(ctx)
		if err != nil {
			return err
		}
		if installed {
			log.Info().Msgf("Dependency '%s' installed", dep.Bin)
		}
	}
	return nil
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if err := deps.CheckDependencies(context.Background(), externalDependencies...); err != nil {
		log.Warn().Msg(fmt.Errorf("Dependency not installed - %w", err).Error())
	}
}
