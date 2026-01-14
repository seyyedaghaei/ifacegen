## Development

### Requirements

- Go (see `go.mod`)
- `make` (optional, but recommended)

### Quickstart

```bash
make ci
```

### Useful targets

- `make fmt`: format code
- `make test`: run tests
- `make vet`: run `go vet`
- `make lint`: run `golangci-lint` (installs if missing)
- `make build`: build `./ifacegen` from `./cmd/ifacegen`

### Notes

- The CLI entrypoint is `cmd/ifacegen`.
- `--version` / `-v` is populated via `-ldflags` in the release workflow.

### Generator tests (golden files)

`internal/generator/generator_golden_test.go` uses small fixture packages under `testdata/`
and compares the generated `iface_gen.go` output against golden files in:

- `internal/generator/testdata/golden/`

To (re)generate the golden files locally, run:

```bash
UPDATE_GOLDEN=1 go test ./...
```

