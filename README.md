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

Full per-function documentation with examples is in [docs/FUNCTIONS.md](docs/FUNCTIONS.md).

### extarray — Array utilities

| JSONata function | Description |
|---|---|
| `$first(array)` | First element |
| `$last(array)` | Last element |
| `$take(array, n)` | First *n* elements |
| `$skip(array, n)` | All but the first *n* elements |
| `$slice(array, start [, end])` | Sub-array (negative indices supported) |
| `$flatten(array [, depth])` | Flatten nested arrays |
| `$chunk(array, size)` | Split into fixed-size chunks |
| `$union(a, b)` | Set union (deduplicated) |
| `$intersection(a, b)` | Set intersection |
| `$difference(a, b)` | Elements in *a* not in *b* |
| `$symmetricDifference(a, b)` | Elements in either but not both |
| `$range(start, end [, step])` | Numeric range (exclusive end) |
| `$zipLongest(a, b [, fill])` | Zip two arrays, padding the shorter one |
| `$window(array, size [, step])` | Sliding windows |

### extcrypto — Cryptography

| JSONata function | Description |
|---|---|
| `$uuid()` | Random UUID v4 |
| `$hash(algorithm, value)` | Hex-encoded hash (md5/sha1/sha256/sha384/sha512) |
| `$hmac(algorithm, key, value)` | Hex-encoded HMAC |

### extdatetime — Date & time

All timestamps are **Unix milliseconds** (`float64`), matching JSONata's `$millis()` / `$toMillis()` convention.

| JSONata function | Description |
|---|---|
| `$dateAdd(timestamp, amount, unit)` | Add/subtract a duration |
| `$dateDiff(t1, t2, unit)` | Difference between two timestamps |
| `$dateComponents(timestamp)` | Map of year/month/day/hour/… |
| `$dateStartOf(timestamp, unit)` | Start of the given unit period |
| `$dateEndOf(timestamp, unit)` | End of the given unit period |

Supported units: `year`, `month`, `day`, `hour`, `minute`, `second`, `millisecond` (plural forms also accepted).

### extformat — Formatting

| JSONata function | Description |
|---|---|
| `$csv(text)` | Parse CSV text → array of objects (first row = header) |
| `$toCSV(array)` | Serialize array of objects → CSV text |
| `$template(str, vars)` | Replace `{{key}}` placeholders |

### extnumeric — Extended math

| JSONata function | Description |
|---|---|
| `$log(n [, base])` | Logarithm (natural or given base) |
| `$sign(n)` | −1, 0, or 1 |
| `$trunc(n)` | Truncate toward zero |
| `$clamp(n, min, max)` | Clamp *n* between *min* and *max* |
| `$sin(n)` / `$cos(n)` / `$tan(n)` | Trigonometric functions |
| `$asin(n)` / `$acos(n)` / `$atan(n)` | Inverse trigonometric functions |
| `$atan2(y, x)` | Two-argument arctangent |
| `$pi()` | π constant |
| `$e()` | Euler's number |
| `$median(array)` | Statistical median |
| `$variance(array)` | Population variance |
| `$stddev(array)` | Population standard deviation |
| `$percentile(array, p)` | *p*-th percentile (0–100) |
| `$mode(array)` | Most frequent value(s) |

### extobject — Object utilities

| JSONata function | Description |
|---|---|
| `$values(object)` | Array of object values |
| `$pairs(object)` | Array of `[key, value]` pairs |
| `$fromPairs(pairs)` | Object from `[[key, value], …]` |
| `$pick(object, keys)` | Keep only the specified keys |
| `$omit(object, keys)` | Remove the specified keys |
| `$deepMerge(target, source)` | Recursively merge *source* into *target* |
| `$invert(object)` | Swap keys and string values |
| `$size(object)` | Number of own keys |
| `$rename(object, oldKey, newKey)` | Rename a single key |

### extstring — String utilities

| JSONata function | Description |
|---|---|
| `$startsWith(str, prefix)` | True if *str* starts with *prefix* |
| `$endsWith(str, suffix)` | True if *str* ends with *suffix* |
| `$indexOf(str, search [, start])` | First index of *search* (−1 if not found) |
| `$lastIndexOf(str, search)` | Last index of *search* (−1 if not found) |
| `$capitalize(str)` | Uppercase first char, lowercase rest |
| `$titleCase(str)` | Title-case every word |
| `$camelCase(str)` | camelCase |
| `$snakeCase(str)` | snake_case |
| `$kebabCase(str)` | kebab-case |
| `$repeat(str, n)` | Repeat *str* *n* times |
| `$words(str)` | Split into array of words |
| `$template(str, vars)` | Replace `{{key}}` placeholders |

### exttypes — Type inspection

| JSONata function | Description |
|---|---|
| `$isString(v)` | True if *v* is a string |
| `$isNumber(v)` | True if *v* is a number |
| `$isBoolean(v)` | True if *v* is a boolean |
| `$isArray(v)` | True if *v* is an array |
| `$isObject(v)` | True if *v* is an object |
| `$isNull(v)` | True if *v* is null |
| `$isUndefined(v)` | True if *v* is undefined |
| `$isEmpty(v)` | True for nil, `""`, `[]`, `{}` |
| `$default(v, d)` | *v* if non-nil, otherwise *d* |
| `$identity(v)` | Returns *v* unchanged |

### extpath — Dot-path access

| JSONata function | Description |
|---|---|
| `$get(obj, path [, default])` | Read a nested value by dot-path |
| `$set(obj, path, value)` | Immutable write at dot-path |
| `$del(obj, path)` | Immutable delete at dot-path |
| `$has(obj, path)` | True if path exists and is non-nil |
| `$flattenObj(obj [, sep])` | `{"a":{"b":1}}` → `{"a.b":1}` |
| `$expandObj(obj [, sep])` | `{"a.b":1}` → `{"a":{"b":1}}` |

### extvalidate — Input validation

| JSONata function | Description |
|---|---|
| `$isEmail(str)` | RFC 5322 simplified email format |
| `$isURL(str)` | Valid http/https/ftp URL |
| `$isUUID(str)` | UUID v1–v5 format |
| `$isIPv4(str)` | IPv4 address |
| `$isIPv6(str)` | IPv6 address |
| `$isAlpha(str)` | Only Unicode letters |
| `$isAlphanumeric(str)` | Only Unicode letters and digits |
| `$isNumericStr(str)` | Parses as a number |
| `$matchesRegex(str, pattern)` | Matches RE2 pattern |
| `$inSet(v, set)` | Value is in array set |
| `$minLen(str, n)` | Rune length ≥ n |
| `$maxLen(str, n)` | Rune length ≤ n |
| `$minItems(arr, n)` | Array length ≥ n |
| `$maxItems(arr, n)` | Array length ≤ n |

### extjson — JSON operations

| JSONata function | Description |
|---|---|
| `$jsonParse(str)` | Parse JSON string into a value |
| `$jsonStringify(v [, indent])` | Serialise value to JSON string |
| `$jsonDiff(a, b)` | Differences as a JSON Patch array |
| `$jsonPatch(obj, ops)` | Apply RFC 6902 JSON Patch operations |
| `$jsonPointer(obj, pointer)` | Resolve RFC 6901 JSON Pointer |

### extgeo — Geospatial utilities

All calculations use the WGS-84 mean Earth radius (6 371 km). Coordinates are decimal degrees.

| JSONata function | Description |
|---|---|
| `$haversine(lat1, lon1, lat2, lon2)` | Great-circle distance in km |
| `$bearing(lat1, lon1, lat2, lon2)` | Initial bearing in degrees (0–360) |
| `$geoFormat(lat, lon [, format])` | Format as `"decimal"` or `"dms"` |
| `$geoParse(str)` | Parse `"lat, lon"` string → `{lat, lon}` |
| `$inBoundingBox(lat, lon, minLat, minLon, maxLat, maxLon)` | Point-in-bbox test |
| `$geoDistance(point, points)` | Array of distances from point to each in points (km) |

### extnet — Network utilities

| JSONata function | Description |
|---|---|
| `$ipVersion(str)` | Returns `4`, `6`, or `-1` |
| `$isPrivateIP(str)` | RFC1918 / loopback / link-local |
| `$ipToInt(str)` | IPv4 → uint32 as float64 |
| `$intToIP(n)` | uint32 as float64 → IPv4 string |
| `$ipInCIDR(ip, cidr)` | True if ip is in CIDR block |
| `$expandCIDR(cidr)` | Network info object for CIDR block |

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
