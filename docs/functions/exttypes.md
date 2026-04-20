# exttypes — Type Checking & Coercion

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/exttypes`

---

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

### `$toArray(v)`

Wraps *v* in a single-element array if it is not already an array. `null`/nil returns `[]`.

```jsonata
$toArray(1)         /* → [1] */
$toArray([1, 2])    /* → [1, 2] */
$toArray(null)      /* → [] */
```

---

### `$defined(v)`

Returns `true` if *v* is not `null`/nil. Use `$exists(path)` for path-based existence tests in standard JSONata.

```jsonata
$defined("hello")   /* → true */
$defined(null)      /* → false */
```

---

### `$nullish(v, default)`

Returns *v* if *v* is not `null`/nil, otherwise returns *default*. Equivalent to `$default`.

```jsonata
$nullish(null, "fallback")   /* → "fallback" */
$nullish(0, "fallback")      /* → 0 */
```

---

### `$typeOf(v)`

Returns the type of *v* as a string. Possible values: `"string"`, `"number"`, `"boolean"`, `"array"`, `"object"`, `"null"`.

> **JSONata native alternative:** `$type(v)` — identical for defined values. `$type` returns `undefined` for undefined values; `$typeOf` returns `"null"`. See [guides/jsonata-overlap.md](../guides/jsonata-overlap.md).

```jsonata
$typeOf("hello")    /* → "string" */
$typeOf([1, 2])     /* → "array" */
$typeOf(null)       /* → "null" */
```

---

### `$toNumber(v)`

Coerces *v* to a number. Booleans become 0/1; strings are parsed as decimal floats; arrays/objects return an error.

> **JSONata native alternative:** `$number(v)` — equivalent, plus supports hex/octal/binary string literals. Prefer the built-in.

```jsonata
$toNumber("3.14")   /* → 3.14 */
$toNumber(true)     /* → 1 */
$toNumber(false)    /* → 0 */
```

---

### `$toString(v)`

Coerces *v* to a string. Objects and arrays are JSON-serialised.

> **JSONata native alternative:** `$string(v)` — equivalent. Prefer the built-in.

```jsonata
$toString(42)          /* → "42" */
$toString(true)        /* → "true" */
$toString({"a": 1})    /* → "{\"a\":1}" */
```

---

### `$toBool(v)`

Coerces *v* to a boolean using JSONata casting rules.

> **JSONata native equivalent:** `$boolean(v)` — identical semantics. Prefer the built-in.

```jsonata
$toBool(1)        /* → true */
$toBool(0)        /* → false */
$toBool("")       /* → false */
$toBool("hello")  /* → true */
```

---
