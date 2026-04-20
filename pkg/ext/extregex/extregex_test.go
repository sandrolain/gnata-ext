package extregex_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extregex"
)

// ---------- RegexAll ----------

func TestRegexAll_Basic(t *testing.T) {
	fn := extregex.RegexAll()
	got, err := fn([]any{"abc123def456", `\d+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	arr := got.([]any)
	if len(arr) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(arr))
	}
	if arr[0] != "123" || arr[1] != "456" {
		t.Errorf("unexpected matches: %v", arr)
	}
}

func TestRegexAll_NoMatch(t *testing.T) {
	fn := extregex.RegexAll()
	got, err := fn([]any{"abc", `\d+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	arr := got.([]any)
	if len(arr) != 0 {
		t.Errorf("expected empty array, got %v", arr)
	}
}

func TestRegexAll_InvalidPattern(t *testing.T) {
	fn := extregex.RegexAll()
	_, err := fn([]any{"abc", `[invalid`}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexAll_NoArgs(t *testing.T) {
	fn := extregex.RegexAll()
	_, err := fn([]any{"only-one"}, nil)
	if err == nil {
		t.Error("expected error for missing pattern")
	}
}

func TestRegexAll_WrongType(t *testing.T) {
	fn := extregex.RegexAll()
	_, err := fn([]any{42, `\d+`}, nil)
	if err == nil {
		t.Error("expected error for non-string input")
	}
}

// ---------- RegexNamedGroups ----------

func TestRegexNamedGroups_Basic(t *testing.T) {
	fn := extregex.RegexNamedGroups()
	got, err := fn([]any{"2024-04-19", `(?P<y>\d{4})-(?P<m>\d{2})-(?P<d>\d{2})`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if obj["y"] != "2024" || obj["m"] != "04" || obj["d"] != "19" {
		t.Errorf("unexpected groups: %v", obj)
	}
}

func TestRegexNamedGroups_NoMatch(t *testing.T) {
	fn := extregex.RegexNamedGroups()
	got, err := fn([]any{"hello", `(?P<y>\d{4})`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	if len(obj) != 0 {
		t.Errorf("expected empty map, got %v", obj)
	}
}

func TestRegexNamedGroups_NoNamedGroups(t *testing.T) {
	fn := extregex.RegexNamedGroups()
	got, err := fn([]any{"hello123", `(\d+)`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	obj := got.(map[string]any)
	// No named groups — should return empty object
	if len(obj) != 0 {
		t.Errorf("expected empty object for unnamed groups, got %v", obj)
	}
}

func TestRegexNamedGroups_InvalidPattern(t *testing.T) {
	fn := extregex.RegexNamedGroups()
	_, err := fn([]any{"hello", `[bad`}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexNamedGroups_NoArgs(t *testing.T) {
	fn := extregex.RegexNamedGroups()
	_, err := fn([]any{"only"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- RegexSplit ----------

func TestRegexSplit_Basic(t *testing.T) {
	fn := extregex.RegexSplit()
	got, err := fn([]any{"one,two;;three", `[,;]+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	arr := got.([]any)
	if len(arr) != 3 {
		t.Fatalf("expected 3 parts, got %d: %v", len(arr), arr)
	}
	if arr[0] != "one" || arr[1] != "two" || arr[2] != "three" {
		t.Errorf("unexpected parts: %v", arr)
	}
}

func TestRegexSplit_NoSep(t *testing.T) {
	fn := extregex.RegexSplit()
	got, _ := fn([]any{"abc", `\d+`}, nil)
	arr := got.([]any)
	if len(arr) != 1 || arr[0] != "abc" {
		t.Errorf("expected ['abc'], got %v", arr)
	}
}

func TestRegexSplit_InvalidPattern(t *testing.T) {
	fn := extregex.RegexSplit()
	_, err := fn([]any{"abc", `[bad`}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexSplit_NoArgs(t *testing.T) {
	fn := extregex.RegexSplit()
	_, err := fn([]any{"s"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- RegexReplaceAll ----------

func TestRegexReplaceAll_Basic(t *testing.T) {
	fn := extregex.RegexReplaceAll()
	got, err := fn([]any{"hello world", `\bworld\b`, "Go"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello Go" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRegexReplaceAll_Backreference(t *testing.T) {
	fn := extregex.RegexReplaceAll()
	got, err := fn([]any{"John Smith", `(\w+) (\w+)`, "$2 $1"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Smith John" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestRegexReplaceAll_NoMatch(t *testing.T) {
	fn := extregex.RegexReplaceAll()
	got, _ := fn([]any{"hello", `\d+`, "X"}, nil)
	if got != "hello" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestRegexReplaceAll_InvalidPattern(t *testing.T) {
	fn := extregex.RegexReplaceAll()
	_, err := fn([]any{"abc", `[bad`, "x"}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexReplaceAll_MissingArgs(t *testing.T) {
	fn := extregex.RegexReplaceAll()
	_, err := fn([]any{"abc", `\d+`}, nil)
	if err == nil {
		t.Error("expected error for missing repl")
	}
}

func TestRegexReplaceAll_WrongType(t *testing.T) {
	fn := extregex.RegexReplaceAll()
	_, err := fn([]any{42, `\d+`, "x"}, nil)
	if err == nil {
		t.Error("expected error for non-string input")
	}
}

// ---------- RegexCount ----------

func TestRegexCount_Basic(t *testing.T) {
	fn := extregex.RegexCount()
	got, err := fn([]any{"abc123def456ghi", `\d+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != float64(2) {
		t.Errorf("expected 2, got %v", got)
	}
}

func TestRegexCount_Zero(t *testing.T) {
	fn := extregex.RegexCount()
	got, _ := fn([]any{"no-digits-here", `\d+`}, nil)
	if got != float64(0) {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestRegexCount_InvalidPattern(t *testing.T) {
	fn := extregex.RegexCount()
	_, err := fn([]any{"abc", `[bad`}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexCount_NoArgs(t *testing.T) {
	fn := extregex.RegexCount()
	_, err := fn([]any{"x"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- RegexTest ----------

func TestRegexTest_Match(t *testing.T) {
	fn := extregex.RegexTest()
	got, err := fn([]any{"hello123", `\d+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != true {
		t.Error("expected true")
	}
}

func TestRegexTest_NoMatch(t *testing.T) {
	fn := extregex.RegexTest()
	got, _ := fn([]any{"hello", `\d+`}, nil)
	if got != false {
		t.Error("expected false")
	}
}

func TestRegexTest_InvalidPattern(t *testing.T) {
	fn := extregex.RegexTest()
	_, err := fn([]any{"abc", `[bad`}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexTest_NoArgs(t *testing.T) {
	fn := extregex.RegexTest()
	_, err := fn([]any{"x"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- RegexExtract ----------

func TestRegexExtract_FullMatch(t *testing.T) {
	fn := extregex.RegexExtract()
	got, err := fn([]any{"abc-123-def", `\d+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "123" {
		t.Errorf("expected '123', got %v", got)
	}
}

func TestRegexExtract_Group(t *testing.T) {
	fn := extregex.RegexExtract()
	got, err := fn([]any{"2024-04-19", `(\d{4})-(\d{2})-(\d{2})`, float64(2)}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != "04" {
		t.Errorf("expected '04', got %v", got)
	}
}

func TestRegexExtract_NoMatch(t *testing.T) {
	fn := extregex.RegexExtract()
	got, err := fn([]any{"no-digits", `\d+`}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestRegexExtract_InvalidPattern(t *testing.T) {
	fn := extregex.RegexExtract()
	_, err := fn([]any{"abc", `[bad`}, nil)
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestRegexExtract_GroupOutOfRange(t *testing.T) {
	fn := extregex.RegexExtract()
	_, err := fn([]any{"abc123", `(\d+)`, float64(5)}, nil)
	if err == nil {
		t.Error("expected error for out-of-range group")
	}
}

func TestRegexExtract_NoArgs(t *testing.T) {
	fn := extregex.RegexExtract()
	_, err := fn([]any{"x"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRegexExtract_WrongGroupType(t *testing.T) {
	fn := extregex.RegexExtract()
	_, err := fn([]any{"abc123", `(\d+)`, "not-a-number"}, nil)
	if err == nil {
		t.Error("expected error for non-numeric group index")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := extregex.All()
	expected := []string{"regexAll", "regexNamedGroups", "regexSplit", "regexReplaceAll", "regexCount", "regexTest", "regexExtract"}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
