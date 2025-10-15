# CLAUDE.md

Guidance for Claude Code when contributing to the Armis Centrix Terraform/OpenTofu provider.

## Repository Overview

- **Focus**: Provider that manages Armis Centrix resources via the Armis API.
- **Primary libs**: Terraform Plugin SDK v2, HashiCorp testing harness, internal `armis/` SDK.
- **Go version**: `1.25` as documented in `README.md`.

## Local Environment

- **Install Task**: `brew install go-task` for Taskfile support.
- **Init tools**: `task install-tools` installs `tfplugindocs`, `goreleaser`, and `tfproviderlint` when required.
- **Environment vars**: set `ARMIS_API_KEY` and `ARMIS_API_URL` (see `docs/Contributing.md`).

## Repository Layout

- **`armis/`**: Armis API client that configures the HTTP client and exposes the structs plus CRUD helpers used by Terraform resources and data sources.
- **`internal/provider/`**: Terraform resource and data source implementations plus acceptance tests.
- **`docs/`**: Generated provider docs (`task docs` refreshes from code descriptions).
- **`examples/`**: Usage samples consumed in documentation.
- **`tools/`**: Doc generation helpers executed by `task docs`.

## Taskfile Reference

- **`task build`**: Runs `go build -v ./...` and produces `terraform-provider-armis-centrix`.
- **`task fmt`**: Invokes `gofumpt -w .`; run after code edits.
- **`task lint`**: Executes `golangci-lint`, `tfproviderlint`, ensures `gofmt` cleanliness, and runs `go vet`.
- **`task test`**: Acceptance suite (`go test ./internal/provider/... -v -timeout=30m`) followed by `task sweep` cleanup.
- **`task docs`**: Regenerates Terraform Registry docs; depends on `install-tools`.
- **`task install`/`task uninstall`**: Manage local `~/.terraformrc` overrides for testing a locally built provider.
- **`task prep`**: Composite `fmt`, `lint`, `docs`, `install`, and `go mod tidy` for release readiness.

## Coding Standards

- **Formatting**: Always run `task fmt`; CI enforces `gofumpt`, `gofmt`, and `golangci-lint`.
- **Imports**: Group as standard library, third-party, repository local; maintain alphabetical order.
- **Error handling**: Wrap errors with context using `fmt.Errorf("%w", err)`; leverage sentinel errors in `internal/provider` when available.
- **Comments**: Follow Go doc comment conventions; keep exported symbol docs in sync with generated docs (`task docs`).
- **APIs**: Use helper functions in `armis/` for REST interactions; avoid duplicating HTTP logic.

## Testing Guidance

- **Acceptance tests**: Require `ARMIS_API_KEY` and `ARMIS_API_URL`; `task test` must pass locally before commits land to keep the provider healthy.
- **Unit tests**: Add table-driven tests alongside implementations in `*_test.go` files; rely on mockable interfaces where possible.
- **Cleanup**: Extend `task sweep` when resources need explicit teardown to prevent drift in shared environments.
- **Pre-submit**: Run `go test ./...` before opening a PR; ensure acceptance tests pass locally when touching provider logic.

## Documentation Workflow

- **Generation**: Update schema descriptions/comments inside resource or data source definitions, then run `task docs` (also invoked automatically by the pre-commit hook via `task prep`).
- **Manual docs**: Edit `docs/` markdown if additional context is required; keep formatting consistent with generated structure.
- **Examples**: Sync Terraform snippets in `examples/` with documented usage.

## Pull Requests

- **Checklist**:
  - **`task fmt`** and **`task lint`** succeed locally.
  - **Tests**: Relevant `go test` and acceptance suites executed; include command output in PR if non-trivial.
  - **Docs**: Regenerate via `task docs` when resource descriptions change.
  - **Changelog**: Follow repository release process if a new version is required (check `.github` workflows).
- **Template**: Complete `.github/PULL_REQUEST_TEMPLATE.md` sections (`what`, `why`, `references`).

## Release & Distribution

- **Local validation**: Use `task install` to drop provider binary into `$GOPATH/bin` for Terraform override testing.
- **Versioning**: Update `terraform-registry-manifest.json` and tags as directed by maintainers; automated releases use `goreleaser` per `.goreleaser.yml`.
- **Registry docs**: Ensure `docs/index.md` edits preserve provider naming fix applied by the sed command in `task docs`.

## Additional Notes

- **Dependencies**: Run `go mod tidy` (part of `task prep`) after adding/removing imports.
- **API changes**: Coordinate with the Armis Centrix API team before introducing breaking changes in `armis/` models.
- **Security**: Never log sensitive API tokens; scrub before committing.
- **Pre-commit**: Hooks run `task prep -f` automatically for Go files under `main.go` and `internal/`; do not bypass with `--no-verify`.
- **Hook contents**: `task prep` formats (`gofumpt`), lints (`golangci-lint`, `tfproviderlint`, `go vet`), regenerates docs, installs the provider locally, and tidies modules.
- **Failures**: Address hook failures locally, rerun `task prep`, and recommit once clean.
