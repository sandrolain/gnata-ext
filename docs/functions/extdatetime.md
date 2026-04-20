# extdatetime — Date & Time Utilities

**Import path:** `github.com/sandrolain/gnata-ext/pkg/ext/extdatetime`

All functions accept and return timestamps as **Unix milliseconds** (`float64`), compatible with JSONata's built-in `$millis()` and `$toMillis()`.

Accepted unit strings (singular and plural forms both work):
`year` · `month` · `day` · `hour` · `minute` · `second` · `millisecond`

> **Note on JSONata native alternatives:** `$dateFormat` / `$dateParse` overlap with `$fromMillis` / `$toMillis` but use **Go time layout strings** rather than XPath F&O picture strings. See [guides/jsonata-overlap.md](../guides/jsonata-overlap.md) for guidance.

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

### `$dateFormat(timestamp, layout)`

Formats *timestamp* (Unix milliseconds) into a string using a **Go time layout** (e.g. `"2006-01-02T15:04:05Z07:00"`).

> **JSONata native alternative:** `$fromMillis(timestamp, picture)` — same purpose, but uses XPath F&O picture strings (e.g. `"[Y0001]-[M01]-[D01]"`). See [guides/jsonata-overlap.md](../guides/jsonata-overlap.md).

```jsonata
$dateFormat(1705319400000, "2006-01-02")           /* → "2024-01-15" */
$dateFormat(1705319400000, "Jan 2, 2006 15:04")    /* → "Jan 15, 2024 12:30" */
```

---

### `$dateParse(str, layout)`

Parses a date/time string using a **Go time layout** and returns the result as Unix milliseconds.

> **JSONata native alternative:** `$toMillis(str, picture)` — same purpose, but uses XPath F&O picture strings. See [guides/jsonata-overlap.md](../guides/jsonata-overlap.md).

```jsonata
$dateParse("2024-01-15", "2006-01-02")   /* → 1705276800000 */
```

---

### `$dateIsBefore(t1, t2)`

Returns `true` if timestamp *t1* is strictly before *t2*.

```jsonata
$dateIsBefore(1705319400000, 1705319500000)   /* → true */
```

---

### `$dateIsAfter(t1, t2)`

Returns `true` if timestamp *t1* is strictly after *t2*.

```jsonata
$dateIsAfter(1705319500000, 1705319400000)   /* → true */
```

---

### `$dateIsBetween(t, start, end)`

Returns `true` if *t* is greater than or equal to *start* and less than or equal to *end*.

```jsonata
$dateIsBetween(1705319400000, 1705000000000, 1705500000000)   /* → true */
```

---

### `$dateWeek(timestamp)`

Returns the ISO 8601 week number (1–53) for the week containing *timestamp*.

```jsonata
$dateWeek(1705319400000)   /* → 3  (week 3 of 2024) */
```

---

### `$dateQuarter(timestamp)`

Returns the calendar quarter (1–4) in which *timestamp* falls.

```jsonata
$dateQuarter(1705319400000)   /* → 1  (January 2024) */
$dateQuarter(1718000000000)   /* → 2  (June 2024) */
```

---

### `$dateDayOfYear(timestamp)`

Returns the day of year (1–366) of *timestamp*.

```jsonata
$dateDayOfYear(1705319400000)   /* → 15  (15 January 2024) */
```

---

### `$isLeapYear(timestamp)`

Returns `true` if the year of *timestamp* is a leap year.

```jsonata
$isLeapYear(1705319400000)   /* → true  (2024 is a leap year) */
$isLeapYear(1672531200000)   /* → false (2023 is not) */
```

---

### `$daysInMonthOf(timestamp)`

Returns the number of days in the calendar month containing *timestamp*.

```jsonata
$daysInMonthOf(1705319400000)   /* → 31  (January 2024) */
$daysInMonthOf(1706745600000)   /* → 29  (February 2024 – leap year) */
```

---
