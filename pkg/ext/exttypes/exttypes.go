// Package exttypes provides type-inspection and value-coercion functions for gnata.
//
// Functions
//
//   - $isString(v)    – true if v is a string
//   - $isNumber(v)    – true if v is a number (float64, int, int64)
//   - $isBoolean(v)   – true if v is a boolean
//   - $isArray(v)     – true if v is an array
//   - $isObject(v)    – true if v is an object
//   - $isNull(v)      – true if v is null/nil
//   - $isUndefined(v) – true if v is undefined (nil in gnata)
//   - $isEmpty(v)     – true if v is nil, empty string, empty array, or empty object
//   - $default(v, d)  – returns v if non-nil, otherwise d
//   - $identity(v)    – returns v unchanged
package exttypes

import (
	"github.com/recolabs/gnata"
)

// All returns a map of all extended type functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"isString":    IsString(),
		"isNumber":    IsNumber(),
		"isBoolean":   IsBoolean(),
		"isArray":     IsArray(),
		"isObject":    IsObject(),
		"isNull":      IsNull(),
		"isUndefined": IsUndefined(),
		"isEmpty":     IsEmpty(),
		"default":     Default(),
		"identity":    Identity(),
	}
}

// IsString returns the CustomFunc for $isString(v).
func IsString() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		_, ok := args[0].(string)
		return ok, nil
	}
}

// IsNumber returns the CustomFunc for $isNumber(v).
func IsNumber() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		switch args[0].(type) {
		case float64, int, int64:
			return true, nil
		default:
			return false, nil
		}
	}
}

// IsBoolean returns the CustomFunc for $isBoolean(v).
func IsBoolean() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		_, ok := args[0].(bool)
		return ok, nil
	}
}

// IsArray returns the CustomFunc for $isArray(v).
func IsArray() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		_, ok := args[0].([]any)
		return ok, nil
	}
}

// IsObject returns the CustomFunc for $isObject(v).
func IsObject() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		_, ok := args[0].(map[string]any)
		return ok, nil
	}
}

// IsNull returns the CustomFunc for $isNull(v).
func IsNull() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return true, nil
		}
		return args[0] == nil, nil
	}
}

// IsUndefined returns the CustomFunc for $isUndefined(v).
// In gnata, undefined is represented as nil.
func IsUndefined() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return true, nil
		}
		return args[0] == nil, nil
	}
}

// IsEmpty returns the CustomFunc for $isEmpty(v).
// Returns true for nil, empty string, empty array, and empty object.
func IsEmpty() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return true, nil
		}
		v := args[0]
		if v == nil {
			return true, nil
		}
		switch val := v.(type) {
		case string:
			return val == "", nil
		case []any:
			return len(val) == 0, nil
		case map[string]any:
			return len(val) == 0, nil
		default:
			return false, nil
		}
	}
}

// Default returns the CustomFunc for $default(v, d).
// Returns v if non-nil; returns d otherwise.
func Default() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			if len(args) == 1 {
				return args[0], nil
			}
			return nil, nil
		}
		if args[0] != nil {
			return args[0], nil
		}
		return args[1], nil
	}
}

// Identity returns the CustomFunc for $identity(v).
// Returns v unchanged.
func Identity() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, nil
		}
		return args[0], nil
	}
}
