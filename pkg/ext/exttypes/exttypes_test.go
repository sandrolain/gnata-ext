package exttypes_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

func TestIsString(t *testing.T) {
	f := exttypes.IsString()
	cases := []struct {
		in   any
		want bool
	}{
		{"hello", true},
		{42.0, false},
		{nil, false},
		{true, false},
	}
	for _, c := range cases {
		got, err := invoke(f, c.in)
		if err != nil || got.(bool) != c.want {
			t.Errorf("isString(%v): got %v, %v; want %v", c.in, got, err, c.want)
		}
	}
}

func TestIsNumber(t *testing.T) {
	f := exttypes.IsNumber()
	cases := []struct {
		in   any
		want bool
	}{
		{42.0, true},
		{int(5), true},
		{int64(10), true},
		{"nope", false},
		{nil, false},
	}
	for _, c := range cases {
		got, err := invoke(f, c.in)
		if err != nil || got.(bool) != c.want {
			t.Errorf("isNumber(%v): got %v, %v; want %v", c.in, got, err, c.want)
		}
	}
}

func TestIsBoolean(t *testing.T) {
	f := exttypes.IsBoolean()
	got, err := invoke(f, true)
	if err != nil || !got.(bool) {
		t.Errorf("isBoolean(true): got %v, %v", got, err)
	}
	got, err = invoke(f, 1.0)
	if err != nil || got.(bool) {
		t.Errorf("isBoolean(1.0): got %v, %v", got, err)
	}
}

func TestIsArray(t *testing.T) {
	f := exttypes.IsArray()
	got, err := invoke(f, []any{1, 2})
	if err != nil || !got.(bool) {
		t.Errorf("isArray([1,2]): got %v, %v", got, err)
	}
	got, err = invoke(f, "array")
	if err != nil || got.(bool) {
		t.Errorf("isArray(string): got %v, %v", got, err)
	}
}

func TestIsObject(t *testing.T) {
	f := exttypes.IsObject()
	got, err := invoke(f, map[string]any{"k": "v"})
	if err != nil || !got.(bool) {
		t.Errorf("isObject(map): got %v, %v", got, err)
	}
	got, err = invoke(f, "object")
	if err != nil || got.(bool) {
		t.Errorf("isObject(string): got %v, %v", got, err)
	}
}

func TestIsNull(t *testing.T) {
	f := exttypes.IsNull()
	got, err := invoke(f, nil)
	if err != nil || !got.(bool) {
		t.Errorf("isNull(nil): got %v, %v", got, err)
	}
	got, err = invoke(f, "val")
	if err != nil || got.(bool) {
		t.Errorf("isNull(string): got %v, %v", got, err)
	}
}

func TestIsEmpty(t *testing.T) {
	f := exttypes.IsEmpty()
	cases := []struct {
		in   any
		want bool
	}{
		{nil, true},
		{"", true},
		{[]any{}, true},
		{map[string]any{}, true},
		{"hello", false},
		{[]any{1}, false},
		{map[string]any{"k": "v"}, false},
	}
	for _, c := range cases {
		got, err := invoke(f, c.in)
		if err != nil || got.(bool) != c.want {
			t.Errorf("isEmpty(%v): got %v, %v; want %v", c.in, got, err, c.want)
		}
	}
}

func TestDefault(t *testing.T) {
	f := exttypes.Default()
	got, err := invoke(f, nil, "fallback")
	if err != nil || got.(string) != "fallback" {
		t.Errorf("default(nil, fallback): got %v, %v", got, err)
	}
	got, err = invoke(f, "value", "fallback")
	if err != nil || got.(string) != "value" {
		t.Errorf("default(value, fallback): got %v, %v", got, err)
	}
}

func TestIdentity(t *testing.T) {
	f := exttypes.Identity()
	obj := map[string]any{"x": 1.0}
	got, err := invoke(f, obj)
	if err != nil {
		t.Fatalf("identity: unexpected error: %v", err)
	}
	m := got.(map[string]any)
	if m["x"].(float64) != 1.0 {
		t.Errorf("identity: got %v", got)
	}
}

func TestAll(t *testing.T) {
	all := exttypes.All()
	expected := []string{"isString", "isNumber", "isBoolean", "isArray", "isObject",
		"isNull", "isUndefined", "isEmpty", "default", "identity"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing function %q", name)
		}
	}
}
