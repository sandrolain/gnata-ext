// Package extregex provides advanced regular expression utilities for gnata.
//
// All patterns use Go RE2 syntax (no look-ahead/look-behind).
// Compiled patterns are cached internally to avoid repeated compilation overhead.
//
// Functions:
//
//   - $regexAll(str, pattern)               – all non-overlapping matches
//   - $regexNamedGroups(str, pattern)       – named capture groups as object
//   - $regexSplit(str, pattern)             – split by regex delimiter
//   - $regexReplaceAll(str, pattern, repl)  – replace all matches
//   - $regexCount(str, pattern)             – number of matches
//   - $regexTest(str, pattern)              – true if pattern matches
//   - $regexExtract(str, pattern [, group]) – first match or specific group
package extregex

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/recolabs/gnata"
)

// cache stores compiled regexps to avoid re-compilation on every call.
var (
	cacheMu sync.RWMutex
	cache   = map[string]*regexp.Regexp{}
)

func compile(pattern string) (*regexp.Regexp, error) {
	cacheMu.RLock()
	re, ok := cache[pattern]
	cacheMu.RUnlock()
	if ok {
		return re, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	cacheMu.Lock()
	cache[pattern] = re
	cacheMu.Unlock()
	return re, nil
}

// All returns a map of all extregex functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"regexAll":        RegexAll(),
		"regexNamedGroups": RegexNamedGroups(),
		"regexSplit":      RegexSplit(),
		"regexReplaceAll": RegexReplaceAll(),
		"regexCount":      RegexCount(),
		"regexTest":       RegexTest(),
		"regexExtract":    RegexExtract(),
	}
}

// RegexAll returns the CustomFunc for $regexAll(str, pattern).
// Returns an array of all non-overlapping matches (full match string).
func RegexAll() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		str, pattern, err := twoStrArgs("$regexAll", args)
		if err != nil {
			return nil, err
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexAll: invalid pattern: %w", err)
		}
		matches := re.FindAllString(str, -1)
		if matches == nil {
			return []any{}, nil
		}
		out := make([]any, len(matches))
		for i, m := range matches {
			out[i] = m
		}
		return out, nil
	}
}

// RegexNamedGroups returns the CustomFunc for $regexNamedGroups(str, pattern).
// Returns an object mapping named capture group names to their matched values.
// Groups that do not participate in the match are omitted.
func RegexNamedGroups() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		str, pattern, err := twoStrArgs("$regexNamedGroups", args)
		if err != nil {
			return nil, err
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexNamedGroups: invalid pattern: %w", err)
		}
		names := re.SubexpNames()
		match := re.FindStringSubmatch(str)
		if match == nil {
			return map[string]any{}, nil
		}
		obj := make(map[string]any)
		for i, name := range names {
			if name != "" && i < len(match) && match[i] != "" {
				obj[name] = match[i]
			}
		}
		return obj, nil
	}
}

// RegexSplit returns the CustomFunc for $regexSplit(str, pattern).
// Splits str at each occurrence of pattern.
func RegexSplit() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		str, pattern, err := twoStrArgs("$regexSplit", args)
		if err != nil {
			return nil, err
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexSplit: invalid pattern: %w", err)
		}
		parts := re.Split(str, -1)
		out := make([]any, len(parts))
		for i, p := range parts {
			out[i] = p
		}
		return out, nil
	}
}

// RegexReplaceAll returns the CustomFunc for $regexReplaceAll(str, pattern, repl).
// Replaces all matches of pattern in str with repl.
// Supports $1, $2, … and ${name} back-references in repl.
func RegexReplaceAll() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$regexReplaceAll: requires 3 arguments (str, pattern, repl)")
		}
		str, ok1 := args[0].(string)
		pattern, ok2 := args[1].(string)
		repl, ok3 := args[2].(string)
		if !ok1 || !ok2 || !ok3 {
			return nil, fmt.Errorf("$regexReplaceAll: all arguments must be strings")
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexReplaceAll: invalid pattern: %w", err)
		}
		return re.ReplaceAllString(str, repl), nil
	}
}

// RegexCount returns the CustomFunc for $regexCount(str, pattern).
// Returns the number of non-overlapping matches.
func RegexCount() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		str, pattern, err := twoStrArgs("$regexCount", args)
		if err != nil {
			return nil, err
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexCount: invalid pattern: %w", err)
		}
		return float64(len(re.FindAllString(str, -1))), nil
	}
}

// RegexTest returns the CustomFunc for $regexTest(str, pattern).
// Returns true if the pattern matches anywhere in str.
func RegexTest() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		str, pattern, err := twoStrArgs("$regexTest", args)
		if err != nil {
			return nil, err
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexTest: invalid pattern: %w", err)
		}
		return re.MatchString(str), nil
	}
}

// RegexExtract returns the CustomFunc for $regexExtract(str, pattern [, group]).
// Returns the first match string, or the contents of the specified capture group (0-based index).
// Returns nil if there is no match.
func RegexExtract() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$regexExtract: requires at least 2 arguments (str, pattern)")
		}
		str, ok1 := args[0].(string)
		pattern, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$regexExtract: first two arguments must be strings")
		}
		re, err := compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("$regexExtract: invalid pattern: %w", err)
		}

		groupIdx := 0
		if len(args) >= 3 {
			switch n := args[2].(type) {
			case float64:
				groupIdx = int(n)
			case int:
				groupIdx = n
			case int64:
				groupIdx = int(n)
			default:
				return nil, fmt.Errorf("$regexExtract: group argument must be a number")
			}
		}

		match := re.FindStringSubmatch(str)
		if match == nil {
			return nil, nil
		}
		if groupIdx >= len(match) {
			return nil, fmt.Errorf("$regexExtract: group index %d out of range (match has %d groups)", groupIdx, len(match)-1)
		}
		return match[groupIdx], nil
	}
}

// --- helpers ---

func twoStrArgs(name string, args []any) (string, string, error) {
	if len(args) < 2 {
		return "", "", fmt.Errorf("%s: requires 2 arguments (str, pattern)", name)
	}
	s, ok1 := args[0].(string)
	p, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("%s: arguments must be strings", name)
	}
	return s, p, nil
}
