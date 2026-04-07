package extjson_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extjson"
)

func TestJsonParse(t *testing.T) {
	fn := extjson.JsonParse()

	got, err := fn([]any{`{"a":1,"b":true}`}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if m["a"] != float64(1) {
		t.Errorf("a: got %v, want 1", m["a"])
	}
	if m["b"] != true {
		t.Errorf("b: got %v, want true", m["b"])
	}

	_, err = fn([]any{"not json"}, nil)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestJsonStringify(t *testing.T) {
	fn := extjson.JsonStringify()

	got, err := fn([]any{map[string]any{"x": float64(1)}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("expected string, got %T", got)
	}
	if s != `{"x":1}` {
		t.Errorf("got %q, want %q", s, `{"x":1}`)
	}

	// with indent
	got, err = fn([]any{map[string]any{"x": float64(1)}, "  "}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.(string) == `{"x":1}` {
		t.Error("expected indented output")
	}
}

func TestJsonDiff(t *testing.T) {
	fn := extjson.JsonDiff()

	a := map[string]any{"x": float64(1), "y": float64(2)}
	b := map[string]any{"x": float64(1), "y": float64(3), "z": float64(4)}

	got, err := fn([]any{a, b}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ops, ok := got.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", got)
	}
	if len(ops) == 0 {
		t.Error("expected at least one diff op")
	}

	// no diff
	got, _ = fn([]any{a, a}, nil)
	if len(got.([]any)) != 0 {
		t.Errorf("expected empty diff, got %v", got)
	}
}

func TestJsonPatch(t *testing.T) {
	fn := extjson.JsonPatch()

	doc := map[string]any{"a": float64(1), "b": float64(2)}
	ops := []any{
		map[string]any{"op": "add", "path": "/c", "value": float64(3)},
		map[string]any{"op": "remove", "path": "/b"},
		map[string]any{"op": "replace", "path": "/a", "value": float64(10)},
	}

	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if m["a"] != float64(10) {
		t.Errorf("a: got %v, want 10", m["a"])
	}
	if _, exists := m["b"]; exists {
		t.Error("b should have been removed")
	}
	if m["c"] != float64(3) {
		t.Errorf("c: got %v, want 3", m["c"])
	}
}

func TestJsonPatchMove(t *testing.T) {
	fn := extjson.JsonPatch()

	doc := map[string]any{"a": float64(1), "b": float64(2)}
	ops := []any{
		map[string]any{"op": "move", "from": "/a", "path": "/d"},
	}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := got.(map[string]any)
	if _, exists := m["a"]; exists {
		t.Error("a should be removed after move")
	}
	if m["d"] != float64(1) {
		t.Errorf("d: got %v, want 1", m["d"])
	}
}

func TestJsonPatchCopy(t *testing.T) {
	fn := extjson.JsonPatch()

	doc := map[string]any{"a": float64(5)}
	ops := []any{
		map[string]any{"op": "copy", "from": "/a", "path": "/b"},
	}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := got.(map[string]any)
	if m["a"] != float64(5) || m["b"] != float64(5) {
		t.Errorf("expected a=5, b=5, got %v", m)
	}
}

func TestJsonPatchTest(t *testing.T) {
	fn := extjson.JsonPatch()

	doc := map[string]any{"a": float64(1)}
	// passing test
	ops := []any{
		map[string]any{"op": "test", "path": "/a", "value": float64(1)},
	}
	_, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("test op should pass: %v", err)
	}
	// failing test
	ops2 := []any{
		map[string]any{"op": "test", "path": "/a", "value": float64(99)},
	}
	_, err = fn([]any{doc, ops2}, nil)
	if err == nil {
		t.Error("test op should fail with mismatched value")
	}
}

func TestJsonPointer(t *testing.T) {
	fn := extjson.JsonPointerFunc()

	doc := map[string]any{
		"a": map[string]any{
			"b": []any{float64(10), float64(20), float64(30)},
		},
	}

	got, err := fn([]any{doc, "/a/b/1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != float64(20) {
		t.Errorf("got %v, want 20", got)
	}

	_, err = fn([]any{doc, "/a/nonexistent"}, nil)
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestAll(t *testing.T) {
	all := extjson.All()
	expected := []string{"jsonParse", "jsonStringify", "jsonDiff", "jsonPatch", "jsonPointer"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All() missing function: %q", name)
		}
	}
}
