package extobject_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extobject"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

var sampleObj = map[string]any{"a": "1", "b": "2", "c": "3"}

func TestValues(t *testing.T) {
	f := extobject.Values()
	got, err := invoke(f, map[string]any{"x": 1.0})
	if err != nil {
		t.Fatalf("values: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 1 || arr[0].(float64) != 1.0 {
		t.Errorf("values: got %v", arr)
	}
}

func TestPairs(t *testing.T) {
	f := extobject.Pairs()
	got, err := invoke(f, map[string]any{"k": "v"})
	if err != nil {
		t.Fatalf("pairs: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 1 {
		t.Errorf("pairs: got %v", arr)
	}
	pair := arr[0].([]any)
	if pair[0].(string) != "k" || pair[1].(string) != "v" {
		t.Errorf("pairs: pair = %v", pair)
	}
}

func TestFromPairs(t *testing.T) {
	f := extobject.FromPairs()
	got, err := invoke(f, []any{[]any{"x", 42.0}, []any{"y", "hello"}})
	if err != nil {
		t.Fatalf("fromPairs: %v", err)
	}
	obj := got.(map[string]any)
	if obj["x"].(float64) != 42.0 || obj["y"].(string) != "hello" {
		t.Errorf("fromPairs: got %v", obj)
	}
}

func TestPick(t *testing.T) {
	f := extobject.Pick()
	got, err := invoke(f, sampleObj, []any{"a", "c"})
	if err != nil {
		t.Fatalf("pick: %v", err)
	}
	obj := got.(map[string]any)
	if len(obj) != 2 || obj["a"] != "1" || obj["c"] != "3" {
		t.Errorf("pick: got %v", obj)
	}
}

func TestOmit(t *testing.T) {
	f := extobject.Omit()
	got, err := invoke(f, sampleObj, []any{"b"})
	if err != nil {
		t.Fatalf("omit: %v", err)
	}
	obj := got.(map[string]any)
	if len(obj) != 2 || obj["b"] != nil {
		t.Errorf("omit: got %v", obj)
	}
	if _, ok := obj["b"]; ok {
		t.Errorf("omit: key 'b' should be absent")
	}
}

func TestDeepMerge(t *testing.T) {
	f := extobject.DeepMerge()
	target := map[string]any{
		"a": map[string]any{"x": 1.0, "y": 2.0},
		"b": "keep",
	}
	source := map[string]any{
		"a": map[string]any{"y": 99.0, "z": 3.0},
		"c": "new",
	}
	got, err := invoke(f, target, source)
	if err != nil {
		t.Fatalf("deepMerge: %v", err)
	}
	obj := got.(map[string]any)
	inner := obj["a"].(map[string]any)
	if inner["x"].(float64) != 1.0 || inner["y"].(float64) != 99.0 || inner["z"].(float64) != 3.0 {
		t.Errorf("deepMerge inner: got %v", inner)
	}
	if obj["b"].(string) != "keep" || obj["c"].(string) != "new" {
		t.Errorf("deepMerge outer: got %v", obj)
	}
}

func TestInvert(t *testing.T) {
	f := extobject.Invert()
	got, err := invoke(f, map[string]any{"k1": "v1", "k2": "v2"})
	if err != nil {
		t.Fatalf("invert: %v", err)
	}
	obj := got.(map[string]any)
	if obj["v1"].(string) != "k1" || obj["v2"].(string) != "k2" {
		t.Errorf("invert: got %v", obj)
	}
}

func TestSize(t *testing.T) {
	f := extobject.Size()
	got, err := invoke(f, sampleObj)
	if err != nil {
		t.Fatalf("size: %v", err)
	}
	if got.(float64) != 3.0 {
		t.Errorf("size: got %v", got)
	}
}

func TestRename(t *testing.T) {
	f := extobject.Rename()
	got, err := invoke(f, map[string]any{"old": "val", "other": "x"}, "old", "new")
	if err != nil {
		t.Fatalf("rename: %v", err)
	}
	obj := got.(map[string]any)
	if obj["new"].(string) != "val" {
		t.Errorf("rename: missing new key, got %v", obj)
	}
	if _, ok := obj["old"]; ok {
		t.Errorf("rename: old key should be gone")
	}
}

func TestAll(t *testing.T) {
	all := extobject.All()
	expected := []string{"values", "pairs", "fromPairs", "pick", "omit", "deepMerge", "invert", "size", "rename"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing function %q", name)
		}
	}
}
