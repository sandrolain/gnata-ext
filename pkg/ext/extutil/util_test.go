package extutil_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

func TestToFloat(t *testing.T) {
	cases := []struct {
		input any
		want  float64
		ok    bool
	}{
		{float64(3.14), 3.14, true},
		{int(5), 5.0, true},
		{int64(10), 10.0, true},
		{"bad", 0, false},
	}
	for _, c := range cases {
		got, err := extutil.ToFloat(c.input)
		if c.ok && (err != nil || got != c.want) {
			t.Errorf("ToFloat(%v): got %v, %v; want %v", c.input, got, err, c.want)
		} else if !c.ok && err == nil {
			t.Errorf("ToFloat(%v): expected error", c.input)
		}
	}
}

func TestToInt(t *testing.T) {
	cases := []struct {
		input any
		want  int
		ok    bool
	}{
		{float64(7.0), 7, true},
		{int(3), 3, true},
		{int64(9), 9, true},
		{"bad", 0, false},
	}
	for _, c := range cases {
		got, ok := extutil.ToInt(c.input)
		if ok != c.ok || (ok && got != c.want) {
			t.Errorf("ToInt(%v): got (%v, %v); want (%v, %v)", c.input, got, ok, c.want, c.ok)
		}
	}
}

func TestToFloatSlice(t *testing.T) {
	got, err := extutil.ToFloatSlice([]any{1.0, 2.0, 3.0})
	if err != nil || len(got) != 3 || got[0] != 1.0 {
		t.Errorf("ToFloatSlice: got %v, %v", got, err)
	}
	_, err = extutil.ToFloatSlice("bad")
	if err == nil {
		t.Error("ToFloatSlice with string: expected error")
	}
}

func TestToArray(t *testing.T) {
	arr := []any{1.0, "two"}
	got, err := extutil.ToArray(arr)
	if err != nil || len(got) != 2 {
		t.Errorf("ToArray: got %v, %v", got, err)
	}
	_, err = extutil.ToArray("not array")
	if err == nil {
		t.Error("ToArray with string: expected error")
	}
}

func TestToObject(t *testing.T) {
	obj := map[string]any{"k": "v"}
	got, err := extutil.ToObject(obj)
	if err != nil || got["k"] != "v" {
		t.Errorf("ToObject: got %v, %v", got, err)
	}
	_, err = extutil.ToObject("not obj")
	if err == nil {
		t.Error("ToObject with string: expected error")
	}
}
