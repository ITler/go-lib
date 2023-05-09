//go:build mage

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/itler/go-lib/magefiles/deps"

	"github.com/itler/go-lib/magefiles/run/golang"
	"github.com/magefile/mage/mg"
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
	err := golang.RunLint()
	if err != nil {
		return err
	}
	return golang.RunVet()
}

// Test validates static site configuration
func (Code) Test(ctx context.Context) error {
	return golang.RunTest()
}

// Ci installs dependencies in a quick way, suitable for temporary pipeline runners
func (Deps) Ci(ctx context.Context) error {
	return deps.InstallDependencies(ctx, externalDependencies...)
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if err := deps.CheckDependencies(context.Background(), externalDependencies...); err != nil {
		log.Warn().Msg(fmt.Errorf("Dependency not installed - %w", err).Error())
	}
}
