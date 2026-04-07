// Package extjson provides JSON serialisation, deserialisation, and structural
// manipulation functions for gnata.
//
// Functions
//
//   - $jsonParse(str)               – parse JSON string into a gnata value
//   - $jsonStringify(v [, indent])  – serialise value to JSON string
//   - $jsonDiff(a, b)               – compute differences as a JSON Patch array
//   - $jsonPatch(obj, ops)          – apply RFC 6902 JSON Patch operations
//   - $jsonPointer(obj, pointer)    – resolve RFC 6901 JSON Pointer
package extjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/recolabs/gnata"
)

// All returns a map of all JSON functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"jsonParse":     JsonParse(),
		"jsonStringify": JsonStringify(),
		"jsonDiff":      JsonDiff(),
		"jsonPatch":     JsonPatch(),
		"jsonPointer":   JsonPointerFunc(),
	}
}

// JsonParse returns the CustomFunc for $jsonParse(str).
func JsonParse() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$jsonParse: requires 1 argument (str)")
		}
		s, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("$jsonParse: argument must be a string")
		}
		var result any
		if err := json.Unmarshal([]byte(s), &result); err != nil {
			return nil, fmt.Errorf("$jsonParse: %w", err)
		}
		return result, nil
	}
}

// JsonStringify returns the CustomFunc for $jsonStringify(v [, indent]).
func JsonStringify() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$jsonStringify: requires at least 1 argument")
		}
		var b []byte
		var err error
		if len(args) >= 2 {
			indent, ok := args[1].(string)
			if !ok {
				return nil, fmt.Errorf("$jsonStringify: indent must be a string")
			}
			b, err = json.MarshalIndent(args[0], "", indent)
		} else {
			b, err = json.Marshal(args[0])
		}
		if err != nil {
			return nil, fmt.Errorf("$jsonStringify: %w", err)
		}
		return string(b), nil
	}
}

// JsonDiff returns the CustomFunc for $jsonDiff(a, b).
// Returns a JSON Patch-compatible array of {op, path, value} objects.
func JsonDiff() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$jsonDiff: requires 2 arguments (a, b)")
		}
		ops := buildDiff(args[0], args[1], "")
		return ops, nil
	}
}

// buildDiff recursively computes the diff between a and b, accumulating
// JSON Patch operations into a slice. path uses JSON Pointer notation.
func buildDiff(a, b any, path string) []any {
	var ops []any

	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)

	if aIsMap && bIsMap {
		// removals
		for k := range aMap {
			if _, exists := bMap[k]; !exists {
				ops = append(ops, map[string]any{
					"op":   "remove",
					"path": path + "/" + jsonPointerEscape(k),
				})
			}
		}
		// additions and changes
		for k, bv := range bMap {
			childPath := path + "/" + jsonPointerEscape(k)
			if av, exists := aMap[k]; !exists {
				ops = append(ops, map[string]any{
					"op":    "add",
					"path":  childPath,
					"value": bv,
				})
			} else {
				ops = append(ops, buildDiff(av, bv, childPath)...)
			}
		}
		return ops
	}

	// For non-map types (including slices and primitives), compare with DeepEqual.
	if !reflect.DeepEqual(a, b) {
		p := path
		if p == "" {
			p = "/"
		}
		ops = append(ops, map[string]any{
			"op":    "replace",
			"path":  p,
			"value": b,
		})
	}
	return ops
}

// jsonPointerEscape escapes a key for use in a JSON Pointer path segment.
func jsonPointerEscape(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}

// JsonPatch returns the CustomFunc for $jsonPatch(obj, ops).
// Applies RFC 6902 JSON Patch operations to obj (immutable – returns a new value).
func JsonPatch() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$jsonPatch: requires 2 arguments (obj, ops)")
		}
		ops, ok := args[1].([]any)
		if !ok {
			return nil, fmt.Errorf("$jsonPatch: ops must be an array")
		}
		doc := deepCopy(args[0])
		var err error
		for i, opAny := range ops {
			op, ok := opAny.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("$jsonPatch: operation %d must be an object", i)
			}
			opType, _ := op["op"].(string)
			path, _ := op["path"].(string)
			segs := pointerSegments(path)

			switch opType {
			case "add":
				doc, err = patchAdd(doc, segs, op["value"])
				if err != nil {
					return nil, fmt.Errorf("$jsonPatch op %d (add): %w", i, err)
				}
			case "remove":
				doc, err = patchRemove(doc, segs)
				if err != nil {
					return nil, fmt.Errorf("$jsonPatch op %d (remove): %w", i, err)
				}
			case "replace":
				doc, err = patchReplace(doc, segs, op["value"])
				if err != nil {
					return nil, fmt.Errorf("$jsonPatch op %d (replace): %w", i, err)
				}
			case "move":
				from, _ := op["from"].(string)
				fromSegs := pointerSegments(from)
				val, ok2 := resolvePointer(doc, fromSegs)
				if !ok2 {
					return nil, fmt.Errorf("$jsonPatch op %d (move): from path not found", i)
				}
				doc, err = patchRemove(doc, fromSegs)
				if err != nil {
					return nil, fmt.Errorf("$jsonPatch op %d (move remove): %w", i, err)
				}
				doc, err = patchAdd(doc, segs, val)
				if err != nil {
					return nil, fmt.Errorf("$jsonPatch op %d (move add): %w", i, err)
				}
			case "copy":
				from, _ := op["from"].(string)
				val, ok2 := resolvePointer(doc, pointerSegments(from))
				if !ok2 {
					return nil, fmt.Errorf("$jsonPatch op %d (copy): from path not found", i)
				}
				doc, err = patchAdd(doc, segs, deepCopy(val))
				if err != nil {
					return nil, fmt.Errorf("$jsonPatch op %d (copy): %w", i, err)
				}
			case "test":
				val, ok2 := resolvePointer(doc, segs)
				if !ok2 {
					return nil, fmt.Errorf("$jsonPatch op %d (test): path not found", i)
				}
				if !reflect.DeepEqual(val, op["value"]) {
					return nil, fmt.Errorf("$jsonPatch op %d (test): value mismatch", i)
				}
			default:
				return nil, fmt.Errorf("$jsonPatch op %d: unknown operation %q", i, opType)
			}
		}
		return doc, nil
	}
}

// JsonPointerFunc returns the CustomFunc for $jsonPointer(obj, pointer).
// Resolves an RFC 6901 JSON Pointer (e.g. "/a/b/0") against obj.
func JsonPointerFunc() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$jsonPointer: requires 2 arguments (obj, pointer)")
		}
		ptr, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$jsonPointer: pointer must be a string")
		}
		segs := pointerSegments(ptr)
		val, found := resolvePointer(args[0], segs)
		if !found {
			return nil, fmt.Errorf("$jsonPointer: path not found: %q", ptr)
		}
		return val, nil
	}
}

// pointerSegments splits an RFC 6901 JSON Pointer into unescaped path segments.
func pointerSegments(path string) []string {
	if path == "" || path == "/" {
		return nil
	}
	raw := path
	if strings.HasPrefix(raw, "/") {
		raw = raw[1:]
	}
	parts := strings.Split(raw, "/")
	for i, p := range parts {
		// unescape ~1 before ~0 per RFC 6901
		p = strings.ReplaceAll(p, "~1", "/")
		p = strings.ReplaceAll(p, "~0", "~")
		parts[i] = p
	}
	return parts
}

// resolvePointer recursively traverses obj following path segments.
func resolvePointer(obj any, segs []string) (any, bool) {
	if len(segs) == 0 {
		return obj, true
	}
	key := segs[0]
	switch v := obj.(type) {
	case map[string]any:
		child, ok := v[key]
		if !ok {
			return nil, false
		}
		return resolvePointer(child, segs[1:])
	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(v) {
			return nil, false
		}
		return resolvePointer(v[idx], segs[1:])
	default:
		return nil, false
	}
}

// patchAdd sets value at segs path, creating parent objects as needed.
// On an array target, "-" appends; an integer index inserts.
func patchAdd(doc any, segs []string, value any) (any, error) {
	if len(segs) == 0 {
		return value, nil
	}
	if len(segs) == 1 {
		return patchAddLeaf(doc, segs[0], value)
	}
	m, ok := doc.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cannot traverse non-object at %q", segs[0])
	}
	out := shallowCopyMap(m)
	child, exists := out[segs[0]]
	if !exists {
		child = map[string]any{}
	}
	var err error
	out[segs[0]], err = patchAdd(child, segs[1:], value)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func patchAddLeaf(doc any, key string, value any) (any, error) {
	switch v := doc.(type) {
	case map[string]any:
		out := shallowCopyMap(v)
		out[key] = value
		return out, nil
	case []any:
		if key == "-" {
			return append(append([]any{}, v...), value), nil
		}
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx > len(v) {
			return nil, fmt.Errorf("invalid array index %q", key)
		}
		out := make([]any, len(v)+1)
		copy(out, v[:idx])
		out[idx] = value
		copy(out[idx+1:], v[idx:])
		return out, nil
	default:
		return nil, fmt.Errorf("cannot add to type %T", doc)
	}
}

// patchRemove removes the element at segs path.
func patchRemove(doc any, segs []string) (any, error) {
	if len(segs) == 0 {
		return nil, fmt.Errorf("cannot remove root")
	}
	if len(segs) == 1 {
		return patchRemoveLeaf(doc, segs[0])
	}
	m, ok := doc.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cannot traverse non-object at %q", segs[0])
	}
	out := shallowCopyMap(m)
	child, exists := out[segs[0]]
	if !exists {
		return nil, fmt.Errorf("key not found: %q", segs[0])
	}
	var err error
	out[segs[0]], err = patchRemove(child, segs[1:])
	if err != nil {
		return nil, err
	}
	return out, nil
}

func patchRemoveLeaf(doc any, key string) (any, error) {
	switch v := doc.(type) {
	case map[string]any:
		if _, exists := v[key]; !exists {
			return nil, fmt.Errorf("key not found: %q", key)
		}
		out := shallowCopyMap(v)
		delete(out, key)
		return out, nil
	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(v) {
			return nil, fmt.Errorf("invalid array index %q", key)
		}
		out := make([]any, 0, len(v)-1)
		out = append(out, v[:idx]...)
		out = append(out, v[idx+1:]...)
		return out, nil
	default:
		return nil, fmt.Errorf("cannot remove from type %T", doc)
	}
}

// patchReplace replaces the value at segs path (the path must exist).
func patchReplace(doc any, segs []string, value any) (any, error) {
	if len(segs) == 0 {
		return value, nil
	}
	if len(segs) == 1 {
		return patchReplaceLeaf(doc, segs[0], value)
	}
	m, ok := doc.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("cannot traverse non-object at %q", segs[0])
	}
	out := shallowCopyMap(m)
	child, exists := out[segs[0]]
	if !exists {
		return nil, fmt.Errorf("key not found: %q", segs[0])
	}
	var err error
	out[segs[0]], err = patchReplace(child, segs[1:], value)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func patchReplaceLeaf(doc any, key string, value any) (any, error) {
	switch v := doc.(type) {
	case map[string]any:
		if _, exists := v[key]; !exists {
			return nil, fmt.Errorf("key not found for replace: %q", key)
		}
		out := shallowCopyMap(v)
		out[key] = value
		return out, nil
	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(v) {
			return nil, fmt.Errorf("invalid array index %q", key)
		}
		out := make([]any, len(v))
		copy(out, v)
		out[idx] = value
		return out, nil
	default:
		return nil, fmt.Errorf("cannot replace in type %T", doc)
	}
}

// shallowCopyMap returns a shallow copy of a map[string]any.
func shallowCopyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// deepCopy recursively copies maps and slices; other values are returned as-is.
func deepCopy(v any) any {
	switch vt := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(vt))
		for k, val := range vt {
			out[k] = deepCopy(val)
		}
		return out
	case []any:
		out := make([]any, len(vt))
		for i, val := range vt {
			out[i] = deepCopy(val)
		}
		return out
	default:
		return v
	}
}
