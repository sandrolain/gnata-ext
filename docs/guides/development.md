# Development

Contributing to `gnata-ext` — building, testing, and linting.

---

## Prerequisites

- Go 1.21+
- `github.com/recolabs/gnata v0.2.1`
- Optional: `golangci-lint` for linting

---

## Running Tests

### All tests

```sh
go test ./...
```

Output shows results for all packages (13 extension packages + infrastructure).

### Integration tests only

```sh
go test ./pkg/ext/ -run TestIntegration -v
```

Runs end-to-end tests that drive full gnata expression evaluation.

### Race detector

```sh
go test -race ./...
```

Detects data races in concurrent code.

---

## Linting

```sh
golangci-lint run
```

Checks code style, formatting, and common errors. Configure via `.golangci.yml`.

---

## Building

```sh
go build ./...
```

Builds all packages in the module.

### CLI tool

```sh
go build -o gnata-ext-cli ./cmd/gnata-ext-cli
./gnata-ext-cli eval '$uuid()'
```

Or install globally:

```sh
go install ./cmd/gnata-ext-cli
gnata-ext-cli eval '$first([1,2,3])'
```

---

## Project Structure

```
pkg/ext/
  extarray/           Extension package
  extcrypto/
  ...
  builder.go          EnvBuilder implementation
  catalog.go          FuncCatalog implementation
  middleware/         Cross-cutting wrappers
  presets/            Pre-built environments
  testing/            Test helper (exttesting)

cmd/gnata-ext-cli/    CLI tool source code

docs/
  FUNCTIONS.md        Function reference
  functions/          Per-package documentation
  guides/             Usage guides and patterns
  USAGE.md            Integration patterns

examples/             Runnable example projects
```

---

## Common Development Tasks

### Add a new extension function

1. Edit the relevant `pkg/ext/ext<name>/funcs.go`
2. Add function implementation
3. Update `All()` export
4. Add tests in the same directory
5. Update `docs/functions/ext<name>.md`
6. Run `go test ./pkg/ext/ext<name>`

### Modify pre-built environments

1. Edit `pkg/ext/presets/presets.go`
2. Update tests in `presets_test.go`
3. Build and test: `go test ./pkg/ext/presets`

### Update documentation

- Per-function docs: `docs/functions/*.md`
- Usage patterns: `docs/guides/*.md`
- API reference: `docs/FUNCTIONS.md`
- Integration examples: `docs/USAGE.md`

### Test a specific package

```sh
go test ./pkg/ext/extstring -v
go test ./pkg/ext/middleware -v
```

---

## Continuous Integration

Tests run on every push via GitHub Actions (if configured). Ensure:

1. `go test ./...` passes
2. `go test -race ./...` passes
3. `golangci-lint run` has no errors
