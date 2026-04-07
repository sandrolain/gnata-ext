# extobject — Object Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extobject`

---

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
