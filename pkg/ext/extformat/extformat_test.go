package extformat_test

import (
	"strings"
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extformat"
)

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

func TestParseCSV(t *testing.T) {
	f := extformat.ParseCSV()
	csvText := "name,age\nAlice,30\nBob,25\n"
	got, err := invoke(f, csvText)
	if err != nil {
		t.Fatalf("csv: %v", err)
	}
	rows := got.([]any)
	if len(rows) != 2 {
		t.Fatalf("csv: expected 2 rows, got %d", len(rows))
	}
	row0 := rows[0].(map[string]any)
	if row0["name"].(string) != "Alice" || row0["age"].(string) != "30" {
		t.Errorf("csv row0: got %v", row0)
	}
	row1 := rows[1].(map[string]any)
	if row1["name"].(string) != "Bob" {
		t.Errorf("csv row1: got %v", row1)
	}
}

func TestParseCSVEmpty(t *testing.T) {
	f := extformat.ParseCSV()
	got, err := invoke(f, "")
	if err != nil {
		t.Fatalf("csv empty: %v", err)
	}
	rows := got.([]any)
	if len(rows) != 0 {
		t.Errorf("csv empty: got %v", rows)
	}
}

func TestToCSV(t *testing.T) {
	f := extformat.ToCSV()
	rows := []any{
		map[string]any{"name": "Alice", "age": "30"},
		map[string]any{"name": "Bob", "age": "25"},
	}
	got, err := invoke(f, rows)
	if err != nil {
		t.Fatalf("toCSV: %v", err)
	}
	s := got.(string)
	if !strings.Contains(s, "Alice") || !strings.Contains(s, "Bob") {
		t.Errorf("toCSV: got %q", s)
	}
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) != 3 { // header + 2 rows
		t.Errorf("toCSV: expected 3 lines, got %d: %q", len(lines), s)
	}
}

func TestToCSVEmpty(t *testing.T) {
	f := extformat.ToCSV()
	got, err := invoke(f, []any{})
	if err != nil {
		t.Fatalf("toCSV empty: %v", err)
	}
	if got.(string) != "" {
		t.Errorf("toCSV empty: got %q", got)
	}
}

func TestTemplate(t *testing.T) {
	f := extformat.Template()
	tmpl := "Hello, {{name}}! You are {{age}} years old."
	vars := map[string]any{"name": "Alice", "age": 30.0}
	got, err := invoke(f, tmpl, vars)
	if err != nil {
		t.Fatalf("template: %v", err)
	}
	s := got.(string)
	if !strings.Contains(s, "Alice") || !strings.Contains(s, "30") {
		t.Errorf("template: got %q", s)
	}
}

func TestTemplateUnknownKey(t *testing.T) {
	f := extformat.Template()
	got, err := invoke(f, "Hello, {{unknown}}!", map[string]any{})
	if err != nil {
		t.Fatalf("template unknown: %v", err)
	}
	// Unknown placeholders should be left as-is
	if got.(string) != "Hello, {{unknown}}!" {
		t.Errorf("template unknown: got %q", got)
	}
}

func TestAll(t *testing.T) {
	all := extformat.All()
	expected := []string{"csv", "toCSV", "template"}
	for _, name := range expected {
		if _, ok := all[name]; !ok {
			t.Errorf("All(): missing function %q", name)
		}
	}
}
