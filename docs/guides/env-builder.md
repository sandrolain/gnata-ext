# EnvBuilder Pattern

Build custom function environments with a fluent API.

---

## Basic Usage

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

---

## Available Methods

- `WithAllFuncs()` — Select every function in one call
- `WithArrayFuncs()` — `extarray` functions
- `WithCryptoFuncs()` — `extcrypto` functions
- `WithDateTimeFuncs()` — `extdatetime` functions
- `WithFormatFuncs()` — `extformat` functions
- `WithGeoFuncs()` — `extgeo` functions
- `WithJSONFuncs()` — `extjson` functions
- `WithNetFuncs()` — `extnet` functions
- `WithNumericFuncs()` — `extnumeric` functions
- `WithObjectFuncs()` — `extobject` functions
- `WithPathFuncs()` — `extpath` functions
- `WithStringFuncs()` — `extstring` functions
- `WithTypesFuncs()` — `exttypes` functions
- `WithValidateFuncs()` — `extvalidate` functions
- `WithPackage(pkg map[string]gnata.CustomFunc)` — Add a custom package
- `WithFunc(name string, fn gnata.CustomFunc)` — Add a single function
- `Without(names ...string)` — Remove specific functions by name
- `Build()` or `Funcs()` — Return the accumulated map

---

## Composition with Middleware

Chain with middleware wrappers for observability:

```go
import (
    "log/slog"
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
    "github.com/sandrolain/gnata-ext/pkg/ext/middleware"
)

funcs := ext.NewEnvBuilder().
    WithStringFuncs().
    WithArrayFuncs().
    Build()

// Add logging and memoization
funcs = middleware.WithLogging(funcs, slog.Default())
funcs = middleware.WithMemoize(funcs, "uuid")

env := gnata.NewCustomEnv(funcs)
```
