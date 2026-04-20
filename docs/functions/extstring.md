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

### `$padStart(str, width, char)`

Left-pads *str* to *width* total characters using *char* (default `" "`). If *str* is already at least *width* chars, it is returned unchanged.

> **JSONata native alternative:** `$pad(str, -width [, char])` — negative width left-pads in standard JSONata.

```jsonata
$padStart("5", 3, "0")   /* → "005" */
$padStart("hi", 5)       /* → "   hi" */
```

---

### `$padEnd(str, width, char)`

Right-pads *str* to *width* total characters using *char* (default `" "`).

> **JSONata native alternative:** `$pad(str, width [, char])` — positive width right-pads in standard JSONata.

```jsonata
$padEnd("hi", 5, ".")   /* → "hi..." */
```

---

### `$truncate(str, max [, suffix])`

Returns *str* truncated to at most *max* characters. When truncation occurs, *suffix* (default `"..."`) is appended and the total length stays within *max*.

```jsonata
$truncate("Hello World", 8)         /* → "Hello..." */
$truncate("Hello World", 8, "…")    /* → "Hello W…" */
$truncate("Hi", 10)                 /* → "Hi" */
```

---

### `$slugify(str)`

Converts *str* to a URL-safe lowercase slug: lowercases, replaces non-alphanumeric sequences with `-`, and strips leading/trailing dashes.

```jsonata
$slugify("Hello World!")     /* → "hello-world" */
$slugify("  Foo  Bar  ")     /* → "foo-bar" */
```

---

### `$countOccurrences(str, sub)`

Returns the number of non-overlapping occurrences of *sub* in *str*.

```jsonata
$countOccurrences("abcabcabc", "abc")   /* → 3 */
$countOccurrences("aaa", "aa")          /* → 1 */
```

---

### `$initials(str)`

Extracts the uppercase initial letter of each word.

```jsonata
$initials("John Doe")          /* → "JD" */
$initials("hello world foo")   /* → "HWF" */
```

---

### `$escapeHTML(str)`

Replaces `&`, `<`, `>`, `"`, and `'` with their HTML entity equivalents.

```jsonata
$escapeHTML("<b>Hello & 'World'</b>")
/* → "&lt;b&gt;Hello &amp; &#39;World&#39;&lt;/b&gt;" */
```

---

### `$unescapeHTML(str)`

Decodes HTML entities back to their characters.

```jsonata
$unescapeHTML("&lt;b&gt;Hello&lt;/b&gt;")   /* → "<b>Hello</b>" */
```

---

### `$reverseWords(str)`

Reverses the order of words in *str*.

```jsonata
$reverseWords("hello world foo")   /* → "foo world hello" */
```

---

### `$levenshtein(a, b)`

Returns the Levenshtein edit distance between strings *a* and *b*.

```jsonata
$levenshtein("kitten", "sitting")   /* → 3 */
$levenshtein("abc", "abc")          /* → 0 */
```

---

### `$longestCommonPrefix(array)`

Returns the longest string that is a prefix of every element in *array*. Returns `""` if the array is empty or there is no common prefix.

```jsonata
$longestCommonPrefix(["flower", "flow", "flight"])   /* → "fl" */
$longestCommonPrefix(["dog", "racecar", "car"])       /* → "" */
```

---
