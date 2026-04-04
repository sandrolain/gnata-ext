// Package extutil provides shared type-conversion helpers for the gnata-ext
// extension sub-packages.
package extutil

import "fmt"

// ToFloat converts v to float64. Supports float64, int, and int64.
func ToFloat(v any) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("expected a number, got %T", v)
	}
}

// ToInt converts v to int. Returns (0, false) when v is not numeric.
func ToInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}

// ToFloatSlice converts an []any of numbers to []float64.
func ToFloatSlice(v any) ([]float64, error) {
	arr, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", v)
	}
	out := make([]float64, len(arr))
	for i, elem := range arr {
		f, err := ToFloat(elem)
		if err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		out[i] = f
	}
	return out, nil
}

// ToArray asserts that v is []any and returns it.
func ToArray(v any) ([]any, error) {
	arr, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", v)
	}
	return arr, nil
}

// ToObject asserts that v is map[string]any and returns it.
func ToObject(v any) (map[string]any, error) {
	obj, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected object, got %T", v)
	}
	return obj, nil
}
