# extarray — Array Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extarray`

---

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
