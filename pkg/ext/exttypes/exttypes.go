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
//   - $toArray(v)     – wraps v in array if not already one (nil → [])
//   - $defined(v)     – true if v is not nil
//   - $nullish(v, d)  – returns v if not nil, otherwise d (alias for $default)
//   - $typeOf(v)      – returns type name as string
//   - $toNumber(v)    – coerces v to number
//   - $toString(v)    – coerces v to string
//   - $toBool(v)      – coerces v to boolean
package exttypes

import (
	"fmt"
	"strconv"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
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
		"toArray":     ToArray(),
		"defined":     Defined(),
		"nullish":     Nullish(),
		"typeOf":      TypeOf(),
		"toNumber":    ToNumber(),
		"toString":    ToString(),
		"toBool":      ToBool(),
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

// ToArray returns the CustomFunc for $toArray(v).
// Wraps v in a slice if not already one; nil → empty slice.
func ToArray() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return []any{}, nil
		}
		v := args[0]
		if v == nil {
			return []any{}, nil
		}
		if arr, ok := extutil.ToArray(v); ok == nil {
			return arr, nil
		}
		return []any{v}, nil
	}
}

// Defined returns the CustomFunc for $defined(v).
// Returns true if v is not nil.
func Defined() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		return args[0] != nil, nil
	}
}

// Nullish returns the CustomFunc for $nullish(v, d).
// Returns v if not nil, otherwise d.
func Nullish() gnata.CustomFunc {
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

// TypeOf returns the CustomFunc for $typeOf(v).
// Returns "string", "number", "boolean", "array", "object", or "null".
func TypeOf() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return "null", nil
		}
		switch v := args[0].(type) {
		case nil:
			return "null", nil
		case string:
			_ = v
			return "string", nil
		case float64, int, int64, float32:
			return "number", nil
		case bool:
			return "boolean", nil
		case []any:
			return "array", nil
		case map[string]any:
			return "object", nil
		default:
			return "unknown", nil
		}
	}
}

// ToNumber returns the CustomFunc for $toNumber(v).
// Coerces string or bool to float64.
func ToNumber() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$toNumber: requires 1 argument")
		}
		v := args[0]
		switch val := v.(type) {
		case float64:
			return val, nil
		case int:
			return float64(val), nil
		case int64:
			return float64(val), nil
		case bool:
			if val {
				return float64(1), nil
			}
			return float64(0), nil
		case string:
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("$toNumber: cannot convert %q to number", val)
			}
			return n, nil
		default:
			return nil, fmt.Errorf("$toNumber: cannot convert type %T", v)
		}
	}
}

// ToString returns the CustomFunc for $toString(v).
// Coerces number, bool, or nil to string.
func ToString() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return "", nil
		}
		return fmt.Sprintf("%v", args[0]), nil
	}
}

// ToBool returns the CustomFunc for $toBool(v).
// Returns false for nil, false, 0, "", and empty collections.
func ToBool() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return false, nil
		}
		switch v := args[0].(type) {
		case nil:
			return false, nil
		case bool:
			return v, nil
		case float64:
			return v != 0, nil
		case int:
			return v != 0, nil
		case int64:
			return v != 0, nil
		case string:
			return v != "", nil
		case []any:
			return len(v) > 0, nil
		case map[string]any:
			return len(v) > 0, nil
		default:
			return true, nil
		}
	}
}
