// Package extnumeric provides extended numeric/math functions for gnata beyond
// the official JSONata 2.x specification.
//
// Functions
//
//   - $log(n [, base])          – logarithm (natural or specified base)
//   - $sign(n)                  – -1, 0, or 1
//   - $trunc(n)                 – truncate toward zero
//   - $clamp(n, min, max)       – clamp n between min and max
//   - $sin(n) / $cos(n) / $tan(n)            – trig functions
//   - $asin(n) / $acos(n) / $atan(n)         – inverse trig
//   - $atan2(y, x)              – two-argument arctangent
//   - $pi()                     – π constant
//   - $e()                      – Euler's number
//   - $median(array)            – statistical median
//   - $variance(array)          – population variance
//   - $stddev(array)            – population standard deviation
//   - $percentile(array, p)     – p-th percentile (0–100)
//   - $mode(array)              – most frequent value(s)
package extnumeric

import (
	"fmt"
	"math"
	"sort"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extended numeric functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"log":        Log(),
		"sign":       Sign(),
		"trunc":      Trunc(),
		"clamp":      Clamp(),
		"sin":        mathFunc1("sin", math.Sin),
		"cos":        mathFunc1("cos", math.Cos),
		"tan":        mathFunc1("tan", math.Tan),
		"asin":       mathFunc1("asin", math.Asin),
		"acos":       mathFunc1("acos", math.Acos),
		"atan":       mathFunc1("atan", math.Atan),
		"atan2":      Atan2(),
		"pi":         Pi(),
		"e":          E(),
		"median":     Median(),
		"variance":   Variance(),
		"stddev":     Stddev(),
		"percentile": Percentile(),
		"mode":       Mode(),
	}
}

// Log returns the CustomFunc for $log(n [, base]).
func Log() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$log: requires at least 1 argument")
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$log: %w", err)
		}
		if n <= 0 {
			return nil, fmt.Errorf("$log: argument must be positive")
		}
		if len(args) >= 2 && args[1] != nil {
			base, err := extutil.ToFloat(args[1])
			if err != nil {
				return nil, fmt.Errorf("$log: %w", err)
			}
			if base <= 0 || base == 1 {
				return nil, fmt.Errorf("$log: base must be positive and not 1")
			}
			return math.Log(n) / math.Log(base), nil
		}
		return math.Log(n), nil
	}
}

// Sign returns the CustomFunc for $sign(n).
func Sign() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$sign: requires 1 argument")
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$sign: %w", err)
		}
		switch {
		case n < 0:
			return float64(-1), nil
		case n > 0:
			return float64(1), nil
		default:
			return float64(0), nil
		}
	}
}

// Trunc returns the CustomFunc for $trunc(n).
func Trunc() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$trunc: requires 1 argument")
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$trunc: %w", err)
		}
		return math.Trunc(n), nil
	}
}

// Clamp returns the CustomFunc for $clamp(n, min, max).
func Clamp() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$clamp: requires 3 arguments")
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$clamp: %w", err)
		}
		minV, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$clamp: %w", err)
		}
		maxV, err := extutil.ToFloat(args[2])
		if err != nil {
			return nil, fmt.Errorf("$clamp: %w", err)
		}
		if n < minV {
			return minV, nil
		}
		if n > maxV {
			return maxV, nil
		}
		return n, nil
	}
}

// Atan2 returns the CustomFunc for $atan2(y, x).
func Atan2() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$atan2: requires 2 arguments")
		}
		y, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$atan2: %w", err)
		}
		x, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$atan2: %w", err)
		}
		return math.Atan2(y, x), nil
	}
}

// Pi returns the CustomFunc for $pi().
func Pi() gnata.CustomFunc {
	return func(_ []any, _ any) (any, error) {
		return math.Pi, nil
	}
}

// E returns the CustomFunc for $e().
func E() gnata.CustomFunc {
	return func(_ []any, _ any) (any, error) {
		return math.E, nil
	}
}

// Median returns the CustomFunc for $median(array).
func Median() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$median: requires 1 argument")
		}
		nums, err := extutil.ToFloatSlice(args[0])
		if err != nil {
			return nil, fmt.Errorf("$median: %w", err)
		}
		if len(nums) == 0 {
			return nil, nil
		}
		sorted := make([]float64, len(nums))
		copy(sorted, nums)
		sort.Float64s(sorted)
		mid := len(sorted) / 2
		if len(sorted)%2 == 0 {
			return (sorted[mid-1] + sorted[mid]) / 2, nil
		}
		return sorted[mid], nil
	}
}

// Variance returns the CustomFunc for $variance(array).
func Variance() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$variance: requires 1 argument")
		}
		nums, err := extutil.ToFloatSlice(args[0])
		if err != nil {
			return nil, fmt.Errorf("$variance: %w", err)
		}
		return calcVariance(nums)
	}
}

// Stddev returns the CustomFunc for $stddev(array).
func Stddev() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$stddev: requires 1 argument")
		}
		nums, err := extutil.ToFloatSlice(args[0])
		if err != nil {
			return nil, fmt.Errorf("$stddev: %w", err)
		}
		v, err := calcVariance(nums)
		if err != nil || v == nil {
			return v, err
		}
		return math.Sqrt(v.(float64)), nil
	}
}

// Percentile returns the CustomFunc for $percentile(array, p).
// p is in range [0, 100].
func Percentile() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$percentile: requires 2 arguments")
		}
		nums, err := extutil.ToFloatSlice(args[0])
		if err != nil {
			return nil, fmt.Errorf("$percentile: %w", err)
		}
		p, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$percentile: %w", err)
		}
		if p < 0 || p > 100 {
			return nil, fmt.Errorf("$percentile: p must be between 0 and 100")
		}
		if len(nums) == 0 {
			return nil, nil
		}
		sorted := make([]float64, len(nums))
		copy(sorted, nums)
		sort.Float64s(sorted)
		idx := p / 100 * float64(len(sorted)-1)
		lo := int(math.Floor(idx))
		hi := int(math.Ceil(idx))
		if lo == hi {
			return sorted[lo], nil
		}
		frac := idx - float64(lo)
		return sorted[lo]*(1-frac) + sorted[hi]*frac, nil
	}
}

// Mode returns the CustomFunc for $mode(array).
// Returns the most frequent value; multiple values if there is a tie.
func Mode() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$mode: requires 1 argument")
		}
		nums, err := extutil.ToFloatSlice(args[0])
		if err != nil {
			return nil, fmt.Errorf("$mode: %w", err)
		}
		if len(nums) == 0 {
			return nil, nil
		}
		counts := make(map[float64]int)
		for _, n := range nums {
			counts[n]++
		}
		maxCount := 0
		for _, c := range counts {
			if c > maxCount {
				maxCount = c
			}
		}
		var modes []any
		for _, n := range nums {
			if counts[n] == maxCount {
				found := false
				for _, m := range modes {
					if m.(float64) == n {
						found = true
						break
					}
				}
				if !found {
					modes = append(modes, n)
				}
			}
		}
		if len(modes) == 1 {
			return modes[0], nil
		}
		return modes, nil
	}
}

// mathFunc1 creates a CustomFunc wrapping a single-argument math function.
func mathFunc1(name string, fn func(float64) float64) gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$%s: requires 1 argument", name)
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$%s: %w", name, err)
		}
		return fn(n), nil
	}
}

func calcVariance(nums []float64) (any, error) {
	if len(nums) == 0 {
		return nil, nil
	}
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	mean := sum / float64(len(nums))
	varSum := 0.0
	for _, n := range nums {
		d := n - mean
		varSum += d * d
	}
	return varSum / float64(len(nums)), nil
}
