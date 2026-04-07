# Function Reference

Complete documentation for every function provided by `gnata-ext`.

All functions are registered as JSONata custom functions and are called with a leading `$` in expressions (e.g. `$uuid()`). Timestamps are always **Unix milliseconds** (`float64`) unless stated otherwise.

---

## Table of Contents

- [extarray — Array utilities](#extarray)
- [extcrypto — Cryptography](#extcrypto)
- [extdatetime — Date & time](#extdatetime)
- [extformat — Formatting](#extformat)
- [extnumeric — Extended math](#extnumeric)
- [extobject — Object utilities](#extobject)
- [extstring — String utilities](#extstring)
- [exttypes — Type inspection](#exttypes)
- [extpath — Dot-path access](#extpath)
- [extvalidate — Input validation](#extvalidate)
- [extjson — JSON operations](#extjson)
- [extgeo — Geospatial utilities](#extgeo)
- [extnet — Network utilities](#extnet)

---

## extarray

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extarray`

### `$first(array)`

Returns the first element of an array, or `null` if the array is empty.

```jsonata
$first([10, 20, 30])          /* → 10 */
$first([])                    /* → null */
```

---

### `$last(array)`

Returns the last element of an array, or `null` if the array is empty.

```jsonata
$last([10, 20, 30])           /* → 30 */
```

---

### `$take(array, n)`

Returns the first *n* elements. If *n* exceeds the length, the full array is returned.

```jsonata
$take([1, 2, 3, 4, 5], 3)    /* → [1, 2, 3] */
$take([1, 2], 10)             /* → [1, 2] */
```

---

### `$skip(array, n)`

Returns all elements after the first *n*. If *n* exceeds the length, returns `[]`.

```jsonata
$skip([1, 2, 3, 4, 5], 2)    /* → [3, 4, 5] */
```

---

### `$slice(array, start [, end])`

Returns a sub-array from index *start* (inclusive) to *end* (exclusive). Negative indices count from the end. *end* defaults to the array length.

```jsonata
$slice([1, 2, 3, 4, 5], 1, 4)   /* → [2, 3, 4] */
$slice([1, 2, 3, 4, 5], -2)     /* → [4, 5] */
```

---

### `$flatten(array [, depth])`

Recursively flattens nested arrays. *depth* limits how many levels to descend; omit or pass `null` for unlimited flattening.

```jsonata
$flatten([[1, 2], [3, [4, 5]]])     /* → [1, 2, 3, 4, 5] */
$flatten([[1, [2, [3]]], 4], 1)     /* → [1, 2, [3], 4] */
```

---

### `$chunk(array, size)`

Splits *array* into consecutive sub-arrays of length *size*. The last chunk may be shorter.

```jsonata
$chunk([1, 2, 3, 4, 5], 2)   /* → [[1, 2], [3, 4], [5]] */
```

---

### `$union(a, b)`

Returns a de-duplicated array containing all elements from both *a* and *b*, preserving order of first occurrence.

```jsonata
$union([1, 2, 3], [2, 3, 4])  /* → [1, 2, 3, 4] */
```

---

### `$intersection(a, b)`

Returns elements present in **both** *a* and *b*.

```jsonata
$intersection([1, 2, 3], [2, 3, 4])   /* → [2, 3] */
```

---

### `$difference(a, b)`

Returns elements in *a* that are **not** in *b*.

```jsonata
$difference([1, 2, 3], [2])   /* → [1, 3] */
```

---

### `$symmetricDifference(a, b)`

Returns elements that are in *either* array but **not in both**.

```jsonata
$symmetricDifference([1, 2, 3], [2, 4])   /* → [1, 3, 4] */
```

---

### `$range(start, end [, step])`

Generates a numeric range from *start* (inclusive) to *end* (exclusive), advancing by *step* (default `1`). Negative steps are supported for descending ranges. Maximum 100,000 items.

```jsonata
$range(1, 5)           /* → [1, 2, 3, 4] */
$range(0, 1, 0.25)     /* → [0, 0.25, 0.5, 0.75] */
$range(5, 0, -1)       /* → [5, 4, 3, 2, 1] */
```

---

### `$zipLongest(a, b [, fill])`

Pairs elements of *a* and *b* positionally. When arrays differ in length, the shorter one is padded with *fill* (default `null`).

```jsonata
$zipLongest([1, 2], [3, 4, 5], 0)
/* → [[1, 3], [2, 4], [0, 5]] */
```

---

### `$window(array, size [, step])`

Produces an array of sliding windows of the given *size*. *step* (default `1`) controls how many positions to advance between windows.

```jsonata
$window([1, 2, 3, 4], 2)        /* → [[1, 2], [2, 3], [3, 4]] */
$window([1, 2, 3, 4], 2, 2)     /* → [[1, 2], [3, 4]] */
```

---

## extcrypto

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extcrypto`

### `$uuid()`

Generates a random UUID v4.

```jsonata
$uuid()   /* → "550e8400-e29b-41d4-a716-446655440000" */
```

---

### `$hash(algorithm, value)`

Returns a lower-case hex-encoded digest of *value*.

Supported algorithms: `md5`, `sha1`, `sha256`, `sha384`, `sha512`.

```jsonata
$hash("sha256", "hello")
/* → "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" */

$hash("md5", "gnata-ext")
```

---

### `$hmac(algorithm, key, value)`

Returns a lower-case hex-encoded HMAC of *value* signed with *key*.

Supported algorithms: `md5`, `sha1`, `sha256`, `sha384`, `sha512`.

```jsonata
$hmac("sha256", "secret", "message")
```

---

## extdatetime

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extdatetime`

All functions accept and return timestamps as **Unix milliseconds** (`float64`), compatible with JSONata's built-in `$millis()` and `$toMillis()`.

Accepted unit strings (singular and plural forms both work):
`year` · `month` · `day` · `hour` · `minute` · `second` · `millisecond`

---

### `$dateAdd(timestamp, amount, unit)`

Adds *amount* units to a timestamp. Use a negative *amount* to subtract.

```jsonata
/* 2024-01-15T12:30:00Z = 1705319400000 ms */
$dateAdd(1705319400000, 1, "day")     /* +1 day  */
$dateAdd(1705319400000, -2, "month")  /* -2 months */
$dateAdd(1705319400000, 3, "hour")
```

---

### `$dateDiff(t1, t2, unit)`

Returns the whole-unit difference `t2 − t1`.

```jsonata
$dateDiff(1705319400000, 1705405800000, "hour")   /* → 24 */
$dateDiff(1705319400000, 1705405800000, "day")    /* → 1 */
```

---

### `$dateComponents(timestamp)`

Returns an object with the UTC components of the timestamp.

```jsonata
$dateComponents(1705319400000)
/* → {
     "year": 2024, "month": 1, "day": 15,
     "hour": 12, "minute": 30, "second": 0,
     "millisecond": 0, "weekday": 1
   } */
```

`weekday` follows Go's `time.Weekday`: Sunday=0, Monday=1, …, Saturday=6.

---

### `$dateStartOf(timestamp, unit)`

Returns the timestamp at the start of the given period.

```jsonata
$dateStartOf(1705319400000, "day")    /* 2024-01-15T00:00:00.000Z */
$dateStartOf(1705319400000, "month")  /* 2024-01-01T00:00:00.000Z */
$dateStartOf(1705319400000, "year")   /* 2024-01-01T00:00:00.000Z */
```

---

### `$dateEndOf(timestamp, unit)`

Returns the timestamp at the last millisecond of the given period.

```jsonata
$dateEndOf(1705319400000, "day")      /* 2024-01-15T23:59:59.999Z */
$dateEndOf(1705319400000, "month")    /* 2024-01-31T23:59:59.999Z */
```

---

## extformat

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extformat`

### `$csv(text)`

Parses a CSV string into an array of objects. The first row is the header; each subsequent row becomes an object with header keys.

```jsonata
$csv("name,age\nAlice,30\nBob,25")
/* → [
     {"name": "Alice", "age": "30"},
     {"name": "Bob",   "age": "25"}
   ] */
```

All values are returned as strings, consistent with standard CSV semantics.

---

### `$toCSV(array)`

Serializes an array of objects to CSV text. The header row is derived from the keys of the **first** object; missing keys in subsequent rows are written as empty strings.

```jsonata
$toCSV([
  {"name": "Alice", "age": "30"},
  {"name": "Bob",   "age": "25"}
])
/* → "name,age\nAlice,30\nBob,25\n" */
```

---

### `$template(str, vars)`

Replaces `{{key}}` placeholders in *str* with values from the *vars* object. Unknown placeholders are left unchanged.

```jsonata
$template("Hello, {{name}}! You are {{age}}.", {"name": "Alice", "age": 30})
/* → "Hello, Alice! You are 30." */

$template("{{greeting}}, world!", {})
/* → "{{greeting}}, world!" */
```

> **Note:** `extstring` also exposes `$template` with identical behaviour. When both packages are registered, the last one wins; register only the package you need if this matters.

---

## extnumeric

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extnumeric`

### `$log(n [, base])`

Returns the logarithm of *n*. Without *base*, returns the natural logarithm. *n* must be positive; *base* must be positive and ≠ 1.

```jsonata
$log(1)         /* → 0 */
$log(100, 10)   /* → 2 */
$log($e())      /* → 1 */
```

---

### `$sign(n)`

Returns `−1`, `0`, or `1` depending on the sign of *n*.

```jsonata
$sign(-42)  /* → -1 */
$sign(0)    /* →  0 */
$sign(7)    /* →  1 */
```

---

### `$trunc(n)`

Truncates *n* toward zero (removes the fractional part without rounding).

```jsonata
$trunc(3.9)    /* → 3 */
$trunc(-3.9)   /* → -3 */
```

---

### `$clamp(n, min, max)`

Returns *n* clamped to the inclusive range [*min*, *max*].

```jsonata
$clamp(5,  1, 10)   /* → 5 */
$clamp(-1, 1, 10)   /* → 1 */
$clamp(99, 1, 10)   /* → 10 */
```

---

### Trigonometric functions

All accept/return values in **radians**.

| Function | Description |
|---|---|
| `$sin(n)` | Sine |
| `$cos(n)` | Cosine |
| `$tan(n)` | Tangent |
| `$asin(n)` | Arc-sine |
| `$acos(n)` | Arc-cosine |
| `$atan(n)` | Arc-tangent |
| `$atan2(y, x)` | `atan(y/x)` with correct quadrant |

```jsonata
$sin(0)          /* → 0 */
$cos(0)          /* → 1 */
$atan2(1, 1)     /* → π/4 ≈ 0.7854 */
```

---

### `$pi()`

Returns π (≈ 3.141592653589793).

```jsonata
$pi()   /* → 3.141592653589793 */
```

---

### `$e()`

Returns Euler's number *e* (≈ 2.718281828459045).

```jsonata
$e()    /* → 2.718281828459045 */
```

---

### `$median(array)`

Returns the median value of a numeric array. For even-length arrays, returns the average of the two middle values.

```jsonata
$median([1, 2, 3])        /* → 2 */
$median([1, 2, 3, 4])     /* → 2.5 */
```

---

### `$variance(array)`

Returns the population variance of a numeric array.

```jsonata
$variance([2, 4, 4, 4, 5, 5, 7, 9])   /* → 4 */
```

---

### `$stddev(array)`

Returns the population standard deviation (√variance).

```jsonata
$stddev([2, 4, 4, 4, 5, 5, 7, 9])   /* → 2 */
```

---

### `$percentile(array, p)`

Returns the *p*-th percentile of *array*, where *p* is in [0, 100]. Uses linear interpolation between adjacent values.

```jsonata
$percentile([1, 2, 3, 4, 5], 50)   /* → 3 */
$percentile([1, 2, 3, 4, 5], 75)   /* → 4 */
```

---

### `$mode(array)`

Returns the most frequently occurring value. When multiple values share the highest frequency, returns an array of all of them.

```jsonata
$mode([1, 2, 2, 3])        /* → 2 */
$mode([1, 1, 2, 2, 3])     /* → [1, 2] */
```

---

## extobject

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extobject`

### `$values(object)`

Returns an array of the object's own values.

```jsonata
$values({"a": 1, "b": 2})   /* → [1, 2] */
```

---

### `$pairs(object)`

Returns an array of `[key, value]` pairs.

```jsonata
$pairs({"x": 10, "y": 20})
/* → [["x", 10], ["y", 20]] */
```

---

### `$fromPairs(pairs)`

Constructs an object from an array of `[key, value]` pairs. Also accepts objects with `"key"` and `"value"` fields.

```jsonata
$fromPairs([["a", 1], ["b", 2]])
/* → {"a": 1, "b": 2} */
```

---

### `$pick(object, keys)`

Returns a new object containing only the listed keys.

```jsonata
$pick({"a": 1, "b": 2, "c": 3}, ["a", "c"])
/* → {"a": 1, "c": 3} */
```

---

### `$omit(object, keys)`

Returns a new object with the listed keys removed.

```jsonata
$omit({"a": 1, "b": 2, "c": 3}, ["b"])
/* → {"a": 1, "c": 3} */
```

---

### `$deepMerge(target, source)`

Returns a new object with *source* merged recursively into *target*. Nested objects are merged; all other values in *source* overwrite *target*.

```jsonata
$deepMerge(
  {"a": {"x": 1, "y": 2}, "b": "keep"},
  {"a": {"y": 99, "z": 3}, "c": "new"}
)
/* → {"a": {"x": 1, "y": 99, "z": 3}, "b": "keep", "c": "new"} */
```

---

### `$invert(object)`

Swaps keys and values. All values must be strings.

```jsonata
$invert({"a": "alpha", "b": "beta"})
/* → {"alpha": "a", "beta": "b"} */
```

---

### `$size(object)`

Returns the number of own keys.

```jsonata
$size({"a": 1, "b": 2, "c": 3})   /* → 3 */
```

---

### `$rename(object, oldKey, newKey)`

Returns a copy of *object* with *oldKey* renamed to *newKey*. If *oldKey* does not exist, the object is returned unchanged.

```jsonata
$rename({"firstName": "Alice", "age": 30}, "firstName", "name")
/* → {"name": "Alice", "age": 30} */
```

---

## extstring

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extstring`

### `$startsWith(str, prefix)`

Returns `true` if *str* starts with *prefix*.

```jsonata
$startsWith("hello world", "hello")   /* → true */
$startsWith("hello world", "world")   /* → false */
```

---

### `$endsWith(str, suffix)`

Returns `true` if *str* ends with *suffix*.

```jsonata
$endsWith("hello world", "world")   /* → true */
```

---

### `$indexOf(str, search [, start])`

Returns the first index of *search* in *str*, or `−1` if not found. Optional *start* offsets where to begin the search.

```jsonata
$indexOf("abcabc", "b")      /* → 1 */
$indexOf("abcabc", "b", 2)   /* → 4 */
$indexOf("abcabc", "x")      /* → -1 */
```

---

### `$lastIndexOf(str, search)`

Returns the last index of *search* in *str*, or `−1` if not found.

```jsonata
$lastIndexOf("abcabc", "b")   /* → 4 */
```

---

### `$capitalize(str)`

Uppercases the first character and lowercases the rest.

```jsonata
$capitalize("hELLO")   /* → "Hello" */
```

---

### `$titleCase(str)`

Title-cases every word (first letter uppercase, rest lowercase).

```jsonata
$titleCase("hello world foo")   /* → "Hello World Foo" */
```

---

### `$camelCase(str)`

Converts to camelCase. Splits on whitespace, hyphens, underscores, and case boundaries.

```jsonata
$camelCase("hello world")    /* → "helloWorld" */
$camelCase("foo-bar-baz")    /* → "fooBarBaz" */
$camelCase("my_variable")    /* → "myVariable" */
```

---

### `$snakeCase(str)`

Converts to snake_case.

```jsonata
$snakeCase("helloWorld")     /* → "hello_world" */
$snakeCase("Hello World")    /* → "hello_world" */
```

---

### `$kebabCase(str)`

Converts to kebab-case.

```jsonata
$kebabCase("helloWorld")     /* → "hello-world" */
$kebabCase("Hello World")    /* → "hello-world" */
```

---

### `$repeat(str, n)`

Returns *str* repeated *n* times.

```jsonata
$repeat("ab", 3)    /* → "ababab" */
$repeat("x", 0)     /* → "" */
```

---

### `$words(str)`

Splits *str* into an array of words, splitting on whitespace, hyphens, underscores, and camelCase boundaries.

```jsonata
$words("hello world")      /* → ["hello", "world"] */
$words("camelCaseWord")    /* → ["camel", "Case", "Word"] */
$words("foo-bar_baz")      /* → ["foo", "bar", "baz"] */
```

---

### `$template(str, vars)`

Replaces `{{key}}` placeholders in *str* with values from the *vars* object. See also [`extformat.$template`](#template).

```jsonata
$template("Hi {{name}}", {"name": "Bob"})   /* → "Hi Bob" */
```

---

## exttypes

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/exttypes`

### `$isString(v)`

```jsonata
$isString("hello")   /* → true */
$isString(42)        /* → false */
```

---

### `$isNumber(v)`

Accepts `float64`, `int`, and `int64`.

```jsonata
$isNumber(42)        /* → true */
$isNumber("42")      /* → false */
```

---

### `$isBoolean(v)`

```jsonata
$isBoolean(true)     /* → true */
$isBoolean(1)        /* → false */
```

---

### `$isArray(v)`

```jsonata
$isArray([1, 2])     /* → true */
$isArray("abc")      /* → false */
```

---

### `$isObject(v)`

```jsonata
$isObject({"k": 1})  /* → true */
$isObject([1, 2])    /* → false */
```

---

### `$isNull(v)`

```jsonata
$isNull(null)        /* → true */
$isNull("")          /* → false */
```

---

### `$isUndefined(v)`

In gnata, `undefined` is represented as `nil`, which is indistinguishable from `null`. This function behaves identically to `$isNull`.

```jsonata
$isUndefined(null)   /* → true */
```

---

### `$isEmpty(v)`

Returns `true` for `null`, the empty string `""`, an empty array `[]`, or an empty object `{}`.

```jsonata
$isEmpty(null)       /* → true */
$isEmpty("")         /* → true */
$isEmpty([])         /* → true */
$isEmpty({})         /* → true */
$isEmpty("x")        /* → false */
```

---

### `$default(v, d)`

Returns *v* if it is non-null; returns *d* otherwise.

```jsonata
$default(null, "fallback")   /* → "fallback" */
$default("value", "other")   /* → "value" */
```

---

### `$identity(v)`

Returns *v* unchanged. Useful as a no-op placeholder in pipelines.

```jsonata
$identity(42)          /* → 42 */
$identity({"k": "v"})  /* → {"k": "v"} */
```

---

## extpath

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extpath`

Provides immutable read/write access to nested objects via dot-path strings (e.g. `"a.b.c"`).

### `$get(obj, path [, default])`

Reads a nested value by dot-path. Returns *default* (or `null`) if any segment is missing.

```jsonata
$get({"a": {"b": 42}}, "a.b")          /* → 42 */
$get({"a": 1}, "a.b.c", "missing")     /* → "missing" */
```

---

### `$set(obj, path, value)`

Returns a new object with *value* written at *path*. Intermediate objects are created as needed. The original object is not mutated.

```jsonata
$set({"a": 1}, "b.c", 2)
/* → {"a": 1, "b": {"c": 2}} */

$set({"a": {"x": 1}}, "a.y", 2)
/* → {"a": {"x": 1, "y": 2}} */
```

---

### `$del(obj, path)`

Returns a new object with the value at *path* removed. The original is not mutated.

```jsonata
$del({"a": 1, "b": 2}, "b")       /* → {"a": 1} */
$del({"a": {"b": 1, "c": 2}}, "a.b")
/* → {"a": {"c": 2}} */
```

---

### `$has(obj, path)`

Returns `true` if the path exists in *obj* and its value is non-null.

```jsonata
$has({"a": {"b": 1}}, "a.b")   /* → true */
$has({"a": 1}, "a.b")          /* → false */
```

---

### `$flattenObj(obj [, sep])`

Flattens a nested object to a single level, joining keys with *sep* (default `"."`).

```jsonata
$flattenObj({"a": {"b": {"c": 1}, "d": 2}})
/* → {"a.b.c": 1, "a.d": 2} */

$flattenObj({"a": {"b": 1}}, "/")
/* → {"a/b": 1} */
```

---

### `$expandObj(obj [, sep])`

Expands a flat object with compound keys back to a nested structure. Inverse of `$flattenObj`.

```jsonata
$expandObj({"a.b.c": 1, "a.d": 2})
/* → {"a": {"b": {"c": 1}, "d": 2}} */
```

---

## extvalidate

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extvalidate`

All functions return `true`/`false`. On invalid input types (non-string when a string is expected, etc.) they return `false` rather than erroring, making them safe in conditional expressions.

> **Regex safety:** `$matchesRegex` uses Go's RE2-based `regexp` package, which guarantees linear-time matching and is inherently safe against ReDoS attacks.

### `$isEmail(str)`

Returns `true` if *str* matches a simplified RFC 5322 email format.

```jsonata
$isEmail("user@example.com")          /* → true */
$isEmail("user+tag@sub.example.co")   /* → true */
$isEmail("notanemail")                /* → false */
```

---

### `$isURL(str)`

Returns `true` if *str* is a valid URL with scheme `http`, `https`, or `ftp`.

```jsonata
$isURL("https://example.com/path?q=1")   /* → true */
$isURL("ftp://files.example.com")        /* → true */
$isURL("not-a-url")                      /* → false */
```

---

### `$isUUID(str)`

Returns `true` if *str* matches a UUID v1–v5 format (case-insensitive).

```jsonata
$isUUID("550e8400-e29b-41d4-a716-446655440000")   /* → true */
$isUUID("not-a-uuid")                             /* → false */
```

---

### `$isIPv4(str)`

```jsonata
$isIPv4("192.168.1.1")   /* → true */
$isIPv4("::1")           /* → false */
```

---

### `$isIPv6(str)`

```jsonata
$isIPv6("2001:db8::1")     /* → true */
$isIPv6("192.168.1.1")     /* → false */
```

---

### `$isAlpha(str)`

Returns `true` if *str* contains only Unicode letters (empty string → `false`).

```jsonata
$isAlpha("Hello")      /* → true */
$isAlpha("hello123")   /* → false */
```

---

### `$isAlphanumeric(str)`

Returns `true` if *str* contains only Unicode letters and digits.

```jsonata
$isAlphanumeric("hello123")    /* → true */
$isAlphanumeric("hello 123")   /* → false */
```

---

### `$isNumericStr(str)`

Returns `true` if *str* can be parsed as a number.

```jsonata
$isNumericStr("3.14")     /* → true */
$isNumericStr("-1.5e10")  /* → true */
$isNumericStr("abc")      /* → false */
```

---

### `$matchesRegex(str, pattern)`

Returns `true` if *str* fully or partially matches the RE2 *pattern*.

```jsonata
$matchesRegex("hello123", "^\w+$")   /* → true */
$matchesRegex("hello world", "^\w+$") /* → false */
```

---

### `$inSet(v, set)`

Returns `true` if *v* is present in the array *set* (strict equality).

```jsonata
$inSet("b", ["a", "b", "c"])   /* → true */
$inSet("z", ["a", "b", "c"])   /* → false */
```

---

### `$minLen(str, n)`

Returns `true` if the rune (Unicode character) length of *str* is ≥ *n*.

```jsonata
$minLen("hello", 3)   /* → true */
$minLen("hi", 5)      /* → false */
```

---

### `$maxLen(str, n)`

Returns `true` if the rune length of *str* is ≤ *n*.

```jsonata
$maxLen("hi", 5)       /* → true */
$maxLen("toolong", 3)  /* → false */
```

---

### `$minItems(arr, n)`

Returns `true` if the array length of *arr* is ≥ *n*.

```jsonata
$minItems([1, 2, 3], 2)   /* → true */
$minItems([1], 2)          /* → false */
```

---

### `$maxItems(arr, n)`

Returns `true` if the array length of *arr* is ≤ *n*.

```jsonata
$maxItems([1, 2], 5)    /* → true */
$maxItems([1, 2, 3, 4, 5, 6], 5)   /* → false */
```

---

## extjson

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extjson`

### `$jsonParse(str)`

Parses a JSON string and returns the corresponding gnata value. Errors on invalid JSON.

```jsonata
$jsonParse("{\"a\":1,\"b\":true}")
/* → {"a": 1, "b": true} */

$jsonParse("[1, 2, 3]")
/* → [1, 2, 3] */
```

---

### `$jsonStringify(v [, indent])`

Serialises *v* to a JSON string. If *indent* is provided (a string, e.g. `"  "`), the output is pretty-printed.

```jsonata
$jsonStringify({"x": 1})          /* → "{\"x\":1}" */
$jsonStringify({"x": 1}, "  ")    /* pretty-printed */
```

---

### `$jsonDiff(a, b)`

Returns a JSON Patch-compatible array of `{op, path, value}` operations describing what changed between *a* and *b*. Uses RFC 6902 operation names (`add`, `remove`, `replace`). Paths use JSON Pointer notation.

```jsonata
$jsonDiff({"x": 1, "y": 2}, {"x": 1, "y": 3, "z": 4})
/* → [
     {"op": "replace", "path": "/y", "value": 3},
     {"op": "add",     "path": "/z", "value": 4}
   ] */
```

---

### `$jsonPatch(obj, ops)`

Applies RFC 6902 JSON Patch operations to *obj* (immutable — returns a new value). Supported operations: `add`, `remove`, `replace`, `move`, `copy`, `test`.

```jsonata
$jsonPatch(
  {"a": 1, "b": 2},
  [
    {"op": "replace", "path": "/a", "value": 10},
    {"op": "remove",  "path": "/b"},
    {"op": "add",     "path": "/c", "value": 3}
  ]
)
/* → {"a": 10, "c": 3} */
```

The `"test"` operation throws an error if the path value does not match, which causes the entire patch to fail.

---

### `$jsonPointer(obj, pointer)`

Resolves an RFC 6901 JSON Pointer against *obj*. Supports traversal of both objects and arrays.

```jsonata
$jsonPointer({"a": {"b": [10, 20, 30]}}, "/a/b/1")   /* → 20 */
$jsonPointer({"a": 1}, "/a")                          /* → 1 */
```

---

## extgeo

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extgeo`

All calculations use the WGS-84 mean Earth radius **6 371 km**. Coordinates are always **decimal degrees**. No external dependencies — pure `math` stdlib.

### `$haversine(lat1, lon1, lat2, lon2)`

Returns the great-circle distance between two points in **kilometres**.

```jsonata
/* London → Paris ≈ 340 km */
$haversine(51.5074, -0.1278, 48.8566, 2.3522)
/* → ~340.6 */

/* Same point */
$haversine(48.0, 2.0, 48.0, 2.0)   /* → 0 */
```

---

### `$bearing(lat1, lon1, lat2, lon2)`

Returns the initial bearing from point 1 to point 2 in **degrees**, clockwise from north (0–360).

```jsonata
$bearing(0, 0, 0, 1)    /* → ~90  (east) */
$bearing(0, 0, 1, 0)    /* → ~0   (north) */
```

---

### `$geoFormat(lat, lon [, format])`

Formats a coordinate pair as a string.

- `"decimal"` (default): `"48.8566, 2.3522"` — 4 decimal places
- `"dms"`: degrees/minutes/seconds with cardinal directions, e.g. `"48°51'23.76\"N 2°21'7.92\"E"`

```jsonata
$geoFormat(48.8566, 2.3522)           /* → "48.8566, 2.3522" */
$geoFormat(48.8566, 2.3522, "dms")    /* → "48°51'23.76\"N 2°21'7.92\"E" */
```

---

### `$geoParse(str)`

Parses a `"lat, lon"` decimal string and returns `{lat, lon}`.

```jsonata
$geoParse("48.8566, 2.3522")
/* → {"lat": 48.8566, "lon": 2.3522} */
```

---

### `$inBoundingBox(lat, lon, minLat, minLon, maxLat, maxLon)`

Returns `true` if the point (`lat`, `lon`) lies within the axis-aligned bounding box.

```jsonata
/* Paris inside Europe bbox */
$inBoundingBox(48.8566, 2.3522, 36.0, -10.0, 71.0, 40.0)   /* → true */
```

---

### `$geoDistance(point, points)`

Computes the haversine distance from *point* to each element of *points* and returns an array of distances in **kilometres**.

- *point*: `{"lat": float, "lon": float}`
- *points*: array of `{"lat": float, "lon": float}`

```jsonata
$geoDistance(
  {"lat": 51.5074, "lon": -0.1278},
  [
    {"lat": 48.8566, "lon": 2.3522},
    {"lat": 52.5200, "lon": 13.4050}
  ]
)
/* → [~340.6, ~930.9] */
```

---

## extnet

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extnet`

### `$ipVersion(str)`

Returns `4` for IPv4, `6` for IPv6, or `-1` if *str* is not a valid IP address.

```jsonata
$ipVersion("192.168.1.1")   /* → 4 */
$ipVersion("::1")           /* → 6 */
$ipVersion("not-an-ip")     /* → -1 */
```

---

### `$isPrivateIP(str)`

Returns `true` if *str* is a private, loopback, or link-local address:

- RFC 1918: `10/8`, `172.16/12`, `192.168/16`
- Loopback: `127/8`, `::1`
- Link-local: `169.254/16`, `fe80::/10`

```jsonata
$isPrivateIP("192.168.1.1")   /* → true */
$isPrivateIP("8.8.8.8")       /* → false */
$isPrivateIP("::1")           /* → true */
```

---

### `$ipToInt(str)`

Converts an IPv4 address string to its 32-bit unsigned integer representation (returned as `float64`). Only IPv4 is supported.

```jsonata
$ipToInt("0.0.0.0")           /* → 0 */
$ipToInt("255.255.255.255")   /* → 4294967295 */
$ipToInt("192.168.1.1")       /* → 3232235777 */
```

---

### `$intToIP(n)`

Converts a 32-bit unsigned integer (as `float64`) back to an IPv4 address string.

```jsonata
$intToIP(0)           /* → "0.0.0.0" */
$intToIP(3232235777)  /* → "192.168.1.1" */
```

---

### `$ipInCIDR(ip, cidr)`

Returns `true` if *ip* is contained within the *cidr* block.

```jsonata
$ipInCIDR("192.168.1.100", "192.168.1.0/24")   /* → true */
$ipInCIDR("10.0.0.1",      "192.168.1.0/24")   /* → false */
```

---

### `$expandCIDR(cidr)`

Returns a network-info object for the CIDR block.

**IPv4** — returns `{network, broadcast, first, last, count}`:

```jsonata
$expandCIDR("192.168.1.0/24")
/* → {
     "network":   "192.168.1.0",
     "broadcast": "192.168.1.255",
     "first":     "192.168.1.1",
     "last":      "192.168.1.254",
     "count":     256
   } */
```

**IPv6** — returns `{network, first, last}` (no broadcast, count omitted due to scale):

```jsonata
$expandCIDR("2001:db8::/32")
/* → {
     "network": "2001:db8::",
     "first":   "2001:db8::",
     "last":    "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff"
   } */
```
