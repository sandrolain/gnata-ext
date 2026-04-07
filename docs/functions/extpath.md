# extpath — Nested Object Access

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extpath`

Provides immutable read/write access to nested objects via dot-path strings (e.g. `"a.b.c"`).

---

### `$get(obj, path [, default])`

Reads a nested value by dot-path. Returns *default* (or `null`) if any segment is missing.

```jsonata
$get({"a": {"b": 42}}, "a.b")          /* → 42 */
$get({"a": 1}, "a.b.c", "missing")     /* → "missing" */
```

---

### `$set(obj, path, value)`

Returns a new object with *value* written at *path*. Intermediate objects are created as needed. The original object is not mutated.

```jsonata
$set({"a": 1}, "b.c", 2)
/* → {"a": 1, "b": {"c": 2}} */

$set({"a": {"x": 1}}, "a.y", 2)
/* → {"a": {"x": 1, "y": 2}} */
```

---

### `$del(obj, path)`

Returns a new object with the value at *path* removed. The original is not mutated.

```jsonata
$del({"a": 1, "b": 2}, "b")       /* → {"a": 1} */
$del({"a": {"b": 1, "c": 2}}, "a.b")
/* → {"a": {"c": 2}} */
```

---

### `$has(obj, path)`

Returns `true` if the path exists in *obj* and its value is non-null.

```jsonata
$has({"a": {"b": 1}}, "a.b")   /* → true */
$has({"a": 1}, "a.b")          /* → false */
```

---

### `$flattenObj(obj [, sep])`

Flattens a nested object to a single level, joining keys with *sep* (default `"."`).

```jsonata
$flattenObj({"a": {"b": {"c": 1}, "d": 2}})
/* → {"a.b.c": 1, "a.d": 2} */

$flattenObj({"a": {"b": 1}}, "/")
/* → {"a/b": 1} */
```

---

### `$expandObj(obj [, sep])`

Expands a flat object with compound keys back to a nested structure. Inverse of `$flattenObj`.

```jsonata
$expandObj({"a.b.c": 1, "a.d": 2})
/* → {"a": {"b": {"c": 1}, "d": 2}} */
```

---
