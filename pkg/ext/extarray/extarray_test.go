package extarray_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extarray"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

func TestFirst(t *testing.T) {
	f := extarray.First()
	got, err := invoke(f, []any{1.0, 2.0, 3.0})
	if err != nil || got.(float64) != 1.0 {
		t.Errorf("first: got %v, %v", got, err)
	}
	got, err = invoke(f, []any{})
	if err != nil || got != nil {
		t.Errorf("first([]): got %v, %v", got, err)
	}
}

func TestLast(t *testing.T) {
	f := extarray.Last()
	got, err := invoke(f, []any{1.0, 2.0, 3.0})
	if err != nil || got.(float64) != 3.0 {
		t.Errorf("last: got %v, %v", got, err)
	}
}

func TestTake(t *testing.T) {
	f := extarray.Take()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0}, 2.0)
	if err != nil {
		t.Fatalf("take: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 || arr[0].(float64) != 1.0 || arr[1].(float64) != 2.0 {
		t.Errorf("take(2): got %v", arr)
	}
}

func TestSkip(t *testing.T) {
	f := extarray.Skip()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0}, 2.0)
	if err != nil {
		t.Fatalf("skip: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 || arr[0].(float64) != 3.0 {
		t.Errorf("skip(2): got %v", arr)
	}
}

func TestSlice(t *testing.T) {
	f := extarray.Slice()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0}, 1.0, 3.0)
	if err != nil {
		t.Fatalf("slice: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 || arr[0].(float64) != 2.0 {
		t.Errorf("slice(1,3): got %v", arr)
	}
}

func TestFlatten(t *testing.T) {
	f := extarray.Flatten()
	got, err := invoke(f, []any{1.0, []any{2.0, []any{3.0}}})
	if err != nil {
		t.Fatalf("flatten: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 3 {
		t.Errorf("flatten: got %v", arr)
	}
	// depth=1
	got, err = invoke(f, []any{1.0, []any{2.0, []any{3.0}}}, 1.0)
	if err != nil {
		t.Fatalf("flatten depth=1: %v", err)
	}
	arr = got.([]any)
	if len(arr) != 3 {
		t.Errorf("flatten depth=1: got %v", arr)
	}
	nested, ok := arr[2].([]any)
	if !ok || nested[0].(float64) != 3.0 {
		t.Errorf("flatten depth=1: inner should remain array, got %v", arr[2])
	}
}

func TestChunk(t *testing.T) {
	f := extarray.Chunk()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0, 5.0}, 2.0)
	if err != nil {
		t.Fatalf("chunk: %v", err)
	}
	chunks := got.([]any)
	if len(chunks) != 3 {
		t.Errorf("chunk(2): got %v chunks", len(chunks))
	}
	if len(chunks[2].([]any)) != 1 {
		t.Errorf("chunk(2): last chunk should have 1 item")
	}
}

func TestUnion(t *testing.T) {
	f := extarray.Union()
	got, err := invoke(f, []any{1.0, 2.0}, []any{2.0, 3.0})
	if err != nil {
		t.Fatalf("union: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 3 {
		t.Errorf("union: got %v", arr)
	}
}

func TestIntersection(t *testing.T) {
	f := extarray.Intersection()
	got, err := invoke(f, []any{1.0, 2.0, 3.0}, []any{2.0, 3.0, 4.0})
	if err != nil {
		t.Fatalf("intersection: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 {
		t.Errorf("intersection: got %v", arr)
	}
}

func TestDifference(t *testing.T) {
	f := extarray.Difference()
	got, err := invoke(f, []any{1.0, 2.0, 3.0}, []any{2.0})
	if err != nil {
		t.Fatalf("difference: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 {
		t.Errorf("difference: got %v", arr)
	}
}

func TestSymmetricDifference(t *testing.T) {
	f := extarray.SymmetricDifference()
	got, err := invoke(f, []any{1.0, 2.0, 3.0}, []any{2.0, 4.0})
	if err != nil {
		t.Fatalf("symmetricDifference: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 3 { // 1, 3, 4
		t.Errorf("symmetricDifference: got %v", arr)
	}
}

func TestRange(t *testing.T) {
	f := extarray.Range()
	got, err := invoke(f, 1.0, 5.0)
	if err != nil {
		t.Fatalf("range: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 4 {
		t.Errorf("range(1,5): got %v", arr)
	}
	// with step
	got, err = invoke(f, 0.0, 10.0, 2.0)
	if err != nil {
		t.Fatalf("range step: %v", err)
	}
	arr = got.([]any)
	if len(arr) != 5 {
		t.Errorf("range(0,10,2): got %v", arr)
	}
}

func TestZipLongest(t *testing.T) {
	f := extarray.ZipLongest()
	got, err := invoke(f, []any{1.0, 2.0}, []any{3.0, 4.0, 5.0}, nil)
	if err != nil {
		t.Fatalf("zipLongest: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 3 {
		t.Errorf("zipLongest: got %v", arr)
	}
}

func TestWindow(t *testing.T) {
	f := extarray.Window()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0}, 2.0)
	if err != nil {
		t.Fatalf("window: %v", err)
	}
	windows := got.([]any)
	if len(windows) != 3 {
		t.Errorf("window(size=2): got %v windows", len(windows))
	}
}

func TestAll(t *testing.T) {
	all := extarray.All()
	expected := []string{
		"first", "last", "take", "skip", "slice", "flatten", "chunk",
		"union", "intersection", "difference", "symmetricDifference",
		"range", "zipLongest", "window",
	}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing function %q", name)
		}
	}
}

func TestFirstErrors(t *testing.T) {
	f := extarray.First()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("first: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array"); err == nil {
		t.Error("first: expected error for non-array")
	}
}

func TestLastEmpty(t *testing.T) {
	f := extarray.Last()
	got, err := invoke(f, []any{})
	if err != nil || got != nil {
		t.Errorf("last([]): expected nil, got %v, %v", got, err)
	}
}

func TestLastErrors(t *testing.T) {
	f := extarray.Last()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("last: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array"); err == nil {
		t.Error("last: expected error for non-array")
	}
}

func TestTakeEdgeCases(t *testing.T) {
	f := extarray.Take()
	// n < 0 -> empty
	got, _ := invoke(f, []any{1.0, 2.0}, -1.0)
	if len(got.([]any)) != 0 {
		t.Errorf("take(-1): expected empty, got %v", got)
	}
	// n > len -> full array
	got, _ = invoke(f, []any{1.0, 2.0}, 10.0)
	if len(got.([]any)) != 2 {
		t.Errorf("take(10): expected 2, got %v", got)
	}
}

func TestTakeErrors(t *testing.T) {
	f := extarray.Take()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("take: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", 2.0); err == nil {
		t.Error("take: expected error for non-array")
	}
	if _, err := invoke(f, []any{1.0}, "bad"); err == nil {
		t.Error("take: expected error for non-numeric n")
	}
}

func TestSkipEdgeCases(t *testing.T) {
	f := extarray.Skip()
	// n < 0 -> full array
	got, _ := invoke(f, []any{1.0, 2.0}, -1.0)
	if len(got.([]any)) != 2 {
		t.Errorf("skip(-1): expected 2, got %v", got)
	}
	// n > len -> empty
	got, _ = invoke(f, []any{1.0, 2.0}, 10.0)
	if len(got.([]any)) != 0 {
		t.Errorf("skip(10): expected empty, got %v", got)
	}
}

func TestSkipErrors(t *testing.T) {
	f := extarray.Skip()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("skip: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", 2.0); err == nil {
		t.Error("skip: expected error for non-array")
	}
	if _, err := invoke(f, []any{1.0}, "bad"); err == nil {
		t.Error("skip: expected error for non-numeric n")
	}
}

func TestSliceNegativeIndex(t *testing.T) {
	f := extarray.Slice()
	// negative start: -1 normalises to len-1=3
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0}, -1.0)
	if err != nil {
		t.Fatalf("slice(-1): %v", err)
	}
	arr := got.([]any)
	if len(arr) != 1 || arr[0].(float64) != 4.0 {
		t.Errorf("slice(-1): got %v", arr)
	}
	// start > end -> empty
	got, _ = invoke(f, []any{1.0, 2.0, 3.0}, 3.0, 1.0)
	if len(got.([]any)) != 0 {
		t.Errorf("slice(3,1): expected empty, got %v", got)
	}
}

func TestSliceErrors(t *testing.T) {
	f := extarray.Slice()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("slice: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", 0.0); err == nil {
		t.Error("slice: expected error for non-array")
	}
}

func TestFlattenErrors(t *testing.T) {
	f := extarray.Flatten()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("flatten: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array"); err == nil {
		t.Error("flatten: expected error for non-array")
	}
	if _, err := invoke(f, []any{1.0}, "bad"); err == nil {
		t.Error("flatten: expected error for non-numeric depth")
	}
}

func TestChunkEmpty(t *testing.T) {
	f := extarray.Chunk()
	got, err := invoke(f, []any{}, 2.0)
	if err != nil || len(got.([]any)) != 0 {
		t.Errorf("chunk([]): expected empty chunks, got %v, %v", got, err)
	}
}

func TestChunkErrors(t *testing.T) {
	f := extarray.Chunk()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("chunk: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", 2.0); err == nil {
		t.Error("chunk: expected error for non-array")
	}
	if _, err := invoke(f, []any{1.0}, 0.0); err == nil {
		t.Error("chunk: expected error for size=0")
	}
	if _, err := invoke(f, []any{1.0}, "bad"); err == nil {
		t.Error("chunk: expected error for non-numeric size")
	}
}

func TestUnionEmpty(t *testing.T) {
	f := extarray.Union()
	got, _ := invoke(f, []any{}, []any{})
	if len(got.([]any)) != 0 {
		t.Errorf("union(empty): expected [], got %v", got)
	}
}

func TestUnionErrors(t *testing.T) {
	f := extarray.Union()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("union: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", []any{}); err == nil {
		t.Error("union: expected error for non-array a")
	}
	if _, err := invoke(f, []any{}, "not-array"); err == nil {
		t.Error("union: expected error for non-array b")
	}
}

func TestIntersectionEmpty(t *testing.T) {
	f := extarray.Intersection()
	got, _ := invoke(f, []any{1.0}, []any{2.0})
	if len(got.([]any)) != 0 {
		t.Errorf("intersection(no common): expected [], got %v", got)
	}
}

func TestIntersectionErrors(t *testing.T) {
	f := extarray.Intersection()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("intersection: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", []any{}); err == nil {
		t.Error("intersection: expected error for non-array a")
	}
	if _, err := invoke(f, []any{}, "not-array"); err == nil {
		t.Error("intersection: expected error for non-array b")
	}
}

func TestDifferenceEmpty(t *testing.T) {
	f := extarray.Difference()
	got, _ := invoke(f, []any{1.0}, []any{1.0})
	if len(got.([]any)) != 0 {
		t.Errorf("difference(a subset b): expected [], got %v", got)
	}
}

func TestDifferenceErrors(t *testing.T) {
	f := extarray.Difference()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("difference: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", []any{}); err == nil {
		t.Error("difference: expected error for non-array a")
	}
	if _, err := invoke(f, []any{}, "not-array"); err == nil {
		t.Error("difference: expected error for non-array b")
	}
}

func TestSymmetricDifferenceEmpty(t *testing.T) {
	f := extarray.SymmetricDifference()
	got, _ := invoke(f, []any{1.0, 2.0}, []any{1.0, 2.0})
	if len(got.([]any)) != 0 {
		t.Errorf("symmetricDiff(identical): expected [], got %v", got)
	}
}

func TestSymmetricDifferenceErrors(t *testing.T) {
	f := extarray.SymmetricDifference()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("symmetricDifference: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", []any{}); err == nil {
		t.Error("symmetricDifference: expected error for non-array a")
	}
	if _, err := invoke(f, []any{}, "not-array"); err == nil {
		t.Error("symmetricDifference: expected error for non-array b")
	}
}

func TestRangeEdgeCases(t *testing.T) {
	f := extarray.Range()
	// empty range (start == end)
	got, _ := invoke(f, 5.0, 5.0)
	if len(got.([]any)) != 0 {
		t.Errorf("range(5,5): expected [], got %v", got)
	}
	// negative step
	got, err := invoke(f, 5.0, 1.0, -1.0)
	if err != nil {
		t.Fatalf("range negative step: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 4 { // 5,4,3,2
		t.Errorf("range(5,1,-1): expected 4 items, got %v", arr)
	}
}

func TestRangeErrors(t *testing.T) {
	f := extarray.Range()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("range: expected error for 0 args")
	}
	if _, err := invoke(f, "bad", 5.0); err == nil {
		t.Error("range: expected error for non-numeric start")
	}
	if _, err := invoke(f, 0.0, "bad"); err == nil {
		t.Error("range: expected error for non-numeric end")
	}
	if _, err := invoke(f, 0.0, 5.0, 0.0); err == nil {
		t.Error("range: expected error for step=0")
	}
	if _, err := invoke(f, 0.0, 5.0, "bad"); err == nil {
		t.Error("range: expected error for non-numeric step")
	}
}

func TestZipLongestErrors(t *testing.T) {
	f := extarray.ZipLongest()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("zipLongest: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", []any{}); err == nil {
		t.Error("zipLongest: expected error for non-array a")
	}
	if _, err := invoke(f, []any{}, "not-array"); err == nil {
		t.Error("zipLongest: expected error for non-array b")
	}
}

func TestWindowWithStep(t *testing.T) {
	f := extarray.Window()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0, 5.0}, 2.0, 2.0)
	if err != nil {
		t.Fatalf("window step=2: %v", err)
	}
	windows := got.([]any)
	if len(windows) != 2 { // [1,2], [3,4]
		t.Errorf("window(size=2,step=2): expected 2, got %v", len(windows))
	}
}

func TestWindowSizeExceedsArray(t *testing.T) {
	f := extarray.Window()
	got, err := invoke(f, []any{1.0, 2.0}, 5.0)
	if err != nil || len(got.([]any)) != 0 {
		t.Errorf("window(size>len): expected [], got %v, %v", got, err)
	}
}

func TestWindowErrors(t *testing.T) {
	f := extarray.Window()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("window: expected error for 0 args")
	}
	if _, err := invoke(f, "not-array", 2.0); err == nil {
		t.Error("window: expected error for non-array")
	}
	if _, err := invoke(f, []any{1.0}, 0.0); err == nil {
		t.Error("window: expected error for size=0")
	}
	if _, err := invoke(f, []any{1.0}, 1.0, 0.0); err == nil {
		t.Error("window: expected error for step=0")
	}
	if _, err := invoke(f, []any{1.0}, 1.0, "bad"); err == nil {
		t.Error("window: expected error for non-numeric step")
	}
}

func TestCompact(t *testing.T) {
	fn := extarray.Compact()
	got, err := fn([]any{[]any{float64(1), nil, false, float64(0), "", float64(2)}}, nil)
	if err != nil {
		t.Errorf("compact: unexpected error: %v", err)
	}
	arr := got.([]any)
	if len(arr) != 2 || arr[0] != float64(1) || arr[1] != float64(2) {
		t.Errorf("compact: got %v, want [1 2]", arr)
	}
}

func TestCompactErrors(t *testing.T) {
	fn := extarray.Compact()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("compact: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array"}, nil); err == nil {
		t.Error("compact: expected error for non-array")
	}
}

func TestGroupByKey(t *testing.T) {
	fn := extarray.GroupByKey()
	arr := []any{
		map[string]any{"type": "a", "val": float64(1)},
		map[string]any{"type": "b", "val": float64(2)},
		map[string]any{"type": "a", "val": float64(3)},
	}
	got, err := fn([]any{arr, "type"}, nil)
	if err != nil {
		t.Errorf("groupByKey: unexpected error: %v", err)
	}
	obj := got.(map[string]any)
	if len(obj["a"].([]any)) != 2 {
		t.Errorf("groupByKey: expected 2 items in group 'a'")
	}
	if len(obj["b"].([]any)) != 1 {
		t.Errorf("groupByKey: expected 1 item in group 'b'")
	}
}

func TestGroupByKeyErrors(t *testing.T) {
	fn := extarray.GroupByKey()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("groupByKey: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", "key"}, nil); err == nil {
		t.Error("groupByKey: expected error for non-array")
	}
	if _, err := fn([]any{[]any{}, 42}, nil); err == nil {
		t.Error("groupByKey: expected error for non-string key")
	}
	if _, err := fn([]any{[]any{"not-obj"}, "key"}, nil); err == nil {
		t.Error("groupByKey: expected error for non-object element")
	}
}

func TestSortBy(t *testing.T) {
	fn := extarray.SortBy()
	arr := []any{
		map[string]any{"name": "charlie"},
		map[string]any{"name": "alice"},
		map[string]any{"name": "bob"},
	}
	got, err := fn([]any{arr, "name"}, nil)
	if err != nil {
		t.Errorf("sortBy: unexpected error: %v", err)
	}
	result := got.([]any)
	names := []string{"alice", "bob", "charlie"}
	for i, n := range names {
		obj := result[i].(map[string]any)
		if obj["name"] != n {
			t.Errorf("sortBy[%d]: got %v, want %v", i, obj["name"], n)
		}
	}
	// desc
	got2, _ := fn([]any{arr, "name", true}, nil)
	result2 := got2.([]any)
	names2 := []string{"charlie", "bob", "alice"}
	for i, n := range names2 {
		obj := result2[i].(map[string]any)
		if obj["name"] != n {
			t.Errorf("sortBy desc[%d]: got %v, want %v", i, obj["name"], n)
		}
	}
}

func TestSortByErrors(t *testing.T) {
	fn := extarray.SortBy()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("sortBy: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", "key"}, nil); err == nil {
		t.Error("sortBy: expected error for non-array")
	}
	if _, err := fn([]any{[]any{}, 42}, nil); err == nil {
		t.Error("sortBy: expected error for non-string key")
	}
	if _, err := fn([]any{[]any{}, "key", "bad"}, nil); err == nil {
		t.Error("sortBy: expected error for non-bool desc")
	}
}

func TestUniqueBy(t *testing.T) {
	fn := extarray.UniqueBy()
	arr := []any{
		map[string]any{"id": float64(1), "val": "a"},
		map[string]any{"id": float64(1), "val": "b"},
		map[string]any{"id": float64(2), "val": "c"},
	}
	got, err := fn([]any{arr, "id"}, nil)
	if err != nil {
		t.Errorf("uniqueBy: unexpected error: %v", err)
	}
	result := got.([]any)
	if len(result) != 2 {
		t.Errorf("uniqueBy: got %d elements, want 2", len(result))
	}
}

func TestUniqueByErrors(t *testing.T) {
	fn := extarray.UniqueBy()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("uniqueBy: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", "key"}, nil); err == nil {
		t.Error("uniqueBy: expected error for non-array")
	}
	if _, err := fn([]any{[]any{}, 42}, nil); err == nil {
		t.Error("uniqueBy: expected error for non-string key")
	}
	if _, err := fn([]any{[]any{"not-obj"}, "key"}, nil); err == nil {
		t.Error("uniqueBy: expected error for non-object element")
	}
}

func TestSumByKey(t *testing.T) {
	fn := extarray.SumByKey()
	arr := []any{
		map[string]any{"type": "a", "amount": float64(10)},
		map[string]any{"type": "b", "amount": float64(5)},
		map[string]any{"type": "a", "amount": float64(3)},
	}
	got, err := fn([]any{arr, "type", "amount"}, nil)
	if err != nil {
		t.Errorf("sumByKey: unexpected error: %v", err)
	}
	obj := got.(map[string]any)
	if obj["a"] != float64(13) {
		t.Errorf("sumByKey a: got %v, want 13", obj["a"])
	}
	if obj["b"] != float64(5) {
		t.Errorf("sumByKey b: got %v, want 5", obj["b"])
	}
}

func TestSumByKeyErrors(t *testing.T) {
	fn := extarray.SumByKey()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("sumByKey: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", "k", "v"}, nil); err == nil {
		t.Error("sumByKey: expected error for non-array")
	}
	if _, err := fn([]any{[]any{}, 42, "v"}, nil); err == nil {
		t.Error("sumByKey: expected error for non-string key")
	}
	if _, err := fn([]any{[]any{"not-obj"}, "k", "v"}, nil); err == nil {
		t.Error("sumByKey: expected error for non-object element")
	}
}

func TestCountByKey(t *testing.T) {
	fn := extarray.CountByKey()
	arr := []any{
		map[string]any{"type": "a"},
		map[string]any{"type": "b"},
		map[string]any{"type": "a"},
	}
	got, err := fn([]any{arr, "type"}, nil)
	if err != nil {
		t.Errorf("countByKey: unexpected error: %v", err)
	}
	obj := got.(map[string]any)
	if obj["a"] != float64(2) {
		t.Errorf("countByKey a: got %v, want 2", obj["a"])
	}
	if obj["b"] != float64(1) {
		t.Errorf("countByKey b: got %v, want 1", obj["b"])
	}
}

func TestCountByKeyErrors(t *testing.T) {
	fn := extarray.CountByKey()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("countByKey: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", "key"}, nil); err == nil {
		t.Error("countByKey: expected error for non-array")
	}
	if _, err := fn([]any{[]any{}, 42}, nil); err == nil {
		t.Error("countByKey: expected error for non-string key")
	}
	if _, err := fn([]any{[]any{"not-obj"}, "key"}, nil); err == nil {
		t.Error("countByKey: expected error for non-object element")
	}
}

func TestRotate(t *testing.T) {
	fn := extarray.Rotate()
	cases := []struct {
		arr  []any
		n    float64
		want []any
	}{
		{[]any{float64(1), float64(2), float64(3), float64(4)}, 1, []any{float64(2), float64(3), float64(4), float64(1)}},
		{[]any{float64(1), float64(2), float64(3)}, -1, []any{float64(3), float64(1), float64(2)}},
		{[]any{float64(1), float64(2)}, 0, []any{float64(1), float64(2)}},
		{[]any{}, 2, []any{}},
	}
	for _, c := range cases {
		got, err := fn([]any{c.arr, c.n}, nil)
		if err != nil {
			t.Errorf("rotate %v by %v: unexpected error: %v", c.arr, c.n, err)
			continue
		}
		arr := got.([]any)
		for i, w := range c.want {
			if arr[i] != w {
				t.Errorf("rotate[%d]: got %v, want %v", i, arr[i], w)
			}
		}
	}
}

func TestRotateErrors(t *testing.T) {
	fn := extarray.Rotate()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("rotate: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", float64(1)}, nil); err == nil {
		t.Error("rotate: expected error for non-array")
	}
	if _, err := fn([]any{[]any{}, "bad"}, nil); err == nil {
		t.Error("rotate: expected error for non-integer n")
	}
}

func TestIndexof(t *testing.T) {
	fn := extarray.Indexof()
	arr := []any{float64(10), float64(20), float64(30)}
	cases := []struct {
		val  any
		want float64
	}{
		{float64(20), 1},
		{float64(99), -1},
	}
	for _, c := range cases {
		got, err := fn([]any{arr, c.val}, nil)
		if err != nil {
			t.Errorf("indexof: unexpected error: %v", err)
		}
		if got != c.want {
			t.Errorf("indexof(%v): got %v, want %v", c.val, got, c.want)
		}
	}
}

func TestIndexofErrors(t *testing.T) {
	fn := extarray.Indexof()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("indexof: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array", float64(1)}, nil); err == nil {
		t.Error("indexof: expected error for non-array")
	}
}

func TestTranspose(t *testing.T) {
	fn := extarray.Transpose()
	matrix := []any{
		[]any{float64(1), float64(2), float64(3)},
		[]any{float64(4), float64(5), float64(6)},
	}
	got, err := fn([]any{matrix}, nil)
	if err != nil {
		t.Errorf("transpose: unexpected error: %v", err)
	}
	result := got.([]any)
	if len(result) != 3 {
		t.Fatalf("transpose: got %d cols, want 3", len(result))
	}
	want := [][]float64{{1, 4}, {2, 5}, {3, 6}}
	for c, col := range want {
		arr := result[c].([]any)
		for r, v := range col {
			if arr[r] != v {
				t.Errorf("transpose[%d][%d]: got %v, want %v", c, r, arr[r], v)
			}
		}
	}
	// empty
	got2, _ := fn([]any{[]any{}}, nil)
	if len(got2.([]any)) != 0 {
		t.Errorf("transpose empty: got non-empty")
	}
}

func TestTransposeErrors(t *testing.T) {
	fn := extarray.Transpose()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("transpose: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array"}, nil); err == nil {
		t.Error("transpose: expected error for non-array")
	}
	if _, err := fn([]any{[]any{"not-row"}}, nil); err == nil {
		t.Error("transpose: expected error for non-array row")
	}
}

func TestAdjacentPairs(t *testing.T) {
	fn := extarray.AdjacentPairs()
	arr := []any{float64(1), float64(2), float64(3), float64(4)}
	got, err := fn([]any{arr}, nil)
	if err != nil {
		t.Errorf("adjacentPairs: unexpected error: %v", err)
	}
	result := got.([]any)
	if len(result) != 3 {
		t.Fatalf("adjacentPairs: got %d pairs, want 3", len(result))
	}
	pairs := [][2]float64{{1, 2}, {2, 3}, {3, 4}}
	for i, p := range pairs {
		pair := result[i].([]any)
		if pair[0] != p[0] || pair[1] != p[1] {
			t.Errorf("adjacentPairs[%d]: got %v, want %v", i, pair, p)
		}
	}
	// single element
	got2, _ := fn([]any{[]any{float64(1)}}, nil)
	if len(got2.([]any)) != 0 {
		t.Errorf("adjacentPairs single: got non-empty")
	}
}

func TestAdjacentPairsErrors(t *testing.T) {
	fn := extarray.AdjacentPairs()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("adjacentPairs: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array"}, nil); err == nil {
		t.Error("adjacentPairs: expected error for non-array")
	}
}
