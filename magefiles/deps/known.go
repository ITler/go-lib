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
	// DartSass describes the Dart Sass dependency (https://sass-lang.com).
	// Version is intentionally unset — the latest published release is resolved
	// automatically via the GitHub API. Pin it in consumer code when reproducibility
	// is required:
	//
	//	deps.DartSass.GithubRelease.Version = "1.103.1"
	DartSass = &Dependency{
		Bin: "sass",
		GithubRelease: &GithubRelease{
			Repo:         "sass/dart-sass",
			AssetPattern: "dart-sass-{{.Version}}-{{.OS}}-{{.Arch}}.tar.gz",
			BundleDir:    "dart-sass",
			BinPath:      "dart-sass/sass",
		},
	}
)
