# `jn` — CLI Tool

`jn` is a command-line JSONata processor inspired by `jq`. It evaluates
JSONata expressions with all gnata-ext extension functions loaded and reads
JSON from files or stdin.

---

## Installation

```sh
# From source (requires Go 1.21+)
go install github.com/sandrolain/gnata-ext/cmd/jn@latest

# Or download a pre-built binary from GitHub Releases:
# https://github.com/sandrolain/gnata-ext/releases
```

---

## Basic Usage

```
jn [flags] [expr] [file...]
```

- **`expr`** — a JSONata expression (default: `$`, returns input unchanged)
- **`file...`** — one or more JSON files; stdin is used when omitted

---

## Quick Examples

```sh
# Identity — pretty-print stdin
echo '{"name":"Alice","age":30}' | jn '$'

# Field access
echo '{"name":"Alice"}' | jn '$.name'

# Extension function
echo '"hello world"' | jn '$camelCase($)'

# No input needed
jn -n '$uuid()'
jn -n '$dateAdd($millis(), 7, "day")'

# From a file
jn '$.users.$count($)' data.json

# Multiple files (processed sequentially)
jn '$sum($.prices)' jan.json feb.json mar.json

# Compact output
echo '[3,1,2]' | jn -c '$sort($)'

# Raw string output (no JSON quotes)
echo '{"msg":"hello"}' | jn -r '$.msg'

# Slurp all JSON values into one array, then evaluate
cat records.ndjson | jn -s '$count($)'
cat records.ndjson | jn -s '$filter($, function($v){ $v.active })'

# Read expression from a file
jn -f transform.jsonata data.json

# Bind a string variable
jn --arg 'prefix=hello' '$prefix & " " & $.name' data.json

# Bind a JSON variable
jn --argjson 'limit=10' '$take($.items, $limit)' data.json

# Exit status: exit 5 if output is null/false
jn -e '$.active' record.json && echo "active"
```

---

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--compact` | `-c` | Compact JSON output (no indentation) |
| `--raw-output` | `-r` | Output raw strings without JSON quotes |
| `--raw-output0` | | Like `-r` but use NUL (`\0`) as record separator instead of newline |
| `--raw-input` | `-R` | Read each input line as a raw string |
| `--null-input` | `-n` | Use null as input (evaluate without data) |
| `--exit-status` | `-e` | Exit 5 if last output is null or false |
| `--from-file <file>` | `-f` | Read JSONata expression from a file |
| `--slurp` | `-s` | Slurp all JSON values into an array first |
| `--join-output` | `-j` | No trailing newline after each output |
| `--tab` | | Use tab indentation (overrides `--indent`) |
| `--indent <n>` | | Indentation width in spaces (0–7, default 2) |
| `--sort-keys` | `-S` | Sort object keys (no-op: Go already sorts) |
| `--color-output` | `-C` | Force colorized output even when not writing to a terminal |
| `--monochrome-output` | `-M` | Disable colorized output |
| `--unbuffered` | | Flush output after each JSON object is printed |
| `--data <json>` | | Inline JSON input data string |
| `--data-file <file>` | | Path to JSON input file |
| `--arg name=value` | | Bind `$name` to a string value |
| `--argjson name=json` | | Bind `$name` to a parsed JSON value |

---

## Subcommands

### `jn list` — List extension functions

```sh
# All functions
jn list

# Filter by package
jn list --package extarray
jn list -p extnumeric
```

Output columns: `FUNCTION`, `PACKAGE`, `DESCRIPTION`.

---

### `jn describe` — Show function details

```sh
jn describe haversine
jn describe $chunk     # leading $ is optional
```

Output:

```
Function:    $haversine
Package:     extgeo
Signature:   $haversine(lat1, lon1, lat2, lon2: number) → number
Description: Great-circle distance in kilometres.
```

---

### `jn version` — Show version

```sh
jn version
# jn version v1.2.3
```

---

## Color Output

`jn` automatically colorizes JSON output (similar to `jq`) when writing to a
terminal, using colors close to jq's default scheme. Color is disabled
automatically when stdout is redirected to a pipe or file.

| Flag | Behavior |
|------|----------|
| _(default)_ | Color when stdout is a TTY, monochrome otherwise |
| `-C` / `--color-output` | Force color even when not writing to a terminal |
| `-M` / `--monochrome-output` | Disable color unconditionally |
| `NO_COLOR=1` (env) | Disables color unless `-C` is also given |

### `JN_COLORS` — Custom color scheme

Set the `JN_COLORS` environment variable to a colon-separated list of eight
ANSI escape code fragments (same format as jq's `JQ_COLORS`):

```
null:false:true:numbers:strings:arrays:objects:object-keys
```

Each value is an ANSI SGR parameter string such as `"1;31"` (bright red).
Fields left empty keep the default color.

```sh
# Example: red strings, cyan numbers
export JN_COLORS="::::::1;31:0;36"
echo '{"x":42,"y":"hello"}' | jn -C '$'
```

---

## Advanced Examples

```sh
# Date arithmetic
jn -n '$fromMillis($dateAdd($millis(), 7, "day"))'

# String manipulation
echo '["hello world","foo bar"]' | jn '$map($, $camelCase)'

# Crypto
jn -n '$hash("sha256", "hello")'
jn -n '$uuid()'

# Array operations
echo '[1,2,3,4,5,6,7,8,9,10]' | jn '$chunk($, 3)'
echo '[1,2,3,4,5]' | jn '$window($, 3)'

# Geo distance
jn -n '$haversine(51.5, -0.1, 48.8, 2.3)'

# Validate input
echo '"user@example.com"' | jn '$isEmail($)'

# CSV parsing
echo '"a,b,c\n1,2,3\n4,5,6"' | jn '$csv($)'

# Complex transform from file
jn '
  $map($.orders, function($o) {
    {
      "id":    $o.id,
      "total": $sum($o.items.price),
      "date":  $fromMillis($o.timestamp)
    }
  })
' orders.json
```

---

## Using `--arg` and `--argjson`

Variable bindings are injected as JSONata block assignments wrapping the expression:

```sh
# String binding
jn --arg 'sep=, ' '$join($.tags, $sep)' item.json
# Equivalent expression: ($sep := ", "; $join($.tags, $sep))

# JSON binding
jn --argjson 'n=5' '$take($sort($.scores), $n)' data.json
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Usage or flag error |
| `2` | Compile error (invalid expression) |
| `5` | `--exit-status` and last output was null/false |

---

## Comparison with jq

| Feature | `jq` | `jn` |
|---------|------|------|
| Query language | jq filter syntax | JSONata 2.x |
| Extension functions | built-in + modules | 110+ gnata-ext functions |
| Null input (`-n`) | ✓ | ✓ |
| Raw output (`-r`) | ✓ | ✓ |
| Raw input (`-R`) | ✓ | ✓ |
| Slurp (`-s`) | ✓ | ✓ |
| Compact (`-c`) | ✓ | ✓ |
| Exit status (`-e`) | ✓ | ✓ |
| From file (`-f`) | ✓ | ✓ |
| `--arg` / `--argjson` | ✓ | ✓ (name=value format) |
| Streaming large files | ✓ | ✓ (json.Decoder) |
| Color output (`-C`/`-M`) | ✓ | ✓ (auto-detected; `JN_COLORS` env) |
| `--raw-output0` | ✓ | ✓ |
| `--unbuffered` | ✓ | ✓ |
| `@base64`, `@csv` formats | ✓ | — (use extension functions) |
