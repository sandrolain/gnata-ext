package extvalidate_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extvalidate"
)

func TestIsEmail(t *testing.T) {
	fn := extvalidate.IsEmail()
	cases := []struct {
		input any
		want  bool
	}{
		{"user@example.com", true},
		{"user+tag@sub.domain.co.uk", true},
		{"notanemail", false},
		{"@no-user.com", false},
		{"no-domain@", false},
		{"", false},
		{42, false},
	}
	for _, c := range cases {
		got, err := fn([]any{c.input}, nil)
		if err != nil {
			t.Errorf("IsEmail(%v) unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("IsEmail(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsURL(t *testing.T) {
	fn := extvalidate.IsURL()
	cases := []struct {
		input any
		want  bool
	}{
		{"https://example.com", true},
		{"http://sub.example.com/path?q=1", true},
		{"ftp://files.example.com", true},
		{"not-a-url", false},
		{"file:///local", false},
		{"", false},
	}
	for _, c := range cases {
		got, err := fn([]any{c.input}, nil)
		if err != nil {
			t.Errorf("IsURL(%v) unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("IsURL(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsUUID(t *testing.T) {
	fn := extvalidate.IsUUID()
	cases := []struct {
		input any
		want  bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true}, // uppercase
		{"not-a-uuid", false},
		{"", false},
	}
	for _, c := range cases {
		got, err := fn([]any{c.input}, nil)
		if err != nil {
			t.Errorf("IsUUID(%v) unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("IsUUID(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsIPv4(t *testing.T) {
	fn := extvalidate.IsIPv4()
	cases := []struct {
		input any
		want  bool
	}{
		{"192.168.1.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"::1", false},
		{"not-an-ip", false},
	}
	for _, c := range cases {
		got, _ := fn([]any{c.input}, nil)
		if got != c.want {
			t.Errorf("IsIPv4(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsIPv6(t *testing.T) {
	fn := extvalidate.IsIPv6()
	cases := []struct {
		input any
		want  bool
	}{
		{"::1", true},
		{"2001:db8::1", true},
		{"192.168.1.1", false},
		{"not-an-ip", false},
	}
	for _, c := range cases {
		got, _ := fn([]any{c.input}, nil)
		if got != c.want {
			t.Errorf("IsIPv6(%v): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	fn := extvalidate.IsAlpha()
	cases := []struct {
		input any
		want  bool
	}{
		{"hello", true},
		{"Hello", true},
		{"hello123", false},
		{"", false},
		{"hello world", false},
	}
	for _, c := range cases {
		got, _ := fn([]any{c.input}, nil)
		if got != c.want {
			t.Errorf("IsAlpha(%q): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsAlphanumeric(t *testing.T) {
	fn := extvalidate.IsAlphanumeric()
	cases := []struct {
		input any
		want  bool
	}{
		{"hello123", true},
		{"Hello", true},
		{"hello world", false},
		{"hello!", false},
		{"", false},
	}
	for _, c := range cases {
		got, _ := fn([]any{c.input}, nil)
		if got != c.want {
			t.Errorf("IsAlphanumeric(%q): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestIsNumericStr(t *testing.T) {
	fn := extvalidate.IsNumericStr()
	cases := []struct {
		input any
		want  bool
	}{
		{"42", true},
		{"3.14", true},
		{"-1.5e10", true},
		{"abc", false},
		{"", false},
	}
	for _, c := range cases {
		got, _ := fn([]any{c.input}, nil)
		if got != c.want {
			t.Errorf("IsNumericStr(%q): got %v, want %v", c.input, got, c.want)
		}
	}
}

func TestMatchesRegex(t *testing.T) {
	fn := extvalidate.MatchesRegex()
	got, err := fn([]any{"hello123", `^\w+$`}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("expected true, got %v", got)
	}

	got, err = fn([]any{"hello world", `^\w+$`}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != false {
		t.Errorf("expected false, got %v", got)
	}

	_, err = fn([]any{"s", `[invalid`}, nil)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestInSet(t *testing.T) {
	fn := extvalidate.InSet()
	got, _ := fn([]any{"b", []any{"a", "b", "c"}}, nil)
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
	got, _ = fn([]any{"z", []any{"a", "b", "c"}}, nil)
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

func TestMinLen(t *testing.T) {
	fn := extvalidate.MinLen()
	got, _ := fn([]any{"hello", float64(3)}, nil)
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
	got, _ = fn([]any{"hi", float64(5)}, nil)
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

func TestMaxLen(t *testing.T) {
	fn := extvalidate.MaxLen()
	got, _ := fn([]any{"hi", float64(5)}, nil)
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
	got, _ = fn([]any{"toolong", float64(3)}, nil)
	if got != false {
		t.Errorf("expected false, got %v", got)
	}
}

func TestMinItems(t *testing.T) {
	fn := extvalidate.MinItems()
	got, _ := fn([]any{[]any{1, 2, 3}, float64(2)}, nil)
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

func TestMaxItems(t *testing.T) {
	fn := extvalidate.MaxItems()
	got, _ := fn([]any{[]any{1, 2}, float64(5)}, nil)
	if got != true {
		t.Errorf("expected true, got %v", got)
	}
}

func TestAll(t *testing.T) {
	all := extvalidate.All()
	expected := []string{
		"isEmail", "isURL", "isUUID", "isIPv4", "isIPv6",
		"isAlpha", "isAlphanumeric", "isNumericStr", "matchesRegex",
		"inSet", "minLen", "maxLen", "minItems", "maxItems",
	}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All() missing function: %q", name)
		}
	}
}

func TestIsEmailNoArgs(t *testing.T) {
	f := extvalidate.IsEmail()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isEmail(no args): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isEmail(non-string): expected false")
	}
}

func TestIsURLEdgeCases(t *testing.T) {
	f := extvalidate.IsURL()
	// no args
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isURL(no args): expected false")
	}
	// non-string
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isURL(non-string): expected false")
	}
	// non-http scheme
	got, _ = f([]any{"file:///etc/passwd"}, nil)
	if got.(bool) {
		t.Error("isURL(file://): expected false")
	}
}

func TestIsUUIDNoArgs(t *testing.T) {
	f := extvalidate.IsUUID()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isUUID(no args): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isUUID(non-string): expected false")
	}
}

func TestIsIPv4NoArgs(t *testing.T) {
	f := extvalidate.IsIPv4()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isIPv4(no args): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isIPv4(non-string): expected false")
	}
}

func TestIsIPv6NoArgs(t *testing.T) {
	f := extvalidate.IsIPv6()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isIPv6(no args): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isIPv6(non-string): expected false")
	}
}

func TestIsAlphaEdges(t *testing.T) {
	f := extvalidate.IsAlpha()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isAlpha(no args): expected false")
	}
	got, _ = f([]any{""}, nil)
	if got.(bool) {
		t.Error("isAlpha(empty): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isAlpha(non-string): expected false")
	}
}

func TestIsAlphanumericEdges(t *testing.T) {
	f := extvalidate.IsAlphanumeric()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isAlphanumeric(no args): expected false")
	}
	got, _ = f([]any{""}, nil)
	if got.(bool) {
		t.Error("isAlphanumeric(empty): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isAlphanumeric(non-string): expected false")
	}
}

func TestIsNumericStrEdges(t *testing.T) {
	f := extvalidate.IsNumericStr()
	got, _ := f([]any{}, nil)
	if got.(bool) {
		t.Error("isNumericStr(no args): expected false")
	}
	got, _ = f([]any{42}, nil)
	if got.(bool) {
		t.Error("isNumericStr(non-string): expected false")
	}
}

func TestMatchesRegexErrors(t *testing.T) {
	f := extvalidate.MatchesRegex()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("matchesRegex: expected error for 0 args")
	}
	if _, err := f([]any{42, `\d+`}, nil); err == nil {
		t.Error("matchesRegex: expected error for non-string str")
	}
	if _, err := f([]any{"hello", 42}, nil); err == nil {
		t.Error("matchesRegex: expected error for non-string pattern")
	}
	if _, err := f([]any{"hello", `[invalid`}, nil); err == nil {
		t.Error("matchesRegex: expected error for invalid regex")
	}
}

func TestInSetErrors(t *testing.T) {
	f := extvalidate.InSet()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("inSet: expected error for 0 args")
	}
	if _, err := f([]any{"x", "not-array"}, nil); err == nil {
		t.Error("inSet: expected error for non-array set")
	}
}

func TestMinLenErrors(t *testing.T) {
	f := extvalidate.MinLen()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("minLen: expected error for 0 args")
	}
	if _, err := f([]any{42, 3.0}, nil); err == nil {
		t.Error("minLen: expected error for non-string")
	}
	if _, err := f([]any{"hello", "bad"}, nil); err == nil {
		t.Error("minLen: expected error for non-numeric n")
	}
}

func TestMaxLenErrors(t *testing.T) {
	f := extvalidate.MaxLen()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("maxLen: expected error for 0 args")
	}
	if _, err := f([]any{42, 3.0}, nil); err == nil {
		t.Error("maxLen: expected error for non-string")
	}
	if _, err := f([]any{"hello", "bad"}, nil); err == nil {
		t.Error("maxLen: expected error for non-numeric n")
	}
}

func TestMinItemsFalse(t *testing.T) {
	f := extvalidate.MinItems()
	got, _ := f([]any{[]any{1.0}, float64(5)}, nil)
	if got.(bool) {
		t.Error("minItems: expected false when len < n")
	}
}

func TestMinItemsErrors(t *testing.T) {
	f := extvalidate.MinItems()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("minItems: expected error for 0 args")
	}
	if _, err := f([]any{"not-array", 2.0}, nil); err == nil {
		t.Error("minItems: expected error for non-array")
	}
	if _, err := f([]any{[]any{1.0}, "bad"}, nil); err == nil {
		t.Error("minItems: expected error for non-numeric n")
	}
}

func TestMaxItemsFalse(t *testing.T) {
	f := extvalidate.MaxItems()
	got, _ := f([]any{[]any{1.0, 2.0, 3.0}, float64(2)}, nil)
	if got.(bool) {
		t.Error("maxItems: expected false when len > n")
	}
}

func TestMaxItemsErrors(t *testing.T) {
	f := extvalidate.MaxItems()
	if _, err := f([]any{}, nil); err == nil {
		t.Error("maxItems: expected error for 0 args")
	}
	if _, err := f([]any{"not-array", 2.0}, nil); err == nil {
		t.Error("maxItems: expected error for non-array")
	}
	if _, err := f([]any{[]any{1.0}, "bad"}, nil); err == nil {
		t.Error("maxItems: expected error for non-numeric n")
	}
}
