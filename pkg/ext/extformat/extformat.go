// Package extformat provides CSV parsing/serialisation and template-rendering
// functions for gnata.
//
// Functions
//
//   - $csv(text)            – parse CSV text into an array of objects (first row = header)
//   - $toCSV(array)        – serialise an array of objects to CSV text
//   - $template(str, vars) – replace {{key}} placeholders with values from vars
package extformat

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

var templateRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

// All returns a map of all extended format functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"csv":      ParseCSV(),
		"toCSV":    ToCSV(),
		"template": Template(),
	}
}

// ParseCSV returns the CustomFunc for $csv(text).
// The first row of the CSV is treated as the header. Returns []any of map[string]any.
func ParseCSV() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$csv: requires 1 argument")
		}
		text, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$csv: argument must be a string")
		}
		r := csv.NewReader(strings.NewReader(text))
		records, err := r.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("$csv: %w", err)
		}
		if len(records) == 0 {
			return []any{}, nil
		}
		headers := records[0]
		result := make([]any, 0, len(records)-1)
		for _, row := range records[1:] {
			obj := make(map[string]any, len(headers))
			for i, h := range headers {
				if i < len(row) {
					obj[h] = row[i]
				} else {
					obj[h] = ""
				}
			}
			result = append(result, obj)
		}
		return result, nil
	}
}

// ToCSV returns the CustomFunc for $toCSV(array).
// Accepts an array of objects. The header row is derived from the keys of the
// first object. Missing keys in subsequent rows are written as empty strings.
func ToCSV() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$toCSV: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$toCSV: %w", err)
		}
		if len(arr) == 0 {
			return "", nil
		}
		firstObj, err := extutil.ToObject(arr[0])
		if err != nil {
			return nil, fmt.Errorf("$toCSV: first element: %w", err)
		}
		// Collect headers from first object.
		headers := make([]string, 0, len(firstObj))
		for k := range firstObj {
			headers = append(headers, k)
		}
		var buf bytes.Buffer
		w := csv.NewWriter(&buf)
		if err := w.Write(headers); err != nil {
			return nil, fmt.Errorf("$toCSV: %w", err)
		}
		for i, item := range arr {
			obj, err := extutil.ToObject(item)
			if err != nil {
				return nil, fmt.Errorf("$toCSV: row %d: %w", i, err)
			}
			row := make([]string, len(headers))
			for j, h := range headers {
				if v, ok := obj[h]; ok {
					row[j] = fmt.Sprintf("%v", v)
				}
			}
			if err := w.Write(row); err != nil {
				return nil, fmt.Errorf("$toCSV: %w", err)
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return nil, fmt.Errorf("$toCSV: %w", err)
		}
		return buf.String(), nil
	}
}

// Template returns the CustomFunc for $template(str, vars).
// Replaces {{key}} placeholders in str with corresponding values from vars.
func Template() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$template: requires 2 arguments (template, vars)")
		}
		tmpl, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$template: first argument must be a string")
		}
		vars, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$template: %w", err)
		}
		result := templateRe.ReplaceAllStringFunc(tmpl, func(match string) string {
			key := match[2 : len(match)-2] // strip {{ and }}
			if v, ok := vars[key]; ok {
				return fmt.Sprintf("%v", v)
			}
			return match
		})
		return result, nil
	}
}
