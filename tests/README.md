# AxisML Lite black-box tests

The suite exercises the Lite Compose deployment through its HTTP contracts and
the Platform UI. It uses committed OpenAPI-generated Python clients and treats
the Go implementation as a black box.

## Setup

```sh
cd tests
uv sync
uv run playwright install chromium
uv run test-setup
```

`test-setup` runs the repository's `make lite-up`, waits for Platform readiness,
and resets the test admin password. The sibling `axisml` checkout is required by
the local Go `replace` directives while the shared modules remain unpublished.

## Run

```sh
uv run pytest api
uv run pytest e2e
uv run pytest api/compute_service -k mlrun -v
```

Tests marked `standard_only` are skipped because this repository owns only the
Lite deployment form. Unmarked and `lite_only` tests run against `axisml-core`
on `http://localhost:8090` and Platform on `http://localhost:8080` by default.

## Regenerate clients

```sh
make -C .. client-gen
```

Client generation reads System and Platform OpenAPI specs from the sibling
`../axisml` checkout. Override its location with `AXISML_REPO` when necessary.

## Teardown

```sh
uv run test-teardown
uv run test-teardown --clean
```
