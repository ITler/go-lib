package deps

var (
	// Docker describes the docker dependency
	Docker = &Dependency{
		Bin: "docker",
	}
	// Yamllint describes the yamllint dependency
	Yamllint = &Dependency{
		Bin: "yamllint",
	}
	// Golint describes the golint dependency
	Golint = &Dependency{
		Bin:       "golint",
		GoInstall: []string{"golang.org/x/lint/golint@latest"},
	}
	// Htmltest describes the htmltest dependency
	Htmltest = &Dependency{
		Bin:       "htmltest",
		GoInstall: []string{"github.com/wjdp/htmltest@latest"},
	}
	// Hugo describes the hugo dependency
	Hugo = &Dependency{
		Bin:       "hugo",
		Env:       map[string]string{"CGO_ENABLED": "1"},
		GoInstall: []string{"-tags", "extended", "github.com/gohugoio/hugo@latest"},
	}
)
