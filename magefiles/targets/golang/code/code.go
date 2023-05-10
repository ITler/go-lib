package code

import (
	"fmt"
	"strings"

	"github.com/itler/go-lib/magefiles/deps"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// BuildBinary triggers go build for main package in repository root
func BuildBinary(name string) error {
	return BuildBinaryFromDir(name, ".")
}

// BuildBinaryFromDir triggers go build for package in specified dir
func BuildBinaryFromDir(name, dir string) error {
	return sh.RunV(mg.GoCmd(), "build", "-o", name,
		"-ldflags", "-w", "-ldflags", "-s", dir)
}

// GenerateMocks creates mock code via vektra/mockery
func GenerateMocks() error {
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

// Lint runs go Lint for all modules
func Lint() error {
	return sh.RunV(deps.Golint.Bin, strings.Split("-set_exit_status ./...", " ")...)
}

// Test runs go Test for all modules
func Test() error {
	return sh.RunV(mg.GoCmd(), "test", "./...", "-short", "-v", "-race",
		"-coverprofile=coverage.out", "-covermode=atomic", "-tags=\"\"")
}

// Vet runs go Vet for all modules
func Vet() error {
	return sh.RunV(mg.GoCmd(), "vet", "./...")
}
