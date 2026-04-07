# Middleware

Add cross-cutting concerns like logging, memoization, and rate limiting.

---

## Basic Usage

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

---

## WithLogging

Emits a structured log line for every function call and result:

```go
import (
    "log/slog"
    "github.com/sandrolain/gnata-ext/pkg/ext"
    "github.com/sandrolain/gnata-ext/pkg/ext/middleware"
)

funcs := middleware.WithLogging(ext.AllFuncs(), slog.Default())
// Each call logs: level="INFO" name="uuid" args=[...] result="..." err=<nil>
```

---

## WithMemoize

Caches results keyed on function name + serialised arguments. Pass function names to exclude non-deterministic functions:

```go
funcs := middleware.WithMemoize(ext.AllFuncs(), "uuid", "millis", "now")
// $first([1,2,3]) result cached, but $uuid() called fresh each time
```

---

## WithRateLimit

Enforces a global rate limit across all wrapped functions (requests per second):

```go
funcs := middleware.WithRateLimit(ext.AllFuncs(), 100) // 100 calls/sec
// RPS < 1 (e.g. 0) is a no-op
```

---

## Composition

Stack multiple wrappers:

```go
funcs := ext.AllFuncs()
funcs = middleware.WithLogging(funcs, slog.Default())
funcs = middleware.WithMemoize(funcs)
funcs = middleware.WithRateLimit(funcs, 200)

env := gnata.NewCustomEnv(funcs)
```

Order matters — in the example above, logging wraps memoization, which wraps rate limiting.

---

## Selective Application

Apply middleware only to a subset of functions:

```go
base := ext.NewEnvBuilder().
    WithStringFuncs().
    WithArrayFuncs().
    Build()

// Only wrap string functions with logging
stringFuncs := make(map[string]gnata.CustomFunc)
for k, v := range base {
    if k != "uuid" && k != "hash" { // string-related funcs
        stringFuncs[k] = v
    }
}

stringFuncs = middleware.WithLogging(stringFuncs, slog.Default())
env := gnata.NewCustomEnv(stringFuncs)
```
