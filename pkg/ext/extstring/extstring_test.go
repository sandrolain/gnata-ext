package extstring_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extstring"
)

func TestStartsWith(t *testing.T) {
	fn := extstring.StartsWith()
	got, _ := fn([]any{"hello world", "hello"}, nil)
	if got != true {
		t.Errorf("StartsWith true: got %v", got)
	}
	got, _ = fn([]any{"hello world", "world"}, nil)
	if got != false {
		t.Errorf("StartsWith false: got %v", got)
	}
}

func TestEndsWith(t *testing.T) {
	fn := extstring.EndsWith()
	got, _ := fn([]any{"hello world", "world"}, nil)
	if got != true {
		t.Errorf("EndsWith true: got %v", got)
	}
	got, _ = fn([]any{"hello world", "hello"}, nil)
	if got != false {
		t.Errorf("EndsWith false: got %v", got)
	}
}

func TestIndexOf(t *testing.T) {
	fn := extstring.IndexOf()
	got, _ := fn([]any{"hello world", "world"}, nil)
	if got != float64(6) {
		t.Errorf("IndexOf: got %v", got)
	}
	got, _ = fn([]any{"hello world", "xyz"}, nil)
	if got != float64(-1) {
		t.Errorf("IndexOf not found: got %v", got)
	}
}

func TestLastIndexOf(t *testing.T) {
	fn := extstring.LastIndexOf()
	got, _ := fn([]any{"abcabc", "b"}, nil)
	if got != float64(4) {
		t.Errorf("LastIndexOf: got %v", got)
	}
}

func TestCapitalize(t *testing.T) {
	fn := extstring.Capitalize()
	got, _ := fn([]any{"hello WORLD"}, nil)
	if got != "Hello world" {
		t.Errorf("Capitalize: got %v", got)
	}
}

func TestTitleCase(t *testing.T) {
	fn := extstring.TitleCase()
	got, _ := fn([]any{"hello world"}, nil)
	if got != "Hello World" {
		t.Errorf("TitleCase: got %v", got)
	}
}

func TestCamelCase(t *testing.T) {
	fn := extstring.CamelCase()
	got, _ := fn([]any{"hello_world_foo"}, nil)
	if got != "helloWorldFoo" {
		t.Errorf("CamelCase: got %v", got)
	}
}

func TestSnakeCase(t *testing.T) {
	fn := extstring.SnakeCase()
	got, _ := fn([]any{"helloWorldFoo"}, nil)
	if got != "hello_world_foo" {
		t.Errorf("SnakeCase: got %v", got)
	}
}

func TestKebabCase(t *testing.T) {
	fn := extstring.KebabCase()
	got, _ := fn([]any{"Hello World Foo"}, nil)
	if got != "hello-world-foo" {
		t.Errorf("KebabCase: got %v", got)
	}
}

func TestRepeat(t *testing.T) {
	fn := extstring.Repeat()
	got, _ := fn([]any{"ab", 3.0}, nil)
	if got != "ababab" {
		t.Errorf("Repeat: got %v", got)
	}
}

func TestWords(t *testing.T) {
	fn := extstring.Words()
	got, _ := fn([]any{"hello world foo"}, nil)
	arr := got.([]any)
	if len(arr) != 3 || arr[0] != "hello" {
		t.Errorf("Words: got %v", arr)
	}
}

func TestTemplate(t *testing.T) {
	fn := extstring.Template()
	got, _ := fn([]any{"Hello, {{name}}!", map[string]any{"name": "Alice"}}, nil)
	if got != "Hello, Alice!" {
		t.Errorf("Template: got %v", got)
	}
}

func TestAllKeys(t *testing.T) {
	m := extstring.All()
	for _, k := range []string{"startsWith", "endsWith", "indexOf", "lastIndexOf",
		"capitalize", "titleCase", "camelCase", "snakeCase", "kebabCase",
		"repeat", "words", "template"} {
		if _, ok := m[k]; !ok {
			t.Errorf("All() missing: %s", k)
		}
	}
}
