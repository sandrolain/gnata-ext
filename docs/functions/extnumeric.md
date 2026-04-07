# extnumeric — Numeric Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extnumeric`

---

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
