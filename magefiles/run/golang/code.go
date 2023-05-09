package golang

import (
	"fmt"
	"strings"

	"github.com/itler/go-lib/magefiles/deps"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// GenerateVektraMocks creates mock files via vektra/mockery
func GenerateVektraMocks() error {
	userID, err := sh.Output("id", "-u")
	if err != nil {
		return fmt.Errorf("Unable to determine users id: %w", err)
	}
	groupID, err := sh.Output("id", "-g")
	if err != nil {
		return fmt.Errorf("Unable to determine users group id: %w", err)
	}
	user := fmt.Sprintf("%s:%s", userID, groupID)
	return sh.RunV(deps.Docker.Bin, strings.Split(
		"run --user "+user+" -v ${PWD}:/src -w /src vektra/mockery --all", " ")...)

}

// RunLintDefault runs go lint for all modules
func RunLintDefault() error {
	return sh.RunV(deps.Golint.Bin, strings.Split("-set_exit_status ./...", " ")...)
}

// RunTestDefault runs go test for all modules
func RunTestDefault() error {
	return sh.RunV(mg.GoCmd(), "test", "./...", "-short", "-v", "-race",
		"-coverprofile=coverage.out", "-covermode=atomic", "-tags=\"\"")
}

// RunVetDefault runs go vet for all modules
func RunVetDefault() error {
	return sh.RunV(mg.GoCmd(), "vet", "./...")
}
