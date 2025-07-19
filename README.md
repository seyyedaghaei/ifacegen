# ifacegen

`ifacegen` is a Go code generation tool that scans Go packages to automatically generate interfaces based on struct method sets. It's particularly useful for creating abstractions such as `Service`, `Repository`, or similar patterns.

## Features

- Generates interfaces for structs matching configurable name patterns (e.g., `*Service`, `*Repository`).
- Skips specific structs or methods via `// ifacegen:skip` comment.
- Includes structs explicitly marked with `// ifacegen:generate`, even if their name doesn't match patterns.
- Preserves method comments in the generated interfaces.
- Resolves types using `go/types` for accurate imports and signatures.
- Handles cross-package imports with proper aliasing.
- Concurrently processes multiple packages for improved performance.
- Outputs one file per package (default name: `iface_gen.go`).
- Automatically formats and organizes the output with `goimports`.

## Installation

```bash
go install github.com/seyyedaghaei/ifacegen@latest
```

## Usage

```bash
ifacegen -match '*Service,*Repository' -output iface_gen.go ./...
```

### Flags

| Flag         | Description                                                                 |
|--------------|-----------------------------------------------------------------------------|
| `-match`     | Comma-separated list of glob patterns to match struct names (e.g., `*Service`). |
| `-output`    | Output filename used inside each package (default: `iface_gen.go`).         |
| `-name`      | Pattern for naming interfaces, using `{}` as the placeholder (default: `I{}`). |
| `-help`      | Show usage information.                                                     |

## Skipping and Including Structs or Methods

You can **exclude** specific structs or methods from interface generation using the `// ifacegen:skip` comment:

```go
// ifacegen:skip
type AuthService struct {
	// ...
}
```

```go
// ifacegen:skip
func (s *UserService) InternalLogic() {}
```

You can **explicitly include** structs (even if they don't match `-match`) with the `// ifacegen:generate` comment:

```go
// ifacegen:generate
type Special struct {
	// ...
}
```

## Example

Given this struct:

```go
// UserService handles user-related operations.
type UserService struct {}

func (s *UserService) CreateUser(name string) error {
	return nil
}
```

Running:

```bash
ifacegen -match '*Service' ./...
```

Generates:

```go
// UserService handles user-related operations.
type IUserService interface {
	CreateUser(name string) error
}
```

## Output Location

The generated file is written **once per package**, in the same directory as the structs, using the filename from `-output` (default `iface_gen.go`). If the content hasn't changed, the file won't be overwritten.

## License

MIT © Morteza SeyyedAgahei
