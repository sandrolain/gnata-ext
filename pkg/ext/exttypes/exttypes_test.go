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
func TestToArray(t *testing.T) {
	f := exttypes.ToArray()
	// nil -> empty slice
	got, _ := f([]any{nil}, nil)
	if len(got.([]any)) != 0 {
		t.Errorf("toArray nil: expected empty slice")
	}
	// no args -> empty slice
	got, _ = f([]any{}, nil)
	if len(got.([]any)) != 0 {
		t.Errorf("toArray no args: expected empty slice")
	}
	// scalar -> wrapped
	got, _ = f([]any{"hello"}, nil)
	arr := got.([]any)
	if len(arr) != 1 || arr[0] != "hello" {
		t.Errorf("toArray scalar: got %v", arr)
	}
	// array -> passthrough
	input := []any{float64(1), float64(2)}
	got, _ = f([]any{input}, nil)
	if len(got.([]any)) != 2 {
		t.Errorf("toArray array: got %v", got)
	}
}

func TestDefined(t *testing.T) {
	f := exttypes.Defined()
	got, _ := f([]any{nil}, nil)
	if got != false {
		t.Error("defined(nil): expected false")
	}
	got, _ = f([]any{"x"}, nil)
	if got != true {
		t.Error("defined(x): expected true")
	}
	got, _ = f([]any{}, nil)
	if got != false {
		t.Error("defined(no args): expected false")
	}
}

func TestNullish(t *testing.T) {
	f := exttypes.Nullish()
	// nil -> fallback
	got, _ := f([]any{nil, "fallback"}, nil)
	if got != "fallback" {
		t.Errorf("nullish nil: got %v", got)
	}
	// non-nil -> first arg
	got, _ = f([]any{"value", "fallback"}, nil)
	if got != "value" {
		t.Errorf("nullish non-nil: got %v", got)
	}
	// no args -> nil
	got, _ = f([]any{}, nil)
	if got != nil {
		t.Errorf("nullish no args: got %v", got)
	}
}

func TestTypeOf(t *testing.T) {
	f := exttypes.TypeOf()
	tests := []struct {
		v    any
		want string
	}{
		{nil, "null"},
		{"hello", "string"},
		{float64(1), "number"},
		{true, "boolean"},
		{[]any{}, "array"},
		{map[string]any{}, "object"},
	}
	for _, tc := range tests {
		got, _ := f([]any{tc.v}, nil)
		if got != tc.want {
			t.Errorf("typeOf(%v): got %v, want %v", tc.v, got, tc.want)
		}
	}
	// no args -> null
	got, _ := f([]any{}, nil)
	if got != "null" {
		t.Errorf("typeOf no args: got %v", got)
	}
}

func TestToNumber(t *testing.T) {
	f := exttypes.ToNumber()
	tests := []struct {
		v    any
		want float64
	}{
		{float64(3.14), 3.14},
		{int(5), 5},
		{int64(7), 7},
		{true, 1},
		{false, 0},
		{"42", 42},
	}
	for _, tc := range tests {
		got, err := f([]any{tc.v}, nil)
		if err != nil {
			t.Errorf("toNumber(%v): unexpected error: %v", tc.v, err)
		}
		if got.(float64) != tc.want {
			t.Errorf("toNumber(%v): got %v, want %v", tc.v, got, tc.want)
		}
	}
	// error cases
	if _, err := f([]any{}, nil); err == nil {
		t.Error("toNumber: expected error for 0 args")
	}
	if _, err := f([]any{"not-a-number"}, nil); err == nil {
		t.Error("toNumber: expected error for invalid string")
	}
	if _, err := f([]any{map[string]any{}}, nil); err == nil {
		t.Error("toNumber: expected error for object")
	}
}

func TestToString(t *testing.T) {
	f := exttypes.ToString()
	got, _ := f([]any{float64(42)}, nil)
	if got != "42" {
		t.Errorf("toString float64: got %v", got)
	}
	got, _ = f([]any{true}, nil)
	if got != "true" {
		t.Errorf("toString bool: got %v", got)
	}
	got, _ = f([]any{nil}, nil)
	if got != "<nil>" {
		t.Errorf("toString nil: got %v", got)
	}
	got, _ = f([]any{}, nil)
	if got != "" {
		t.Errorf("toString no args: got %v", got)
	}
}

func TestToBool(t *testing.T) {
	f := exttypes.ToBool()
	tests := []struct {
		v    any
		want bool
	}{
		{nil, false},
		{false, false},
		{float64(0), false},
		{"", false},
		{[]any{}, false},
		{map[string]any{}, false},
		{true, true},
		{float64(1), true},
		{"x", true},
		{[]any{float64(1)}, true},
	}
	for _, tc := range tests {
		got, _ := f([]any{tc.v}, nil)
		if got != tc.want {
			t.Errorf("toBool(%v): got %v, want %v", tc.v, got, tc.want)
		}
	}
	// no args
	got, _ := f([]any{}, nil)
	if got != false {
		t.Errorf("toBool no args: got %v", got)
	}
}
