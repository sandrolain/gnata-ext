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
	cases := []struct{ input any; want bool }{
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
	cases := []struct{ input any; want bool }{
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
	cases := []struct{ input any; want bool }{
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
