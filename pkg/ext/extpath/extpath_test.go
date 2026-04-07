package extpath_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extpath"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

func TestGet(t *testing.T) {
	f := extpath.Get()
	obj := map[string]any{"a": map[string]any{"b": "hello"}}

	got, err := invoke(f, obj, "a.b")
	if err != nil || got != "hello" {
		t.Errorf("get a.b: got %v, %v", got, err)
	}

	got, err = invoke(f, obj, "a.c", "default")
	if err != nil || got != "default" {
		t.Errorf("get missing with default: got %v, %v", got, err)
	}

	got, err = invoke(f, obj, "x")
	if err != nil || got != nil {
		t.Errorf("get missing no default: got %v, %v", got, err)
	}
}

func TestSet(t *testing.T) {
	f := extpath.Set()
	obj := map[string]any{"a": map[string]any{"b": 1.0}}

	got, err := invoke(f, obj, "a.b", "new")
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	m := got.(map[string]any)
	inner := m["a"].(map[string]any)
	if inner["b"] != "new" {
		t.Errorf("set a.b: got %v", inner["b"])
	}
	// original must be unchanged
	if obj["a"].(map[string]any)["b"] != 1.0 {
		t.Error("set must be immutable")
	}

	// create nested
	got, err = invoke(f, map[string]any{}, "x.y.z", 42.0)
	if err != nil {
		t.Fatalf("set nested create: %v", err)
	}
	v, _ := extpath.Get()([]any{got, "x.y.z"}, nil)
	if v != 42.0 {
		t.Errorf("set nested create: got %v", v)
	}
}

func TestDel(t *testing.T) {
	f := extpath.Del()
	obj := map[string]any{"a": map[string]any{"b": 1.0, "c": 2.0}, "d": 3.0}

	got, err := invoke(f, obj, "a.b")
	if err != nil {
		t.Fatalf("del: %v", err)
	}
	m := got.(map[string]any)
	inner := m["a"].(map[string]any)
	if _, ok := inner["b"]; ok {
		t.Error("del a.b: key should be removed")
	}
	if inner["c"] != 2.0 {
		t.Error("del a.b: key c should remain")
	}
}

func TestHas(t *testing.T) {
	f := extpath.Has()
	obj := map[string]any{"a": map[string]any{"b": "v"}}

	got, err := invoke(f, obj, "a.b")
	if err != nil || got != true {
		t.Errorf("has a.b: got %v, %v", got, err)
	}

	got, err = invoke(f, obj, "a.x")
	if err != nil || got != false {
		t.Errorf("has a.x: got %v, %v", got, err)
	}
}

func TestFlattenObj(t *testing.T) {
	f := extpath.FlattenObj()
	obj := map[string]any{
		"a": map[string]any{
			"b": 1.0,
			"c": map[string]any{"d": 2.0},
		},
		"e": 3.0,
	}
	got, err := invoke(f, obj)
	if err != nil {
		t.Fatalf("flattenObj: %v", err)
	}
	m := got.(map[string]any)
	if m["a.b"] != 1.0 {
		t.Errorf("flattenObj a.b: got %v", m["a.b"])
	}
	if m["a.c.d"] != 2.0 {
		t.Errorf("flattenObj a.c.d: got %v", m["a.c.d"])
	}
	if m["e"] != 3.0 {
		t.Errorf("flattenObj e: got %v", m["e"])
	}
}

func TestExpandObj(t *testing.T) {
	f := extpath.ExpandObj()
	obj := map[string]any{
		"a.b": 1.0,
		"a.c": 2.0,
		"e":   3.0,
	}
	got, err := invoke(f, obj)
	if err != nil {
		t.Fatalf("expandObj: %v", err)
	}
	m := got.(map[string]any)
	inner, ok := m["a"].(map[string]any)
	if !ok {
		t.Fatalf("expandObj: 'a' not a map: %T", m["a"])
	}
	if inner["b"] != 1.0 || inner["c"] != 2.0 {
		t.Errorf("expandObj a: got %v", inner)
	}
	if m["e"] != 3.0 {
		t.Errorf("expandObj e: got %v", m["e"])
	}
}

func TestAll(t *testing.T) {
	all := extpath.All()
	for _, name := range []string{"get", "set", "del", "has", "flattenObj", "expandObj"} {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing %q", name)
		}
	}
}
