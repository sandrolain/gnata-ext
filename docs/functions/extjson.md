# extjson — JSON Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extjson`

---

### `$jsonParse(str)`

Parses a JSON string and returns the corresponding gnata value. Errors on invalid JSON.

```jsonata
$jsonParse("{\"a\":1,\"b\":true}")
/* → {"a": 1, "b": true} */

$jsonParse("[1, 2, 3]")
/* → [1, 2, 3] */
```

---

### `$jsonStringify(v [, indent])`

Serialises *v* to a JSON string. If *indent* is provided (a string, e.g. `"  "`), the output is pretty-printed.

```jsonata
$jsonStringify({"x": 1})          /* → "{\"x\":1}" */
$jsonStringify({"x": 1}, "  ")    /* pretty-printed */
```

---

### `$jsonDiff(a, b)`

Returns a JSON Patch-compatible array of `{op, path, value}` operations describing what changed between *a* and *b*. Uses RFC 6902 operation names (`add`, `remove`, `replace`). Paths use JSON Pointer notation.

```jsonata
$jsonDiff({"x": 1, "y": 2}, {"x": 1, "y": 3, "z": 4})
/* → [
     {"op": "replace", "path": "/y", "value": 3},
     {"op": "add",     "path": "/z", "value": 4}
   ] */
```

---

### `$jsonPatch(obj, ops)`

Applies RFC 6902 JSON Patch operations to *obj* (immutable — returns a new value). Supported operations: `add`, `remove`, `replace`, `move`, `copy`, `test`.

```jsonata
$jsonPatch(
  {"a": 1, "b": 2},
  [
    {"op": "replace", "path": "/a", "value": 10},
    {"op": "remove",  "path": "/b"},
    {"op": "add",     "path": "/c", "value": 3}
  ]
)
/* → {"a": 10, "c": 3} */
```

The `"test"` operation throws an error if the path value does not match, which causes the entire patch to fail.

---

### `$jsonPointer(obj, pointer)`

Resolves an RFC 6901 JSON Pointer against *obj*. Supports traversal of both objects and arrays.

```jsonata
$jsonPointer({"a": {"b": [10, 20, 30]}}, "/a/b/1")   /* → 20 */
$jsonPointer({"a": 1}, "/a")                          /* → 1 */
```

---
