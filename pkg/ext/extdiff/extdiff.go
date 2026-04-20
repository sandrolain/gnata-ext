// Package extdiff provides structural diff and deep equality functions for gnata.
//
// Functions:
//
//   - $diff(a, b)          – structural diff {added, removed, changed}
//   - $patch(obj, diff)    – apply a diff to reconstruct b from a
//   - $changed(a, b, key)  – true if a specific key changed
//   - $addedKeys(a, b)     – keys present in b but not in a
//   - $removedKeys(a, b)   – keys present in a but not in b
//   - $arrayDiff(a, b)     – {added, removed} between two arrays
//   - $deepEqual(a, b)     – recursive deep equality
package extdiff

import (
	"fmt"
	"reflect"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extdiff functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"diff":        Diff(),
		"patch":       Patch(),
		"changed":     Changed(),
		"addedKeys":   AddedKeys(),
		"removedKeys": RemovedKeys(),
		"arrayDiff":   ArrayDiff(),
		"deepEqual":   DeepEqual(),
	}
}

// Diff returns the CustomFunc for $diff(a, b).
// Both arguments must be objects (map[string]any).
// Returns {added: object, removed: object, changed: object}.
// "added" contains keys in b but not in a.
// "removed" contains keys in a but not in b.
// "changed" maps key → {from: oldVal, to: newVal} for keys whose value changed.
func Diff() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$diff: requires 2 arguments (a, b)")
		}
		a, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$diff: first argument must be an object: %w", err)
		}
		b, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$diff: second argument must be an object: %w", err)
		}

		added := map[string]any{}
		removed := map[string]any{}
		changed := map[string]any{}

		for k, bv := range b {
			if av, exists := a[k]; !exists {
				added[k] = bv
			} else if !deepEqual(av, bv) {
				changed[k] = map[string]any{"from": av, "to": bv}
			}
		}
		for k, av := range a {
			if _, exists := b[k]; !exists {
				removed[k] = av
			}
		}

		return map[string]any{
			"added":   added,
			"removed": removed,
			"changed": changed,
		}, nil
	}
}

// Patch returns the CustomFunc for $patch(obj, diff).
// Applies a diff (as produced by $diff) to obj, returning the modified object.
// Keys in diff.added are added; keys in diff.removed are deleted; keys in diff.changed are updated.
func Patch() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$patch: requires 2 arguments (obj, diff)")
		}
		obj, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$patch: first argument must be an object: %w", err)
		}
		d, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$patch: second argument must be an object: %w", err)
		}

		// Build a copy.
		result := make(map[string]any, len(obj))
		for k, v := range obj {
			result[k] = v
		}

		if added, ok := d["added"].(map[string]any); ok {
			for k, v := range added {
				result[k] = v
			}
		}
		if removed, ok := d["removed"].(map[string]any); ok {
			for k := range removed {
				delete(result, k)
			}
		}
		if changed, ok := d["changed"].(map[string]any); ok {
			for k, entry := range changed {
				if e, ok := entry.(map[string]any); ok {
					if to, exists := e["to"]; exists {
						result[k] = to
					}
				}
			}
		}

		return result, nil
	}
}

// Changed returns the CustomFunc for $changed(a, b, key).
// Returns true if the value of key differs between objects a and b.
func Changed() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$changed: requires 3 arguments (a, b, key)")
		}
		a, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$changed: first argument must be an object: %w", err)
		}
		b, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$changed: second argument must be an object: %w", err)
		}
		key, ok := args[2].(string)
		if !ok {
			return nil, fmt.Errorf("$changed: key must be a string")
		}
		return !deepEqual(a[key], b[key]), nil
	}
}

// AddedKeys returns the CustomFunc for $addedKeys(a, b).
// Returns keys present in b but not in a.
func AddedKeys() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$addedKeys: requires 2 arguments (a, b)")
		}
		a, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$addedKeys: first argument must be an object: %w", err)
		}
		b, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$addedKeys: second argument must be an object: %w", err)
		}
		var out []any
		for k := range b {
			if _, exists := a[k]; !exists {
				out = append(out, k)
			}
		}
		if out == nil {
			out = []any{}
		}
		return out, nil
	}
}

// RemovedKeys returns the CustomFunc for $removedKeys(a, b).
// Returns keys present in a but not in b.
func RemovedKeys() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$removedKeys: requires 2 arguments (a, b)")
		}
		a, err := extutil.ToObject(args[0])
		if err != nil {
			return nil, fmt.Errorf("$removedKeys: first argument must be an object: %w", err)
		}
		b, err := extutil.ToObject(args[1])
		if err != nil {
			return nil, fmt.Errorf("$removedKeys: second argument must be an object: %w", err)
		}
		var out []any
		for k := range a {
			if _, exists := b[k]; !exists {
				out = append(out, k)
			}
		}
		if out == nil {
			out = []any{}
		}
		return out, nil
	}
}

// ArrayDiff returns the CustomFunc for $arrayDiff(a, b).
// Returns {added: array, removed: array} representing elements present in only one array.
// Comparison is deep equality.
func ArrayDiff() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$arrayDiff: requires 2 arguments (a, b)")
		}
		a, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$arrayDiff: first argument must be an array: %w", err)
		}
		b, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$arrayDiff: second argument must be an array: %w", err)
		}

		added := filterNotIn(b, a)
		removed := filterNotIn(a, b)
		return map[string]any{
			"added":   added,
			"removed": removed,
		}, nil
	}
}

// DeepEqual returns the CustomFunc for $deepEqual(a, b).
// Returns true if a and b are recursively equal.
func DeepEqual() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$deepEqual: requires 2 arguments (a, b)")
		}
		return deepEqual(args[0], args[1]), nil
	}
}

// --- helpers ---

// deepEqual performs recursive deep equality comparison.
func deepEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

// filterNotIn returns elements of src that are not deeply equal to any element in other.
func filterNotIn(src, other []any) []any {
	var out []any
	for _, s := range src {
		found := false
		for _, o := range other {
			if deepEqual(s, o) {
				found = true
				break
			}
		}
		if !found {
			out = append(out, s)
		}
	}
	if out == nil {
		out = []any{}
	}
	return out
}
