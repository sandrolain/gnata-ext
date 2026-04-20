# extcrypto — Cryptographic Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extcrypto`

---

### `$uuid()`

Generates a random UUID v4.

```jsonata
$uuid()   /* → "550e8400-e29b-41d4-a716-446655440000" */
```

---

### `$hash(algorithm, value)`

Returns a lower-case hex-encoded digest of *value*.

Supported algorithms: `md5`, `sha1`, `sha256`, `sha384`, `sha512`.

```jsonata
$hash("sha256", "hello")
/* → "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" */

$hash("md5", "gnata-ext")
```

---

### `$hmac(algorithm, key, value)`

Returns a lower-case hex-encoded HMAC of *value* signed with *key*.

Supported algorithms: `md5`, `sha1`, `sha256`, `sha384`, `sha512`.

```jsonata
$hmac("sha256", "secret", "message")
```

---

### `$randomBytes(n)`

Returns *n* cryptographically random bytes encoded as a lowercase hexadecimal string. The result is always 2×*n* characters long.

```jsonata
$randomBytes(8)   /* → e.g. "3f9a1b2c4e5d6a7b" */
```

---

### `$base64url(str)`

Encodes *str* using RFC 4648 §5 URL-safe Base64 alphabet with **no padding characters**.

> **Note:** This is different from `$base64encode(str)`, which uses the standard alphabet with `=` padding. The two encodings are **not interchangeable**.

```jsonata
$base64url("hello world")   /* → "aGVsbG8gd29ybGQ" */
```

---

### `$unbase64url(str)`

Decodes a URL-safe Base64 (no padding) string produced by `$base64url`.

> **Note:** This function does not decode output from `$base64decode`; use the matching `$base64url` / `$unbase64url` pair consistently.

```jsonata
$unbase64url("aGVsbG8gd29ybGQ")   /* → "hello world" */
```

---
