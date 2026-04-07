# Quick Start

Get up and running with `gnata-ext` in minutes.

---

## Register all extensions at once

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

---

## Register only a specific package

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

---

## Use with StreamEvaluator

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

Individual functions can also be retrieved directly via their exported Go constructors (e.g. `extcrypto.UUID()`, `extnumeric.Median()`).
