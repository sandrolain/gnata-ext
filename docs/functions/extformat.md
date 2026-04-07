# extformat — Formatting & Serialisation

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extformat`

---

### `$csv(text)`

Parses a CSV string into an array of objects. The first row is the header; each subsequent row becomes an object with header keys.

```jsonata
$csv("name,age\nAlice,30\nBob,25")
/* → [
     {"name": "Alice", "age": "30"},
     {"name": "Bob",   "age": "25"}
   ] */
```

All values are returned as strings, consistent with standard CSV semantics.

---

### `$toCSV(array)`

Serializes an array of objects to CSV text. The header row is derived from the keys of the **first** object; missing keys in subsequent rows are written as empty strings.

```jsonata
$toCSV([
  {"name": "Alice", "age": "30"},
  {"name": "Bob",   "age": "25"}
])
/* → "name,age\nAlice,30\nBob,25\n" */
```

---

### `$template(str, vars)`

Replaces `{{key}}` placeholders in *str* with values from the *vars* object. Unknown placeholders are left unchanged.

```jsonata
$template("Hello, {{name}}! You are {{age}}.", {"name": "Alice", "age": 30})
/* → "Hello, Alice! You are 30." */

$template("{{greeting}}, world!", {})
/* → "{{greeting}}, world!" */
```

> **Note:** `extstring` also exposes `$template` with identical behaviour. When both packages are registered, the last one wins; register only the package you need if this matters.

---
