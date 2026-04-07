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
