# CLI Tool

Command-line interface for evaluating expressions and discovering functions.

---

## Installation

```sh
go install github.com/sandrolain/gnata-ext/cmd/gnata-ext-cli@latest
```

---

## Commands

### `eval` — Evaluate expressions

Evaluate a JSONata expression with all extension functions available:

```sh
# Simple expression
gnata-ext-cli eval '$first([1,2,3])'
# Output: 1

# With inline data
gnata-ext-cli eval '$hash("sha256","hello")' --data '{"msg":"world"}'

# With data from file
gnata-ext-cli eval '$dateAdd($millis(),1,"day")' --data-file payload.json
```

---

### `list` — List functions

Show all registered functions (or filter by package):

```sh
# All functions
gnata-ext-cli list
# extarray: first, last, take, skip, ...
# extcrypto: uuid, hash, hmac
# ...

# Filter by package
gnata-ext-cli list --package extarray
# $first(array) — First element
# $last(array) — Last element
# ...
```

---

### `describe` — Function details

Show signature and description for a specific function:

```sh
gnata-ext-cli describe haversine
# Output:
# Package: extgeo
# Signature: haversine(lat1, lon1, lat2, lon2)
# Great-circle distance in km

gnata-ext-cli describe chunk
# Output:
# Package: extarray
# Signature: chunk(array, size)
# Split into fixed-size chunks
```

---

## Flags

| Flag | Description | Example |
|---|---|---|
| `--data` | JSON data to pass to expression | `--data '{"x":1}'` |
| `--data-file` | Path to JSON file | `--data-file input.json` |
| `--package` | Filter by package name (for `list`) | `--package extstring` |

---

## Examples

```sh
# Date arithmetic
gnata-ext-cli eval '$dateAdd($millis(), 7, "day")'

# String manipulation
gnata-ext-cli eval '$camelCase("hello world")'

# Array operations
gnata-ext-cli eval '$chunk($range(1, 11), 3)'

# Complex expression with data
gnata-ext-cli eval '$uuid() & ":" & $hash("sha256", email)' \
    --data '{"email":"user@example.com"}'

# Read from file
gnata-ext-cli eval '$toCSV(records)' --data-file records.json

# List all validation functions
gnata-ext-cli list --package extvalidate

# Check if a string is an email
gnata-ext-cli eval '$isEmail("test@example.com")'
```
