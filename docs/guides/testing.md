# Testing

Simplify table-driven tests with the `exttesting` helper.

---

## Basic Usage

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

---

## Methods

### `New(funcs map[string]gnata.CustomFunc, opts ...Option)`

Create a test environment with custom functions and options.

### `Eval(t testing.TB, expr string, data any) any`

Evaluate an expression and return the result (fails test on error).

### `AssertEqual(t testing.TB, expr string, data any, want any)`

Assert that the expression result equals the expected value.

### `AssertError(t testing.TB, expr string, data any)`

Assert that the expression raises an error (fails if no error occurs).

---

## Options

### `WithFrozenTime(ts int64)`

Pin `$millis()` and `$now()` to a fixed Unix milliseconds timestamp:

```go
env := exttesting.New(ext.AllFuncs(),
    exttesting.WithFrozenTime(1705319400000), // Jan 15, 2024
)

result := env.Eval(t, `$millis()`, nil)
// Always returns 1705319400000
```

### `WithDeterministicUUID(fixed string)`

Pin `$uuid()` to a fixed UUID string:

```go
env := exttesting.New(ext.AllFuncs(),
    exttesting.WithDeterministicUUID("00000000-0000-0000-0000-000000000000"),
)

result := env.Eval(t, `$uuid()`, nil)
// Always returns "00000000-0000-0000-0000-000000000000"
```

### `WithExtraFuncs(funcs map[string]gnata.CustomFunc)`

Inject additional custom functions for this environment only:

```go
extra := map[string]gnata.CustomFunc{
    "custom": func(args []any, focus any) (any, error) {
        return "custom-result", nil
    },
}

env := exttesting.New(ext.AllFuncs(),
    exttesting.WithExtraFuncs(extra),
)

result := env.Eval(t, `$custom()`, nil)
// Returns "custom-result"
```

---

## Table-Driven Tests

```go
func TestStringFuncs(t *testing.T) {
    env := exttesting.New(ext.AllFuncs())

    tests := []struct {
        name    string
        expr    string
        data    any
        want    any
        wantErr bool
    }{
        {"camelCase", `$camelCase("hello world")`, nil, "helloWorld", false},
        {"capitalize", `$capitalize("alice")`, nil, "Alice", false},
        {"startsWith", `$startsWith("hello", "he")`, nil, true, false},
        {"invalidFunc", `$unknownFunc()`, nil, nil, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.wantErr {
                env.AssertError(t, tt.expr, tt.data)
            } else {
                env.AssertEqual(t, tt.expr, tt.data, tt.want)
            }
        })
    }
}
```

---

## Direct Evaluation

For complex assertions, use `Eval` to get the raw result:

```go
result := env.Eval(t, `[1, 2, 3].$first($.range(1, 5))`, nil)
if result != float64(1) {
    t.Fatalf("got %v, want 1", result)
}
```
