// Package extstring provides extended string functions for gnata beyond the
// official JSONata 2.x specification.
//
// Functions
//
//   - $startsWith(str, prefix)            – true if str starts with prefix
//   - $endsWith(str, suffix)              – true if str ends with suffix
//   - $indexOf(str, search [, start])     – first index of search (-1 if not found)
//   - $lastIndexOf(str, search)           – last index of search (-1 if not found)
//   - $capitalize(str)                    – uppercase first char, lowercase rest
//   - $titleCase(str)                     – uppercase first char of each word
//   - $camelCase(str)                     – convert to camelCase
//   - $snakeCase(str)                     – convert to snake_case
//   - $kebabCase(str)                     – convert to kebab-case
//   - $repeat(str, n)                     – repeat str n times
//   - $words(str)                         – split into array of words
//   - $template(str, bindings)            – replace {{key}} placeholders
//   - $padStart(str, len [, fill])        – left-pad to length
//   - $padEnd(str, len [, fill])          – right-pad to length
//   - $truncate(str, len [, suffix])      – truncate to len chars with optional suffix
//   - $slugify(str)                       – URL-friendly slug
//   - $countOccurrences(str, sub)         – count non-overlapping occurrences
//   - $initials(str [, sep])              – initials from words, joined by sep
//   - $escapeHTML(str)                    – escape HTML special characters
//   - $unescapeHTML(str)                  – unescape HTML entities
//   - $reverseWords(str)                  – reverse the order of words
//   - $levenshtein(a, b)                  – edit distance between two strings
//   - $longestCommonPrefix(strs)          – longest common prefix of an array
package extstring

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

var (
	splitWordsRe    = regexp.MustCompile(`[_\-\s]+|([a-z])([A-Z])`)
	templateRe      = regexp.MustCompile(`\{\{(\w+)\}\}`)
	slugifyNonAlNum = regexp.MustCompile(`[^a-z0-9]+`)
)

// All returns a map of all extended string functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"startsWith":          StartsWith(),
		"endsWith":            EndsWith(),
		"indexOf":             IndexOf(),
		"lastIndexOf":         LastIndexOf(),
		"capitalize":          Capitalize(),
		"titleCase":           TitleCase(),
		"camelCase":           CamelCase(),
		"snakeCase":           SnakeCase(),
		"kebabCase":           KebabCase(),
		"repeat":              Repeat(),
		"words":               Words(),
		"template":            Template(),
		"padStart":            PadStart(),
		"padEnd":              PadEnd(),
		"truncate":            Truncate(),
		"slugify":             Slugify(),
		"countOccurrences":    CountOccurrences(),
		"initials":            Initials(),
		"escapeHTML":          EscapeHTML(),
		"unescapeHTML":        UnescapeHTML(),
		"reverseWords":        ReverseWords(),
		"levenshtein":         Levenshtein(),
		"longestCommonPrefix": LongestCommonPrefix(),
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

// PadStart returns the CustomFunc for $padStart(str, len [, fill]).
// Left-pads str to length using fill (default " ").
func PadStart() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$padStart: requires at least 2 arguments")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$padStart: first argument must be a string")
		}
		padLen, ok := extutil.ToInt(args[1])
		if !ok || padLen < 0 {
			return nil, fmt.Errorf("$padStart: len must be a non-negative integer")
		}
		fill := " "
		if len(args) >= 3 && args[2] != nil {
			f, ok := args[2].(string)
			if !ok || f == "" {
				return nil, fmt.Errorf("$padStart: fill must be a non-empty string")
			}
			fill = f
		}
		runes := []rune(str)
		needed := padLen - len(runes)
		if needed <= 0 {
			return str, nil
		}
		fillRunes := []rune(fill)
		var pad []rune
		for len(pad) < needed {
			pad = append(pad, fillRunes...)
		}
		return string(pad[:needed]) + str, nil
	}
}

// PadEnd returns the CustomFunc for $padEnd(str, len [, fill]).
// Right-pads str to length using fill (default " ").
func PadEnd() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$padEnd: requires at least 2 arguments")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$padEnd: first argument must be a string")
		}
		padLen, ok := extutil.ToInt(args[1])
		if !ok || padLen < 0 {
			return nil, fmt.Errorf("$padEnd: len must be a non-negative integer")
		}
		fill := " "
		if len(args) >= 3 && args[2] != nil {
			f, ok := args[2].(string)
			if !ok || f == "" {
				return nil, fmt.Errorf("$padEnd: fill must be a non-empty string")
			}
			fill = f
		}
		runes := []rune(str)
		needed := padLen - len(runes)
		if needed <= 0 {
			return str, nil
		}
		fillRunes := []rune(fill)
		var pad []rune
		for len(pad) < needed {
			pad = append(pad, fillRunes...)
		}
		return str + string(pad[:needed]), nil
	}
}

// Truncate returns the CustomFunc for $truncate(str, len [, suffix]).
// Truncates str to at most len runes; if truncated, appends suffix (default "…").
func Truncate() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$truncate: requires at least 2 arguments")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$truncate: first argument must be a string")
		}
		truncLen, ok := extutil.ToInt(args[1])
		if !ok || truncLen < 0 {
			return nil, fmt.Errorf("$truncate: len must be a non-negative integer")
		}
		suffix := "…"
		if len(args) >= 3 && args[2] != nil {
			s, ok := args[2].(string)
			if !ok {
				return nil, fmt.Errorf("$truncate: suffix must be a string")
			}
			suffix = s
		}
		runes := []rune(str)
		if len(runes) <= truncLen {
			return str, nil
		}
		return string(runes[:truncLen]) + suffix, nil
	}
}

// Slugify returns the CustomFunc for $slugify(str).
// Converts str to a URL-friendly lowercase slug (letters, digits, hyphens).
func Slugify() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$slugify: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$slugify: argument must be a string")
		}
		slug := strings.ToLower(str)
		slug = slugifyNonAlNum.ReplaceAllString(slug, "-")
		slug = strings.Trim(slug, "-")
		return slug, nil
	}
}

// CountOccurrences returns the CustomFunc for $countOccurrences(str, sub).
// Returns the number of non-overlapping occurrences of sub in str.
func CountOccurrences() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$countOccurrences: requires 2 arguments")
		}
		str, ok1 := args[0].(string)
		sub, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$countOccurrences: arguments must be strings")
		}
		if sub == "" {
			return float64(0), nil
		}
		return float64(strings.Count(str, sub)), nil
	}
}

// Initials returns the CustomFunc for $initials(str [, sep]).
// Returns the initials of each word joined by sep (default "").
func Initials() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$initials: requires at least 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$initials: argument must be a string")
		}
		sep := ""
		if len(args) >= 2 && args[1] != nil {
			s, ok := args[1].(string)
			if !ok {
				return nil, fmt.Errorf("$initials: sep must be a string")
			}
			sep = s
		}
		words := strings.Fields(str)
		inits := make([]string, 0, len(words))
		for _, w := range words {
			runes := []rune(w)
			if len(runes) > 0 {
				inits = append(inits, string(unicode.ToUpper(runes[0])))
			}
		}
		return strings.Join(inits, sep), nil
	}
}

// EscapeHTML returns the CustomFunc for $escapeHTML(str).
// Escapes &, <, >, ", ' to HTML entities.
func EscapeHTML() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$escapeHTML: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$escapeHTML: argument must be a string")
		}
		return html.EscapeString(str), nil
	}
}

// UnescapeHTML returns the CustomFunc for $unescapeHTML(str).
// Unescapes HTML entities back to their characters.
func UnescapeHTML() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$unescapeHTML: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$unescapeHTML: argument must be a string")
		}
		return html.UnescapeString(str), nil
	}
}

// ReverseWords returns the CustomFunc for $reverseWords(str).
// Reverses the order of whitespace-separated words.
func ReverseWords() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$reverseWords: requires 1 argument")
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$reverseWords: argument must be a string")
		}
		words := strings.Fields(str)
		for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
			words[i], words[j] = words[j], words[i]
		}
		return strings.Join(words, " "), nil
	}
}

// Levenshtein returns the CustomFunc for $levenshtein(a, b).
// Returns the edit distance (insertions, deletions, substitutions) between a and b.
func Levenshtein() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$levenshtein: requires 2 arguments")
		}
		a, ok1 := args[0].(string)
		b, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$levenshtein: arguments must be strings")
		}
		ra, rb := []rune(a), []rune(b)
		la, lb := len(ra), len(rb)
		if la == 0 {
			return float64(lb), nil
		}
		if lb == 0 {
			return float64(la), nil
		}
		prev := make([]int, lb+1)
		curr := make([]int, lb+1)
		for j := 0; j <= lb; j++ {
			prev[j] = j
		}
		for i := 1; i <= la; i++ {
			curr[0] = i
			for j := 1; j <= lb; j++ {
				cost := 1
				if ra[i-1] == rb[j-1] {
					cost = 0
				}
				del := prev[j] + 1
				ins := curr[j-1] + 1
				sub := prev[j-1] + cost
				m := del
				if ins < m {
					m = ins
				}
				if sub < m {
					m = sub
				}
				curr[j] = m
			}
			prev, curr = curr, prev
		}
		return float64(prev[lb]), nil
	}
}

// LongestCommonPrefix returns the CustomFunc for $longestCommonPrefix(strs).
// Returns the longest string that is a prefix of all strings in the array.
func LongestCommonPrefix() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$longestCommonPrefix: requires 1 argument")
		}
		arr, ok := args[0].([]any)
		if !ok {
			return nil, fmt.Errorf("$longestCommonPrefix: argument must be an array")
		}
		if len(arr) == 0 {
			return "", nil
		}
		strings_ := make([]string, 0, len(arr))
		for _, v := range arr {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("$longestCommonPrefix: all elements must be strings")
			}
			strings_ = append(strings_, s)
		}
		prefix := []rune(strings_[0])
		for _, s := range strings_[1:] {
			sr := []rune(s)
			max := len(prefix)
			if len(sr) < max {
				max = len(sr)
			}
			prefix = prefix[:max]
			for i := range prefix {
				if prefix[i] != sr[i] {
					prefix = prefix[:i]
					break
				}
			}
		}
		return string(prefix), nil
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
