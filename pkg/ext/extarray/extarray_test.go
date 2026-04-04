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
