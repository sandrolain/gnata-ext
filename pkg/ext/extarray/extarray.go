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
//   - $compact(array)                      – remove falsy values (nil/false/0/"")
//   - $groupByKey(array, key)              – group objects by a field value
//   - $sortBy(array, key [, desc])         – stable sort objects by a field
//   - $uniqueBy(array, key)               – remove duplicate objects by a field
//   - $sumByKey(array, key, valueKey)      – sum valueKey per key group
//   - $countByKey(array, key)             – count elements per key group
//   - $rotate(array, n)                   – rotate array by n positions
//   - $indexof(array, value)              – first index of value, or -1
//   - $transpose(matrix)                  – transpose a 2D array
//   - $adjacentPairs(array)               – array of [a[i], a[i+1]] pairs
package extarray

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

const maxRangeItems = 100000

// All returns a map of all extended array functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"first":               First(),
		"last":                Last(),
		"take":                Take(),
		"skip":                Skip(),
		"slice":               Slice(),
		"flatten":             Flatten(),
		"chunk":               Chunk(),
		"union":               Union(),
		"intersection":        Intersection(),
		"difference":          Difference(),
		"symmetricDifference": SymmetricDifference(),
		"range":               Range(),
		"zipLongest":          ZipLongest(),
		"window":              Window(),
		"compact":             Compact(),
		"groupByKey":          GroupByKey(),
		"sortBy":              SortBy(),
		"uniqueBy":            UniqueBy(),
		"sumByKey":            SumByKey(),
		"countByKey":          CountByKey(),
		"rotate":              Rotate(),
		"indexof":             Indexof(),
		"transpose":           Transpose(),
		"adjacentPairs":       AdjacentPairs(),
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

// Compact returns the CustomFunc for $compact(array).
// Removes falsy values: nil, false, float64(0), "".
func Compact() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$compact: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$compact: %w", err)
		}
		result := make([]any, 0, len(arr))
		for _, v := range arr {
			if !isFalsy(v) {
				result = append(result, v)
			}
		}
		return result, nil
	}
}

// GroupByKey returns the CustomFunc for $groupByKey(array, key).
// Groups an array of objects by the value of a specified key.
func GroupByKey() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$groupByKey: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$groupByKey: %w", err)
		}
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$groupByKey: key must be a string")
		}
		result := map[string]any{}
		order := []string{}
		for _, item := range arr {
			obj, err := extutil.ToObject(item)
			if err != nil {
				return nil, fmt.Errorf("$groupByKey: each element must be an object")
			}
			groupVal := fmt.Sprint(obj[key])
			if _, exists := result[groupVal]; !exists {
				result[groupVal] = []any{}
				order = append(order, groupVal)
			}
			result[groupVal] = append(result[groupVal].([]any), item)
		}
		_ = order
		return result, nil
	}
}

// SortBy returns the CustomFunc for $sortBy(array, key [, desc]).
// Stable-sorts an array of objects by a string field.
func SortBy() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$sortBy: requires at least 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$sortBy: %w", err)
		}
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$sortBy: key must be a string")
		}
		desc := false
		if len(args) >= 3 && args[2] != nil {
			d, ok := args[2].(bool)
			if !ok {
				return nil, fmt.Errorf("$sortBy: desc must be a boolean")
			}
			desc = d
		}
		result := make([]any, len(arr))
		copy(result, arr)
		sort.SliceStable(result, func(i, j int) bool {
			oi, _ := extutil.ToObject(result[i])
			oj, _ := extutil.ToObject(result[j])
			vi := fmt.Sprint(oi[key])
			vj := fmt.Sprint(oj[key])
			if desc {
				return vi > vj
			}
			return vi < vj
		})
		return result, nil
	}
}

// UniqueBy returns the CustomFunc for $uniqueBy(array, key).
// Returns the first element seen for each unique value of key.
func UniqueBy() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$uniqueBy: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$uniqueBy: %w", err)
		}
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$uniqueBy: key must be a string")
		}
		seen := map[string]bool{}
		result := make([]any, 0)
		for _, item := range arr {
			obj, err := extutil.ToObject(item)
			if err != nil {
				return nil, fmt.Errorf("$uniqueBy: each element must be an object")
			}
			k := fmt.Sprint(obj[key])
			if !seen[k] {
				seen[k] = true
				result = append(result, item)
			}
		}
		return result, nil
	}
}

// SumByKey returns the CustomFunc for $sumByKey(array, key, valueKey).
// Returns an object mapping each group (by key) to the sum of valueKey.
func SumByKey() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$sumByKey: requires 3 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$sumByKey: %w", err)
		}
		key, ok1 := args[1].(string)
		valueKey, ok2 := args[2].(string)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$sumByKey: key and valueKey must be strings")
		}
		result := map[string]any{}
		for _, item := range arr {
			obj, err := extutil.ToObject(item)
			if err != nil {
				return nil, fmt.Errorf("$sumByKey: each element must be an object")
			}
			groupVal := fmt.Sprint(obj[key])
			v, err := extutil.ToFloat(obj[valueKey])
			if err != nil {
				v = 0
			}
			if cur, exists := result[groupVal]; exists {
				result[groupVal] = cur.(float64) + v
			} else {
				result[groupVal] = v
			}
		}
		return result, nil
	}
}

// CountByKey returns the CustomFunc for $countByKey(array, key).
// Returns an object mapping each group (by key) to its element count.
func CountByKey() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$countByKey: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$countByKey: %w", err)
		}
		key, ok := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("$countByKey: key must be a string")
		}
		result := map[string]any{}
		for _, item := range arr {
			obj, err := extutil.ToObject(item)
			if err != nil {
				return nil, fmt.Errorf("$countByKey: each element must be an object")
			}
			groupVal := fmt.Sprint(obj[key])
			if cur, exists := result[groupVal]; exists {
				result[groupVal] = cur.(float64) + 1
			} else {
				result[groupVal] = float64(1)
			}
		}
		return result, nil
	}
}

// Rotate returns the CustomFunc for $rotate(array, n).
// Rotates array left by n positions (negative n rotates right).
func Rotate() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$rotate: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$rotate: %w", err)
		}
		n, ok := extutil.ToInt(args[1])
		if !ok {
			return nil, fmt.Errorf("$rotate: n must be an integer")
		}
		if len(arr) == 0 {
			return arr, nil
		}
		l := len(arr)
		n = ((n % l) + l) % l
		result := make([]any, l)
		copy(result, arr[n:])
		copy(result[l-n:], arr[:n])
		return result, nil
	}
}

// Indexof returns the CustomFunc for $indexof(array, value).
// Returns the first index at which value appears, or -1 if not found.
func Indexof() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$indexof: requires 2 arguments")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$indexof: %w", err)
		}
		target := args[1]
		for i, v := range arr {
			if reflect.DeepEqual(v, target) {
				return float64(i), nil
			}
		}
		return float64(-1), nil
	}
}

// Transpose returns the CustomFunc for $transpose(matrix).
// Transposes a 2D array (matrix), swapping rows and columns.
func Transpose() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$transpose: requires 1 argument")
		}
		matrix, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$transpose: %w", err)
		}
		if len(matrix) == 0 {
			return []any{}, nil
		}
		rows := make([][]any, len(matrix))
		maxCols := 0
		for i, row := range matrix {
			r, err := extutil.ToArray(row)
			if err != nil {
				return nil, fmt.Errorf("$transpose: row %d must be an array", i)
			}
			rows[i] = r
			if len(r) > maxCols {
				maxCols = len(r)
			}
		}
		result := make([]any, maxCols)
		for c := 0; c < maxCols; c++ {
			col := make([]any, len(rows))
			for r, row := range rows {
				if c < len(row) {
					col[r] = row[c]
				}
			}
			result[c] = col
		}
		return result, nil
	}
}

// AdjacentPairs returns the CustomFunc for $adjacentPairs(array).
// Returns [[a[0],a[1]], [a[1],a[2]], ...] for consecutive elements.
func AdjacentPairs() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$adjacentPairs: requires 1 argument")
		}
		arr, err := extutil.ToArray(args[0])
		if err != nil {
			return nil, fmt.Errorf("$adjacentPairs: %w", err)
		}
		if len(arr) < 2 {
			return []any{}, nil
		}
		result := make([]any, len(arr)-1)
		for i := 0; i < len(arr)-1; i++ {
			result[i] = []any{arr[i], arr[i+1]}
		}
		return result, nil
	}
}

// isFalsy returns true for nil, false, 0.0, and "".
func isFalsy(v any) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case bool:
		return !t
	case float64:
		return t == 0
	case string:
		return t == ""
	}
	return false
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
