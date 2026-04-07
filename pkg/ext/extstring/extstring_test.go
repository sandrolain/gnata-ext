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

func TestStartsWithErrors(t *testing.T) {
	fn := extstring.StartsWith()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("startsWith: expected error for 0 args")
	}
	if _, err := fn([]any{1, "x"}, nil); err == nil {
		t.Error("startsWith: expected error for non-string first arg")
	}
	if _, err := fn([]any{"hello", 2}, nil); err == nil {
		t.Error("startsWith: expected error for non-string second arg")
	}
}

func TestEndsWithErrors(t *testing.T) {
	fn := extstring.EndsWith()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("endsWith: expected error for 0 args")
	}
	if _, err := fn([]any{1, "x"}, nil); err == nil {
		t.Error("endsWith: expected error for non-string first arg")
	}
}

func TestIndexOfWithStart(t *testing.T) {
	fn := extstring.IndexOf()
	// with start parameter
	got, _ := fn([]any{"abcabc", "b", 2.0}, nil)
	if got != float64(4) {
		t.Errorf("IndexOf with start: got %v", got)
	}
	// negative start treated as 0
	got, _ = fn([]any{"abcabc", "a", -5.0}, nil)
	if got != float64(0) {
		t.Errorf("IndexOf negative start: got %v", got)
	}
	// start >= len returns -1
	got, _ = fn([]any{"abc", "a", 10.0}, nil)
	if got != float64(-1) {
		t.Errorf("IndexOf start>=len: got %v", got)
	}
}

func TestIndexOfErrors(t *testing.T) {
	fn := extstring.IndexOf()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("indexOf: expected error for 0 args")
	}
	if _, err := fn([]any{42, "x"}, nil); err == nil {
		t.Error("indexOf: expected error for non-string first arg")
	}
	if _, err := fn([]any{"hello", "o", "bad"}, nil); err == nil {
		t.Error("indexOf: expected error for non-numeric start")
	}
}

func TestLastIndexOfErrors(t *testing.T) {
	fn := extstring.LastIndexOf()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("lastIndexOf: expected error for 0 args")
	}
	if _, err := fn([]any{42, "x"}, nil); err == nil {
		t.Error("lastIndexOf: expected error for non-string first arg")
	}
}

func TestCapitalizeEdge(t *testing.T) {
	fn := extstring.Capitalize()
	got, _ := fn([]any{""}, nil)
	if got != "" {
		t.Errorf("Capitalize empty: got %v", got)
	}
}

func TestCapitalizeErrors(t *testing.T) {
	fn := extstring.Capitalize()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("capitalize: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("capitalize: expected error for non-string")
	}
}

func TestTitleCaseErrors(t *testing.T) {
	fn := extstring.TitleCase()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("titleCase: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("titleCase: expected error for non-string")
	}
}

func TestCamelCaseErrors(t *testing.T) {
	fn := extstring.CamelCase()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("camelCase: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("camelCase: expected error for non-string")
	}
	// empty string returns ""
	got, _ := fn([]any{""}, nil)
	if got != "" {
		t.Errorf("camelCase empty: got %v", got)
	}
}

func TestSnakeCaseErrors(t *testing.T) {
	fn := extstring.SnakeCase()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("snakeCase: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("snakeCase: expected error for non-string")
	}
}

func TestKebabCaseErrors(t *testing.T) {
	fn := extstring.KebabCase()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("kebabCase: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("kebabCase: expected error for non-string")
	}
}

func TestRepeatErrors(t *testing.T) {
	fn := extstring.Repeat()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("repeat: expected error for 0 args")
	}
	if _, err := fn([]any{42, 3.0}, nil); err == nil {
		t.Error("repeat: expected error for non-string first arg")
	}
	if _, err := fn([]any{"ab", -1.0}, nil); err == nil {
		t.Error("repeat: expected error for negative n")
	}
	if _, err := fn([]any{"ab", "bad"}, nil); err == nil {
		t.Error("repeat: expected error for non-numeric n")
	}
}

func TestWordsEmpty(t *testing.T) {
	fn := extstring.Words()
	got, err := fn([]any{""}, nil)
	if err != nil || got != nil {
		t.Errorf("Words empty: expected nil result, got %v %v", got, err)
	}
}

func TestWordsErrors(t *testing.T) {
	fn := extstring.Words()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("words: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("words: expected error for non-string")
	}
}

func TestTemplateMissingKey(t *testing.T) {
	fn := extstring.Template()
	got, _ := fn([]any{"Hello, {{missing}}!", map[string]any{}}, nil)
	if got != "Hello, {{missing}}!" {
		t.Errorf("Template missing key: got %v", got)
	}
}

func TestTemplateErrors(t *testing.T) {
	fn := extstring.Template()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("template: expected error for 0 args")
	}
	if _, err := fn([]any{42, map[string]any{}}, nil); err == nil {
		t.Error("template: expected error for non-string first arg")
	}
	if _, err := fn([]any{"hello {{x}}", "not-an-object"}, nil); err == nil {
		t.Error("template: expected error for non-object bindings")
	}
}
