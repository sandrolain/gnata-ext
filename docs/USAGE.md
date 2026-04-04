# Usage Guide

Practical patterns for integrating `gnata-ext` into Go projects.

---

## Contents

- [Installing the module](#installing-the-module)
- [Registering all functions](#registering-all-functions)
- [Registering a single package](#registering-a-single-package)
- [Merging selected packages](#merging-selected-packages)
- [Using individual functions directly in Go](#using-individual-functions-directly-in-go)
- [StreamEvaluator integration](#streamevaluator-integration)
- [Writing expressions with extension functions](#writing-expressions-with-extension-functions)
- [Conflict resolution and name overrides](#conflict-resolution-and-name-overrides)
- [Passing the environment to EvalWithCustomFuncs](#passing-the-environment-to-evalwithcustomfuncs)
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

`ext.AllFuncs()` merges all sub-packages. The merge order is:

1. `extarray`
2. `extcrypto`
3. `extdatetime`
4. `extformat`
5. `extnumeric`
6. `extobject`
7. `extstring`
8. `exttypes`

Both `extformat` and `extstring` register a `template` function with identical behaviour. Because `extstring` is merged after `extformat`, the `extstring` implementation wins when you use `ext.AllFuncs()`. If you need a specific implementation, build the map manually:

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
