# extvalidate — Validation Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extvalidate`

All functions return `true`/`false`. On invalid input types (non-string when a string is expected, etc.) they return `false` rather than erroring, making them safe in conditional expressions.

> **Regex safety:** `$matchesRegex` uses Go's RE2-based `regexp` package, which guarantees linear-time matching and is inherently safe against ReDoS attacks.

---

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
$matchesRegex("hello123", "^\w+$")    /* → true */
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
