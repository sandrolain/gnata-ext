# Usage Guide

Practical patterns for integrating `gnata-ext` into Go projects.

---

## Contents

- [Usage Guide](#usage-guide)
  - [Contents](#contents)
  - [Installing the module](#installing-the-module)
  - [Registering all functions](#registering-all-functions)
  - [Registering a single package](#registering-a-single-package)
  - [EnvBuilder — fluent construction](#envbuilder--fluent-construction)
  - [Presets — pre-built function sets](#presets--pre-built-function-sets)
  - [Merging selected packages](#merging-selected-packages)
  - [Using individual functions directly in Go](#using-individual-functions-directly-in-go)
  - [StreamEvaluator integration](#streamevaluator-integration)
  - [Writing expressions with extension functions](#writing-expressions-with-extension-functions)
  - [Conflict resolution and name overrides](#conflict-resolution-and-name-overrides)
  - [Passing the environment to EvalWithCustomFuncs](#passing-the-environment-to-evalwithcustomfuncs)
  - [Middleware — cross-cutting concerns](#middleware--cross-cutting-concerns)
  - [FuncCatalog — function discovery](#funccatalog--function-discovery)
  - [Testing helper](#testing-helper)
  - [CLI tool](#cli-tool)
  - [Examples](#examples)
  - [Integration tests](#integration-tests)

---

## Installing the module

```sh
go get github.com/sandrolain/gnata-ext
```

---

## Registering all functions

Use `ext.AllFuncs()` to get a flat `map[string]gnata.CustomFunc` containing every extension function, then pass it to `gnata.NewCustomEnv`.

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

    expr, err := gnata.Compile(`$camelCase("hello world") & " / " & $uuid()`)
    if err != nil {
        panic(err)
    }

    result, err := expr.EvalWithCustomFuncs(context.Background(), nil, env)
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
    // e.g.: "helloWorld / 550e8400-e29b-41d4-a716-446655440000"
}
```

---

## Registering a single package

Each sub-package exposes `All() map[string]gnata.CustomFunc`. This is the preferred approach when you want a minimal footprint or need to avoid name collisions.

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
)

env := gnata.NewCustomEnv(extcrypto.All())

expr, _ := gnata.Compile(`$hash("sha256", payload)`)
result, _ := expr.EvalWithCustomFuncs(ctx, map[string]any{"payload": "hello"}, env)
```

---

## EnvBuilder — fluent construction

`EnvBuilder` provides a fluent API for assembling function sets. It is more concise than manual map merging and supports explicit exclusion:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

funcs := ext.NewEnvBuilder().
    WithStringFuncs().
    WithNumericFuncs().
    WithArrayFuncs().
    Without("range", "mode").  // exclude by name
    Build()

env := gnata.NewCustomEnv(funcs)
```

Available `With<Package>Funcs()` methods: `WithArrayFuncs`, `WithCryptoFuncs`, `WithDateTimeFuncs`, `WithFormatFuncs`, `WithGeoFuncs`, `WithJSONFuncs`, `WithNetFuncs`, `WithNumericFuncs`, `WithObjectFuncs`, `WithPathFuncs`, `WithStringFuncs`, `WithTypesFuncs`, `WithValidateFuncs`.

`WithAllFuncs()` selects the full set in one call. `Build()` and `Funcs()` are equivalent — both return a `map[string]gnata.CustomFunc` copy.

---

## Presets — pre-built function sets

The `presets` package provides named environments for common use cases:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/presets"
)

env := gnata.NewCustomEnv(presets.DataEnv())
```

| Function | Included packages | Typical use |
|---|---|---|
| `DataEnv()` | `extarray`, `extobject`, `exttypes`, `extnumeric`, `extpath` | Data transformation pipelines |
| `TextEnv()` | `extstring`, `extformat`, `exttypes` | Text processing, templating |
| `SecureEnv()` | All packages, `$uuid` removed | Environments where non-determinism is undesirable |
| `AnalyticsEnv()` | `extdatetime`, `extarray`, `extnumeric` | Aggregation and time-series work |

The returned value is a `map[string]gnata.CustomFunc`; pass it to `gnata.NewCustomEnv` or further customise it with `EnvBuilder.WithPackage()` before use.

---

## Merging selected packages

Manually merge any combination of packages:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/extarray"
    "github.com/sandrolain/gnata-ext/pkg/ext/extobject"
    "github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
)

funcs := make(map[string]gnata.CustomFunc)
for k, v := range extarray.All()  { funcs[k] = v }
for k, v := range extobject.All() { funcs[k] = v }
for k, v := range exttypes.All()  { funcs[k] = v }

env := gnata.NewCustomEnv(funcs)
```

---

## Using individual functions directly in Go

Every exported constructor (`extcrypto.UUID()`, `extnumeric.Median()`, etc.) returns a plain `gnata.CustomFunc` — a `func(args []any, focus any) (any, error)` — which you can call directly from Go without going through gnata at all:

```go
import "github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"

median := extnumeric.Median()
result, err := median([]any{[]any{1.0, 2.0, 3.0}}, nil)
// result → 2.0
```

This is particularly useful when unit-testing application logic that depends on these functions.

---

## StreamEvaluator integration

For high-throughput batch-evaluation use `gnata.WithCustomFunctions`:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

expressions := []gnata.Expression{expr1, expr2}

se := gnata.NewStreamEvaluator(expressions,
    gnata.WithCustomFunctions(ext.AllFuncs()),
)

results, err := se.Eval(ctx, inputData)
```

---

## Writing expressions with extension functions

Extension functions are available in JSONata expressions exactly like built-in functions, with a leading `$`:

```jsonata
/* Array operations */
$chunk($range(1, 11), 3)
/* → [[1,2,3],[4,5,6],[7,8,9],[10]] */

/* Date arithmetic */
$dateAdd($millis(), 30, "day")

/* String transforms */
$snakeCase($trim("  Hello World  "))
/* → "hello_world" */

/* Object manipulation */
$omit(order, ["internalId", "createdAt"])

/* Type guards */
$isArray(items) ? $flatten(items) : [$default(items, [])]

/* Combine several packages */
$toCSV(
  $chunk(users, 100).$[0]
)
```

---

## Conflict resolution and name overrides

`ext.AllFuncs()` merges all sub-packages in this order:

1. `extarray`
2. `extcrypto`
3. `extdatetime`
4. `extformat`
5. `extgeo`
6. `extjson`
7. `extnet`
8. `extnumeric`
9. `extobject`
10. `extpath`
11. `extstring`
12. `exttypes`
13. `extvalidate`

Both `extformat` and `extstring` register a `template` function with identical behaviour. Because `extstring` is merged after `extformat`, the `extstring` implementation wins when you use `ext.AllFuncs()`.

If you need a specific function implementation from a particular package, build the map manually:

```go
funcs := extformat.All()
// Override "template" with a custom implementation:
funcs["template"] = func(args []any, focus any) (any, error) {
    // ...
    return "", nil
}
env := gnata.NewCustomEnv(funcs)
```

---

## Passing the environment to EvalWithCustomFuncs

`gnata.NewCustomEnv` returns `*evaluator.Environment` from gnata's internal package. This type is not importable externally, but it can be stored as `any` and passed back to `EvalWithCustomFuncs`:

```go
var env any = gnata.NewCustomEnv(ext.AllFuncs())

// Later:
result, err := expr.EvalWithCustomFuncs(ctx, data, env)
```

The `ext.NewEnv()` helper is provided for this pattern:

```go
env := ext.NewEnv() // returns any wrapping *evaluator.Environment

result, err := expr.EvalWithCustomFuncs(ctx, data, env)
```

---

## Middleware — cross-cutting concerns

The `middleware` package wraps a `map[string]gnata.CustomFunc` to add logging, memoization, or rate limiting without touching any expression or function implementation:

```go
import (
    "log/slog"
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
    "github.com/sandrolain/gnata-ext/pkg/ext/middleware"
)

funcs := ext.AllFuncs()

// Emit a structured log line for every function call and result
funcs = middleware.WithLogging(funcs, slog.Default())

// Cache results keyed on function name + serialised arguments
// Pass names to exclude non-deterministic functions
funcs = middleware.WithMemoize(funcs, "uuid", "millis", "now")

// Allow at most 200 calls per second across all wrapped functions
funcs = middleware.WithRateLimit(funcs, 200)

env := gnata.NewCustomEnv(funcs)
```

Wrappers are composable and can be applied to any map, including the output of `EnvBuilder.Build()` or a preset.

---

## FuncCatalog — function discovery

`Catalog()` and `CatalogByPackage()` return metadata for all registered extension functions:

```go
import (
    "fmt"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

// All functions, sorted by name
for _, f := range ext.Catalog() {
    fmt.Printf("$%s %s\n  package: %s\n  %s\n\n",
        f.Name, f.Signature, f.Package, f.Description)
}

// Grouped by package
byPkg := ext.CatalogByPackage()
for _, meta := range byPkg["extarray"] {
    fmt.Printf("  $%s — %s\n", meta.Name, meta.Description)
}
```

`FuncMeta` fields: `Name string`, `Package string`, `Signature string`, `Description string`.

---

## Testing helper

`pkg/ext/testing` (package `exttesting`) simplifies table-driven tests against JSONata expressions:

```go
package mypackage_test

import (
    "testing"

    "github.com/sandrolain/gnata-ext/pkg/ext"
    exttesting "github.com/sandrolain/gnata-ext/pkg/ext/testing"
)

func TestExpressions(t *testing.T) {
    env := exttesting.New(ext.AllFuncs(),
        exttesting.WithFrozenTime(1705319400000),        // pins $millis()
        exttesting.WithDeterministicUUID("fixed-uuid"), // pins $uuid()
    )

    env.AssertEqual(t, `$uuid()`, nil, "fixed-uuid")
    env.AssertEqual(t, `$first([1,2,3])`, nil, float64(1))
    env.AssertError(t, `$hash("unsupported-algo","x")`, nil)

    // Raw evaluation
    result := env.Eval(t, `$camelCase("hello world")`, nil)
    _ = result
}
```

`WithExtraFuncs(funcs)` injects additional custom functions into the environment without affecting the base set.

---

## CLI tool

`jn` is a command-line JSONata processor (inspired by `jq`) available under `cmd/jn`:

```sh
# Install
go install github.com/sandrolain/gnata-ext/cmd/jn@latest

# Evaluate an expression — all extension functions are available
echo '[1,2,3]' | jn '$first($)'
jn -n '$hash("sha256","hello")'
echo '{"name":"hello world"}' | jn '$camelCase($.name)'
jn -n '$dateAdd($millis(),7,"day")'

# List all registered functions
jn list

# Filter by package name
jn list --package extarray

# Show signature and description for a specific function
jn describe haversine
jn describe chunk
```

See [docs/guides/cli.md](guides/cli.md) for the full flag reference and jq-compatible options (`-c`, `-r`, `-n`, `-s`, `-f`, `--arg`, `--argjson`, etc.).

---

## Examples

Runnable examples for each domain package are in the [`examples/`](../examples/) directory. Each subdirectory is a standalone `main` package:

| Example | Package(s) used |
|---|---|
| [`examples/allFuncs`](../examples/allFuncs/main.go) | `ext.AllFuncs()` |
| [`examples/strings`](../examples/strings/main.go) | `extstring` |
| [`examples/arrays`](../examples/arrays/main.go) | `extarray` |
| [`examples/objects`](../examples/objects/main.go) | `extobject` |
| [`examples/numeric`](../examples/numeric/main.go) | `extnumeric` |
| [`examples/datetime`](../examples/datetime/main.go) | `extdatetime` |
| [`examples/crypto`](../examples/crypto/main.go) | `extcrypto` |
| [`examples/stream`](../examples/stream/main.go) | `ext.AllFuncs` + `gnata.StreamEvaluator` |
| [`examples/types`](../examples/types/main.go) | `exttypes` |
| [`examples/format`](../examples/format/main.go) | `extformat` |

Run any example directly:

```sh
go run ./examples/allFuncs
go run ./examples/stream
```

---

## Integration tests

End-to-end integration tests that drive full gnata expression evaluation are in [`pkg/ext/integration_test.go`](../pkg/ext/integration_test.go). They cover:

- `TestIntegration_AllFuncs` — `ext.AllFuncs()` + `gnata.NewCustomEnv` + `EvalWithCustomFuncs`
- `TestIntegration_SinglePackage` — single package isolation (`extcrypto`)
- `TestIntegration_SelectivePackages` — manual merge of multiple packages
- `TestIntegration_StreamEvaluator` — `gnata.NewStreamEvaluator` with `gnata.WithCustomFunctions`
- `TestIntegration_Expressions` — table-driven, covering extstring / extnumeric / exttypes
- `TestIntegration_DataDriven` — expressions that reference fields in the input data
- `TestIntegration_DirectFuncCall` — calling a sub-package function as a plain Go function

Run only integration tests:

```sh
go test ./pkg/ext/ -run TestIntegration -v
```
