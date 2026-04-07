// Package extpath provides deep path-access functions for gnata using dot-notation.
//
// Paths use "." as separator by default: "a.b.c" traverses {a:{b:{c:value}}}.
// All mutations return new objects (immutable); the input is never modified.
//
// Functions
//
//   - $get(obj, path [, default])  – value at path; returns default (or nil) if absent
//   - $set(obj, path, value)       – new object with value set at path (creates nested maps)
//   - $del(obj, path)              – new object with the key at path removed
//   - $has(obj, path)              – true if path exists and its value is non-nil
//   - $flattenObj(obj [, sep])     – {a:{b:1}} → {"a.b":1}
//   - $expandObj(obj [, sep])      – {"a.b":1} → {a:{b:1}}
package extpath

import (
	"fmt"
	"strings"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all deep-path functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"get":        Get(),
		"set":        Set(),
		"del":        Del(),
		"has":        Has(),
		"flattenObj": FlattenObj(),
		"expandObj":  ExpandObj(),
	}
}

// parsePath splits a dot-notation path into segments.
func parsePath(path string) []string {
	if path == "" {
		return nil
	}
	return strings.Split(path, ".")
}

// getPath traverses obj following keys, returning (value, true) when found.
func getPath(obj map[string]any, keys []string) (any, bool) {
	if len(keys) == 0 {
		return obj, true
	}
	v, ok := obj[keys[0]]
	if !ok {
		return nil, false
	}
	if len(keys) == 1 {
		return v, true
	}
	child, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}
	return getPath(child, keys[1:])
}

// copyMap makes a shallow copy of m.
func copyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// setPath returns a new map with value written at keys depth.
func setPath(obj map[string]any, keys []string, value any) map[string]any {
	out := copyMap(obj)
	if len(keys) == 1 {
		out[keys[0]] = value
		return out
	}
	child, _ := obj[keys[0]].(map[string]any)
	if child == nil {
		child = map[string]any{}
	}
	out[keys[0]] = setPath(child, keys[1:], value)
	return out
}

// delPath returns a new map with the leaf key removed.
func delPath(obj map[string]any, keys []string) map[string]any {
	out := copyMap(obj)
	if len(keys) == 1 {
		delete(out, keys[0])
		return out
	}
	child, ok := obj[keys[0]].(map[string]any)
	if !ok {
		return out
	}
	out[keys[0]] = delPath(child, keys[1:])
	return out
}

// Get returns the CustomFunc for $get(obj, path [, default]).
func Get() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$get: requires at least 2 arguments (obj, path)")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$get: %w", err)
		}
		path, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$get: path must be a string")
		}
		keys := parsePath(path)
		if len(keys) == 0 {
			return obj, nil
		}
		v, found := getPath(obj, keys)
		if !found || v == nil {
			if len(args) >= 3 {
				return args[2], nil
			}
			return nil, nil
		}
		return v, nil
	}
}

// Set returns the CustomFunc for $set(obj, path, value).
func Set() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$set: requires 3 arguments (obj, path, value)")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$set: %w", err)
		}
		path, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$set: path must be a string")
		}
		keys := parsePath(path)
		if len(keys) == 0 {
			return nil, fmt.Errorf("$set: path must not be empty")
		}
		return setPath(obj, keys, args[2]), nil
	}
}

// Del returns the CustomFunc for $del(obj, path).
func Del() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$del: requires 2 arguments (obj, path)")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$del: %w", err)
		}
		path, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$del: path must be a string")
		}
		keys := parsePath(path)
		if len(keys) == 0 {
			return copyMap(obj), nil
		}
		return delPath(obj, keys), nil
	}
}

// Has returns the CustomFunc for $has(obj, path).
func Has() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$has: requires 2 arguments (obj, path)")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$has: %w", err)
		}
		path, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$has: path must be a string")
		}
		keys := parsePath(path)
		if len(keys) == 0 {
			return true, nil
		}
		v, found := getPath(obj, keys)
		return found && v != nil, nil
	}
}

// FlattenObj returns the CustomFunc for $flattenObj(obj [, sep]).
// Converts {a:{b:1}} to {"a.b":1}.
func FlattenObj() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$flattenObj: requires at least 1 argument")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$flattenObj: %w", err)
		}
		sep := "."
		if len(args) >= 2 {
			if s, ok := args[1].(string); ok && s != "" {
				sep = s
			}
		}
		result := make(map[string]any)
		flattenInto(result, obj, "", sep)
		return result, nil
	}
}

func flattenInto(out map[string]any, obj map[string]any, prefix, sep string) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		if child, ok := v.(map[string]any); ok {
			flattenInto(out, child, key, sep)
		} else {
			out[key] = v
		}
	}
}

// ExpandObj returns the CustomFunc for $expandObj(obj [, sep]).
// Converts {"a.b":1} to {a:{b:1}}.
func ExpandObj() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$expandObj: requires at least 1 argument")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$expandObj: %w", err)
		}
		sep := "."
		if len(args) >= 2 {
			if s, ok := args[1].(string); ok && s != "" {
				sep = s
			}
		}
		result := map[string]any{}
		for k, v := range obj {
			keys := strings.Split(k, sep)
			result = setPath(result, keys, v)
		}
		return result, nil
	}
}
