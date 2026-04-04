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
