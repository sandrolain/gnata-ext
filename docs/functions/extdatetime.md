# extdatetime — Date & Time Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extdatetime`

All functions accept and return timestamps as **Unix milliseconds** (`float64`), compatible with JSONata's built-in `$millis()` and `$toMillis()`.

Accepted unit strings (singular and plural forms both work):
`year` · `month` · `day` · `hour` · `minute` · `second` · `millisecond`

---

### `$dateAdd(timestamp, amount, unit)`

Adds *amount* units to a timestamp. Use a negative *amount* to subtract.

```jsonata
/* 2024-01-15T12:30:00Z = 1705319400000 ms */
$dateAdd(1705319400000, 1, "day")     /* +1 day  */
$dateAdd(1705319400000, -2, "month")  /* -2 months */
$dateAdd(1705319400000, 3, "hour")
```

---

### `$dateDiff(t1, t2, unit)`

Returns the whole-unit difference `t2 − t1`.

```jsonata
$dateDiff(1705319400000, 1705405800000, "hour")   /* → 24 */
$dateDiff(1705319400000, 1705405800000, "day")    /* → 1 */
```

---

### `$dateComponents(timestamp)`

Returns an object with the UTC components of the timestamp.

```jsonata
$dateComponents(1705319400000)
/* → {
     "year": 2024, "month": 1, "day": 15,
     "hour": 12, "minute": 30, "second": 0,
     "millisecond": 0, "weekday": 1
   } */
```

`weekday` follows Go's `time.Weekday`: Sunday=0, Monday=1, …, Saturday=6.

---

### `$dateStartOf(timestamp, unit)`

Returns the timestamp at the start of the given period.

```jsonata
$dateStartOf(1705319400000, "day")    /* 2024-01-15T00:00:00.000Z */
$dateStartOf(1705319400000, "month")  /* 2024-01-01T00:00:00.000Z */
$dateStartOf(1705319400000, "year")   /* 2024-01-01T00:00:00.000Z */
```

---

### `$dateEndOf(timestamp, unit)`

Returns the timestamp at the last millisecond of the given period.

```jsonata
$dateEndOf(1705319400000, "day")      /* 2024-01-15T23:59:59.999Z */
$dateEndOf(1705319400000, "month")    /* 2024-01-31T23:59:59.999Z */
```

---
