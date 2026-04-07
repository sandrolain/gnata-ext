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
