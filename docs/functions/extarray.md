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

### `$compact(array)`

Returns a new array with all falsy values removed. Falsy values are: `null`, `false`, `0`, and `""`.

```jsonata
$compact([1, null, 2, false, 0, "", 3])   /* → [1, 2, 3] */
```

---

### `$groupByKey(array, key)`

Groups an array of objects by the value of *key*. Returns an object where each key is a distinct value of *key* and the value is an array of matching objects.

```jsonata
$groupByKey([{"type":"a","v":1},{"type":"b","v":2},{"type":"a","v":3}], "type")
/* → {"a":[{"type":"a","v":1},{"type":"a","v":3}], "b":[{"type":"b","v":2}]} */
```

---

### `$sortBy(array, key)`

Returns a new array sorted in ascending order by the value of *key*. Uses a stable sort. Numeric values are compared numerically; all others are compared as strings.

> **JSONata native alternative:** `$sort(array, function($l, $r){ $l.key > $r.key })` — more flexible but requires a lambda.

```jsonata
$sortBy([{"name":"c"},{"name":"a"},{"name":"b"}], "name")
/* → [{"name":"a"},{"name":"b"},{"name":"c"}] */
```

---

### `$uniqueBy(array, key)`

Returns an array with duplicates removed, keeping the first occurrence for each value of *key*.

```jsonata
$uniqueBy([{"id":1,"v":"a"},{"id":2,"v":"b"},{"id":1,"v":"c"}], "id")
/* → [{"id":1,"v":"a"},{"id":2,"v":"b"}] */
```

---

### `$sumByKey(array, key)`

Returns the sum of the numeric values at *key* across all objects in *array*.

```jsonata
$sumByKey([{"price":10},{"price":20},{"price":5}], "price")   /* → 35 */
```

---

### `$countByKey(array, key)`

Returns an object counting how many times each distinct value of *key* appears.

```jsonata
$countByKey([{"type":"a"},{"type":"b"},{"type":"a"}], "type")
/* → {"a": 2, "b": 1} */
```

---

### `$rotate(array, n)`

Rotates *array* right by *n* positions. A negative *n* rotates left.

```jsonata
$rotate([1, 2, 3, 4, 5], 2)    /* → [4, 5, 1, 2, 3] */
$rotate([1, 2, 3, 4, 5], -1)   /* → [2, 3, 4, 5, 1] */
```

---

### `$indexof(array, value)`

Returns the first index of *value* in *array* using deep equality, or `−1` if not found.

```jsonata
$indexof([1, 2, 3], 2)                         /* → 1 */
$indexof([{"a":1},{"a":2}], {"a":2})           /* → 1 */
$indexof([1, 2, 3], 99)                        /* → -1 */
```

---

### `$transpose(matrix)`

Transposes a 2-D array (matrix). The *i*-th row of the result contains the *i*-th element of each original row.

```jsonata
$transpose([[1, 2, 3], [4, 5, 6]])
/* → [[1, 4], [2, 5], [3, 6]] */
```

---

### `$adjacentPairs(array)`

Returns an array of consecutive `[a, b]` pairs from *array*. An array of length *n* produces *n − 1* pairs.

```jsonata
$adjacentPairs([1, 2, 3, 4])   /* → [[1, 2], [2, 3], [3, 4]] */
```

---
