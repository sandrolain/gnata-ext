package extlogic_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extlogic"
)

// ---------- IfElse ----------

func TestIfElse_True(t *testing.T) {
	fn := extlogic.IfElse()
	got, err := fn([]any{true, "yes", "no"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "yes" {
		t.Errorf("expected 'yes', got %v", got)
	}
}

func TestIfElse_False(t *testing.T) {
	fn := extlogic.IfElse()
	got, _ := fn([]any{false, "yes", "no"}, nil)
	if got != "no" {
		t.Errorf("expected 'no', got %v", got)
	}
}

func TestIfElse_Numeric(t *testing.T) {
	fn := extlogic.IfElse()
	got, _ := fn([]any{float64(0), "yes", "no"}, nil)
	if got != "no" {
		t.Errorf("expected 'no' for 0, got %v", got)
	}
	got, _ = fn([]any{float64(1), "yes", "no"}, nil)
	if got != "yes" {
		t.Errorf("expected 'yes' for 1, got %v", got)
	}
}

func TestIfElse_EmptyString(t *testing.T) {
	fn := extlogic.IfElse()
	got, _ := fn([]any{"", "yes", "no"}, nil)
	if got != "no" {
		t.Errorf("expected 'no' for empty string, got %v", got)
	}
}

func TestIfElse_NilCond(t *testing.T) {
	fn := extlogic.IfElse()
	got, _ := fn([]any{nil, "yes", "no"}, nil)
	if got != "no" {
		t.Errorf("expected 'no' for nil cond, got %v", got)
	}
}

func TestIfElse_MissingArgs(t *testing.T) {
	fn := extlogic.IfElse()
	_, err := fn([]any{true, "yes"}, nil)
	if err == nil {
		t.Error("expected error for missing third arg")
	}
}

func TestIfElse_NilValues(t *testing.T) {
	fn := extlogic.IfElse()
	got, _ := fn([]any{true, nil, "fallback"}, nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

// ---------- When ----------

func TestWhen_True(t *testing.T) {
	fn := extlogic.When()
	got, err := fn([]any{true, "value"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "value" {
		t.Errorf("expected 'value', got %v", got)
	}
}

func TestWhen_False(t *testing.T) {
	fn := extlogic.When()
	got, _ := fn([]any{false, "value"}, nil)
	if got != nil {
		t.Errorf("expected nil when false, got %v", got)
	}
}

func TestWhen_MissingArgs(t *testing.T) {
	fn := extlogic.When()
	_, err := fn([]any{true}, nil)
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}

func TestWhen_EmptySliceFalsy(t *testing.T) {
	fn := extlogic.When()
	got, _ := fn([]any{[]any{}, "value"}, nil)
	if got != nil {
		t.Errorf("expected nil for empty slice condition, got %v", got)
	}
}

// ---------- Unless ----------

func TestUnless_False(t *testing.T) {
	fn := extlogic.Unless()
	got, err := fn([]any{false, "value"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "value" {
		t.Errorf("expected 'value', got %v", got)
	}
}

func TestUnless_True(t *testing.T) {
	fn := extlogic.Unless()
	got, _ := fn([]any{true, "value"}, nil)
	if got != nil {
		t.Errorf("expected nil when true, got %v", got)
	}
}

func TestUnless_MissingArgs(t *testing.T) {
	fn := extlogic.Unless()
	_, err := fn([]any{false}, nil)
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}

// ---------- Switch ----------

func TestSwitch_Match(t *testing.T) {
	fn := extlogic.Switch()
	cases := map[string]any{"active": "A", "inactive": "I"}
	got, err := fn([]any{"active", cases}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "A" {
		t.Errorf("expected 'A', got %v", got)
	}
}

func TestSwitch_NoMatchWithDefault(t *testing.T) {
	fn := extlogic.Switch()
	cases := map[string]any{"active": "A"}
	got, _ := fn([]any{"unknown", cases, "U"}, nil)
	if got != "U" {
		t.Errorf("expected default 'U', got %v", got)
	}
}

func TestSwitch_NoMatchNoDefault(t *testing.T) {
	fn := extlogic.Switch()
	cases := map[string]any{"active": "A"}
	got, _ := fn([]any{"unknown", cases}, nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestSwitch_NumericKey(t *testing.T) {
	fn := extlogic.Switch()
	cases := map[string]any{"1": "one", "2": "two"}
	got, _ := fn([]any{float64(2), cases}, nil)
	if got != "two" {
		t.Errorf("expected 'two', got %v", got)
	}
}

func TestSwitch_MissingArgs(t *testing.T) {
	fn := extlogic.Switch()
	_, err := fn([]any{"x"}, nil)
	if err == nil {
		t.Error("expected error for missing cases arg")
	}
}

func TestSwitch_InvalidCases(t *testing.T) {
	fn := extlogic.Switch()
	_, err := fn([]any{"x", "not-a-map"}, nil)
	if err == nil {
		t.Error("expected error for non-object cases")
	}
}

// ---------- Coalesce ----------

func TestCoalesce_FirstNonNil(t *testing.T) {
	fn := extlogic.Coalesce()
	got, err := fn([]any{nil, "", float64(0), "found", "other"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "found" {
		t.Errorf("expected 'found', got %v", got)
	}
}

func TestCoalesce_AllNil(t *testing.T) {
	fn := extlogic.Coalesce()
	got, _ := fn([]any{nil, nil, nil}, nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestCoalesce_FirstIsNonEmpty(t *testing.T) {
	fn := extlogic.Coalesce()
	got, _ := fn([]any{"first", "second"}, nil)
	if got != "first" {
		t.Errorf("expected 'first', got %v", got)
	}
}

func TestCoalesce_BoolTrue(t *testing.T) {
	fn := extlogic.Coalesce()
	got, _ := fn([]any{false, true, "x"}, nil)
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

func TestCoalesce_MissingArgs(t *testing.T) {
	fn := extlogic.Coalesce()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for empty args")
	}
}

func TestCoalesce_MapNonEmpty(t *testing.T) {
	fn := extlogic.Coalesce()
	m := map[string]any{"a": 1}
	got, _ := fn([]any{nil, m}, nil)
	if got == nil {
		t.Error("expected non-nil map")
	}
}

// ---------- Tap ----------

func TestTap(t *testing.T) {
	fn := extlogic.Tap()
	got, err := fn([]any{"hello"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Errorf("expected 'hello', got %v", got)
	}
}

func TestTap_NilValue(t *testing.T) {
	fn := extlogic.Tap()
	got, _ := fn([]any{nil}, nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestTap_MissingArgs(t *testing.T) {
	fn := extlogic.Tap()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for empty args")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := extlogic.All()
	expected := []string{"ifElse", "when", "unless", "switch", "coalesce", "tap"}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
