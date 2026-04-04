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
