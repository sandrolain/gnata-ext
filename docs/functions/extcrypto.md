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
