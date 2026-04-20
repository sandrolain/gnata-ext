// Package extdatetime provides extended date/time functions for gnata.
//
// All dates are expressed as Unix milliseconds (float64), matching the
// JSONata $millis() / $toMillis() convention.
//
// Functions
//
//   - $dateAdd(timestamp, amount, unit)     – add/subtract duration
//   - $dateDiff(t1, t2, unit)              – difference between two timestamps
//   - $dateComponents(timestamp)           – map of year/month/day/… components
//   - $dateStartOf(timestamp, unit)        – start of the given time unit
//   - $dateEndOf(timestamp, unit)          – end of the given time unit
//   - $dateFormat(timestamp, layout)       – format timestamp using Go layout
//   - $dateParse(str, layout)              – parse string to timestamp
//   - $dateIsBefore(t1, t2)               – true if t1 < t2
//   - $dateIsAfter(t1, t2)               – true if t1 > t2
//   - $dateIsBetween(t, start, end)        – true if start <= t <= end
//   - $dateWeek(timestamp)                – ISO week number
//   - $dateQuarter(timestamp)             – quarter number (1-4)
//   - $dateDayOfYear(timestamp)           – day of year (1-366)
//   - $isLeapYear(timestamp)              – true if the year is a leap year
//   - $daysInMonth(timestamp)             – number of days in the timestamp's month
package extdatetime

import (
	"fmt"
	"strings"
	"time"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extended date/time functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"dateAdd":        DateAdd(),
		"dateDiff":       DateDiff(),
		"dateComponents": DateComponents(),
		"dateStartOf":    DateStartOf(),
		"dateEndOf":      DateEndOf(),
		"dateFormat":     DateFormat(),
		"dateParse":      DateParse(),
		"dateIsBefore":   DateIsBefore(),
		"dateIsAfter":    DateIsAfter(),
		"dateIsBetween":  DateIsBetween(),
		"dateWeek":       DateWeek(),
		"dateQuarter":    DateQuarter(),
		"dateDayOfYear":  DateDayOfYear(),
		"isLeapYear":     IsLeapYear(),
		"daysInMonthOf":  DaysInMonthOf(),
	}
}

// DateAdd returns the CustomFunc for $dateAdd(timestamp, amount, unit).
// unit: "year"/"month"/"day"/"hour"/"minute"/"second"/"millisecond"
func DateAdd() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$dateAdd: requires 3 arguments (timestamp, amount, unit)")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateAdd: timestamp: %w", err)
		}
		amount, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$dateAdd: amount: %w", err)
		}
		unit, ok := args[2].(string)
		if !ok {
			return nil, fmt.Errorf("$dateAdd: unit must be a string")
		}
		t := msToTime(ms)
		n := int(amount)
		switch strings.ToLower(unit) {
		case "year", "years":
			t = t.AddDate(n, 0, 0)
		case "month", "months":
			t = t.AddDate(0, n, 0)
		case "day", "days":
			t = t.AddDate(0, 0, n)
		case "hour", "hours":
			t = t.Add(time.Duration(n) * time.Hour)
		case "minute", "minutes":
			t = t.Add(time.Duration(n) * time.Minute)
		case "second", "seconds":
			t = t.Add(time.Duration(n) * time.Second)
		case "millisecond", "milliseconds":
			t = t.Add(time.Duration(n) * time.Millisecond)
		default:
			return nil, fmt.Errorf("$dateAdd: unknown unit %q", unit)
		}
		return timeToMs(t), nil
	}
}

// DateDiff returns the CustomFunc for $dateDiff(t1, t2, unit).
func DateDiff() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$dateDiff: requires 3 arguments (t1, t2, unit)")
		}
		ms1, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateDiff: t1: %w", err)
		}
		ms2, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$dateDiff: t2: %w", err)
		}
		unit, ok := args[2].(string)
		if !ok {
			return nil, fmt.Errorf("$dateDiff: unit must be a string")
		}
		t1 := msToTime(ms1)
		t2 := msToTime(ms2)
		switch strings.ToLower(unit) {
		case "millisecond", "milliseconds":
			return math64(t2.UnixMilli() - t1.UnixMilli()), nil
		case "second", "seconds":
			return math64(t2.Unix() - t1.Unix()), nil
		case "minute", "minutes":
			return math64((t2.Unix() - t1.Unix()) / 60), nil
		case "hour", "hours":
			return math64((t2.Unix() - t1.Unix()) / 3600), nil
		case "day", "days":
			return math64((t2.Unix() - t1.Unix()) / 86400), nil
		case "month", "months":
			years, months, _ := dateDiffYMD(t1, t2)
			return float64(years*12 + months), nil
		case "year", "years":
			years, _, _ := dateDiffYMD(t1, t2)
			return float64(years), nil
		default:
			return nil, fmt.Errorf("$dateDiff: unknown unit %q", unit)
		}
	}
}

// DateComponents returns the CustomFunc for $dateComponents(timestamp).
func DateComponents() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$dateComponents: requires 1 argument")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateComponents: %w", err)
		}
		t := msToTime(ms).UTC()
		return map[string]any{
			"year":        float64(t.Year()),
			"month":       float64(t.Month()),
			"day":         float64(t.Day()),
			"hour":        float64(t.Hour()),
			"minute":      float64(t.Minute()),
			"second":      float64(t.Second()),
			"millisecond": float64(t.Nanosecond() / 1e6),
			"weekday":     float64(t.Weekday()),
		}, nil
	}
}

// DateStartOf returns the CustomFunc for $dateStartOf(timestamp, unit).
func DateStartOf() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$dateStartOf: requires 2 arguments (timestamp, unit)")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateStartOf: %w", err)
		}
		unit, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$dateStartOf: unit must be a string")
		}
		t := msToTime(ms).UTC()
		var result time.Time
		switch strings.ToLower(unit) {
		case "year":
			result = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		case "month":
			result = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
		case "day":
			result = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		case "hour":
			result = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
		case "minute":
			result = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.UTC)
		case "second":
			result = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
		default:
			return nil, fmt.Errorf("$dateStartOf: unknown unit %q", unit)
		}
		return timeToMs(result), nil
	}
}

// DateEndOf returns the CustomFunc for $dateEndOf(timestamp, unit).
func DateEndOf() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$dateEndOf: requires 2 arguments (timestamp, unit)")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateEndOf: %w", err)
		}
		unit, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$dateEndOf: unit must be a string")
		}
		t := msToTime(ms).UTC()
		var result time.Time
		switch strings.ToLower(unit) {
		case "year":
			result = time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
		case "month":
			// First day of next month minus 1 ns
			first := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			result = first.Add(-time.Nanosecond)
		case "day":
			result = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
		case "hour":
			result = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 999999999, time.UTC)
		case "minute":
			result = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, 999999999, time.UTC)
		case "second":
			result = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 999999999, time.UTC)
		default:
			return nil, fmt.Errorf("$dateEndOf: unknown unit %q", unit)
		}
		return timeToMs(result), nil
	}
}

// msToTime converts Unix milliseconds to time.Time in UTC.
func msToTime(ms float64) time.Time {
	sec := int64(ms) / 1000
	nsec := (int64(ms) % 1000) * int64(time.Millisecond)
	return time.Unix(sec, nsec).UTC()
}

// timeToMs converts time.Time to Unix milliseconds as float64.
func timeToMs(t time.Time) float64 {
	return float64(t.UnixMilli())
}

// math64 converts int64 to float64 for result values.
func math64(n int64) float64 {
	return float64(n)
}

// dateDiffYMD computes the calendar difference in years, months, and days.
func dateDiffYMD(from, to time.Time) (years, months, days int) {
	y1, m1, d1 := from.Date()
	y2, m2, d2 := to.Date()
	years = y2 - y1
	months = int(m2 - m1)
	days = d2 - d1
	if days < 0 {
		months--
		// Days in previous month relative to to
		prev := to.AddDate(0, -1, 0)
		days += daysInMonth(prev.Year(), prev.Month())
	}
	if months < 0 {
		years--
		months += 12
	}
	return
}

// daysInMonth returns the number of days in the given month.
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// DateFormat returns a CustomFunc for $dateFormat(timestamp, layout).
// layout uses Go time layout string (e.g. "2006-01-02").
func DateFormat() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$dateFormat: requires 2 arguments (timestamp, layout)")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateFormat: timestamp: %w", err)
		}
		layout, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$dateFormat: layout must be a string")
		}
		return msToTime(ms).Format(layout), nil
	}
}

// DateParse returns a CustomFunc for $dateParse(str, layout).
// Returns Unix milliseconds as float64.
func DateParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$dateParse: requires 2 arguments (str, layout)")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$dateParse: str must be a string")
		}
		layout, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$dateParse: layout must be a string")
		}
		t, err := time.Parse(layout, str)
		if err != nil {
			return nil, fmt.Errorf("$dateParse: %w", err)
		}
		return timeToMs(t), nil
	}
}

// DateIsBefore returns a CustomFunc for $dateIsBefore(t1, t2).
func DateIsBefore() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$dateIsBefore: requires 2 arguments")
		}
		t1, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateIsBefore: t1: %w", err)
		}
		t2, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$dateIsBefore: t2: %w", err)
		}
		return t1 < t2, nil
	}
}

// DateIsAfter returns a CustomFunc for $dateIsAfter(t1, t2).
func DateIsAfter() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$dateIsAfter: requires 2 arguments")
		}
		t1, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateIsAfter: t1: %w", err)
		}
		t2, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$dateIsAfter: t2: %w", err)
		}
		return t1 > t2, nil
	}
}

// DateIsBetween returns a CustomFunc for $dateIsBetween(t, start, end).
func DateIsBetween() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$dateIsBetween: requires 3 arguments (t, start, end)")
		}
		t, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateIsBetween: t: %w", err)
		}
		start, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$dateIsBetween: start: %w", err)
		}
		end, err := extutil.ToFloat(args[2])
		if err != nil {
			return nil, fmt.Errorf("$dateIsBetween: end: %w", err)
		}
		return t >= start && t <= end, nil
	}
}

// DateWeek returns a CustomFunc for $dateWeek(timestamp).
// Returns ISO week number (1-53).
func DateWeek() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$dateWeek: requires 1 argument")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateWeek: timestamp: %w", err)
		}
		_, week := msToTime(ms).ISOWeek()
		return float64(week), nil
	}
}

// DateQuarter returns a CustomFunc for $dateQuarter(timestamp).
// Returns quarter 1-4.
func DateQuarter() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$dateQuarter: requires 1 argument")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateQuarter: timestamp: %w", err)
		}
		month := int(msToTime(ms).Month())
		return float64((month-1)/3 + 1), nil
	}
}

// DateDayOfYear returns a CustomFunc for $dateDayOfYear(timestamp).
func DateDayOfYear() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$dateDayOfYear: requires 1 argument")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$dateDayOfYear: timestamp: %w", err)
		}
		return float64(msToTime(ms).YearDay()), nil
	}
}

// IsLeapYear returns a CustomFunc for $isLeapYear(timestamp).
func IsLeapYear() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$isLeapYear: requires 1 argument")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$isLeapYear: timestamp: %w", err)
		}
		year := msToTime(ms).Year()
		leap := (year%4 == 0 && year%100 != 0) || year%400 == 0
		return leap, nil
	}
}

// DaysInMonthOf returns a CustomFunc for $daysInMonthOf(timestamp).
func DaysInMonthOf() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$daysInMonthOf: requires 1 argument")
		}
		ms, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$daysInMonthOf: timestamp: %w", err)
		}
		t := msToTime(ms)
		return float64(daysInMonth(t.Year(), t.Month())), nil
	}
}
