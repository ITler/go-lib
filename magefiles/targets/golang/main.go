package golang

import (
	"context"

	"github.com/itler/go-lib/magefiles/targets/golang/code"
	"github.com/magefile/mage/mg"
	"github.com/rs/zerolog/log"
)

// Go encapsulates go related targets to this namespace
type Go mg.Namespace

// Lint runs code linting
func (Go) Lint(ctx context.Context) error {
	err := code.Lint()
	if err != nil {
		return err
	}
	return code.Vet()
}

// GenerateMocks creates mock files via vektra/mockery
func (Go) GenerateMocks() error {
	return code.GenerateMocks()
}

// Test runs unit tests
func (Go) Test(ctx context.Context) error {
	return code.Test()
}

// Test runs all tests
func Test(ctx context.Context) error {
	log.Info().Msg("Run unit tests...")
	mg.Deps(Go.Test)
	log.Info().Msg("Run lint and vet...")
	mg.Deps(Go.Lint)
	log.Info().Msg("...done\n")

	return nil
}
