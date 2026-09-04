# go-lib

A Go library providing opinionated abstractions and convenience utilities: env helpers, validation, API/OAuth2 clients, markdown tools, and shared test helpers. Module path: `github.com/itler/go-lib`. Task automation runs through [Mage](https://magefile.org/).

## Dev environment

Prerequisites: Go ≥1.24, Node.js/npm, and the `mage` binary in `$PATH`.

```sh
npm install          # installs husky, commitlint, prettier, semantic-release
mage ci              # installs Go-based external tools (golint, htmltest, etc.)
```

The `prepare` script in `package.json` installs husky hooks only when `$CI` is unset, so local dev gets git hooks automatically.

## Build & test

```sh
mage test            # runs unit tests + lint + vet (the canonical CI gate)
mage go:test         # unit tests only: go test ./... -short -v -race -coverprofile=coverage.out -covermode=atomic
mage go:lint         # golint -set_exit_status ./...
mage go:lint         # also runs go vet ./... internally
mage go:generateMocks  # vektra/mockery via Docker (requires Docker daemon)
```

Mage targets are namespaced: top-level `Test` fans out to `Go:Test` then `Go:Lint`. All mage targets are listed with `mage -l`.

## Packages

| Path         | Purpose                                                             |
| ------------ | ------------------------------------------------------------------- |
| `ease/`      | Core utilities: env, validation, refs, common helpers, test helpers |
| `api/`       | HTTP client interface, token handling, GitHub and OAuth2 clients    |
| `md/`        | Markdown utilities                                                  |
| `pkg/ease/`  | File tools (supplemental to `ease/`)                                |
| `magefiles/` | Mage build targets and external dependency definitions              |

## Conventions

- **Error handling**: errors are wrapped with `fmt.Errorf("...: %w", err)`; `zerolog` (`github.com/rs/zerolog/log`) is the logger — `log.Fatal()`, `log.Warn()`, `log.Debug()` etc. No `fmt.Print` for logging.
- **Test helpers**: `ease.ReadFile`, `ease.UnmarshalTestValue`, and `ease.RequireXxx` wrappers live in `ease/test_helpers.go` — use them in `_test.go` files rather than duplicating.
- **Commit format**: Angular + Conventional Commits enforced by commitlint. Scope case may be lower or upper. Body/footer lines max 260 chars.
- **Pre-commit**: `prettier --write` runs on `*.css, *.less, *.scss, *.md, *.json` via lint-staged.
- **Releases**: `semantic-release` on `main` branch only; changelog written to `docs/CHANGELOG.md`.
- **Module path**: always `github.com/itler/go-lib/...` — do not use the `signavio` mirror path in imports.

## Pitfalls

- `mage go:generateMocks` requires Docker; it runs `vektra/mockery` as a container — fails silently if the daemon is down.
- `coverage.out` at repo root is a generated artifact — do not hand-edit or commit it.
- `dummy.tfstate` at repo root is a fixture file — leave it alone.
- The `GH_TOKEN` / `GH_USER` env vars are needed for CI private module access (written to `~/.netrc`); without them, `go get` of private modules will fail locally too.
- Husky hooks are skipped when `CI` env var is set (by design); locally they run on commit.
