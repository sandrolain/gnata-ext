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
