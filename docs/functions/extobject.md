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

### `$clean(object)`

Recursively removes all keys whose value is `null`/nil from *object* and any nested objects.

```jsonata
$clean({"a": 1, "b": null, "c": {"d": null, "e": 2}})
/* → {"a": 1, "c": {"e": 2}} */
```

---

### `$defaults(object, defs)`

Returns a copy of *object* with any missing keys filled in from *defs*. Existing keys are **not** overwritten.

```jsonata
$defaults({"a": 1}, {"a": 99, "b": 2, "c": 3})
/* → {"a": 1, "b": 2, "c": 3} */
```

---

### `$transform(object, keyMap)`

Renames multiple keys in *object* according to the `{oldKey: newKey}` mapping *keyMap*. Keys not present in *keyMap* are kept unchanged.

```jsonata
$transform({"firstName": "Alice", "lastName": "Smith"}, {"firstName": "first", "lastName": "last"})
/* → {"first": "Alice", "last": "Smith"} */
```

---

### `$filterKeys(object, pattern)`

Returns a copy of *object* keeping only keys that match the regular expression *pattern*.

```jsonata
$filterKeys({"foo_a": 1, "foo_b": 2, "bar": 3}, "^foo_")
/* → {"foo_a": 1, "foo_b": 2} */
```

---

### `$groupByValue(object)`

Inverts *object*, grouping together keys that share the same value. The result maps each distinct value (as a string) to an array of the original keys.

```jsonata
$groupByValue({"a": "x", "b": "y", "c": "x"})
/* → {"x": ["a", "c"], "y": ["b"]} */
```

---
