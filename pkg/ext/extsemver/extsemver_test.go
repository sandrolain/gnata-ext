package extsemver_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extsemver"
)

// ---------- SemverParse ----------

func TestSemverParse_Basic(t *testing.T) {
	fn := extsemver.SemverParse()
	got, err := fn([]any{"1.2.3"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["major"] != float64(1) || obj["minor"] != float64(2) || obj["patch"] != float64(3) {
		t.Errorf("unexpected parse result: %v", obj)
	}
}

func TestSemverParse_WithPrerelease(t *testing.T) {
	fn := extsemver.SemverParse()
	got, _ := fn([]any{"2.0.0-alpha.1"}, nil)
	obj := got.(map[string]any)
	if obj["prerelease"] != "alpha.1" {
		t.Errorf("expected prerelease=alpha.1, got %v", obj["prerelease"])
	}
}

func TestSemverParse_WithMetadata(t *testing.T) {
	fn := extsemver.SemverParse()
	got, _ := fn([]any{"1.0.0+build.42"}, nil)
	obj := got.(map[string]any)
	if obj["metadata"] != "build.42" {
		t.Errorf("expected metadata=build.42, got %v", obj["metadata"])
	}
}

func TestSemverParse_Invalid(t *testing.T) {
	fn := extsemver.SemverParse()
	_, err := fn([]any{"not-a-version"}, nil)
	if err == nil {
		t.Error("expected error for invalid version")
	}
}

func TestSemverParse_NoArgs(t *testing.T) {
	fn := extsemver.SemverParse()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverParse_WrongType(t *testing.T) {
	fn := extsemver.SemverParse()
	_, err := fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error for non-string")
	}
}

// ---------- SemverCompare ----------

func TestSemverCompare_Less(t *testing.T) {
	fn := extsemver.SemverCompare()
	got, err := fn([]any{"1.0.0", "2.0.0"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != float64(-1) {
		t.Errorf("expected -1, got %v", got)
	}
}

func TestSemverCompare_Equal(t *testing.T) {
	fn := extsemver.SemverCompare()
	got, _ := fn([]any{"1.0.0", "1.0.0"}, nil)
	if got != float64(0) {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestSemverCompare_Greater(t *testing.T) {
	fn := extsemver.SemverCompare()
	got, _ := fn([]any{"2.0.0", "1.0.0"}, nil)
	if got != float64(1) {
		t.Errorf("expected 1, got %v", got)
	}
}

func TestSemverCompare_NoArgs(t *testing.T) {
	fn := extsemver.SemverCompare()
	_, err := fn([]any{"1.0.0"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverCompare_WrongTypeA(t *testing.T) {
	fn := extsemver.SemverCompare()
	_, err := fn([]any{1, "1.0.0"}, nil)
	if err == nil {
		t.Error("expected error for non-string a")
	}
}

func TestSemverCompare_WrongTypeB(t *testing.T) {
	fn := extsemver.SemverCompare()
	_, err := fn([]any{"1.0.0", false}, nil)
	if err == nil {
		t.Error("expected error for non-string b")
	}
}

func TestSemverCompare_InvalidA(t *testing.T) {
	fn := extsemver.SemverCompare()
	_, err := fn([]any{"bad", "1.0.0"}, nil)
	if err == nil {
		t.Error("expected error for invalid version a")
	}
}

func TestSemverCompare_InvalidB(t *testing.T) {
	fn := extsemver.SemverCompare()
	_, err := fn([]any{"1.0.0", "bad"}, nil)
	if err == nil {
		t.Error("expected error for invalid version b")
	}
}

// ---------- SemverSatisfies ----------

func TestSemverSatisfies_True(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	got, err := fn([]any{"1.5.0", ">=1.0.0 <2.0.0"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != true {
		t.Error("expected true")
	}
}

func TestSemverSatisfies_False(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	got, _ := fn([]any{"2.5.0", ">=1.0.0 <2.0.0"}, nil)
	if got != false {
		t.Error("expected false")
	}
}

func TestSemverSatisfies_NoArgs(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	_, err := fn([]any{"1.0.0"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverSatisfies_WrongTypeV(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	_, err := fn([]any{42, ">= 1.0.0"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverSatisfies_WrongTypeC(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	_, err := fn([]any{"1.0.0", 42}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverSatisfies_InvalidVersion(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	_, err := fn([]any{"bad", ">=1.0.0"}, nil)
	if err == nil {
		t.Error("expected error for invalid version")
	}
}

func TestSemverSatisfies_InvalidConstraint(t *testing.T) {
	fn := extsemver.SemverSatisfies()
	_, err := fn([]any{"1.0.0", "not-a-constraint-!!"}, nil)
	if err == nil {
		t.Error("expected error for invalid constraint")
	}
}

// ---------- SemverBump ----------

func TestSemverBump_Patch(t *testing.T) {
	fn := extsemver.SemverBump()
	got, err := fn([]any{"1.2.3", "patch"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.2.4" {
		t.Errorf("expected '1.2.4', got %v", got)
	}
}

func TestSemverBump_Minor(t *testing.T) {
	fn := extsemver.SemverBump()
	got, _ := fn([]any{"1.2.3", "minor"}, nil)
	if got != "1.3.0" {
		t.Errorf("expected '1.3.0', got %v", got)
	}
}

func TestSemverBump_Major(t *testing.T) {
	fn := extsemver.SemverBump()
	got, _ := fn([]any{"1.2.3", "major"}, nil)
	if got != "2.0.0" {
		t.Errorf("expected '2.0.0', got %v", got)
	}
}

func TestSemverBump_InvalidPart(t *testing.T) {
	fn := extsemver.SemverBump()
	_, err := fn([]any{"1.2.3", "build"}, nil)
	if err == nil {
		t.Error("expected error for invalid part")
	}
}

func TestSemverBump_NoArgs(t *testing.T) {
	fn := extsemver.SemverBump()
	_, err := fn([]any{"1.0.0"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverBump_WrongTypeV(t *testing.T) {
	fn := extsemver.SemverBump()
	_, err := fn([]any{true, "patch"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverBump_WrongTypePart(t *testing.T) {
	fn := extsemver.SemverBump()
	_, err := fn([]any{"1.0.0", 42}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverBump_InvalidVersion(t *testing.T) {
	fn := extsemver.SemverBump()
	_, err := fn([]any{"bad", "patch"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- SemverSort ----------

func TestSemverSort_Basic(t *testing.T) {
	fn := extsemver.SemverSort()
	arr := []any{"2.0.0", "1.0.0", "1.5.0"}
	got, err := fn([]any{arr}, nil)
	if err != nil {
		t.Fatal(err)
	}
	sorted := got.([]any)
	if sorted[0] != "1.0.0" || sorted[1] != "1.5.0" || sorted[2] != "2.0.0" {
		t.Errorf("unexpected sort order: %v", sorted)
	}
}

func TestSemverSort_Empty(t *testing.T) {
	fn := extsemver.SemverSort()
	got, _ := fn([]any{[]any{}}, nil)
	arr := got.([]any)
	if len(arr) != 0 {
		t.Error("expected empty")
	}
}

func TestSemverSort_NoArgs(t *testing.T) {
	fn := extsemver.SemverSort()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverSort_NonArray(t *testing.T) {
	fn := extsemver.SemverSort()
	_, err := fn([]any{"1.0.0"}, nil)
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestSemverSort_InvalidElement(t *testing.T) {
	fn := extsemver.SemverSort()
	_, err := fn([]any{[]any{"1.0.0", "bad"}}, nil)
	if err == nil {
		t.Error("expected error for invalid element")
	}
}

func TestSemverSort_NonStringElement(t *testing.T) {
	fn := extsemver.SemverSort()
	_, err := fn([]any{[]any{42}}, nil)
	if err == nil {
		t.Error("expected error for non-string element")
	}
}

// ---------- SemverMax ----------

func TestSemverMax_Basic(t *testing.T) {
	fn := extsemver.SemverMax()
	arr := []any{"1.0.0", "3.0.0", "2.0.0"}
	got, err := fn([]any{arr}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "3.0.0" {
		t.Errorf("expected '3.0.0', got %v", got)
	}
}

func TestSemverMax_Empty(t *testing.T) {
	fn := extsemver.SemverMax()
	_, err := fn([]any{[]any{}}, nil)
	if err == nil {
		t.Error("expected error for empty array")
	}
}

func TestSemverMax_NoArgs(t *testing.T) {
	fn := extsemver.SemverMax()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverMax_NonArray(t *testing.T) {
	fn := extsemver.SemverMax()
	_, err := fn([]any{"1.0.0"}, nil)
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestSemverMax_NonStringElement(t *testing.T) {
	fn := extsemver.SemverMax()
	_, err := fn([]any{[]any{42}}, nil)
	if err == nil {
		t.Error("expected error for non-string element")
	}
}

func TestSemverMax_InvalidElement(t *testing.T) {
	fn := extsemver.SemverMax()
	_, err := fn([]any{[]any{"bad"}}, nil)
	if err == nil {
		t.Error("expected error for invalid element")
	}
}

// ---------- SemverMin ----------

func TestSemverMin_Basic(t *testing.T) {
	fn := extsemver.SemverMin()
	arr := []any{"3.0.0", "1.0.0", "2.0.0"}
	got, err := fn([]any{arr}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.0.0" {
		t.Errorf("expected '1.0.0', got %v", got)
	}
}

func TestSemverMin_Empty(t *testing.T) {
	fn := extsemver.SemverMin()
	_, err := fn([]any{[]any{}}, nil)
	if err == nil {
		t.Error("expected error for empty array")
	}
}

func TestSemverMin_NoArgs(t *testing.T) {
	fn := extsemver.SemverMin()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSemverMin_NonArray(t *testing.T) {
	fn := extsemver.SemverMin()
	_, err := fn([]any{"1.0.0"}, nil)
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestSemverMin_NonStringElement(t *testing.T) {
	fn := extsemver.SemverMin()
	_, err := fn([]any{[]any{true}}, nil)
	if err == nil {
		t.Error("expected error for non-string element")
	}
}

func TestSemverMin_InvalidElement(t *testing.T) {
	fn := extsemver.SemverMin()
	_, err := fn([]any{[]any{"bad"}}, nil)
	if err == nil {
		t.Error("expected error for invalid element")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := extsemver.All()
	expected := []string{
		"semverParse", "semverCompare", "semverSatisfies",
		"semverBump", "semverSort", "semverMax", "semverMin",
	}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
