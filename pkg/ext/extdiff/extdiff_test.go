package extdiff_test

import (
	"sort"
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extdiff"
)

// --- helpers ---

func objKeys(obj map[string]any) []string {
	ks := make([]string, 0, len(obj))
	for k := range obj {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// ---------- Diff ----------

func TestDiff_Added(t *testing.T) {
	fn := extdiff.Diff()
	a := map[string]any{"x": float64(1)}
	b := map[string]any{"x": float64(1), "y": float64(2)}
	got, err := fn([]any{a, b}, nil)
	if err != nil {
		t.Fatal(err)
	}
	d := got.(map[string]any)
	added := d["added"].(map[string]any)
	if _, ok := added["y"]; !ok {
		t.Error("expected 'y' in added")
	}
	if len(d["removed"].(map[string]any)) != 0 {
		t.Error("expected empty removed")
	}
}

func TestDiff_Removed(t *testing.T) {
	fn := extdiff.Diff()
	a := map[string]any{"x": float64(1), "y": float64(2)}
	b := map[string]any{"x": float64(1)}
	got, _ := fn([]any{a, b}, nil)
	d := got.(map[string]any)
	removed := d["removed"].(map[string]any)
	if _, ok := removed["y"]; !ok {
		t.Error("expected 'y' in removed")
	}
}

func TestDiff_Changed(t *testing.T) {
	fn := extdiff.Diff()
	a := map[string]any{"x": float64(1)}
	b := map[string]any{"x": float64(99)}
	got, _ := fn([]any{a, b}, nil)
	d := got.(map[string]any)
	changed := d["changed"].(map[string]any)
	if entry, ok := changed["x"]; !ok {
		t.Error("expected 'x' in changed")
	} else {
		e := entry.(map[string]any)
		if e["from"] != float64(1) || e["to"] != float64(99) {
			t.Errorf("changed.x: unexpected from/to: %v", e)
		}
	}
}

func TestDiff_NoDiff(t *testing.T) {
	fn := extdiff.Diff()
	a := map[string]any{"x": float64(1)}
	got, _ := fn([]any{a, a}, nil)
	d := got.(map[string]any)
	if len(d["added"].(map[string]any)) != 0 {
		t.Error("expected empty added")
	}
	if len(d["removed"].(map[string]any)) != 0 {
		t.Error("expected empty removed")
	}
	if len(d["changed"].(map[string]any)) != 0 {
		t.Error("expected empty changed")
	}
}

func TestDiff_NoArgs(t *testing.T) {
	fn := extdiff.Diff()
	_, err := fn([]any{map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}

func TestDiff_WrongType(t *testing.T) {
	fn := extdiff.Diff()
	_, err := fn([]any{"not-object", map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for non-object a")
	}
	_, err = fn([]any{map[string]any{}, "not-object"}, nil)
	if err == nil {
		t.Error("expected error for non-object b")
	}
}

func TestDiff_DeepValueChange(t *testing.T) {
	fn := extdiff.Diff()
	a := map[string]any{"meta": map[string]any{"v": float64(1)}}
	b := map[string]any{"meta": map[string]any{"v": float64(2)}}
	got, _ := fn([]any{a, b}, nil)
	d := got.(map[string]any)
	changed := d["changed"].(map[string]any)
	if _, ok := changed["meta"]; !ok {
		t.Error("expected 'meta' in changed for deep value change")
	}
}

// ---------- Patch ----------

func TestPatch_Add(t *testing.T) {
	fn := extdiff.Patch()
	obj := map[string]any{"x": float64(1)}
	d := map[string]any{
		"added":   map[string]any{"y": float64(2)},
		"removed": map[string]any{},
		"changed": map[string]any{},
	}
	got, err := fn([]any{obj, d}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result := got.(map[string]any)
	if result["y"] != float64(2) {
		t.Errorf("expected y=2, got %v", result["y"])
	}
}

func TestPatch_Remove(t *testing.T) {
	fn := extdiff.Patch()
	obj := map[string]any{"x": float64(1), "y": float64(2)}
	d := map[string]any{
		"added":   map[string]any{},
		"removed": map[string]any{"y": float64(2)},
		"changed": map[string]any{},
	}
	got, _ := fn([]any{obj, d}, nil)
	result := got.(map[string]any)
	if _, ok := result["y"]; ok {
		t.Error("expected 'y' to be removed")
	}
}

func TestPatch_Change(t *testing.T) {
	fn := extdiff.Patch()
	obj := map[string]any{"x": float64(1)}
	d := map[string]any{
		"added":   map[string]any{},
		"removed": map[string]any{},
		"changed": map[string]any{"x": map[string]any{"from": float64(1), "to": float64(99)}},
	}
	got, _ := fn([]any{obj, d}, nil)
	result := got.(map[string]any)
	if result["x"] != float64(99) {
		t.Errorf("expected x=99, got %v", result["x"])
	}
}

func TestPatch_NoArgs(t *testing.T) {
	fn := extdiff.Patch()
	_, err := fn([]any{map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for missing diff arg")
	}
}

func TestPatch_WrongType(t *testing.T) {
	fn := extdiff.Patch()
	_, err := fn([]any{"not-obj", map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for non-object obj")
	}
}

func TestPatch_NoDiffSections(t *testing.T) {
	fn := extdiff.Patch()
	obj := map[string]any{"x": float64(1)}
	// diff with no added/removed/changed keys — should be identity
	got, err := fn([]any{obj, map[string]any{}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result := got.(map[string]any)
	if result["x"] != float64(1) {
		t.Errorf("expected x=1, got %v", result["x"])
	}
}

// ---------- Changed ----------

func TestChanged_True(t *testing.T) {
	fn := extdiff.Changed()
	a := map[string]any{"x": float64(1)}
	b := map[string]any{"x": float64(2)}
	got, err := fn([]any{a, b, "x"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != true {
		t.Error("expected true")
	}
}

func TestChanged_False(t *testing.T) {
	fn := extdiff.Changed()
	a := map[string]any{"x": float64(1)}
	b := map[string]any{"x": float64(1)}
	got, _ := fn([]any{a, b, "x"}, nil)
	if got != false {
		t.Error("expected false")
	}
}

func TestChanged_MissingKey(t *testing.T) {
	fn := extdiff.Changed()
	a := map[string]any{"x": float64(1)}
	b := map[string]any{"y": float64(2)}
	got, _ := fn([]any{a, b, "x"}, nil)
	// a["x"]=1 vs b["x"]=nil → changed
	if got != true {
		t.Error("expected true when key missing in b")
	}
}

func TestChanged_NoArgs(t *testing.T) {
	fn := extdiff.Changed()
	_, err := fn([]any{map[string]any{}, map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for missing key arg")
	}
}

func TestChanged_WrongKeyType(t *testing.T) {
	fn := extdiff.Changed()
	_, err := fn([]any{map[string]any{}, map[string]any{}, 42}, nil)
	if err == nil {
		t.Error("expected error for non-string key")
	}
}

func TestChanged_WrongObjType(t *testing.T) {
	fn := extdiff.Changed()
	_, err := fn([]any{"not-obj", map[string]any{}, "k"}, nil)
	if err == nil {
		t.Error("expected error for non-object a")
	}
}

// ---------- AddedKeys ----------

func TestAddedKeys_Basic(t *testing.T) {
	fn := extdiff.AddedKeys()
	a := map[string]any{"x": 1}
	b := map[string]any{"x": 1, "y": 2, "z": 3}
	got, err := fn([]any{a, b}, nil)
	if err != nil {
		t.Fatal(err)
	}
	arr := got.([]any)
	if len(arr) != 2 {
		t.Errorf("expected 2 added keys, got %d: %v", len(arr), arr)
	}
}

func TestAddedKeys_None(t *testing.T) {
	fn := extdiff.AddedKeys()
	a := map[string]any{"x": 1}
	got, _ := fn([]any{a, a}, nil)
	arr := got.([]any)
	if len(arr) != 0 {
		t.Errorf("expected empty, got %v", arr)
	}
}

func TestAddedKeys_NoArgs(t *testing.T) {
	fn := extdiff.AddedKeys()
	_, err := fn([]any{map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestAddedKeys_WrongType(t *testing.T) {
	fn := extdiff.AddedKeys()
	_, err := fn([]any{"bad", map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error for non-object")
	}
}

// ---------- RemovedKeys ----------

func TestRemovedKeys_Basic(t *testing.T) {
	fn := extdiff.RemovedKeys()
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"x": 1}
	got, err := fn([]any{a, b}, nil)
	if err != nil {
		t.Fatal(err)
	}
	arr := got.([]any)
	if len(arr) != 1 || arr[0] != "y" {
		t.Errorf("expected [y], got %v", arr)
	}
}

func TestRemovedKeys_None(t *testing.T) {
	fn := extdiff.RemovedKeys()
	a := map[string]any{"x": 1}
	got, _ := fn([]any{a, a}, nil)
	arr := got.([]any)
	if len(arr) != 0 {
		t.Errorf("expected empty, got %v", arr)
	}
}

func TestRemovedKeys_NoArgs(t *testing.T) {
	fn := extdiff.RemovedKeys()
	_, err := fn([]any{map[string]any{}}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRemovedKeys_WrongType(t *testing.T) {
	fn := extdiff.RemovedKeys()
	_, err := fn([]any{map[string]any{}, 42}, nil)
	if err == nil {
		t.Error("expected error for non-object b")
	}
}

// ---------- ArrayDiff ----------

func TestArrayDiff_Added(t *testing.T) {
	fn := extdiff.ArrayDiff()
	a := []any{"apple", "banana"}
	b := []any{"apple", "banana", "cherry"}
	got, err := fn([]any{a, b}, nil)
	if err != nil {
		t.Fatal(err)
	}
	d := got.(map[string]any)
	added := d["added"].([]any)
	if len(added) != 1 || added[0] != "cherry" {
		t.Errorf("expected [cherry], got %v", added)
	}
	removed := d["removed"].([]any)
	if len(removed) != 0 {
		t.Errorf("expected empty removed, got %v", removed)
	}
}

func TestArrayDiff_Removed(t *testing.T) {
	fn := extdiff.ArrayDiff()
	a := []any{"a", "b", "c"}
	b := []any{"a", "b"}
	got, _ := fn([]any{a, b}, nil)
	d := got.(map[string]any)
	removed := d["removed"].([]any)
	if len(removed) != 1 || removed[0] != "c" {
		t.Errorf("expected [c], got %v", removed)
	}
}

func TestArrayDiff_NoDiff(t *testing.T) {
	fn := extdiff.ArrayDiff()
	a := []any{"a", "b"}
	got, _ := fn([]any{a, a}, nil)
	d := got.(map[string]any)
	if len(d["added"].([]any)) != 0 || len(d["removed"].([]any)) != 0 {
		t.Error("expected no diff")
	}
}

func TestArrayDiff_NoArgs(t *testing.T) {
	fn := extdiff.ArrayDiff()
	_, err := fn([]any{[]any{}}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestArrayDiff_WrongType(t *testing.T) {
	fn := extdiff.ArrayDiff()
	_, err := fn([]any{"not-array", []any{}}, nil)
	if err == nil {
		t.Error("expected error for non-array")
	}
}

// ---------- DeepEqual ----------

func TestDeepEqual_Equal(t *testing.T) {
	fn := extdiff.DeepEqual()
	got, err := fn([]any{map[string]any{"a": float64(1)}, map[string]any{"a": float64(1)}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != true {
		t.Error("expected true")
	}
}

func TestDeepEqual_NotEqual(t *testing.T) {
	fn := extdiff.DeepEqual()
	got, _ := fn([]any{float64(1), float64(2)}, nil)
	if got != false {
		t.Error("expected false")
	}
}

func TestDeepEqual_Nested(t *testing.T) {
	fn := extdiff.DeepEqual()
	a := map[string]any{"x": map[string]any{"y": float64(1)}}
	b := map[string]any{"x": map[string]any{"y": float64(1)}}
	got, _ := fn([]any{a, b}, nil)
	if got != true {
		t.Error("expected true for nested equal objects")
	}
}

func TestDeepEqual_NilNil(t *testing.T) {
	fn := extdiff.DeepEqual()
	got, _ := fn([]any{nil, nil}, nil)
	if got != true {
		t.Error("expected true for nil==nil")
	}
}

func TestDeepEqual_NilVsValue(t *testing.T) {
	fn := extdiff.DeepEqual()
	got, _ := fn([]any{nil, "x"}, nil)
	if got != false {
		t.Error("expected false for nil vs non-nil")
	}
}

func TestDeepEqual_Arrays(t *testing.T) {
	fn := extdiff.DeepEqual()
	a := []any{float64(1), float64(2)}
	b := []any{float64(1), float64(2)}
	got, _ := fn([]any{a, b}, nil)
	if got != true {
		t.Error("expected true for equal arrays")
	}
}

func TestDeepEqual_NoArgs(t *testing.T) {
	fn := extdiff.DeepEqual()
	_, err := fn([]any{"only-one"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := extdiff.All()
	expected := []string{"diff", "patch", "changed", "addedKeys", "removedKeys", "arrayDiff", "deepEqual"}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
