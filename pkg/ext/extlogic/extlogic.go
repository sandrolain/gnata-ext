// Package extlogic provides conditional and control-flow functions for gnata.
//
// Functions:
//
//   - $ifElse(cond, then, else)     – ternary conditional
//   - $when(cond, value)            – returns value only if cond is true
//   - $unless(cond, value)          – returns value only if cond is false
//   - $switch(v, cases [, default]) – map-based switch/case
//   - $coalesce(v1, v2, ...)        – first non-nil/non-empty value
//   - $tap(v)                       – returns v unchanged (debug helper)
package extlogic

import (
	"fmt"

	"github.com/recolabs/gnata"
)

// All returns a map of all extlogic functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"ifElse":   IfElse(),
		"when":     When(),
		"unless":   Unless(),
		"switch":   Switch(),
		"coalesce": Coalesce(),
		"tap":      Tap(),
	}
}

// IfElse returns the CustomFunc for $ifElse(cond, then, else).
// Returns the `then` value when cond is truthy, otherwise the `else` value.
func IfElse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$ifElse: requires 3 arguments (cond, then, else)")
		}
		if isTruthy(args[0]) {
			return args[1], nil
		}
		return args[2], nil
	}
}

// When returns the CustomFunc for $when(cond, value).
// Returns value when cond is truthy, otherwise nil.
func When() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$when: requires 2 arguments (cond, value)")
		}
		if isTruthy(args[0]) {
			return args[1], nil
		}
		return nil, nil
	}
}

// Unless returns the CustomFunc for $unless(cond, value).
// Returns value when cond is falsy, otherwise nil.
func Unless() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$unless: requires 2 arguments (cond, value)")
		}
		if !isTruthy(args[0]) {
			return args[1], nil
		}
		return nil, nil
	}
}

// Switch returns the CustomFunc for $switch(v, cases [, default]).
// cases must be a map[string]any; the key matching fmt.Sprintf("%v", v) is returned.
// If no key matches and a default (third argument) is provided, it is returned;
// otherwise nil is returned.
func Switch() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$switch: requires at least 2 arguments (value, cases)")
		}
		cases, ok := args[1].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("$switch: second argument must be an object")
		}
		key := fmt.Sprintf("%v", args[0])
		if val, found := cases[key]; found {
			return val, nil
		}
		if len(args) >= 3 {
			return args[2], nil
		}
		return nil, nil
	}
}

// Coalesce returns the CustomFunc for $coalesce(v1, v2, ...).
// Returns the first argument that is not nil, false, 0, or empty string.
func Coalesce() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("$coalesce: requires at least 1 argument")
		}
		for _, v := range args {
			if !isNilOrEmpty(v) {
				return v, nil
			}
		}
		return nil, nil
	}
}

// Tap returns the CustomFunc for $tap(v).
// Returns v unchanged. Useful as a pass-through in expression chains.
func Tap() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$tap: requires 1 argument")
		}
		return args[0], nil
	}
}

// isTruthy returns true if v is considered truthy:
// non-nil, non-false, non-zero, non-empty-string, non-empty-slice, non-empty-map.
func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	case int64:
		return val != 0
	case string:
		return val != ""
	case []any:
		return len(val) > 0
	case map[string]any:
		return len(val) > 0
	}
	return true
}

// isNilOrEmpty returns true if v is nil, false, 0, "", empty slice, or empty map.
func isNilOrEmpty(v any) bool {
	return !isTruthy(v)
}
