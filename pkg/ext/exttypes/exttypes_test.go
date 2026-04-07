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

func TestIsUndefined(t *testing.T) {
	f := exttypes.IsUndefined()
	// 0 args -> true
	got, err := f([]any{}, nil)
	if err != nil || !got.(bool) {
		t.Errorf("isUndefined(no args): got %v, %v", got, err)
	}
	// nil arg -> true
	got, _ = invoke(f, nil)
	if !got.(bool) {
		t.Errorf("isUndefined(nil): expected true")
	}
	// non-nil arg -> false
	got, _ = invoke(f, "something")
	if got.(bool) {
		t.Errorf("isUndefined(string): expected false")
	}
}

func TestIsStringNoArgs(t *testing.T) {
	f := exttypes.IsString()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isString(no args): expected false")
	}
}

func TestIsNumberNoArgs(t *testing.T) {
	f := exttypes.IsNumber()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isNumber(no args): expected false")
	}
	// int64 variant
	got, _ = invoke(f, int64(5))
	if !got.(bool) {
		t.Error("isNumber(int64): expected true")
	}
	// int variant
	got, _ = invoke(f, int(5))
	if !got.(bool) {
		t.Error("isNumber(int): expected true")
	}
}

func TestIsBooleanNoArgs(t *testing.T) {
	f := exttypes.IsBoolean()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isBoolean(no args): expected false")
	}
}

func TestIsArrayNoArgs(t *testing.T) {
	f := exttypes.IsArray()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isArray(no args): expected false")
	}
}

func TestIsObjectNoArgs(t *testing.T) {
	f := exttypes.IsObject()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isObject(no args): expected false")
	}
}

func TestIsNullNoArgs(t *testing.T) {
	f := exttypes.IsNull()
	got, _ := f([]any{}, nil)
	if !got.(bool) {
		t.Error("isNull(no args): expected true")
	}
}

func TestIsEmptyDefault(t *testing.T) {
	f := exttypes.IsEmpty()
	// non-nil, non-string, non-array, non-object -> false
	got, _ := invoke(f, 42.0)
	if got.(bool) {
		t.Error("isEmpty(42): expected false")
	}
	// no args -> true
	got, _ = f([]any{}, nil)
	if !got.(bool) {
		t.Error("isEmpty(no args): expected true")
	}
}

func TestDefaultEdgeCases(t *testing.T) {
	f := exttypes.Default()
	// 0 args -> nil
	got, err := f([]any{}, nil)
	if err != nil || got != nil {
		t.Errorf("default(0 args): got %v, %v", got, err)
	}
	// 1 arg non-nil -> returns it
	got, _ = f([]any{"val"}, nil)
	if got != "val" {
		t.Errorf("default(1 arg): got %v", got)
	}
}

func TestIdentityNoArgs(t *testing.T) {
	f := exttypes.Identity()
	got, err := f([]any{}, nil)
	if err != nil || got != nil {
		t.Errorf("identity(no args): got %v, %v", got, err)
	}
}
