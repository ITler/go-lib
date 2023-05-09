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

// RunBuildBinary triggers go build for main package in repository root
func RunBuildBinary(binName string) error {
	return RunBuildBinaryFromDir(binName, ".")
}

// RunBuildBinaryFromDir triggers go build for package in specified dir
func RunBuildBinaryFromDir(binName, dir string) error {
	return sh.RunV(mg.GoCmd(), "build", "-o", binName,
		"-ldflags", "-w", "-ldflags", "-s", dir)
}

// RunLint runs go lint for all modules
func RunLint() error {
	return sh.RunV(deps.Golint.Bin, strings.Split("-set_exit_status ./...", " ")...)
}

// RunTest runs go test for all modules
func RunTest() error {
	return sh.RunV(mg.GoCmd(), "test", "./...", "-short", "-v", "-race",
		"-coverprofile=coverage.out", "-covermode=atomic", "-tags=\"\"")
}

// RunVet runs go vet for all modules
func RunVet() error {
	return sh.RunV(mg.GoCmd(), "vet", "./...")
}
