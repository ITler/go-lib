package deps

var (
	// Golint describes the golint dependency
	Golint = &Dependency{
		Bin:       "golint",
		GoInstall: "golang.org/x/lint/golint@latest",
	}
	// Htmltest describes the htmltest dependency
	Htmltest = &Dependency{
		Bin:       "htmltest",
		GoInstall: "github.com/wjdp/htmltest@latest",
	}
	// Hugo describes the hugo dependency
	Hugo = &Dependency{
		Bin:       "hugo",
		GoInstall: "-tags extended github.com/gohugoio/hugo@latest",
	}
)
