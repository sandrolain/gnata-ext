# Function Reference

Complete documentation for every function provided by `gnata-ext`, split by package.

All functions are registered as JSONata custom functions and are called with a leading `$` in expressions (e.g. `$uuid()`). Timestamps are always **Unix milliseconds** (`float64`) unless stated otherwise.

---

## Quick Overview

| Package | Count | Description |
|---|---|---|
| `extarray` | 14 | Array utilities (slicing, set ops, ranges, windows) |
| `extcrypto` | 3 | UUID generation, hash, HMAC |
| `extdatetime` | 5 | Date arithmetic and decomposition |
| `extformat` | 3 | CSV parsing/serialisation, template rendering |
| `extnumeric` | 14 | Math, trigonometry, statistics |
| `extobject` | 9 | Object manipulation |
| `extstring` | 12 | String utilities and case conversion |
| `exttypes` | 10 | Type predicates and coercion helpers |
| `extpath` | 6 | Dot-path read/write on nested objects |
| `extvalidate` | 14 | Validation predicates (email, URL, IP, regex …) |
| `extjson` | 5 | JSON parse/stringify, diff, patch, pointer |
| `extgeo` | 6 | Geographic calculations (haversine, bearing, bounding box) |
| `extnet` | 6 | IP address and CIDR utilities |
| **Total** | **~110** | Across 13 packages |

---

## By Package

### extarray — Array Utilities

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

**Detailed reference:** [functions/extarray.md](functions/extarray.md)

---

### extcrypto — Cryptography

| JSONata function | Description |
|---|---|
| `$uuid()` | Random UUID v4 |
| `$hash(algorithm, value)` | Hex-encoded hash (md5/sha1/sha256/sha384/sha512) |
| `$hmac(algorithm, key, value)` | Hex-encoded HMAC |

**Detailed reference:** [functions/extcrypto.md](functions/extcrypto.md)

---

### extdatetime — Date & Time

All timestamps are **Unix milliseconds** (`float64`), matching JSONata's `$millis()` / `$toMillis()` convention.

| JSONata function | Description |
|---|---|
| `$dateAdd(timestamp, amount, unit)` | Add/subtract a duration |
| `$dateDiff(t1, t2, unit)` | Difference between two timestamps |
| `$dateComponents(timestamp)` | Map of year/month/day/hour/… |
| `$dateStartOf(timestamp, unit)` | Start of the given unit period |
| `$dateEndOf(timestamp, unit)` | End of the given unit period |

Supported units: `year`, `month`, `day`, `hour`, `minute`, `second`, `millisecond` (plural forms also accepted).

**Detailed reference:** [functions/extdatetime.md](functions/extdatetime.md)

---

### extformat — Formatting

| JSONata function | Description |
|---|---|
| `$csv(text)` | Parse CSV text → array of objects (first row = header) |
| `$toCSV(array)` | Serialize array of objects → CSV text |
| `$template(str, vars)` | Replace `{{key}}` placeholders |

**Detailed reference:** [functions/extformat.md](functions/extformat.md)

---

### extnumeric — Extended Math

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

**Detailed reference:** [functions/extnumeric.md](functions/extnumeric.md)

---

### extobject — Object Utilities

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

**Detailed reference:** [functions/extobject.md](functions/extobject.md)

---

### extstring — String Utilities

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

**Detailed reference:** [functions/extstring.md](functions/extstring.md)

---

### exttypes — Type Inspection

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

**Detailed reference:** [functions/exttypes.md](functions/exttypes.md)

---

### extpath — Dot-Path Access

| JSONata function | Description |
|---|---|
| `$get(obj, path [, default])` | Read a nested value by dot-path |
| `$set(obj, path, value)` | Immutable write at dot-path |
| `$del(obj, path)` | Immutable delete at dot-path |
| `$has(obj, path)` | True if path exists and is non-nil |
| `$flattenObj(obj [, sep])` | `{"a":{"b":1}}` → `{"a.b":1}` |
| `$expandObj(obj [, sep])` | `{"a.b":1}` → `{"a":{"b":1}}` |

**Detailed reference:** [functions/extpath.md](functions/extpath.md)

---

### extvalidate — Input Validation

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

**Detailed reference:** [functions/extvalidate.md](functions/extvalidate.md)

---

### extjson — JSON Operations

| JSONata function | Description |
|---|---|
| `$jsonParse(str)` | Parse JSON string into a value |
| `$jsonStringify(v [, indent])` | Serialise value to JSON string |
| `$jsonDiff(a, b)` | Differences as a JSON Patch array |
| `$jsonPatch(obj, ops)` | Apply RFC 6902 JSON Patch operations |
| `$jsonPointer(obj, pointer)` | Resolve RFC 6901 JSON Pointer |

**Detailed reference:** [functions/extjson.md](functions/extjson.md)

---

### extgeo — Geospatial Utilities

All calculations use the WGS-84 mean Earth radius (6 371 km). Coordinates are decimal degrees.

| JSONata function | Description |
|---|---|
| `$haversine(lat1, lon1, lat2, lon2)` | Great-circle distance in km |
| `$bearing(lat1, lon1, lat2, lon2)` | Initial bearing in degrees (0–360) |
| `$geoFormat(lat, lon [, format])` | Format as `"decimal"` or `"dms"` |
| `$geoParse(str)` | Parse `"lat, lon"` string → `{lat, lon}` |
| `$inBoundingBox(lat, lon, minLat, minLon, maxLat, maxLon)` | Point-in-bbox test |
| `$geoDistance(point, points)` | Array of distances from point to each in points (km) |

**Detailed reference:** [functions/extgeo.md](functions/extgeo.md)

---

### extnet — Network Utilities

| JSONata function | Description |
|---|---|
| `$ipVersion(str)` | Returns `4`, `6`, or `-1` |
| `$isPrivateIP(str)` | RFC1918 / loopback / link-local |
| `$ipToInt(str)` | IPv4 → uint32 as float64 |
| `$intToIP(n)` | uint32 as float64 → IPv4 string |
| `$ipInCIDR(ip, cidr)` | True if ip is in CIDR block |
| `$expandCIDR(cidr)` | Network info object for CIDR block |

**Detailed reference:** [functions/extnet.md](functions/extnet.md)
