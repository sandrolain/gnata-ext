// Package extstring provides extended string functions for gnata beyond the
// official JSONata 2.x specification.
//
// Functions
//
//   - $startsWith(str, prefix)   – true if str starts with prefix
//   - $endsWith(str, suffix)     – true if str ends with suffix
//   - $indexOf(str, search [, start]) – first index of search (-1 if not found)
//   - $lastIndexOf(str, search)  – last index of search (-1 if not found)
//   - $capitalize(str)           – uppercase first char, lowercase rest
//   - $titleCase(str)            – uppercase first char of each word
//   - $camelCase(str)            – convert to camelCase
//   - $snakeCase(str)            – convert to snake_case
//   - $kebabCase(str)            – convert to kebab-case
//   - $repeat(str, n)            – repeat str n times
//   - $words(str)                – split into array of words
//   - $template(str, bindings)   – replace {{key}} placeholders
package extstring

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

var (
	splitWordsRe = regexp.MustCompile(`[_\-\s]+|([a-z])([A-Z])`)
	templateRe   = regexp.MustCompile(`\{\{(\w+)\}\}`)
)

// All returns a map of all extended string functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"startsWith":  StartsWith(),
		"endsWith":    EndsWith(),
		"indexOf":     IndexOf(),
		"lastIndexOf": LastIndexOf(),
		"capitalize":  Capitalize(),
		"titleCase":   TitleCase(),
		"camelCase":   CamelCase(),
		"snakeCase":   SnakeCase(),
		"kebabCase":   KebabCase(),
		"repeat":      Repeat(),
		"words":       Words(),
		"template":    Template(),
	}
}

// StartsWith returns the CustomFunc for $startsWith(str, prefix).
func StartsWith() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$startsWith: requires 2 arguments")
		}
		str, ok1 := args[0].(string)
		prefix, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$startsWith: arguments must be strings")
		}
		return strings.HasPrefix(str, prefix), nil
	}
}

// EndsWith returns the CustomFunc for $endsWith(str, suffix).
func EndsWith() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$endsWith: requires 2 arguments")
		}
		str, ok1 := args[0].(string)
		suffix, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$endsWith: arguments must be strings")
		}
		return strings.HasSuffix(str, suffix), nil
	}
}

// IndexOf returns the CustomFunc for $indexOf(str, search [, start]).
// Returns -1 when not found.
func IndexOf() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$indexOf: requires at least 2 arguments")
		}
		str, ok1 := args[0].(string)
		search, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$indexOf: first two arguments must be strings")
		}
		start := 0
		if len(args) >= 3 && args[2] != nil {
			n, ok := extutil.ToInt(args[2])
			if !ok {
				return nil, fmt.Errorf("$indexOf: start must be a number")
			}
			if n < 0 {
				n = 0
			}
			start = n
		}
		if start >= len(str) {
			return float64(-1), nil
		}
		idx := strings.Index(str[start:], search)
		if idx == -1 {
			return float64(-1), nil
		}
		return float64(idx + start), nil
	}
}

// LastIndexOf returns the CustomFunc for $lastIndexOf(str, search).
func LastIndexOf() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$lastIndexOf: requires 2 arguments")
		}
		str, ok1 := args[0].(string)
		search, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$lastIndexOf: arguments must be strings")
		}
		return float64(strings.LastIndex(str, search)), nil
	}
}

// Capitalize returns the CustomFunc for $capitalize(str).
func Capitalize() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$capitalize: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$capitalize: argument must be a string")
		}
		if str == "" {
			return str, nil
		}
		runes := []rune(str)
		runes[0] = unicode.ToUpper(runes[0])
		for i := 1; i < len(runes); i++ {
			runes[i] = unicode.ToLower(runes[i])
		}
		return string(runes), nil
	}
}

// TitleCase returns the CustomFunc for $titleCase(str).
func TitleCase() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$titleCase: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$titleCase: argument must be a string")
		}
		return strings.Title(strings.ToLower(str)), nil //nolint:staticcheck
	}
}

// CamelCase returns the CustomFunc for $camelCase(str).
func CamelCase() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$camelCase: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$camelCase: argument must be a string")
		}
		words := splitIntoWords(str)
		if len(words) == 0 {
			return "", nil
		}
		var b strings.Builder
		b.WriteString(strings.ToLower(words[0]))
		for _, w := range words[1:] {
			if w == "" {
				continue
			}
			runes := []rune(strings.ToLower(w))
			runes[0] = unicode.ToUpper(runes[0])
			b.WriteString(string(runes))
		}
		return b.String(), nil
	}
}

// SnakeCase returns the CustomFunc for $snakeCase(str).
func SnakeCase() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$snakeCase: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$snakeCase: argument must be a string")
		}
		words := splitIntoWords(str)
		for i, w := range words {
			words[i] = strings.ToLower(w)
		}
		return strings.Join(words, "_"), nil
	}
}

// KebabCase returns the CustomFunc for $kebabCase(str).
func KebabCase() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$kebabCase: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$kebabCase: argument must be a string")
		}
		words := splitIntoWords(str)
		for i, w := range words {
			words[i] = strings.ToLower(w)
		}
		return strings.Join(words, "-"), nil
	}
}

// Repeat returns the CustomFunc for $repeat(str, n).
func Repeat() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$repeat: requires 2 arguments")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$repeat: first argument must be a string")
		}
		n, ok := extutil.ToInt(args[1])
		if !ok || n < 0 {
			return nil, fmt.Errorf("$repeat: second argument must be a non-negative integer")
		}
		return strings.Repeat(str, n), nil
	}
}

// Words returns the CustomFunc for $words(str).
func Words() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$words: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$words: argument must be a string")
		}
		parts := strings.Fields(str)
		if len(parts) == 0 {
			return nil, nil
		}
		result := make([]any, len(parts))
		for i, p := range parts {
			result[i] = p
		}
		return result, nil
	}
}

// Template returns the CustomFunc for $template(str, bindings).
// Replaces {{key}} placeholders with values from the bindings object.
func Template() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$template: requires 2 arguments")
		}
		tmpl, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$template: first argument must be a string")
		}
		bindings, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$template: second argument must be an object")
		}
		result := templateRe.ReplaceAllStringFunc(tmpl, func(match string) string {
			key := match[2 : len(match)-2]
			if val, exists := bindings[key]; exists {
				return fmt.Sprint(val)
			}
			return match
		})
		return result, nil
	}
}

func splitIntoWords(str string) []string {
	expanded := splitWordsRe.ReplaceAllStringFunc(str, func(s string) string {
		if len(s) == 2 && s[0] >= 'a' && s[0] <= 'z' {
			return string(s[0]) + " " + string(s[1])
		}
		return " "
	})
	return strings.Fields(expanded)
}
