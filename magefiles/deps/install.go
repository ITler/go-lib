package deps

import "context"

// Installable describes something that can be installed
type Installable interface {
	Check() error
	Install(ctx context.Context) error
	Via(...InstallMethod) Installable
}

type executor struct {
	foo string
}

type Executor interface {
	Run()
}

func (o *executor) Run() {

}

type InstallCmdF func(opts ...interface{}) (bool, error)

type InstallMethod struct {
	args string
}

// func InstallWithOpts(deps []Dependency, opts ...interface{}) (bool, error) {

// }
