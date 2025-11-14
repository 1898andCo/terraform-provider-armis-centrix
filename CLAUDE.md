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

- **`armis/`**: Armis API client with auth (`client.go`, `auth.go`), CRUD operations (`{resource}.go`), and models (`model_{resource}.go`); uses functional options and token caching.
- **`internal/provider/`**: Terraform resource and data source implementations plus acceptance tests.
- **`docs/`**: Generated provider docs (`task docs` refreshes from code descriptions).
- **`examples/`**: Usage samples consumed in documentation.
- **`tools/`**: Doc generation helpers executed by `task docs`.

## Architecture Patterns

### Dual Model System
- **API models** (`armis/model_*.go`): Match Armis API JSON structure exactly; use `json` tags.
- **Terraform models** (resource files, `internal/utils/model_*.go`): Map to Terraform state; use `tfsdk` tags.
- **Conversion**: Helper functions bridge models (e.g., `buildArmisUser()` converts Terraform → API).

### API Client Design
- **Functional options**: Configure client via `armis.WithAPIURL()`, `armis.WithHTTPClient()`, etc.
- **Token caching**: Bearer token auto-refreshes 5 minutes before expiry; `sync.RWMutex` ensures thread safety.
- **CRUD signature**: `Get(ctx, id) (*Model, error)`, `Create(ctx, model) (*Model, error)`, `Update(ctx, model, id)`, `Delete(ctx, id) (bool, error)`.

### Resource Implementation Pattern
Each resource follows this structure:
1. **armis/model_{resource}.go**: Define API response structs (GetResponse, CreateResponse, etc.).
2. **armis/{resource}.go**: Implement CRUD methods; add sentinel errors to `armis/errors.go`.
3. **internal/provider/{resource}_resource.go**:
   - Implement `resource.Resource` + `resource.ResourceWithConfigure` + `resource.ResourceWithImportState`.
   - Define schema with validators (use `stringvalidator`, `int64validator` from framework).
   - `Create/Read/Update/Delete` methods: extract plan/state → call armis client → update state.
   - Handle 404 in Read: `resp.State.RemoveResource(ctx)` for drift detection.
4. **internal/provider/{resource}_resource_test.go**: Write acceptance tests with test fixtures.
5. **internal/sweep/{resource}_resource_test_sweeper.go**: Register sweeper for test cleanup.
6. **provider.go**: Add factory to `Resources()` slice.
7. Run `task docs` to regenerate documentation from schema descriptions.

### Error Handling Strategy
Three-level approach:
- **Validation errors**: Sentinel errors in `armis/errors.go` (e.g., `ErrCollectorInvalidType`).
- **API errors**: `armis.APIError` with StatusCode and Body; use `appendAPIError()` for diagnostics.
- **Context wrapping**: `fmt.Errorf("%w", err)` preserves error chain for debugging.

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
- **Models**: Maintain separate API models (`armis/`) and Terraform models (`internal/provider`, `internal/utils`); use helper functions for conversion.
- **Struct tags**: When defining nested structs in API models, apply `omitempty` tags to individual fields within the nested struct (e.g., `Name string \`json:"name,omitempty"\``), not to the parent struct field itself; this avoids linter warnings about redundant tags while achieving proper JSON marshaling behavior.

## Testing Guidance

- **Acceptance tests**: Require `ARMIS_API_KEY` and `ARMIS_API_URL`; `task test` must pass locally before commits land to keep the provider healthy.
- **Unit tests**: Add table-driven tests alongside implementations in `*_test.go` files; rely on mockable interfaces where possible.
- **Cleanup**: Extend `task sweep` when resources need explicit teardown to prevent drift in shared environments.
- **Pre-submit**: Run `go test ./...` before opening a PR; ensure acceptance tests pass locally when touching provider logic.

## Running Individual Tests

- **Single unit test**: `go test ./internal/utils -v -run TestBuildRoleRequest`
- **Single acceptance test**: `TF_ACC=true go test ./internal/provider -v -run TestAccCollectorResource_basic -timeout=30m`
- **Package-specific tests**: `go test ./armis -v` or `go test ./internal/utils -v`
- **With coverage**: `go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out`
- **Parallel execution**: Tests use `t.Parallel()` for speed; ensure parallel-safe when adding new tests.

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
