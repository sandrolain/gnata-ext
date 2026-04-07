# gnata extensions

Extended JSONata functions for [gnata](https://github.com/RecoLabs/gnata) — the Go port of JSONata 2.x.

`gnata-ext` ports and adapts the extension functions from the [gosonata](https://github.com/blues/gosonata) `pkg/ext` library to gnata's `CustomFunc` API, providing over **110 additional functions** grouped into thirteen domain packages.

---

## Installation

```sh
go get github.com/sandrolain/gnata-ext
```

Requires Go 1.21+ and `github.com/recolabs/gnata v0.2.1`.

---

## Quick Start

### Register all extensions at once

```go
package main

import (
    "context"
    "fmt"

    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

func main() {
    env := gnata.NewCustomEnv(ext.AllFuncs())

    expr, _ := gnata.Compile(`$uuid()`)
    result, _ := expr.EvalWithCustomFuncs(context.Background(), nil, env)
    fmt.Println(result) // e.g. "550e8400-e29b-41d4-a716-446655440000"
}
```

### Register only a specific package

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/extstring"
)

env := gnata.NewCustomEnv(extstring.All())

expr, _ := gnata.Compile(`$camelCase("hello world foo")`)
result, _ := expr.EvalWithCustomFuncs(ctx, nil, env)
// result → "helloWorldFoo"
```

### Use with StreamEvaluator

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

se := gnata.NewStreamEvaluator(exprs,
    gnata.WithCustomFunctions(ext.AllFuncs()),
)
```

---

## Function Reference

Complete function documentation with examples is in **[docs/FUNCTIONS.md](docs/FUNCTIONS.md)**, organized by package. Each package page includes:

- Function signatures and full descriptions
- Usage examples with JSONata expressions
- Edge cases and implementation notes
- Links to related functions

---

## Selective registration

Each sub-package exposes an `All() map[string]gnata.CustomFunc` function, so you can register only what you need and avoid name collisions:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
    "github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
)

funcs := make(map[string]gnata.CustomFunc)
for k, v := range extcrypto.All() { funcs[k] = v }
for k, v := range extnumeric.All() { funcs[k] = v }

env := gnata.NewCustomEnv(funcs)
```

Full per-package documentation is in [docs/FUNCTIONS.md](docs/FUNCTIONS.md).

Individual functions can also be retrieved directly via their exported Go constructors (e.g. `extcrypto.UUID()`, `extnumeric.Median()`).

---

## EnvBuilder — selective environment construction

`EnvBuilder` provides a fluent API for assembling function sets without manual map operations:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

funcs := ext.NewEnvBuilder().
    WithStringFuncs().
    WithNumericFuncs().
    WithArrayFuncs().
    Without("range", "mode"). // exclude specific functions
    Build()

env := gnata.NewCustomEnv(funcs)
```

`WithAllFuncs()` selects every function in one call; individual `With<Package>Funcs()` methods exist for each of the 13 sub-packages. `Without(names...)` removes names from the accumulated set. `Build()` (or `Funcs()`) returns a `map[string]gnata.CustomFunc` copy.

---

## Presets — pre-built function sets

The `presets` package bundles commonly needed sub-packages into named environments:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/presets"
)

env := gnata.NewCustomEnv(presets.DataEnv())
```

| Preset | Included packages |
|---|---|
| `DataEnv()` | `extarray`, `extobject`, `exttypes`, `extnumeric`, `extpath` |
| `TextEnv()` | `extstring`, `extformat`, `exttypes` |
| `SecureEnv()` | All packages, `$uuid` removed |
| `AnalyticsEnv()` | `extdatetime`, `extarray`, `extnumeric` |

---

## FuncCatalog — function discovery

`Catalog()` returns metadata for every registered function; useful for documentation, tooling, or runtime introspection:

```go
import "github.com/sandrolain/gnata-ext/pkg/ext"

// List all functions
for _, f := range ext.Catalog() {
    fmt.Printf("$%s %s — %s\n", f.Name, f.Signature, f.Description)
}

// Group by package
byPkg := ext.CatalogByPackage()
for _, meta := range byPkg["extarray"] {
    fmt.Println(meta.Name, meta.Signature)
}
```

Each `FuncMeta` carries `Name`, `Package`, `Signature`, and `Description` fields.

---

## Middleware — cross-cutting concerns

The `middleware` package wraps any `map[string]gnata.CustomFunc` with observable, cached, or rate-limited behaviour:

```go
import (
    "log/slog"
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
    "github.com/sandrolain/gnata-ext/pkg/ext/middleware"
)

funcs := ext.AllFuncs()

// Structured logging for every call
funcs = middleware.WithLogging(funcs, slog.Default())

// In-process memoization (exclude non-deterministic functions)
funcs = middleware.WithMemoize(funcs, "uuid", "millis", "now")

// Rate-limit to 100 calls/second across all functions
funcs = middleware.WithRateLimit(funcs, 100)

env := gnata.NewCustomEnv(funcs)
```

Wrappers can be combined in any order and applied to a subset of functions by first selecting them with `EnvBuilder`.

---

## Testing helper

The `exttesting` package simplifies writing table-driven tests against JSONata expressions:

```go
import (
    "testing"

    "github.com/sandrolain/gnata-ext/pkg/ext"
    exttesting "github.com/sandrolain/gnata-ext/pkg/ext/testing"
)

func TestMyExpressions(t *testing.T) {
    env := exttesting.New(ext.AllFuncs(),
        exttesting.WithFrozenTime(1705319400000),        // deterministic $millis()
        exttesting.WithDeterministicUUID("fixed-uuid"),  // deterministic $uuid()
    )

    env.AssertEqual(t, `$uuid()`, nil, "fixed-uuid")
    env.AssertEqual(t, `$first([1,2,3])`, nil, float64(1))
    env.AssertError(t, `$hash("bad", "x")`, nil)
}
```

`Eval(t, expr, data)` evaluates and returns the result; `AssertEqual` and `AssertError` call `t.Fatal` on mismatch. `WithExtraFuncs` injects additional custom functions for a single `Env`.

---

## CLI tool

A command-line tool ships under `cmd/gnata-ext-cli` for quick expression evaluation and function discovery:

```sh
# Install
go install github.com/sandrolain/gnata-ext/cmd/gnata-ext-cli@latest

# Evaluate an expression (all extension functions are available)
gnata-ext-cli eval '$first([1,2,3])'
gnata-ext-cli eval '$hash("sha256","hello")' --data '{"msg":"world"}'
gnata-ext-cli eval '$dateAdd($millis(),1,"day")' --data-file payload.json

# List all functions (or filter by package)
gnata-ext-cli list
gnata-ext-cli list --package extarray

# Show signature and description for a function
gnata-ext-cli describe haversine
```

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

## Development

```sh
# Run all tests (unit + integration)
go test ./...

# Run integration tests only
go test ./pkg/ext/ -run TestIntegration -v

# Run with race detector
go test -race ./...

# Lint
golangci-lint run
```

---

## License

See [LICENSE](LICENSE).
