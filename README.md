# gnata extensions

<img src="gopher.png" width="320" />

Extended JSONata functions for [gnata](https://github.com/RecoLabs/gnata) — the Go port of JSONata 2.x.

`gnata-ext` ports and adapts the extension functions from the [gosonata](https://github.com/blues/gosonata) `pkg/ext` library to gnata's `CustomFunc` API, providing over **110 additional functions** grouped into thirteen domain packages.

Also includes **`jn`** — a command-line JSONata processor inspired by `jq`, powered by all gnata-ext functions.

> **No naming conflicts:** all gnata-ext custom functions (`$first`, `$camelCase`, `$uuid`, `$haversine`, …) have distinct names from the built-in JSONata functions (`$string`, `$length`, `$sum`, `$map`, `$merge`, etc.).

---

## Installation

**Go library:**

```sh
go get github.com/sandrolain/gnata-ext
```

Requires Go 1.21+ and `github.com/recolabs/gnata v0.2.1`.

**CLI (`jn`):**

```sh
# Install from source
go install github.com/sandrolain/gnata-ext/cmd/jn@latest

# Or download a pre-built binary from GitHub Releases:
# https://github.com/sandrolain/gnata-ext/releases
```

---

## Quick Start

See [docs/guides/quick-start.md](docs/guides/quick-start.md) for getting up and running in minutes, including:

- Registering all extensions at once
- Registering only a specific package
- Using with `StreamEvaluator`
- Selective registration patterns

---

## Function Reference

See [docs/FUNCTIONS.md](docs/FUNCTIONS.md) for complete function documentation, organized by package with quick overviews and links to detailed references.

---

## Usage Guides

Extended documentation and patterns are in `docs/guides/`:

| Guide | Topic |
|---|---|
| [quick-start.md](docs/guides/quick-start.md) | Getting started; basic registration patterns |
| [env-builder.md](docs/guides/env-builder.md) | `EnvBuilder` — fluent environment construction |
| [presets.md](docs/guides/presets.md) | Pre-built environments for common use cases |
| [catalog.md](docs/guides/catalog.md) | Runtime function discovery and introspection |
| [middleware.md](docs/guides/middleware.md) | Logging, memoization, rate limiting |
| [testing.md](docs/guides/testing.md) | Table-driven tests with `exttesting` |
| [cli.md](docs/guides/cli.md) | `jn` CLI — JSONata processor inspired by jq, with color output |
| [development.md](docs/guides/development.md) | Building, testing, contributing |

---

## Examples

Runnable examples for each domain package are in the [`examples/`](examples/) directory:

| Directory | Description |
|---|---|
| [`examples/allFuncs`](examples/allFuncs/main.go) | Use `ext.AllFuncs()` to register everything at once |
| [`examples/strings`](examples/strings/main.go) | `extstring` — case conversion, search, template |
| [`examples/arrays`](examples/arrays/main.go) | `extarray` — range, chunk, flatten, set ops, sliding windows |
| [`examples/objects`](examples/objects/main.go) | `extobject` — pick, omit, deepMerge, pairs, invert |
| [`examples/numeric`](examples/numeric/main.go) | `extnumeric` — clamp, sign, statistics, trig |
| [`examples/datetime`](examples/datetime/main.go) | `extdatetime` — dateAdd, dateDiff, dateComponents, start/end of period |
| [`examples/crypto`](examples/crypto/main.go) | `extcrypto` — uuid, hash, hmac |
| [`examples/stream`](examples/stream/main.go) | `gnata.StreamEvaluator` with `WithCustomFunctions` |
| [`examples/types`](examples/types/main.go) | `exttypes` — type checks, default, identity |
| [`examples/format`](examples/format/main.go) | `extformat` — CSV parse/serialise, template |

Run any example with:

```sh
go run ./examples/allFuncs
```

---

## Notes on HOF (Higher-Order Functions)

gnata does not expose a `Caller` interface for invoking lambda values from within custom functions. As a result, HOF functions from gosonata (`groupBy`, `countBy`, `sumBy`, `minBy`, `maxBy`, `accumulate`, `mapValues`, `mapKeys`, `pipe`, `memoize`) are **not included** in this library. Use the built-in JSONata operators `$map`, `$filter`, `$reduce`, and `$each` instead.

---

## `jn` — Command-Line Processor

`jn` is a JSONata CLI tool inspired by `jq`. It reads JSON from stdin or files, evaluates a JSONata expression with all gnata-ext functions loaded, and writes results to stdout.

```sh
# Pretty-print input
echo '{"name":"Alice","age":30}' | jn '$'

# Field access
echo '{"name":"Alice"}' | jn '$.name'

# Use an extension function
echo '"hello world"' | jn '$camelCase($)'

# Null input (no data)
jn -n '$uuid()'

# Compact output
echo '[3,1,2]' | jn -c '$sort($)'

# Raw string output
echo '{"msg":"hello"}' | jn -r '$.msg'

# Slurp all JSON values into array
cat records.json | jn -s '$count($)'

# List all extension functions
jn list
jn list --package extarray

# Describe a specific function
jn describe haversine

# Show version
jn version
```

See [docs/guides/cli.md](docs/guides/cli.md) for full flag reference and examples.

---

## Development

Building, testing, and contributing — see [docs/guides/development.md](docs/guides/development.md) for:

- Running unit and integration tests
- Linting with `golangci-lint`
- Building the CLI tool
- Project structure overview
- Common development tasks

```sh
# Run all tests
go test ./...

# Run with race detector
go test -race ./...

# Build CLI locally
go build -o jn ./cmd/jn

# Lint
golangci-lint run
```

---

## License

See [LICENSE](LICENSE).
