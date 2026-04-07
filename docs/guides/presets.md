# Presets

Pre-built environment configurations for common use cases.

---

## Available Presets

The `presets` package bundles commonly needed sub-packages into named environments:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext/presets"
)

env := gnata.NewCustomEnv(presets.DataEnv())
```

| Preset | Included packages | Use case |
|---|---|---|
| `DataEnv()` | `extarray`, `extobject`, `exttypes`, `extnumeric`, `extpath` | Data transformation pipelines |
| `TextEnv()` | `extstring`, `extformat`, `exttypes` | Text processing, templating |
| `SecureEnv()` | All packages, `$uuid` removed | Non-deterministic functions excluded |
| `AnalyticsEnv()` | `extdatetime`, `extarray`, `extnumeric` | Aggregation and time-series work |

---

## Custom Presets

Extend presets with `EnvBuilder`:

```go
import (
    "github.com/recolabs/gnata"
    "github.com/sandrolain/gnata-ext/pkg/ext"
    "github.com/sandrolain/gnata-ext/pkg/ext/presets"
)

funcs := ext.NewEnvBuilder().
    WithPackage(presets.DataEnv()).
    WithDateTimeFuncs(). // add more
    Without("flatten").  // remove if needed
    Build()

env := gnata.NewCustomEnv(funcs)
```
