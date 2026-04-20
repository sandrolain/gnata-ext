package ext

import "sort"

// FuncMeta holds documentation metadata for a single extension function.
type FuncMeta struct {
	Name        string // function name without leading $
	Package     string // sub-package name (e.g. "extarray")
	Signature   string // short signature string (e.g. "(array) → any")
	Description string // one-line description
}

// catalog is the static registry of all extension function metadata.
var catalog = []FuncMeta{
	// extarray
	{Name: "first", Package: "extarray", Signature: "(array) → any", Description: "Returns the first element of an array, or null if empty."},
	{Name: "last", Package: "extarray", Signature: "(array) → any", Description: "Returns the last element of an array, or null if empty."},
	{Name: "take", Package: "extarray", Signature: "(array, n: number) → array", Description: "Returns the first n elements."},
	{Name: "skip", Package: "extarray", Signature: "(array, n: number) → array", Description: "Returns all elements after the first n."},
	{Name: "slice", Package: "extarray", Signature: "(array, start: number [, end: number]) → array", Description: "Returns a sub-array from start (inclusive) to end (exclusive)."},
	{Name: "flatten", Package: "extarray", Signature: "(array [, depth: number]) → array", Description: "Recursively flattens nested arrays."},
	{Name: "chunk", Package: "extarray", Signature: "(array, size: number) → array", Description: "Splits array into sub-arrays of the given size."},
	{Name: "union", Package: "extarray", Signature: "(a: array, b: array) → array", Description: "De-duplicated union of two arrays."},
	{Name: "intersection", Package: "extarray", Signature: "(a: array, b: array) → array", Description: "Elements present in both arrays."},
	{Name: "difference", Package: "extarray", Signature: "(a: array, b: array) → array", Description: "Elements in a that are not in b."},
	{Name: "symmetricDifference", Package: "extarray", Signature: "(a: array, b: array) → array", Description: "Elements in either array but not in both."},
	{Name: "range", Package: "extarray", Signature: "(start: number, end: number [, step: number]) → array", Description: "Generates a numeric range."},
	{Name: "zipLongest", Package: "extarray", Signature: "(a: array, b: array [, fill: any]) → array", Description: "Zips two arrays, padding the shorter one."},
	{Name: "window", Package: "extarray", Signature: "(array, size: number [, step: number]) → array", Description: "Produces sliding windows over an array."},
	// extcolor
	{Name: "colorParse", Package: "extcolor", Signature: "(s: string) → object", Description: "Parse a CSS color string (#rgb, #rrggbb, rgb(), rgba()) into {r,g,b,a}."},
	{Name: "colorToHex", Package: "extcolor", Signature: "(obj: object) → string", Description: "Convert {r,g,b} to a hex color string."},
	{Name: "colorToRGB", Package: "extcolor", Signature: "(obj: object) → string", Description: "Convert {r,g,b,a?} to rgb() or rgba() CSS string."},
	{Name: "colorToHSL", Package: "extcolor", Signature: "(obj: object) → object", Description: "Convert {r,g,b} to {h,s,l} (HSL color model)."},
	{Name: "colorMix", Package: "extcolor", Signature: "(a: object, b: object, t: number) → object", Description: "Linear blend of two colors at factor t (0–1)."},
	{Name: "colorLighten", Package: "extcolor", Signature: "(obj: object, amount: number) → object", Description: "Increase HSL lightness by amount (0–1)."},
	{Name: "colorDarken", Package: "extcolor", Signature: "(obj: object, amount: number) → object", Description: "Decrease HSL lightness by amount (0–1)."},
	{Name: "colorContrast", Package: "extcolor", Signature: "(fg: object, bg: object) → number", Description: "WCAG contrast ratio between two colors."},
	{Name: "colorLuminance", Package: "extcolor", Signature: "(obj: object) → number", Description: "Relative luminance (0–1) per WCAG 2.1."},
	// extcrypto
	{Name: "uuid", Package: "extcrypto", Signature: "() → string", Description: "Generates a random UUID v4."},
	{Name: "hash", Package: "extcrypto", Signature: "(algorithm: string, value: string) → string", Description: "Hex-encoded cryptographic hash."},
	{Name: "hmac", Package: "extcrypto", Signature: "(algorithm: string, key: string, value: string) → string", Description: "Hex-encoded HMAC."},
	// extdatetime
	{Name: "dateAdd", Package: "extdatetime", Signature: "(timestamp: number, amount: number, unit: string) → number", Description: "Adds time units to a Unix-ms timestamp."},
	{Name: "dateDiff", Package: "extdatetime", Signature: "(t1: number, t2: number, unit: string) → number", Description: "Whole-unit difference between two timestamps."},
	{Name: "dateComponents", Package: "extdatetime", Signature: "(timestamp: number) → object", Description: "UTC components of a Unix-ms timestamp."},
	{Name: "dateStartOf", Package: "extdatetime", Signature: "(timestamp: number, unit: string) → number", Description: "Start of the given period."},
	{Name: "dateEndOf", Package: "extdatetime", Signature: "(timestamp: number, unit: string) → number", Description: "End (last ms) of the given period."},
	// extdiff
	{Name: "diff", Package: "extdiff", Signature: "(a: object, b: object) → object", Description: "Structural diff returning {added, removed, changed}."},
	{Name: "patch", Package: "extdiff", Signature: "(obj: object, diff: object) → object", Description: "Applies a diff to reconstruct the second object."},
	{Name: "changed", Package: "extdiff", Signature: "(a: object, b: object, key: string) → boolean", Description: "Returns true if the given key changed between a and b."},
	{Name: "addedKeys", Package: "extdiff", Signature: "(a: object, b: object) → array", Description: "Keys present in b but not in a."},
	{Name: "removedKeys", Package: "extdiff", Signature: "(a: object, b: object) → array", Description: "Keys present in a but not in b."},
	{Name: "arrayDiff", Package: "extdiff", Signature: "(a: array, b: array) → object", Description: "Returns {added, removed} between two arrays."},
	{Name: "deepEqual", Package: "extdiff", Signature: "(a: any, b: any) → boolean", Description: "Recursive deep equality comparison."},
	// extformat
	{Name: "csv", Package: "extformat", Signature: "(text: string) → array", Description: "Parses a CSV string into an array of objects."},
	{Name: "toCSV", Package: "extformat", Signature: "(array) → string", Description: "Serialises an array of objects to CSV text."},
	{Name: "template", Package: "extformat", Signature: "(str: string, vars: object) → string", Description: "Replaces {{key}} placeholders with values from vars."},
	// extgeo
	{Name: "haversine", Package: "extgeo", Signature: "(lat1, lon1, lat2, lon2: number) → number", Description: "Great-circle distance in kilometres."},
	{Name: "bearing", Package: "extgeo", Signature: "(lat1, lon1, lat2, lon2: number) → number", Description: "Initial bearing in degrees clockwise from north."},
	{Name: "geoFormat", Package: "extgeo", Signature: "(lat, lon: number [, format: string]) → string", Description: "Formats a coordinate pair as a string."},
	{Name: "geoParse", Package: "extgeo", Signature: "(str: string) → object", Description: "Parses a decimal coordinate string to {lat, lon}."},
	{Name: "inBoundingBox", Package: "extgeo", Signature: "(lat, lon, minLat, minLon, maxLat, maxLon: number) → boolean", Description: "Returns true if point is inside bounding box."},
	{Name: "geoDistance", Package: "extgeo", Signature: "(point: object, points: array) → array", Description: "Haversine distance from point to each element of points."},
	// extjson
	{Name: "jsonParse", Package: "extjson", Signature: "(str: string) → any", Description: "Parses a JSON string."},
	{Name: "jsonStringify", Package: "extjson", Signature: "(v: any [, indent: string]) → string", Description: "Serialises a value to a JSON string."},
	{Name: "jsonDiff", Package: "extjson", Signature: "(a: any, b: any) → array", Description: "RFC 6902 patch operations describing the diff between a and b."},
	{Name: "jsonPatch", Package: "extjson", Signature: "(obj: any, ops: array) → any", Description: "Applies RFC 6902 JSON Patch operations."},
	{Name: "jsonPointer", Package: "extjson", Signature: "(obj: any, pointer: string) → any", Description: "Resolves an RFC 6901 JSON Pointer."},
	// extlogic
	{Name: "ifElse", Package: "extlogic", Signature: "(cond: any, then: any, else: any) → any", Description: "Returns then if cond is truthy, otherwise else."},
	{Name: "when", Package: "extlogic", Signature: "(cond: any, value: any) → any", Description: "Returns value if cond is truthy, otherwise null."},
	{Name: "unless", Package: "extlogic", Signature: "(cond: any, value: any) → any", Description: "Returns value if cond is falsy, otherwise null."},
	{Name: "switch", Package: "extlogic", Signature: "(v: any, cases: object [, default: any]) → any", Description: "Returns the matching cases entry, or default."},
	{Name: "coalesce", Package: "extlogic", Signature: "(v1: any, ...) → any", Description: "Returns the first non-nil, non-empty argument."},
	{Name: "tap", Package: "extlogic", Signature: "(v: any) → any", Description: "Returns v unchanged (pass-through helper)."},
	// exturi
	{Name: "uriParse", Package: "exturi", Signature: "(url: string) → object", Description: "Decomposes a URL into {scheme, user, password, host, port, path, query, fragment}."},
	{Name: "uriBuild", Package: "exturi", Signature: "(parts: object) → string", Description: "Builds a URL string from a parts object."},
	{Name: "uriJoin", Package: "exturi", Signature: "(base: string, ref: string) → string", Description: "Resolves a relative URL against a base URL."},
	{Name: "queryParse", Package: "exturi", Signature: "(qs: string) → object", Description: "Parses a query string into an object."},
	{Name: "queryBuild", Package: "exturi", Signature: "(obj: object) → string", Description: "Serializes an object to a URL-encoded query string."},
	{Name: "uriGetPath", Package: "exturi", Signature: "(url: string) → string", Description: "Extracts the path component from a URL."},
	{Name: "uriGetQuery", Package: "exturi", Signature: "(url: string) → object", Description: "Extracts and parses the query string from a URL."},
	{Name: "uriSetQuery", Package: "exturi", Signature: "(url: string, params: object) → string", Description: "Replaces the query string of a URL with the given params."},
	// extnet
	{Name: "ipVersion", Package: "extnet", Signature: "(str: string) → number", Description: "Returns 4, 6, or -1 for the IP version."},
	{Name: "isPrivateIP", Package: "extnet", Signature: "(str: string) → boolean", Description: "Returns true for private/loopback/link-local addresses."},
	{Name: "ipToInt", Package: "extnet", Signature: "(str: string) → number", Description: "Converts an IPv4 address to its integer representation."},
	{Name: "intToIP", Package: "extnet", Signature: "(n: number) → string", Description: "Converts an integer back to an IPv4 address string."},
	{Name: "ipInCIDR", Package: "extnet", Signature: "(ip: string, cidr: string) → boolean", Description: "Returns true if ip is contained in the CIDR block."},
	{Name: "expandCIDR", Package: "extnet", Signature: "(cidr: string) → object", Description: "Returns network info for a CIDR block."},
	// extnumeric
	{Name: "log", Package: "extnumeric", Signature: "(n: number [, base: number]) → number", Description: "Natural or base-n logarithm."},
	{Name: "sign", Package: "extnumeric", Signature: "(n: number) → number", Description: "Returns -1, 0, or 1."},
	{Name: "trunc", Package: "extnumeric", Signature: "(n: number) → number", Description: "Truncates toward zero."},
	{Name: "clamp", Package: "extnumeric", Signature: "(n, min, max: number) → number", Description: "Clamps n to [min, max]."},
	{Name: "sin", Package: "extnumeric", Signature: "(n: number) → number", Description: "Sine (radians)."},
	{Name: "cos", Package: "extnumeric", Signature: "(n: number) → number", Description: "Cosine (radians)."},
	{Name: "tan", Package: "extnumeric", Signature: "(n: number) → number", Description: "Tangent (radians)."},
	{Name: "asin", Package: "extnumeric", Signature: "(n: number) → number", Description: "Arc-sine (radians)."},
	{Name: "acos", Package: "extnumeric", Signature: "(n: number) → number", Description: "Arc-cosine (radians)."},
	{Name: "atan", Package: "extnumeric", Signature: "(n: number) → number", Description: "Arc-tangent (radians)."},
	{Name: "atan2", Package: "extnumeric", Signature: "(y, x: number) → number", Description: "atan(y/x) with correct quadrant."},
	{Name: "pi", Package: "extnumeric", Signature: "() → number", Description: "Returns π."},
	{Name: "e", Package: "extnumeric", Signature: "() → number", Description: "Returns Euler's number e."},
	{Name: "median", Package: "extnumeric", Signature: "(array) → number", Description: "Median of a numeric array."},
	{Name: "variance", Package: "extnumeric", Signature: "(array) → number", Description: "Population variance."},
	{Name: "stddev", Package: "extnumeric", Signature: "(array) → number", Description: "Population standard deviation."},
	{Name: "percentile", Package: "extnumeric", Signature: "(array, p: number) → number", Description: "p-th percentile using linear interpolation."},
	{Name: "mode", Package: "extnumeric", Signature: "(array) → any", Description: "Most frequent value(s)."},
	{Name: "product", Package: "extnumeric", Signature: "(array) → number", Description: "Returns the product of all numbers in the array."},
	{Name: "cumSum", Package: "extnumeric", Signature: "(array) → array", Description: "Returns cumulative sums: result[i] = sum of elements 0..i."},
	{Name: "inRange", Package: "extnumeric", Signature: "(n: number, min: number, max: number) → boolean", Description: "Returns true if min <= n <= max."},
	{Name: "roundTo", Package: "extnumeric", Signature: "(n: number, places: number) → number", Description: "Rounds n to the specified number of decimal places."},
	{Name: "normalize", Package: "extnumeric", Signature: "(array) → array", Description: "Min-max normalizes array values to [0, 1]."},
	{Name: "interpolate", Package: "extnumeric", Signature: "(a: number, b: number, t: number) → number", Description: "Linear interpolation: a + t*(b-a)."},
	{Name: "gcd", Package: "extnumeric", Signature: "(a: integer, b: integer) → integer", Description: "Greatest common divisor of a and b."},
	{Name: "lcm", Package: "extnumeric", Signature: "(a: integer, b: integer) → integer", Description: "Least common multiple of a and b."},
	{Name: "isPrime", Package: "extnumeric", Signature: "(n: integer) → boolean", Description: "Returns true if n is a prime number."},
	{Name: "factorial", Package: "extnumeric", Signature: "(n: integer) → integer", Description: "Returns n! for n in [0, 20]."},
	// extobject
	{Name: "values", Package: "extobject", Signature: "(object) → array", Description: "Returns an array of the object's own values."},
	{Name: "pairs", Package: "extobject", Signature: "(object) → array", Description: "Returns [key, value] pairs."},
	{Name: "fromPairs", Package: "extobject", Signature: "(pairs: array) → object", Description: "Constructs an object from key-value pairs."},
	{Name: "pick", Package: "extobject", Signature: "(object, keys: array) → object", Description: "Returns a new object with only the listed keys."},
	{Name: "omit", Package: "extobject", Signature: "(object, keys: array) → object", Description: "Returns a new object with the listed keys removed."},
	{Name: "deepMerge", Package: "extobject", Signature: "(target: object, source: object) → object", Description: "Recursively merges source into target."},
	{Name: "invert", Package: "extobject", Signature: "(object) → object", Description: "Swaps keys and values."},
	{Name: "size", Package: "extobject", Signature: "(object) → number", Description: "Returns the number of own keys."},
	{Name: "rename", Package: "extobject", Signature: "(object, oldKey: string, newKey: string) → object", Description: "Renames a key."},
	// extpath
	{Name: "get", Package: "extpath", Signature: "(obj: object, path: string [, default: any]) → any", Description: "Reads a nested value by dot-path."},
	{Name: "set", Package: "extpath", Signature: "(obj: object, path: string, value: any) → object", Description: "Returns a new object with value at path."},
	{Name: "del", Package: "extpath", Signature: "(obj: object, path: string) → object", Description: "Returns a new object with the path removed."},
	{Name: "has", Package: "extpath", Signature: "(obj: object, path: string) → boolean", Description: "Returns true if path exists and is non-null."},
	{Name: "flattenObj", Package: "extpath", Signature: "(obj: object [, sep: string]) → object", Description: "Flattens a nested object using sep as key separator."},
	{Name: "expandObj", Package: "extpath", Signature: "(obj: object [, sep: string]) → object", Description: "Expands a flat object with compound keys to nested."},
	// extschema
	{Name: "validateSchema", Package: "extschema", Signature: "(data: any, schema: object) → object", Description: "Validates data against a JSON Schema. Returns {valid, errors}."},
	{Name: "isValid", Package: "extschema", Signature: "(data: any, schema: object) → boolean", Description: "Returns true if data is valid according to the JSON Schema."},
	{Name: "schemaErrors", Package: "extschema", Signature: "(data: any, schema: object) → array", Description: "Returns an array of validation error strings."},
	// extsemver
	{Name: "semverParse", Package: "extsemver", Signature: "(v: string) → object", Description: "Parses a semantic version string into an object."},
	{Name: "semverCompare", Package: "extsemver", Signature: "(a: string, b: string) → number", Description: "Compares two version strings: -1, 0, or 1."},
	{Name: "semverSatisfies", Package: "extsemver", Signature: "(v: string, c: string) → boolean", Description: "True if version v satisfies the constraint c."},
	{Name: "semverBump", Package: "extsemver", Signature: "(v: string, part: string) → string", Description: "Bumps the major, minor, or patch component."},
	{Name: "semverSort", Package: "extsemver", Signature: "(arr: array) → array", Description: "Sorts an array of version strings ascending."},
	{Name: "semverMax", Package: "extsemver", Signature: "(arr: array) → string", Description: "Returns the highest version string."},
	{Name: "semverMin", Package: "extsemver", Signature: "(arr: array) → string", Description: "Returns the lowest version string."},
	// extstring
	{Name: "startsWith", Package: "extstring", Signature: "(str: string, prefix: string) → boolean", Description: "Returns true if str starts with prefix."},
	{Name: "endsWith", Package: "extstring", Signature: "(str: string, suffix: string) → boolean", Description: "Returns true if str ends with suffix."},
	{Name: "indexOf", Package: "extstring", Signature: "(str: string, search: string [, start: number]) → number", Description: "First index of search, or -1."},
	{Name: "lastIndexOf", Package: "extstring", Signature: "(str: string, search: string) → number", Description: "Last index of search, or -1."},
	{Name: "capitalize", Package: "extstring", Signature: "(str: string) → string", Description: "Uppercases the first character."},
	{Name: "titleCase", Package: "extstring", Signature: "(str: string) → string", Description: "Title-cases every word."},
	{Name: "camelCase", Package: "extstring", Signature: "(str: string) → string", Description: "Converts to camelCase."},
	{Name: "snakeCase", Package: "extstring", Signature: "(str: string) → string", Description: "Converts to snake_case."},
	{Name: "kebabCase", Package: "extstring", Signature: "(str: string) → string", Description: "Converts to kebab-case."},
	{Name: "repeat", Package: "extstring", Signature: "(str: string, n: number) → string", Description: "Returns str repeated n times."},
	{Name: "words", Package: "extstring", Signature: "(str: string) → array", Description: "Splits str into words."},
	// extstring.template intentionally omitted (duplicate with extformat.template)
	{Name: "padStart", Package: "extstring", Signature: "(str: string, len: number [, fill: string]) → string", Description: "Left-pads str to length using fill (default space)."},
	{Name: "padEnd", Package: "extstring", Signature: "(str: string, len: number [, fill: string]) → string", Description: "Right-pads str to length using fill (default space)."},
	{Name: "truncate", Package: "extstring", Signature: "(str: string, len: number [, suffix: string]) → string", Description: "Truncates str to len chars, appending suffix (default '…') if truncated."},
	{Name: "slugify", Package: "extstring", Signature: "(str: string) → string", Description: "Converts str to a URL-friendly lowercase slug."},
	{Name: "countOccurrences", Package: "extstring", Signature: "(str: string, sub: string) → number", Description: "Counts non-overlapping occurrences of sub in str."},
	{Name: "initials", Package: "extstring", Signature: "(str: string [, sep: string]) → string", Description: "Returns initials of each word joined by sep."},
	{Name: "escapeHTML", Package: "extstring", Signature: "(str: string) → string", Description: "Escapes HTML special characters (&, <, >, \", ')."},
	{Name: "unescapeHTML", Package: "extstring", Signature: "(str: string) → string", Description: "Unescapes HTML entities back to characters."},
	{Name: "reverseWords", Package: "extstring", Signature: "(str: string) → string", Description: "Reverses the order of whitespace-separated words."},
	{Name: "levenshtein", Package: "extstring", Signature: "(a: string, b: string) → number", Description: "Returns the Levenshtein edit distance between a and b."},
	{Name: "longestCommonPrefix", Package: "extstring", Signature: "(strs: array) → string", Description: "Returns the longest string that is a prefix of all strings in the array."},
	// extregex
	{Name: "regexAll", Package: "extregex", Signature: "(str: string, pattern: string) → array", Description: "Returns all non-overlapping matches of pattern in str."},
	{Name: "regexNamedGroups", Package: "extregex", Signature: "(str: string, pattern: string) → object", Description: "Returns named capture groups as an object."},
	{Name: "regexSplit", Package: "extregex", Signature: "(str: string, pattern: string) → array", Description: "Splits str at each occurrence of pattern."},
	{Name: "regexReplaceAll", Package: "extregex", Signature: "(str: string, pattern: string, repl: string) → string", Description: "Replaces all matches of pattern with repl."},
	{Name: "regexCount", Package: "extregex", Signature: "(str: string, pattern: string) → number", Description: "Returns the number of non-overlapping matches."},
	{Name: "regexTest", Package: "extregex", Signature: "(str: string, pattern: string) → boolean", Description: "Returns true if pattern matches anywhere in str."},
	{Name: "regexExtract", Package: "extregex", Signature: "(str: string, pattern: string [, group: number]) → string", Description: "Returns the first match or the specified capture group."},
	// exttypes
	{Name: "isString", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is a string."},
	{Name: "isNumber", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is a number."},
	{Name: "isBoolean", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is a boolean."},
	{Name: "isArray", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is an array."},
	{Name: "isObject", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is an object."},
	{Name: "isNull", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is null."},
	{Name: "isUndefined", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true if v is undefined (null in gnata)."},
	{Name: "isEmpty", Package: "exttypes", Signature: "(v: any) → boolean", Description: "Returns true for null, empty string, empty array, empty object."},
	{Name: "default", Package: "exttypes", Signature: "(v: any, d: any) → any", Description: "Returns v if non-null, else d."},
	{Name: "identity", Package: "exttypes", Signature: "(v: any) → any", Description: "Returns v unchanged."},
	// exttext
	{Name: "wordCount", Package: "exttext", Signature: "(s: string) → number", Description: "Number of words in the string."},
	{Name: "charCount", Package: "exttext", Signature: "(s: string) → number", Description: "Number of Unicode code points in the string."},
	{Name: "sentenceCount", Package: "exttext", Signature: "(s: string) → number", Description: "Approximate number of sentences."},
	{Name: "readingTime", Package: "exttext", Signature: "(s: string) → number", Description: "Estimated reading time in seconds at 200 wpm."},
	{Name: "wordFrequency", Package: "exttext", Signature: "(s: string) → object", Description: "Map of word to occurrence count."},
	{Name: "ngrams", Package: "exttext", Signature: "(s: string, n: number) → array", Description: "Array of n-gram strings."},
	{Name: "excerpt", Package: "exttext", Signature: "(s: string, maxLen: number) → string", Description: "Truncate string to maxLen runes, appending ellipsis if truncated."},
	{Name: "stripHTML", Package: "exttext", Signature: "(s: string) → string", Description: "Removes HTML tags from a string."},
	{Name: "normalizeWhitespace", Package: "exttext", Signature: "(s: string) → string", Description: "Collapses consecutive whitespace to a single space."},
	// extvalidate
	{Name: "isEmail", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str is a valid email address."},
	{Name: "isURL", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str is a valid http/https/ftp URL."},
	{Name: "isUUID", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str is a valid UUID."},
	{Name: "isIPv4", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str is a valid IPv4 address."},
	{Name: "isIPv6", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str is a valid IPv6 address."},
	{Name: "isAlpha", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str contains only letters."},
	{Name: "isAlphanumeric", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str contains only letters and digits."},
	{Name: "isNumericStr", Package: "extvalidate", Signature: "(str: string) → boolean", Description: "Returns true if str can be parsed as a number."},
	{Name: "matchesRegex", Package: "extvalidate", Signature: "(str: string, pattern: string) → boolean", Description: "Returns true if str matches the RE2 pattern."},
	{Name: "inSet", Package: "extvalidate", Signature: "(v: any, set: array) → boolean", Description: "Returns true if v is present in set."},
	{Name: "minLen", Package: "extvalidate", Signature: "(str: string, n: number) → boolean", Description: "Returns true if string length ≥ n."},
	{Name: "maxLen", Package: "extvalidate", Signature: "(str: string, n: number) → boolean", Description: "Returns true if string length ≤ n."},
	{Name: "minItems", Package: "extvalidate", Signature: "(arr: array, n: number) → boolean", Description: "Returns true if array length ≥ n."},
	{Name: "maxItems", Package: "extvalidate", Signature: "(arr: array, n: number) → boolean", Description: "Returns true if array length ≤ n."},
}

// Catalog returns metadata for all extension functions, sorted by name.
func Catalog() []FuncMeta {
	out := make([]FuncMeta, len(catalog))
	copy(out, catalog)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// CatalogByPackage returns extension function metadata grouped by package name.
func CatalogByPackage() map[string][]FuncMeta {
	m := make(map[string][]FuncMeta)
	for _, f := range catalog {
		m[f.Package] = append(m[f.Package], f)
	}
	return m
}
