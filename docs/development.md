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
- `-version` is populated via `-ldflags` in the release workflow.

