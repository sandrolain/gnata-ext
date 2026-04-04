// Package extobject provides extended object functions for gnata beyond the
// official JSONata 2.x specification.
//
// HOF functions (mapValues, mapKeys) are not included because gnata does not
// expose a Caller interface. Use JSONata built-in $each for higher-order operations.
//
// Functions
//
//   - $values(object)               – array of object values
//   - $pairs(object)                – array of [key, value] pairs
//   - $fromPairs(pairs)             – object from [[key, value], …]
//   - $pick(object, keys)           – keep only the specified keys
//   - $omit(object, keys)           – remove the specified keys
//   - $deepMerge(target, source)    – deep-merge source into target
//   - $invert(object)               – swap keys and values
//   - $size(object)                 – number of own keys
//   - $rename(object, oldKey, newKey) – rename a single key
package extobject

import (
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extended object functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"values":    Values(),
		"pairs":     Pairs(),
		"fromPairs": FromPairs(),
		"pick":      Pick(),
		"omit":      Omit(),
		"deepMerge": DeepMerge(),
		"invert":    Invert(),
		"size":      Size(),
		"rename":    Rename(),
	}
}

// Values returns the CustomFunc for $values(object).
func Values() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$values: requires 1 argument")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$values: %w", err)
		}
		result := make([]any, 0, len(obj))
		for _, v := range obj {
			result = append(result, v)
		}
		return result, nil
	}
}

// Pairs returns the CustomFunc for $pairs(object).
// Returns an array of [key, value] pairs.
func Pairs() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$pairs: requires 1 argument")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$pairs: %w", err)
		}
		result := make([]any, 0, len(obj))
		for k, v := range obj {
			result = append(result, []any{k, v})
		}
		return result, nil
	}
}

// FromPairs returns the CustomFunc for $fromPairs(pairs).
// Accepts an array of [key, value] pairs or objects with key/value fields.
func FromPairs() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$fromPairs: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$fromPairs: %w", err)
		}
		result := make(map[string]any, len(arr))
		for i, item := range arr {
			switch pair := item.(type) {
			case []any:
				if len(pair) < 2 {
					return nil, fmt.Errorf("$fromPairs: pair at index %d has fewer than 2 elements", i)
				}
				k, ok := pair[0].(string)
				if !ok {
					return nil, fmt.Errorf("$fromPairs: key at index %d must be a string", i)
				}
				result[k] = pair[1]
			case map[string]any:
				k, ok := pair["key"].(string)
				if !ok {
					return nil, fmt.Errorf("$fromPairs: object at index %d missing string 'key' field", i)
				}
				result[k] = pair["value"]
			default:
				return nil, fmt.Errorf("$fromPairs: unexpected type at index %d: %T", i, item)
			}
		}
		return result, nil
	}
}

// Pick returns the CustomFunc for $pick(object, keys).
func Pick() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$pick: requires 2 arguments")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$pick: %w", err)
		}
		keys, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$pick: %w", err)
		}
		result := make(map[string]any, len(keys))
		for _, k := range keys {
			ks, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("$pick: key must be a string, got %T", k)
			}
			if v, exists := obj[ks]; exists {
				result[ks] = v
			}
		}
		return result, nil
	}
}

// Omit returns the CustomFunc for $omit(object, keys).
func Omit() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$omit: requires 2 arguments")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$omit: %w", err)
		}
		keys, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$omit: %w", err)
		}
		skip := make(map[string]bool, len(keys))
		for _, k := range keys {
			ks, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("$omit: key must be a string, got %T", k)
			}
			skip[ks] = true
		}
		result := make(map[string]any)
		for k, v := range obj {
			if !skip[k] {
				result[k] = v
			}
		}
		return result, nil
	}
}

// DeepMerge returns the CustomFunc for $deepMerge(target, source).
// Returns a new object with source merged into target recursively.
func DeepMerge() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$deepMerge: requires 2 arguments")
		}
		target, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$deepMerge: target: %w", err)
		}
		source, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$deepMerge: source: %w", err)
		}
		result := make(map[string]any, len(target))
		for k, v := range target {
			result[k] = v
		}
		deepMergeInto(result, source)
		return result, nil
	}
}

// Invert returns the CustomFunc for $invert(object).
// Values must be strings to become keys.
func Invert() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$invert: requires 1 argument")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$invert: %w", err)
		}
		result := make(map[string]any, len(obj))
		for k, v := range obj {
			vs, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("$invert: value for key %q must be a string, got %T", k, v)
			}
			result[vs] = k
		}
		return result, nil
	}
}

// Size returns the CustomFunc for $size(object).
func Size() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$size: requires 1 argument")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$size: %w", err)
		}
		return float64(len(obj)), nil
	}
}

// Rename returns the CustomFunc for $rename(object, oldKey, newKey).
func Rename() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$rename: requires 3 arguments (object, oldKey, newKey)")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$rename: %w", err)
		}
		oldKey, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$rename: oldKey must be a string")
		}
		newKey, ok := args[2].(string)
		if !ok {
			return nil, fmt.Errorf("$rename: newKey must be a string")
		}
		result := make(map[string]any, len(obj))
		for k, v := range obj {
			if k == oldKey {
				result[newKey] = v
			} else {
				result[k] = v
			}
		}
		return result, nil
	}
}

// deepMergeInto merges src into dst in place.
func deepMergeInto(dst, src map[string]any) {
	for k, sv := range src {
		if dv, exists := dst[k]; exists {
			dm, dOK := dv.(map[string]any)
			sm, sOK := sv.(map[string]any)
			if dOK && sOK {
				sub := make(map[string]any, len(dm))
				for dk, dval := range dm {
					sub[dk] = dval
				}
				deepMergeInto(sub, sm)
				dst[k] = sub
				continue
			}
		}
		dst[k] = sv
	}
}
