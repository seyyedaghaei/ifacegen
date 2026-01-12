## Contributing

### Setup

```bash
go test ./...
```

If you have `make`, you can use:

```bash
make ci
```

### Guidelines

- Keep changes focused and easy to review.
- Run `make ci` (or `go test ./...` and `golangci-lint run`) before opening a PR.
- Update `README.md` when behavior, flags, or installation steps change.

