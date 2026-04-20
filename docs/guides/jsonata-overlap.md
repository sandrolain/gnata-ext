# Relationship with JSONata Built-in Functions

`gnata-ext` extends JSONata with functions that are not available natively. However, some extension functions duplicate or partially overlap with JSONata's own built-in functions. This page documents every such overlap so you can make an informed choice.

---

## Significant overlaps — prefer native JSONata

These extension functions are largely equivalent to a native JSONata function. You can use either, but prefer the native one for portability and to avoid importing unnecessary packages.

### String

| gnata-ext | JSONata native | Notes |
|---|---|---|
| `$padStart(s, n, c)` | `$pad(s, -n [, c])` | Negative width in `$pad` means left-padding |
| `$padEnd(s, n, c)` | `$pad(s, n [, c])` | Positive width in `$pad` means right-padding |

`gnata-ext` versions accept explicit `padStart`/`padEnd` naming which is more readable; both implementations are functionally equivalent.

### Numeric

| gnata-ext | JSONata native | Notes |
|---|---|---|
| `$roundTo(n, precision)` | `$round(n, precision)` | Identical semantics; `$roundTo` is purely an alias |

### Array

| gnata-ext | JSONata native | Notes |
|---|---|---|
| `$sortBy(array, key)` | `$sort(array, function($l,$r){ $l.key > $r.key })` | `$sortBy` simplifies alphabetical/numerical sorting by a named key; `$sort` is more flexible for custom comparators |

### Types / coercion

| gnata-ext | JSONata native | Notes |
|---|---|---|
| `$typeOf(v)` | `$type(v)` | Same return values (`"string"`, `"number"`, `"boolean"`, `"array"`, `"object"`, `"null"`). `$type` returns `undefined` for undefined; `$typeOf` returns `"null"`. |
| `$toString(v)` | `$string(v)` | Both coerce to string. `$string` uses JSON serialisation for non-scalars; `$toString` uses Go's `fmt.Sprintf`. |
| `$toNumber(v)` | `$number(v)` | Both coerce to number. `$number` supports hex (`0x`), octal (`0o`), binary (`0b`); `$toNumber` does not. |
| `$toBool(v)` | `$boolean(v)` | Identical casting rules. |

### Date / Time

| gnata-ext | JSONata native | Notes |
|---|---|---|
| `$dateFormat(ms, layout)` | `$fromMillis(ms [, picture])` | Same functionality. **Key difference:** `$dateFormat` uses Go time layout strings (`"2006-01-02"`); `$fromMillis` uses XPath F&O picture strings (`"[Y0001]-[M01]-[D01]"`). |
| `$dateParse(str, layout)` | `$toMillis(str [, picture])` | Same functionality. Same format-string difference as above. |

Use `$fromMillis` / `$toMillis` when you need XPath-compatible format strings. Use `$dateFormat` / `$dateParse` when you prefer Go time format strings (more familiar to Go developers and compatible with `time.Format`/`time.Parse`).

---

## Partial overlaps — gnata-ext adds convenience

These native JSONata expressions can approximate the extension function, but the extension provides a cleaner API or additional behaviour.

### Array

| gnata-ext | Native equivalent | Extra value |
|---|---|---|
| `$first(array)` | `array[0]` | Returns `null` for empty arrays; native throws |
| `$last(array)` | `array[-1]` | Returns `null` for empty arrays; native throws |
| `$union(a, b)` | `$distinct($append(a, b))` | Single function call |
| `$flatten(array [, depth])` | `$reduce(a, $append, [])` | Supports depth limiting; native requires a reduce chain |
| `$range(start, end, step)` | `[start..end]` syntax | Supports fractional steps, negative steps, and stop-before-end semantics |
| `$zipLongest(a, b, fill)` | — | `$zip` truncates to shortest; `$zipLongest` pads with `fill` |

### Type predicates

| gnata-ext | Native equivalent | Extra value |
|---|---|---|
| `$isString(v)` | `$type(v) = "string"` | Shorter expression |
| `$isNumber(v)` | `$type(v) = "number"` | Shorter expression |
| `$isBoolean(v)` | `$type(v) = "boolean"` | Shorter expression |
| `$isArray(v)` | `$type(v) = "array"` | Shorter expression |
| `$isObject(v)` | `$type(v) = "object"` | Shorter expression |

### Object

| gnata-ext | Native equivalent | Extra value |
|---|---|---|
| `$deepMerge(a, b)` | `$merge([a, b])` | `$merge` is shallow; `$deepMerge` recursively merges nested objects |

### Boolean / existence

| gnata-ext | Native equivalent | Notes |
|---|---|---|
| `$defined(v)` | `$exists(v)` | `$exists` operates on path expressions; `$defined` tests an already-evaluated value for `null`/nil |

---

## Related but distinct — different variant

These extension functions are related to a native JSONata function but provide a **different encoding or behaviour** that native cannot replicate.

| gnata-ext | Related native | Difference |
|---|---|---|
| `$base64url(str)` | `$base64encode(str)` | `$base64encode` → standard base64 with `=` padding; `$base64url` → URL-safe alphabet, no padding (RFC 4648 §5) |
| `$unbase64url(str)` | `$base64decode(str)` | Corresponding URL-safe decode |

---

## No overlap — unique to gnata-ext

All other extension functions have no JSONata built-in equivalent:

**extstring:** `startsWith`, `endsWith`, `indexOf`, `lastIndexOf`, `capitalize`, `titleCase`, `camelCase`, `snakeCase`, `kebabCase`, `repeat`, `words`, `truncate`, `slugify`, `countOccurrences`, `initials`, `escapeHTML`, `unescapeHTML`, `reverseWords`, `levenshtein`, `longestCommonPrefix`

**extnumeric:** `log`, `sign`, `trunc`, `clamp`, trigonometric functions, `pi`, `e`, `median`, `variance`, `stddev`, `percentile`, `mode`, `product`, `cumSum`, `inRange`, `normalize`, `interpolate`, `gcd`, `lcm`, `isPrime`, `factorial`

**extarray:** `take`, `skip`, `slice`, `chunk`, `intersection`, `difference`, `symmetricDifference`, `window`, `compact`, `groupByKey`, `uniqueBy`, `sumByKey`, `countByKey`, `rotate`, `indexof`, `transpose`, `adjacentPairs`

**extobject:** `values`, `pairs`, `fromPairs`, `pick`, `omit`, `invert`, `size`, `rename`, `clean`, `defaults`, `transform`, `filterKeys`, `groupByValue`

**extdatetime:** `dateAdd`, `dateDiff`, `dateComponents`, `dateStartOf`, `dateEndOf`, `dateIsBefore`, `dateIsAfter`, `dateIsBetween`, `dateWeek`, `dateQuarter`, `dateDayOfYear`, `isLeapYear`, `daysInMonthOf`

**exttypes:** `isEmpty`, `default`, `identity`, `toArray`, `nullish`

**extcrypto:** `uuid`, `hash`, `hmac`, `randomBytes`

**extformat, extpath, extvalidate, extjson, extgeo, extnet:** all functions are unique.
