# Function Reference

Complete documentation for every function provided by `gnata-ext`, split by package.

All functions are registered as JSONata custom functions and are called with a leading `$` in expressions (e.g. `$uuid()`). Timestamps are always **Unix milliseconds** (`float64`) unless stated otherwise.

| Package | Description | Reference |
|---|---|---|
| `extarray` | Array utilities (slicing, set ops, ranges, windows) | [functions/extarray.md](functions/extarray.md) |
| `extcrypto` | UUID generation, hash, HMAC | [functions/extcrypto.md](functions/extcrypto.md) |
| `extdatetime` | Date arithmetic and decomposition | [functions/extdatetime.md](functions/extdatetime.md) |
| `extformat` | CSV parsing/serialisation, template rendering | [functions/extformat.md](functions/extformat.md) |
| `extnumeric` | Math, trigonometry, statistics | [functions/extnumeric.md](functions/extnumeric.md) |
| `extobject` | Object manipulation | [functions/extobject.md](functions/extobject.md) |
| `extstring` | String utilities and case conversion | [functions/extstring.md](functions/extstring.md) |
| `exttypes` | Type predicates and coercion helpers | [functions/exttypes.md](functions/exttypes.md) |
| `extpath` | Dot-path read/write on nested objects | [functions/extpath.md](functions/extpath.md) |
| `extvalidate` | Validation predicates (email, URL, IP, regex …) | [functions/extvalidate.md](functions/extvalidate.md) |
| `extjson` | JSON parse/stringify, diff, patch, pointer | [functions/extjson.md](functions/extjson.md) |
| `extgeo` | Geographic calculations (haversine, bearing, bounding box) | [functions/extgeo.md](functions/extgeo.md) |
| `extnet` | IP address and CIDR utilities | [functions/extnet.md](functions/extnet.md) |
