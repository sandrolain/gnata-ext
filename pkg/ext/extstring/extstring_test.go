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

func TestPadStart(t *testing.T) {
	fn := extstring.PadStart()
	cases := []struct {
		args []any
		want any
	}{
		{[]any{"hi", float64(5)}, "   hi"},
		{[]any{"hi", float64(5), "0"}, "000hi"},
		{[]any{"hello", float64(3)}, "hello"},
		{[]any{"hi", float64(5), "xy"}, "xyxhi"},
	}
	for _, c := range cases {
		got, err := fn(c.args, nil)
		if err != nil {
			t.Errorf("padStart %v: unexpected error: %v", c.args, err)
		}
		if got != c.want {
			t.Errorf("padStart %v: got %v, want %v", c.args, got, c.want)
		}
	}
}

func TestPadStartErrors(t *testing.T) {
	fn := extstring.PadStart()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("padStart: expected error for 0 args")
	}
	if _, err := fn([]any{42, float64(5)}, nil); err == nil {
		t.Error("padStart: expected error for non-string")
	}
	if _, err := fn([]any{"hi", "bad"}, nil); err == nil {
		t.Error("padStart: expected error for non-int len")
	}
	if _, err := fn([]any{"hi", float64(5), ""}, nil); err == nil {
		t.Error("padStart: expected error for empty fill")
	}
}

func TestPadEnd(t *testing.T) {
	fn := extstring.PadEnd()
	cases := []struct {
		args []any
		want any
	}{
		{[]any{"hi", float64(5)}, "hi   "},
		{[]any{"hi", float64(5), "0"}, "hi000"},
		{[]any{"hello", float64(3)}, "hello"},
		{[]any{"hi", float64(5), "xy"}, "hixyx"},
	}
	for _, c := range cases {
		got, err := fn(c.args, nil)
		if err != nil {
			t.Errorf("padEnd %v: unexpected error: %v", c.args, err)
		}
		if got != c.want {
			t.Errorf("padEnd %v: got %v, want %v", c.args, got, c.want)
		}
	}
}

func TestPadEndErrors(t *testing.T) {
	fn := extstring.PadEnd()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("padEnd: expected error for 0 args")
	}
	if _, err := fn([]any{42, float64(5)}, nil); err == nil {
		t.Error("padEnd: expected error for non-string")
	}
	if _, err := fn([]any{"hi", "bad"}, nil); err == nil {
		t.Error("padEnd: expected error for non-int len")
	}
	if _, err := fn([]any{"hi", float64(5), ""}, nil); err == nil {
		t.Error("padEnd: expected error for empty fill")
	}
}

func TestTruncate(t *testing.T) {
	fn := extstring.Truncate()
	cases := []struct {
		args []any
		want string
	}{
		{[]any{"Hello World", float64(5)}, "Hello…"},
		{[]any{"Hi", float64(10)}, "Hi"},
		{[]any{"Hello World", float64(5), "..."}, "Hello..."},
		{[]any{"Hello World", float64(5), ""}, "Hello"},
	}
	for _, c := range cases {
		got, err := fn(c.args, nil)
		if err != nil {
			t.Errorf("truncate %v: unexpected error: %v", c.args, err)
		}
		if got != c.want {
			t.Errorf("truncate %v: got %v, want %v", c.args, got, c.want)
		}
	}
}

func TestTruncateErrors(t *testing.T) {
	fn := extstring.Truncate()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("truncate: expected error for 0 args")
	}
	if _, err := fn([]any{42, float64(5)}, nil); err == nil {
		t.Error("truncate: expected error for non-string")
	}
	if _, err := fn([]any{"hi", "bad"}, nil); err == nil {
		t.Error("truncate: expected error for non-int len")
	}
	if _, err := fn([]any{"hi", float64(5), 42}, nil); err == nil {
		t.Error("truncate: expected error for non-string suffix")
	}
}

func TestSlugify(t *testing.T) {
	fn := extstring.Slugify()
	cases := []struct {
		input string
		want  string
	}{
		{"Hello World", "hello-world"},
		{"  Foo  Bar  ", "foo-bar"},
		{"It's a Test!", "it-s-a-test"},
		{"already-slug", "already-slug"},
	}
	for _, c := range cases {
		got, err := fn([]any{c.input}, nil)
		if err != nil {
			t.Errorf("slugify %q: unexpected error: %v", c.input, err)
		}
		if got != c.want {
			t.Errorf("slugify %q: got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestSlugifyErrors(t *testing.T) {
	fn := extstring.Slugify()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("slugify: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("slugify: expected error for non-string")
	}
}

func TestCountOccurrences(t *testing.T) {
	fn := extstring.CountOccurrences()
	cases := []struct {
		str  string
		sub  string
		want float64
	}{
		{"hello hello hello", "hello", 3},
		{"banana", "an", 2},
		{"abc", "xyz", 0},
		{"abc", "", 0},
	}
	for _, c := range cases {
		got, err := fn([]any{c.str, c.sub}, nil)
		if err != nil {
			t.Errorf("countOccurrences: unexpected error: %v", err)
		}
		if got != c.want {
			t.Errorf("countOccurrences(%q, %q): got %v, want %v", c.str, c.sub, got, c.want)
		}
	}
}

func TestCountOccurrencesErrors(t *testing.T) {
	fn := extstring.CountOccurrences()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("countOccurrences: expected error for 0 args")
	}
	if _, err := fn([]any{42, "x"}, nil); err == nil {
		t.Error("countOccurrences: expected error for non-string first arg")
	}
}

func TestInitials(t *testing.T) {
	fn := extstring.Initials()
	cases := []struct {
		args []any
		want string
	}{
		{[]any{"John Doe"}, "JD"},
		{[]any{"John Doe", "."}, "J.D"},
		{[]any{"alice bob charlie"}, "ABC"},
		{[]any{""}, ""},
	}
	for _, c := range cases {
		got, err := fn(c.args, nil)
		if err != nil {
			t.Errorf("initials %v: unexpected error: %v", c.args, err)
		}
		if got != c.want {
			t.Errorf("initials %v: got %v, want %v", c.args, got, c.want)
		}
	}
}

func TestInitialsErrors(t *testing.T) {
	fn := extstring.Initials()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("initials: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("initials: expected error for non-string")
	}
	if _, err := fn([]any{"hi", 42}, nil); err == nil {
		t.Error("initials: expected error for non-string sep")
	}
}

func TestEscapeHTML(t *testing.T) {
	fn := extstring.EscapeHTML()
	got, err := fn([]any{"<div class=\"test\">Hello & World</div>"}, nil)
	if err != nil {
		t.Errorf("escapeHTML: unexpected error: %v", err)
	}
	want := "&lt;div class=&#34;test&#34;&gt;Hello &amp; World&lt;/div&gt;"
	if got != want {
		t.Errorf("escapeHTML: got %v, want %v", got, want)
	}
}

func TestEscapeHTMLErrors(t *testing.T) {
	fn := extstring.EscapeHTML()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("escapeHTML: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("escapeHTML: expected error for non-string")
	}
}

func TestUnescapeHTML(t *testing.T) {
	fn := extstring.UnescapeHTML()
	got, err := fn([]any{"&lt;b&gt;Hello &amp; World&lt;/b&gt;"}, nil)
	if err != nil {
		t.Errorf("unescapeHTML: unexpected error: %v", err)
	}
	if got != "<b>Hello & World</b>" {
		t.Errorf("unescapeHTML: got %v", got)
	}
}

func TestUnescapeHTMLErrors(t *testing.T) {
	fn := extstring.UnescapeHTML()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("unescapeHTML: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("unescapeHTML: expected error for non-string")
	}
}

func TestReverseWords(t *testing.T) {
	fn := extstring.ReverseWords()
	cases := []struct {
		input string
		want  string
	}{
		{"hello world foo", "foo world hello"},
		{"one", "one"},
		{"", ""},
	}
	for _, c := range cases {
		got, err := fn([]any{c.input}, nil)
		if err != nil {
			t.Errorf("reverseWords %q: unexpected error: %v", c.input, err)
		}
		if got != c.want {
			t.Errorf("reverseWords %q: got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestReverseWordsErrors(t *testing.T) {
	fn := extstring.ReverseWords()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("reverseWords: expected error for 0 args")
	}
	if _, err := fn([]any{42}, nil); err == nil {
		t.Error("reverseWords: expected error for non-string")
	}
}

func TestLevenshtein(t *testing.T) {
	fn := extstring.Levenshtein()
	cases := []struct {
		a, b string
		want float64
	}{
		{"kitten", "sitting", 3},
		{"", "abc", 3},
		{"abc", "", 3},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
	}
	for _, c := range cases {
		got, err := fn([]any{c.a, c.b}, nil)
		if err != nil {
			t.Errorf("levenshtein(%q,%q): unexpected error: %v", c.a, c.b, err)
		}
		if got != c.want {
			t.Errorf("levenshtein(%q,%q): got %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestLevenshteinErrors(t *testing.T) {
	fn := extstring.Levenshtein()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("levenshtein: expected error for 0 args")
	}
	if _, err := fn([]any{42, "b"}, nil); err == nil {
		t.Error("levenshtein: expected error for non-string first arg")
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	fn := extstring.LongestCommonPrefix()
	cases := []struct {
		strs []any
		want string
	}{
		{[]any{"flower", "flow", "flight"}, "fl"},
		{[]any{"dog", "racecar", "car"}, ""},
		{[]any{"abc", "abcd", "abce"}, "abc"},
		{[]any{}, ""},
	}
	for _, c := range cases {
		got, err := fn([]any{c.strs}, nil)
		if err != nil {
			t.Errorf("longestCommonPrefix: unexpected error: %v", err)
		}
		if got != c.want {
			t.Errorf("longestCommonPrefix(%v): got %v, want %v", c.strs, got, c.want)
		}
	}
}

func TestLongestCommonPrefixErrors(t *testing.T) {
	fn := extstring.LongestCommonPrefix()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("longestCommonPrefix: expected error for 0 args")
	}
	if _, err := fn([]any{"not-an-array"}, nil); err == nil {
		t.Error("longestCommonPrefix: expected error for non-array")
	}
	if _, err := fn([]any{[]any{"valid", 42}}, nil); err == nil {
		t.Error("longestCommonPrefix: expected error for non-string element")
	}
}
