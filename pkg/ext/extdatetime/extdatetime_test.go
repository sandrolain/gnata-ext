package extdatetime_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extdatetime"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

// epoch: 2024-01-15T12:30:00.000Z = 1705319400000 ms
const epoch = float64(1705319400000)

func TestDateAdd(t *testing.T) {
	f := extdatetime.DateAdd()
	// +1 day
	got, err := invoke(f, epoch, 1.0, "day")
	if err != nil {
		t.Fatalf("dateAdd day: %v", err)
	}
	want := epoch + 86400000
	if got.(float64) != want {
		t.Errorf("dateAdd +1 day: got %v, want %v", got, want)
	}
	// +1 hour
	got, err = invoke(f, epoch, 1.0, "hour")
	if err != nil {
		t.Fatalf("dateAdd hour: %v", err)
	}
	if got.(float64) != epoch+3600000 {
		t.Errorf("dateAdd +1 hour: got %v, want %v", got, epoch+3600000)
	}
}

func TestDateAddYear(t *testing.T) {
	f := extdatetime.DateAdd()
	got, err := invoke(f, epoch, 1.0, "year")
	if err != nil {
		t.Fatalf("dateAdd year: %v", err)
	}
	// 2025-01-15 should give ms different from epoch by ~365 days
	diff := got.(float64) - epoch
	if diff < 365*86400000 || diff > 366*86400000 {
		t.Errorf("dateAdd +1 year: diff=%v ms", diff)
	}
}

func TestDateDiff(t *testing.T) {
	f := extdatetime.DateDiff()
	t2 := epoch + 2*86400000 // +2 days
	got, err := invoke(f, epoch, t2, "day")
	if err != nil {
		t.Fatalf("dateDiff day: %v", err)
	}
	if got.(float64) != 2 {
		t.Errorf("dateDiff day: got %v, want 2", got)
	}
	got, err = invoke(f, epoch, t2, "millisecond")
	if err != nil {
		t.Fatalf("dateDiff ms: %v", err)
	}
	if got.(float64) != 2*86400000 {
		t.Errorf("dateDiff ms: got %v", got)
	}
}

func TestDateComponents(t *testing.T) {
	f := extdatetime.DateComponents()
	got, err := invoke(f, epoch)
	if err != nil {
		t.Fatalf("dateComponents: %v", err)
	}
	comp := got.(map[string]any)
	if comp["year"].(float64) != 2024 {
		t.Errorf("dateComponents year: got %v", comp["year"])
	}
	if comp["month"].(float64) != 1 {
		t.Errorf("dateComponents month: got %v", comp["month"])
	}
	if comp["day"].(float64) != 15 {
		t.Errorf("dateComponents day: got %v", comp["day"])
	}
}

func TestDateStartOf(t *testing.T) {
	f := extdatetime.DateStartOf()
	got, err := invoke(f, epoch, "day")
	if err != nil {
		t.Fatalf("dateStartOf day: %v", err)
	}
	// 2024-01-15 00:00:00 UTC
	want := float64(1705276800000)
	if got.(float64) != want {
		t.Errorf("dateStartOf day: got %v, want %v", got, want)
	}
}

func TestDateEndOf(t *testing.T) {
	f := extdatetime.DateEndOf()
	got, err := invoke(f, epoch, "day")
	if err != nil {
		t.Fatalf("dateEndOf day: %v", err)
	}
	// 2024-01-15 23:59:59.999... should be >= start of next day - 1ms
	startOfNext := float64(1705276800000 + 86400000)
	if got.(float64) < startOfNext-1001 || got.(float64) >= startOfNext {
		t.Errorf("dateEndOf day: got %v, expected near %v", got, startOfNext-1)
	}
}

func TestAll(t *testing.T) {
	all := extdatetime.All()
	expected := []string{"dateAdd", "dateDiff", "dateComponents", "dateStartOf", "dateEndOf"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing function %q", name)
		}
	}
}

// --- Additional coverage tests ---

func TestDateAddAllUnits(t *testing.T) {
	f := extdatetime.DateAdd()
	tests := []struct {
		amount float64
		unit   string
	}{
		{1, "month"},
		{1, "months"},
		{1, "year"},
		{1, "years"},
		{1, "days"},
		{1, "hours"},
		{1, "minute"},
		{1, "minutes"},
		{1, "second"},
		{1, "seconds"},
		{1, "millisecond"},
		{1, "milliseconds"},
	}
	for _, tc := range tests {
		got, err := invoke(f, epoch, tc.amount, tc.unit)
		if err != nil {
			t.Errorf("dateAdd %q: unexpected error: %v", tc.unit, err)
		}
		if got.(float64) <= epoch {
			t.Errorf("dateAdd %q: result should be > epoch", tc.unit)
		}
	}
}

func TestDateAddErrors(t *testing.T) {
	f := extdatetime.DateAdd()

	// too few args
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	// bad timestamp
	_, err = invoke(f, "notanumber", 1.0, "day")
	if err == nil {
		t.Error("expected error for non-numeric timestamp")
	}
	// bad amount
	_, err = invoke(f, epoch, "notanumber", "day")
	if err == nil {
		t.Error("expected error for non-numeric amount")
	}
	// bad unit type
	_, err = invoke(f, epoch, 1.0, 42)
	if err == nil {
		t.Error("expected error for non-string unit")
	}
	// unknown unit
	_, err = invoke(f, epoch, 1.0, "decade")
	if err == nil {
		t.Error("expected error for unknown unit")
	}
}

func TestDateDiffAllUnits(t *testing.T) {
	f := extdatetime.DateDiff()
	t2 := epoch + 90*86400000 // +90 days (~3 months)

	units := []string{"second", "seconds", "minute", "minutes", "hour", "hours", "month", "months", "year", "years"}
	for _, unit := range units {
		_, err := invoke(f, epoch, t2, unit)
		if err != nil {
			t.Errorf("dateDiff %q: unexpected error: %v", unit, err)
		}
	}
}

func TestDateDiffErrors(t *testing.T) {
	f := extdatetime.DateDiff()

	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad", epoch, "day")
	if err == nil {
		t.Error("expected error for bad t1")
	}
	_, err = invoke(f, epoch, "bad", "day")
	if err == nil {
		t.Error("expected error for bad t2")
	}
	_, err = invoke(f, epoch, epoch, 99)
	if err == nil {
		t.Error("expected error for non-string unit")
	}
	_, err = invoke(f, epoch, epoch, "decade")
	if err == nil {
		t.Error("expected error for unknown unit")
	}
}

func TestDateDiffYMDNegative(t *testing.T) {
	// dateDiffYMD is tested indirectly via DateDiff months/years
	f := extdatetime.DateDiff()
	// t2 < t1 — negative year/month diffs
	earlier := epoch - 400*86400000
	got, err := invoke(f, epoch, earlier, "year")
	if err != nil {
		t.Fatalf("dateDiff year negative: %v", err)
	}
	if got.(float64) >= 0 {
		t.Errorf("expected negative year diff, got %v", got)
	}
	got, err = invoke(f, epoch, earlier, "month")
	if err != nil {
		t.Fatalf("dateDiff month negative: %v", err)
	}
	if got.(float64) >= 0 {
		t.Errorf("expected negative month diff, got %v", got)
	}
}

func TestDateComponentsErrors(t *testing.T) {
	f := extdatetime.DateComponents()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-numeric timestamp")
	}
}

func TestDateStartOfAllUnits(t *testing.T) {
	f := extdatetime.DateStartOf()
	for _, unit := range []string{"year", "month", "hour", "minute", "second"} {
		got, err := invoke(f, epoch, unit)
		if err != nil {
			t.Errorf("dateStartOf %q: %v", unit, err)
		}
		if got.(float64) > epoch {
			t.Errorf("dateStartOf %q: start should be <= epoch", unit)
		}
	}
}

func TestDateStartOfErrors(t *testing.T) {
	f := extdatetime.DateStartOf()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad", "day")
	if err == nil {
		t.Error("expected error for bad timestamp")
	}
	_, err = invoke(f, epoch, 99)
	if err == nil {
		t.Error("expected error for non-string unit")
	}
	_, err = invoke(f, epoch, "decade")
	if err == nil {
		t.Error("expected error for unknown unit")
	}
}

func TestDateEndOfAllUnits(t *testing.T) {
	f := extdatetime.DateEndOf()
	for _, unit := range []string{"year", "month", "hour", "minute", "second"} {
		got, err := invoke(f, epoch, unit)
		if err != nil {
			t.Errorf("dateEndOf %q: %v", unit, err)
		}
		if got.(float64) < epoch {
			t.Errorf("dateEndOf %q: end should be >= epoch", unit)
		}
	}
}

func TestDateEndOfErrors(t *testing.T) {
	f := extdatetime.DateEndOf()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad", "day")
	if err == nil {
		t.Error("expected error for bad timestamp")
	}
	_, err = invoke(f, epoch, 99)
	if err == nil {
		t.Error("expected error for non-string unit")
	}
	_, err = invoke(f, epoch, "decade")
	if err == nil {
		t.Error("expected error for unknown unit")
	}
}
