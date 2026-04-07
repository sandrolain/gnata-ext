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

// --- Additional coverage tests ---

func TestJsonParseErrors(t *testing.T) {
	fn := extjson.JsonParse()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = fn([]any{42}, nil)
	if err == nil {
		t.Error("expected error for non-string arg")
	}
}

func TestJsonStringifyErrors(t *testing.T) {
	fn := extjson.JsonStringify()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	// non-string indent
	_, err = fn([]any{map[string]any{"x": 1}, 99}, nil)
	if err == nil {
		t.Error("expected error for non-string indent")
	}
}

func TestJsonStringifyIndent(t *testing.T) {
	fn := extjson.JsonStringify()
	got, err := fn([]any{map[string]any{"x": float64(1)}, "  "}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := got.(string)
	if len(s) < 5 {
		t.Errorf("expected indented JSON, got %q", s)
	}
}

func TestJsonDiffErrors(t *testing.T) {
	fn := extjson.JsonDiff()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
}

func TestJsonDiffSliceAndPrimitive(t *testing.T) {
	fn := extjson.JsonDiff()

	// two slices — treated as non-map → replace
	a := []any{float64(1), float64(2)}
	b := []any{float64(1), float64(3)}
	got, err := fn([]any{a, b}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ops := got.([]any)
	if len(ops) == 0 {
		t.Error("expected diff op for changed slice")
	}

	// identical primitives → no ops
	got, err = fn([]any{float64(1), float64(1)}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.([]any)) != 0 {
		t.Error("expected empty diff for identical values")
	}
}

func TestJsonPatchErrors(t *testing.T) {
	fn := extjson.JsonPatch()

	// too few args
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	// ops not an array
	_, err = fn([]any{map[string]any{}, "not-an-array"}, nil)
	if err == nil {
		t.Error("expected error for non-array ops")
	}
	// op not an object
	_, err = fn([]any{map[string]any{}, []any{"not-an-op"}}, nil)
	if err == nil {
		t.Error("expected error for non-object op")
	}
	// unknown op
	_, err = fn([]any{map[string]any{}, []any{
		map[string]any{"op": "unknown", "path": "/a"},
	}}, nil)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestJsonPatchAddArray(t *testing.T) {
	fn := extjson.JsonPatch()

	// add to array via index
	doc := []any{float64(1), float64(2), float64(3)}
	ops := []any{map[string]any{"op": "add", "path": "/1", "value": float64(99)}}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr := got.([]any)
	if arr[1] != float64(99) {
		t.Errorf("expected 99 at index 1, got %v", arr[1])
	}

	// add to array via "-" (append)
	ops2 := []any{map[string]any{"op": "add", "path": "/-", "value": float64(4)}}
	got, err = fn([]any{doc, ops2}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr2 := got.([]any)
	if arr2[len(arr2)-1] != float64(4) {
		t.Errorf("expected 4 appended, got %v", arr2)
	}
}

func TestJsonPatchAddInvalidArrayIndex(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := []any{float64(1)}
	ops := []any{map[string]any{"op": "add", "path": "/99", "value": float64(0)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for out-of-range array index")
	}
}

func TestJsonPatchAddNonTraversable(t *testing.T) {
	fn := extjson.JsonPatch()
	// try to traverse a non-object at intermediate path
	doc := map[string]any{"a": float64(1)}
	ops := []any{map[string]any{"op": "add", "path": "/a/b", "value": float64(2)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for non-traversable path")
	}
}

func TestJsonPatchRemoveArray(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := []any{float64(1), float64(2), float64(3)}
	ops := []any{map[string]any{"op": "remove", "path": "/1"}}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 {
		t.Errorf("expected length 2, got %d", len(arr))
	}
}

func TestJsonPatchRemoveErrors(t *testing.T) {
	fn := extjson.JsonPatch()

	// remove from non-existent map key
	doc := map[string]any{"a": float64(1)}
	ops := []any{map[string]any{"op": "remove", "path": "/z"}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for non-existent key")
	}

	// remove from array out of range
	doc2 := []any{float64(1)}
	ops2 := []any{map[string]any{"op": "remove", "path": "/99"}}
	_, err = fn([]any{doc2, ops2}, nil)
	if err == nil {
		t.Error("expected error for out-of-range index")
	}

	// remove root
	ops3 := []any{map[string]any{"op": "remove", "path": ""}}
	_, err = fn([]any{doc, ops3}, nil)
	if err == nil {
		t.Error("expected error for removing root")
	}
}

func TestJsonPatchReplaceArray(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := []any{float64(1), float64(2), float64(3)}
	ops := []any{map[string]any{"op": "replace", "path": "/0", "value": float64(99)}}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr := got.([]any)
	if arr[0] != float64(99) {
		t.Errorf("expected 99, got %v", arr[0])
	}
}

func TestJsonPatchReplaceErrors(t *testing.T) {
	fn := extjson.JsonPatch()

	// replace non-existent key
	doc := map[string]any{"a": float64(1)}
	ops := []any{map[string]any{"op": "replace", "path": "/z", "value": float64(0)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for non-existent key replace")
	}

	// replace array out of range
	doc2 := []any{float64(1)}
	ops2 := []any{map[string]any{"op": "replace", "path": "/99", "value": float64(0)}}
	_, err = fn([]any{doc2, ops2}, nil)
	if err == nil {
		t.Error("expected error for out-of-range replace")
	}
}

func TestJsonPatchMoveErrors(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": float64(1)}

	// from path not found
	ops := []any{map[string]any{"op": "move", "from": "/notexists", "path": "/b"}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for move from non-existent path")
	}
}

func TestJsonPatchCopyErrors(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": float64(1)}

	ops := []any{map[string]any{"op": "copy", "from": "/notexists", "path": "/b"}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for copy from non-existent path")
	}
}

func TestJsonPatchTestNotFound(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": float64(1)}
	ops := []any{map[string]any{"op": "test", "path": "/notexists", "value": nil}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for test path not found")
	}
}

func TestJsonPointerErrors(t *testing.T) {
	fn := extjson.JsonPointerFunc()
	_, err := fn([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = fn([]any{map[string]any{}, 99}, nil)
	if err == nil {
		t.Error("expected error for non-string pointer")
	}
	// pointer to non-existent path
	_, err = fn([]any{map[string]any{"a": float64(1)}, "/missing"}, nil)
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestJsonPointerArrayIndex(t *testing.T) {
	fn := extjson.JsonPointerFunc()
	doc := map[string]any{"arr": []any{float64(10), float64(20)}}
	got, err := fn([]any{doc, "/arr/0"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != float64(10) {
		t.Errorf("got %v, want 10", got)
	}

	// bad index
	_, err = fn([]any{doc, "/arr/bad"}, nil)
	if err == nil {
		t.Error("expected error for bad array index")
	}

	// out of range
	_, err = fn([]any{doc, "/arr/99"}, nil)
	if err == nil {
		t.Error("expected error for out-of-range index")
	}

	// traverse non-map/slice
	_, err = fn([]any{doc, "/arr/0/x"}, nil)
	if err == nil {
		t.Error("expected error for traversing primitive")
	}
}

func TestJsonPointerRoot(t *testing.T) {
	fn := extjson.JsonPointerFunc()
	doc := map[string]any{"a": float64(1)}
	// empty pointer or "/" → entire doc
	got, err := fn([]any{doc, "/"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = got
}

// --- patchRemove nested path coverage ---

func TestJsonPatchRemoveNested(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": map[string]any{"b": float64(1), "c": float64(2)}}
	ops := []any{map[string]any{"op": "remove", "path": "/a/b"}}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := got.(map[string]any)
	inner := m["a"].(map[string]any)
	if _, exists := inner["b"]; exists {
		t.Error("expected a.b to be removed")
	}
	if inner["c"] != float64(2) {
		t.Errorf("a.c: got %v, want 2", inner["c"])
	}
}

func TestJsonPatchRemoveNestedNonObject(t *testing.T) {
	fn := extjson.JsonPatch()
	// intermediate value is not a map
	doc := map[string]any{"a": float64(1)}
	ops := []any{map[string]any{"op": "remove", "path": "/a/b"}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for non-object traversal in remove")
	}
}

func TestJsonPatchRemoveNestedMissingKey(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": map[string]any{"x": float64(1)}}
	ops := []any{map[string]any{"op": "remove", "path": "/z/b"}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for missing key in nested remove")
	}
}

// --- patchReplace nested path coverage ---

func TestJsonPatchReplaceNested(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": map[string]any{"b": float64(1)}}
	ops := []any{map[string]any{"op": "replace", "path": "/a/b", "value": float64(99)}}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := got.(map[string]any)
	inner := m["a"].(map[string]any)
	if inner["b"] != float64(99) {
		t.Errorf("a.b: got %v, want 99", inner["b"])
	}
}

func TestJsonPatchReplaceNestedNonObject(t *testing.T) {
	fn := extjson.JsonPatch()
	// intermediate value is not a map
	doc := map[string]any{"a": float64(1)}
	ops := []any{map[string]any{"op": "replace", "path": "/a/b", "value": float64(2)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for non-object traversal in replace")
	}
}

func TestJsonPatchReplaceNestedMissingKey(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": map[string]any{"x": float64(1)}}
	ops := []any{map[string]any{"op": "replace", "path": "/z/b", "value": float64(0)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for missing key in nested replace")
	}
}

// --- patchAdd nested path: auto-create missing parent ---

func TestJsonPatchAddNestedAutoCreate(t *testing.T) {
	fn := extjson.JsonPatch()
	doc := map[string]any{"a": map[string]any{}}
	ops := []any{map[string]any{"op": "add", "path": "/a/b/c", "value": float64(42)}}
	got, err := fn([]any{doc, ops}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := got.(map[string]any)
	inner := m["a"].(map[string]any)
	nested := inner["b"].(map[string]any)
	if nested["c"] != float64(42) {
		t.Errorf("a.b.c: got %v, want 42", nested["c"])
	}
}

func TestJsonPatchAddLeafNonMapNonArray(t *testing.T) {
	fn := extjson.JsonPatch()
	// top-level doc is a primitive (float): patchAddLeaf default case
	doc := float64(1)
	ops := []any{map[string]any{"op": "add", "path": "/x", "value": float64(2)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for add to primitive")
	}
}

func TestJsonPatchRemoveLeafNonMapNonArray(t *testing.T) {
	fn := extjson.JsonPatch()
	// top-level doc is a primitive: patchRemoveLeaf default case
	doc := float64(1)
	ops := []any{map[string]any{"op": "remove", "path": "/x"}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for remove from primitive")
	}
}

func TestJsonPatchReplaceLeafNonMapNonArray(t *testing.T) {
	fn := extjson.JsonPatch()
	// top-level doc is a primitive: patchReplaceLeaf default case
	doc := float64(1)
	ops := []any{map[string]any{"op": "replace", "path": "/x", "value": float64(2)}}
	_, err := fn([]any{doc, ops}, nil)
	if err == nil {
		t.Error("expected error for replace in primitive")
	}
}

func TestJsonStringifyMarshalError(t *testing.T) {
	fn := extjson.JsonStringify()
	// channels cannot be marshalled to JSON
	_, err := fn([]any{make(chan int)}, nil)
	if err == nil {
		t.Error("expected error for non-marshallable value")
	}
}
