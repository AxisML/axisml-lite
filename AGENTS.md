# Repository Guidelines

## Project Structure

`axisml-core/` is the Go module for the embeddable Lite control plane and its
Standalone Docker runtime. Public embedding APIs live in `axisml-core/pkg/core`;
runtime implementation stays under `axisml-core/internal/runtime/standalone`.
`deploy/` owns Docker Compose and static ResourcePool/Tenant YAML. Generated
OpenAPI output is committed under `docs/apis/` and embedded into `pkg/core`.

AxisML runtime dependencies resolve through released Go module versions; normal
build, test, image, and document-generation workflows require no sibling
`axisml` checkout.

## Build and Test

- `make build`: build `axisml-core`.
- `make test`: run unit tests.
- `make vet`: vet the runtime and tools modules.
- `make fmt`: run gofmt/goimports.
- `make tidy`: tidy both Go modules.
- `make doc-gen` / `make doc-test`: regenerate or verify the composite OpenAPI.
- `make image`: build the standalone axisml-core image.
- `make lite-up` / `make lite-down`: manage the local Compose stack.

Use `GOCACHE=/tmp/axisml-lite-go-cache` if the default macOS Go cache is not
writable. Do not hand-edit `docs/apis/axisml-core.yaml` or
`axisml-core/pkg/core/openapi.gen.yaml`; regenerate them together.

## Style and Commits

Go code must be gofmt/goimports clean and pass `go vet`. Keep package names
short and lowercase. Tests use Go `testing` plus `testify`. Use Conventional
Commit subjects with the `lite` scope where appropriate, for example
`feat(lite): add runtime capability`.

Never commit secrets, Docker socket data, generated runtime state, local
binaries, coverage output, or a developer-specific `go.work` file.
