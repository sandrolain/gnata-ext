// Package exttext provides text analysis and manipulation functions for gnata.
//
// Functions:
//
//   - $wordCount(s)             – number of words
//   - $charCount(s)             – number of characters (Unicode code points)
//   - $sentenceCount(s)         – approximate number of sentences
//   - $readingTime(s)           – estimated reading time in seconds (200 wpm)
//   - $wordFrequency(s)         – {word: count} map
//   - $ngrams(s, n)             – array of n-gram strings
//   - $excerpt(s, maxLen)       – truncated string with trailing "…"
//   - $stripHTML(s)             – removes HTML tags from a string
//   - $normalizeWhitespace(s)   – collapses consecutive whitespace to single space
package exttext

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all exttext functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"wordCount":           WordCount(),
		"charCount":           CharCount(),
		"sentenceCount":       SentenceCount(),
		"readingTime":         ReadingTime(),
		"wordFrequency":       WordFrequency(),
		"ngrams":              Ngrams(),
		"excerpt":             Excerpt(),
		"stripHTML":           StripHTML(),
		"normalizeWhitespace": NormalizeWhitespace(),
	}
}

var (
	reWords      = regexp.MustCompile(`\b\w+\b`)
	reSentences  = regexp.MustCompile(`[.!?]+`)
	reHTMLTag    = regexp.MustCompile(`<[^>]*>`)
	reWhitespace = regexp.MustCompile(`\s+`)
)

// WordCount returns the CustomFunc for $wordCount(s).
func WordCount() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$wordCount")
		if err != nil {
			return nil, err
		}
		words := reWords.FindAllString(s, -1)
		return float64(len(words)), nil
	}
}

// CharCount returns the CustomFunc for $charCount(s).
// Counts Unicode code points (runes), not bytes.
func CharCount() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$charCount")
		if err != nil {
			return nil, err
		}
		return float64(len([]rune(s))), nil
	}
}

// SentenceCount returns the CustomFunc for $sentenceCount(s).
// Approximates sentences by counting sentence-ending punctuation groups.
func SentenceCount() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$sentenceCount")
		if err != nil {
			return nil, err
		}
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return float64(0), nil
		}
		endings := reSentences.FindAllString(trimmed, -1)
		n := len(endings)
		// If the text does not end with punctuation, the last sentence counts too.
		last := rune(trimmed[len(trimmed)-1])
		if last != '.' && last != '!' && last != '?' {
			n++
		}
		return float64(n), nil
	}
}

// ReadingTime returns the CustomFunc for $readingTime(s).
// Returns estimated reading time in seconds at 200 words per minute.
func ReadingTime() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$readingTime")
		if err != nil {
			return nil, err
		}
		words := reWords.FindAllString(s, -1)
		minutes := float64(len(words)) / 200.0
		seconds := math.Ceil(minutes * 60)
		return seconds, nil
	}
}

// WordFrequency returns the CustomFunc for $wordFrequency(s).
// Returns a map[string]any of {lowercased-word: count (float64)}.
func WordFrequency() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$wordFrequency")
		if err != nil {
			return nil, err
		}
		words := reWords.FindAllString(strings.ToLower(s), -1)
		freq := make(map[string]any, len(words))
		for _, w := range words {
			if prev, ok := freq[w].(float64); ok {
				freq[w] = prev + 1
			} else {
				freq[w] = float64(1)
			}
		}
		return freq, nil
	}
}

// Ngrams returns the CustomFunc for $ngrams(s, n).
// Returns an array of n-gram strings (space-separated words).
func Ngrams() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$ngrams: requires 2 arguments (s, n)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$ngrams: first argument must be a string")
		}
		n, ok := extutil.ToInt(args[1])
		if !ok || n < 1 {
			return nil, fmt.Errorf("$ngrams: n must be a positive integer")
		}
		words := reWords.FindAllString(s, -1)
		var out []any
		for i := 0; i <= len(words)-n; i++ {
			out = append(out, strings.Join(words[i:i+n], " "))
		}
		if out == nil {
			out = []any{}
		}
		return out, nil
	}
}

// Excerpt returns the CustomFunc for $excerpt(s, maxLen).
// Returns the string truncated to maxLen runes. If truncated, appends "…".
func Excerpt() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$excerpt: requires 2 arguments (s, maxLen)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$excerpt: first argument must be a string")
		}
		maxLen, ok := extutil.ToInt(args[1])
		if !ok || maxLen < 0 {
			return nil, fmt.Errorf("$excerpt: maxLen must be a non-negative integer")
		}
		runes := []rune(s)
		if len(runes) <= maxLen {
			return s, nil
		}
		return string(runes[:maxLen]) + "…", nil
	}
}

// StripHTML returns the CustomFunc for $stripHTML(s).
// Removes all HTML tags from the string.
func StripHTML() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$stripHTML")
		if err != nil {
			return nil, err
		}
		return reHTMLTag.ReplaceAllString(s, ""), nil
	}
}

// NormalizeWhitespace returns the CustomFunc for $normalizeWhitespace(s).
// Collapses all consecutive whitespace (including tabs/newlines) to a single space
// and trims leading/trailing whitespace.
func NormalizeWhitespace() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		s, err := requireString(args, "$normalizeWhitespace")
		if err != nil {
			return nil, err
		}
		// Replace each whitespace-only rune group with a space.
		result := reWhitespace.ReplaceAllStringFunc(s, func(m string) string {
			// Keep newline if the entire match consists of newline characters.
			for _, r := range m {
				if !unicode.IsSpace(r) {
					return m
				}
			}
			return " "
		})
		return strings.TrimSpace(result), nil
	}
}

// --- helpers ---

func requireString(args []any, name string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("%s: requires 1 argument", name)
	}
	s, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("%s: argument must be a string", name)
	}
	return s, nil
}
