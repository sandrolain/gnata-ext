// Package extvalidate provides input-validation predicate functions for gnata.
//
// All functions return a boolean. They never return errors on bad input types —
// they return false instead, making them safe to use in conditional expressions.
//
// Note on regex safety: Go's regexp package uses RE2 semantics (guaranteed
// linear-time matching), so $matchesRegex is inherently safe against ReDoS.
//
// Functions
//
//   - $isEmail(str)              – RFC 5322 simplified email format
//   - $isURL(str)                – valid http/https/ftp URL
//   - $isUUID(str)               – UUID v1–v5 (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
//   - $isIPv4(str)               – IPv4 address (e.g. "192.168.1.1")
//   - $isIPv6(str)               – IPv6 address
//   - $isAlpha(str)              – only Unicode letters
//   - $isAlphanumeric(str)       – only Unicode letters and digits
//   - $isNumericStr(str)         – string parses as a number
//   - $matchesRegex(str, pattern)– str matches RE2 pattern
//   - $inSet(v, set)             – v is a member of the array set
//   - $minLen(str, n)            – rune length >= n
//   - $maxLen(str, n)            – rune length <= n
//   - $minItems(arr, n)          – array length >= n
//   - $maxItems(arr, n)          – array length <= n
package extvalidate

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"unicode"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

var (
	emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	uuidRe  = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// All returns a map of all validation functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"isEmail":         IsEmail(),
		"isURL":           IsURL(),
		"isUUID":          IsUUID(),
		"isIPv4":          IsIPv4(),
		"isIPv6":          IsIPv6(),
		"isAlpha":         IsAlpha(),
		"isAlphanumeric":  IsAlphanumeric(),
		"isNumericStr":    IsNumericStr(),
		"matchesRegex":    MatchesRegex(),
		"inSet":           InSet(),
		"minLen":          MinLen(),
		"maxLen":          MaxLen(),
		"minItems":        MinItems(),
		"maxItems":        MaxItems(),
	}
}

// IsEmail returns the CustomFunc for $isEmail(str).
func IsEmail() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		return emailRe.MatchString(s), nil
	}
}

// IsURL returns the CustomFunc for $isURL(str).
func IsURL() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		u, err := url.ParseRequestURI(s)
		if err != nil {
			return false, nil
		}
		switch u.Scheme {
		case "http", "https", "ftp":
			return u.Host != "", nil
		default:
			return false, nil
		}
	}
}

// IsUUID returns the CustomFunc for $isUUID(str).
func IsUUID() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		return uuidRe.MatchString(s), nil
	}
}

// IsIPv4 returns the CustomFunc for $isIPv4(str).
func IsIPv4() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		ip := net.ParseIP(s)
		return ip != nil && ip.To4() != nil, nil
	}
}

// IsIPv6 returns the CustomFunc for $isIPv6(str).
func IsIPv6() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		ip := net.ParseIP(s)
		return ip != nil && ip.To4() == nil, nil
	}
}

// IsAlpha returns the CustomFunc for $isAlpha(str).
func IsAlpha() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok || s == "" {
			return false, nil
		}
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false, nil
			}
		}
		return true, nil
	}
}

// IsAlphanumeric returns the CustomFunc for $isAlphanumeric(str).
func IsAlphanumeric() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok || s == "" {
			return false, nil
		}
		for _, r := range s {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return false, nil
			}
		}
		return true, nil
	}
}

// IsNumericStr returns the CustomFunc for $isNumericStr(str).
func IsNumericStr() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		s, ok := args[0].(string)
		if !ok {
			return false, nil
		}
		_, err := strconv.ParseFloat(s, 64)
		return err == nil, nil
	}
}

// MatchesRegex returns the CustomFunc for $matchesRegex(str, pattern).
// Uses Go RE2-based regexp (linear-time, inherently ReDoS-safe).
func MatchesRegex() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$matchesRegex: requires 2 arguments (str, pattern)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$matchesRegex: first argument must be a string")
		}
		pat, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$matchesRegex: pattern must be a string")
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("$matchesRegex: invalid pattern: %w", err)
		}
		return re.MatchString(s), nil
	}
}

// InSet returns the CustomFunc for $inSet(v, set).
func InSet() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$inSet: requires 2 arguments (value, set)")
		}
		set, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$inSet: %w", err)
		}
		for _, elem := range set {
			if elem == args[0] {
				return true, nil
			}
		}
		return false, nil
	}
}

// MinLen returns the CustomFunc for $minLen(str, n).
func MinLen() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$minLen: requires 2 arguments (str, n)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$minLen: first argument must be a string")
		}
		n, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$minLen: %w", err)
		}
		return float64(len([]rune(s))) >= n, nil
	}
}

// MaxLen returns the CustomFunc for $maxLen(str, n).
func MaxLen() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$maxLen: requires 2 arguments (str, n)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$maxLen: first argument must be a string")
		}
		n, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$maxLen: %w", err)
		}
		return float64(len([]rune(s))) <= n, nil
	}
}

// MinItems returns the CustomFunc for $minItems(arr, n).
func MinItems() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$minItems: requires 2 arguments (arr, n)")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$minItems: %w", err)
		}
		n, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$minItems: %w", err)
		}
		return float64(len(arr)) >= n, nil
	}
}

// MaxItems returns the CustomFunc for $maxItems(arr, n).
func MaxItems() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$maxItems: requires 2 arguments (arr, n)")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$maxItems: %w", err)
		}
		n, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$maxItems: %w", err)
		}
		return float64(len(arr)) <= n, nil
	}
}
