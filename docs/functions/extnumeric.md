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

### `$product(array)`

Returns the arithmetic product of all elements in *array*.

```jsonata
$product([1, 2, 3, 4])   /* → 24 */
$product([2, 0.5])        /* → 1 */
```

---

### `$cumSum(array)`

Returns an array of cumulative sums.

```jsonata
$cumSum([1, 2, 3, 4])   /* → [1, 3, 6, 10] */
```

---

### `$inRange(n, min, max)`

Returns `true` if *min* ≤ *n* ≤ *max*.

```jsonata
$inRange(5, 1, 10)    /* → true */
$inRange(11, 1, 10)   /* → false */
```

---

### `$roundTo(n, precision)`

Rounds *n* to *precision* decimal places.

> **JSONata native equivalent:** `$round(n, precision)` — identical semantics. Prefer the built-in.

```jsonata
$roundTo(3.14159, 2)   /* → 3.14 */
$roundTo(123.5, -1)    /* → 120 */
```

---

### `$normalize(array)`

Applies min-max normalisation, scaling all values to the range [0, 1].

```jsonata
$normalize([1, 2, 3, 4, 5])   /* → [0, 0.25, 0.5, 0.75, 1] */
```

---

### `$interpolate(start, end, t)`

Returns the linearly interpolated value between *start* and *end* at position *t* (0 = start, 1 = end).

```jsonata
$interpolate(0, 100, 0.5)   /* → 50 */
$interpolate(10, 20, 0.25)  /* → 12.5 */
```

---

### `$gcd(a, b)`

Returns the greatest common divisor of integers *a* and *b*.

```jsonata
$gcd(12, 8)   /* → 4 */
$gcd(7, 5)    /* → 1 */
```

---

### `$lcm(a, b)`

Returns the least common multiple of integers *a* and *b*.

```jsonata
$lcm(4, 6)    /* → 12 */
$lcm(3, 7)    /* → 21 */
```

---

### `$isPrime(n)`

Returns `true` if *n* is a prime number. *n* must be a positive integer.

```jsonata
$isPrime(7)    /* → true */
$isPrime(8)    /* → false */
$isPrime(1)    /* → false */
```

---

### `$factorial(n)`

Returns *n*! as a float64. *n* must be a non-negative integer ≤ 20 (larger values overflow float64 exactly).

```jsonata
$factorial(5)    /* → 120 */
$factorial(0)    /* → 1 */
```

---
