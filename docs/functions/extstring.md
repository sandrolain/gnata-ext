# extstring — String Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extstring`

---

### `$startsWith(str, prefix)`

Returns `true` if *str* starts with *prefix*.

```jsonata
$startsWith("hello world", "hello")   /* → true */
$startsWith("hello world", "world")   /* → false */
```

---

### `$endsWith(str, suffix)`

Returns `true` if *str* ends with *suffix*.

```jsonata
$endsWith("hello world", "world")   /* → true */
```

---

### `$indexOf(str, search [, start])`

Returns the first index of *search* in *str*, or `−1` if not found. Optional *start* offsets where to begin the search.

```jsonata
$indexOf("abcabc", "b")      /* → 1 */
$indexOf("abcabc", "b", 2)   /* → 4 */
$indexOf("abcabc", "x")      /* → -1 */
```

---

### `$lastIndexOf(str, search)`

Returns the last index of *search* in *str*, or `−1` if not found.

```jsonata
$lastIndexOf("abcabc", "b")   /* → 4 */
```

---

### `$capitalize(str)`

Uppercases the first character and lowercases the rest.

```jsonata
$capitalize("hELLO")   /* → "Hello" */
```

---

### `$titleCase(str)`

Title-cases every word (first letter uppercase, rest lowercase).

```jsonata
$titleCase("hello world foo")   /* → "Hello World Foo" */
```

---

### `$camelCase(str)`

Converts to camelCase. Splits on whitespace, hyphens, underscores, and case boundaries.

```jsonata
$camelCase("hello world")    /* → "helloWorld" */
$camelCase("foo-bar-baz")    /* → "fooBarBaz" */
$camelCase("my_variable")    /* → "myVariable" */
```

---

### `$snakeCase(str)`

Converts to snake_case.

```jsonata
$snakeCase("helloWorld")     /* → "hello_world" */
$snakeCase("Hello World")    /* → "hello_world" */
```

---

### `$kebabCase(str)`

Converts to kebab-case.

```jsonata
$kebabCase("helloWorld")     /* → "hello-world" */
$kebabCase("Hello World")    /* → "hello-world" */
```

---

### `$repeat(str, n)`

Returns *str* repeated *n* times.

```jsonata
$repeat("ab", 3)    /* → "ababab" */
$repeat("x", 0)     /* → "" */
```

---

### `$words(str)`

Splits *str* into an array of words, splitting on whitespace, hyphens, underscores, and camelCase boundaries.

```jsonata
$words("hello world")      /* → ["hello", "world"] */
$words("camelCaseWord")    /* → ["camel", "Case", "Word"] */
$words("foo-bar_baz")      /* → ["foo", "bar", "baz"] */
```

---

### `$template(str, vars)`

Replaces `{{key}}` placeholders in *str* with values from the *vars* object. See also [`extformat.$template`](extformat.md#template).

```jsonata
$template("Hi {{name}}", {"name": "Bob"})   /* → "Hi Bob" */
```

---
