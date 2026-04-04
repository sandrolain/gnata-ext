// Package extarray provides extended array functions for gnata beyond the
// official JSONata 2.x specification.
//
// HOF functions (groupBy, countBy, sumBy, minBy, maxBy, accumulate) are not
// included because gnata does not expose a Caller interface. Use JSONata built-in
// $map, $filter, $reduce, and $each for higher-order operations.
//
// Functions
//
//   - $first(array)                        – first element
//   - $last(array)                         – last element
//   - $take(array, n)                      – first n elements
//   - $skip(array, n)                      – all but first n elements
//   - $slice(array, start [, end])         – sub-array
//   - $flatten(array [, depth])            – flatten nested arrays
//   - $chunk(array, size)                  – split into chunks
//   - $union(a, b)                         – set union
//   - $intersection(a, b)                  – set intersection
//   - $difference(a, b)                    – set difference (a minus b)
//   - $symmetricDifference(a, b)           – elements in either but not both
//   - $range(start, end [, step])          – numeric range
//   - $zipLongest(a, b [, fill])           – zip two arrays, padding shorter one
//   - $window(array, size [, step])        – sliding windows
package extarray

import (
	"fmt"
	"math"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

const maxRangeItems = 100000

// All returns a map of all extended array functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"first":              First(),
		"last":               Last(),
		"take":               Take(),
		"skip":               Skip(),
		"slice":              Slice(),
		"flatten":            Flatten(),
		"chunk":              Chunk(),
		"union":              Union(),
		"intersection":       Intersection(),
		"difference":         Difference(),
		"symmetricDifference": SymmetricDifference(),
		"range":              Range(),
		"zipLongest":         ZipLongest(),
		"window":             Window(),
	}
}

// First returns the CustomFunc for $first(array).
func First() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$first: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$first: %w", err)
		}
		if len(arr) == 0 {
			return nil, nil
		}
		return arr[0], nil
	}
}

// Last returns the CustomFunc for $last(array).
func Last() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$last: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$last: %w", err)
		}
		if len(arr) == 0 {
			return nil, nil
		}
		return arr[len(arr)-1], nil
	}
}

// Take returns the CustomFunc for $take(array, n).
func Take() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$take: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$take: %w", err)
		}
		n, ok := extutil.ToInt(args[1])
		if !ok {
			return nil, fmt.Errorf("$take: n must be a number")
		}
		if n < 0 {
			n = 0
		}
		if n > len(arr) {
			n = len(arr)
		}
		result := make([]any, n)
		copy(result, arr[:n])
		return result, nil
	}
}

// Skip returns the CustomFunc for $skip(array, n).
func Skip() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$skip: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$skip: %w", err)
		}
		n, ok := extutil.ToInt(args[1])
		if !ok {
			return nil, fmt.Errorf("$skip: n must be a number")
		}
		if n < 0 {
			n = 0
		}
		if n > len(arr) {
			n = len(arr)
		}
		result := make([]any, len(arr)-n)
		copy(result, arr[n:])
		return result, nil
	}
}

// Slice returns the CustomFunc for $slice(array, start [, end]).
func Slice() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$slice: requires at least 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$slice: %w", err)
		}
		start := normaliseIndex(mustInt(args[1]), len(arr))
		end := len(arr)
		if len(args) >= 3 && args[2] != nil {
			end = normaliseIndex(mustInt(args[2]), len(arr))
		}
		if start > end {
			start = end
		}
		result := make([]any, end-start)
		copy(result, arr[start:end])
		return result, nil
	}
}

// Flatten returns the CustomFunc for $flatten(array [, depth]).
func Flatten() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$flatten: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$flatten: %w", err)
		}
		depth := -1 // unlimited
		if len(args) >= 2 && args[1] != nil {
			d, ok := extutil.ToInt(args[1])
			if !ok {
				return nil, fmt.Errorf("$flatten: depth must be a number")
			}
			depth = d
		}
		return flattenArray(arr, depth), nil
	}
}

// Chunk returns the CustomFunc for $chunk(array, size).
func Chunk() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$chunk: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$chunk: %w", err)
		}
		size, ok := extutil.ToInt(args[1])
		if !ok || size <= 0 {
			return nil, fmt.Errorf("$chunk: size must be a positive number")
		}
		var chunks []any
		for i := 0; i < len(arr); i += size {
			end := i + size
			if end > len(arr) {
				end = len(arr)
			}
			chunk := make([]any, end-i)
			copy(chunk, arr[i:end])
			chunks = append(chunks, chunk)
		}
		if chunks == nil {
			return []any{}, nil
		}
		return chunks, nil
	}
}

// Union returns the CustomFunc for $union(a, b).
func Union() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$union: requires 2 arguments")
		}
		a, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$union: %w", err)
		}
		b, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$union: %w", err)
		}
		seen := make(map[any]bool)
		var result []any
		for _, v := range append(a, b...) {
			if !seen[v] {
				seen[v] = true
				result = append(result, v)
			}
		}
		if result == nil {
			return []any{}, nil
		}
		return result, nil
	}
}

// Intersection returns the CustomFunc for $intersection(a, b).
func Intersection() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$intersection: requires 2 arguments")
		}
		a, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$intersection: %w", err)
		}
		b, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$intersection: %w", err)
		}
		bSet := make(map[any]bool)
		for _, v := range b {
			bSet[v] = true
		}
		var result []any
		for _, v := range a {
			if bSet[v] {
				result = append(result, v)
			}
		}
		if result == nil {
			return []any{}, nil
		}
		return result, nil
	}
}

// Difference returns the CustomFunc for $difference(a, b) — elements in a not in b.
func Difference() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$difference: requires 2 arguments")
		}
		a, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$difference: %w", err)
		}
		b, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$difference: %w", err)
		}
		bSet := make(map[any]bool)
		for _, v := range b {
			bSet[v] = true
		}
		var result []any
		for _, v := range a {
			if !bSet[v] {
				result = append(result, v)
			}
		}
		if result == nil {
			return []any{}, nil
		}
		return result, nil
	}
}

// SymmetricDifference returns the CustomFunc for $symmetricDifference(a, b).
func SymmetricDifference() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$symmetricDifference: requires 2 arguments")
		}
		a, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$symmetricDifference: %w", err)
		}
		b, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$symmetricDifference: %w", err)
		}
		aSet := make(map[any]bool)
		for _, v := range a {
			aSet[v] = true
		}
		bSet := make(map[any]bool)
		for _, v := range b {
			bSet[v] = true
		}
		var result []any
		for _, v := range a {
			if !bSet[v] {
				result = append(result, v)
			}
		}
		for _, v := range b {
			if !aSet[v] {
				result = append(result, v)
			}
		}
		if result == nil {
			return []any{}, nil
		}
		return result, nil
	}
}

// Range returns the CustomFunc for $range(start, end [, step]).
// The range is exclusive of end. Maximum 100,000 items.
func Range() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$range: requires at least 2 arguments")
		}
		start, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$range: %w", err)
		}
		end, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$range: %w", err)
		}
		step := 1.0
		if len(args) >= 3 && args[2] != nil {
			s, err := extutil.ToFloat(args[2])
			if err != nil {
				return nil, fmt.Errorf("$range: %w", err)
			}
			if s == 0 {
				return nil, fmt.Errorf("$range: step cannot be zero")
			}
			step = s
		}
		var result []any
		for v := start; (step > 0 && v < end) || (step < 0 && v > end); v += step {
			if len(result) >= maxRangeItems {
				return nil, fmt.Errorf("$range: exceeded maximum of %d items", maxRangeItems)
			}
			rounded := math.Round(v*1e10) / 1e10
			result = append(result, rounded)
		}
		if result == nil {
			return []any{}, nil
		}
		return result, nil
	}
}

// ZipLongest returns the CustomFunc for $zipLongest(a, b [, fill]).
func ZipLongest() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$zipLongest: requires at least 2 arguments")
		}
		a, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$zipLongest: %w", err)
		}
		b, err := extutil.ToArray(args[1])
		if err != nil {
			return nil, fmt.Errorf("$zipLongest: %w", err)
		}
		var fill any
		if len(args) >= 3 {
			fill = args[2]
		}
		length := len(a)
		if len(b) > length {
			length = len(b)
		}
		result := make([]any, length)
		for i := 0; i < length; i++ {
			var va, vb any = fill, fill
			if i < len(a) {
				va = a[i]
			}
			if i < len(b) {
				vb = b[i]
			}
			result[i] = []any{va, vb}
		}
		return result, nil
	}
}

// Window returns the CustomFunc for $window(array, size [, step]).
func Window() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$window: requires at least 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$window: %w", err)
		}
		size, ok := extutil.ToInt(args[1])
		if !ok || size <= 0 {
			return nil, fmt.Errorf("$window: size must be a positive number")
		}
		step := 1
		if len(args) >= 3 && args[2] != nil {
			s, ok := extutil.ToInt(args[2])
			if !ok || s <= 0 {
				return nil, fmt.Errorf("$window: step must be a positive number")
			}
			step = s
		}
		var result []any
		for i := 0; i+size <= len(arr); i += step {
			window := make([]any, size)
			copy(window, arr[i:i+size])
			result = append(result, window)
		}
		if result == nil {
			return []any{}, nil
		}
		return result, nil
	}
}

// flattenArray recursively flattens nested arrays up to the given depth.
// depth < 0 means unlimited.
func flattenArray(arr []any, depth int) []any {
	var result []any
	for _, v := range arr {
		sub, ok := v.([]any)
		if ok && depth != 0 {
			next := depth
			if depth > 0 {
				next = depth - 1
			}
			result = append(result, flattenArray(sub, next)...)
		} else {
			result = append(result, v)
		}
	}
	return result
}

// normaliseIndex converts a possibly-negative index to a valid array index.
func normaliseIndex(idx, length int) int {
	if idx < 0 {
		idx = length + idx
	}
	if idx < 0 {
		return 0
	}
	if idx > length {
		return length
	}
	return idx
}

// mustInt converts any numeric value to int, returning 0 on failure.
func mustInt(v any) int {
	n, _ := extutil.ToInt(v)
	return n
}
