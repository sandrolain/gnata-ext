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
//   - $product(array)           – product of all numbers
//   - $cumSum(array)            – cumulative sum array
//   - $inRange(n, min, max)     – true if min <= n <= max
//   - $roundTo(n, places)       – round to specified decimal places
//   - $normalize(array)         – min-max normalize to [0,1]
//   - $interpolate(a, b, t)     – linear interpolation
//   - $gcd(a, b)                – greatest common divisor
//   - $lcm(a, b)                – least common multiple
//   - $isPrime(n)               – true if n is prime
//   - $factorial(n)             – n!
package extnumeric

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extutil"
)

// All returns a map of all extended numeric functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"log":         Log(),
		"sign":        Sign(),
		"trunc":       Trunc(),
		"clamp":       Clamp(),
		"sin":         mathFunc1("sin", math.Sin),
		"cos":         mathFunc1("cos", math.Cos),
		"tan":         mathFunc1("tan", math.Tan),
		"asin":        mathFunc1("asin", math.Asin),
		"acos":        mathFunc1("acos", math.Acos),
		"atan":        mathFunc1("atan", math.Atan),
		"atan2":       Atan2(),
		"pi":          Pi(),
		"e":           E(),
		"median":      Median(),
		"variance":    Variance(),
		"stddev":      Stddev(),
		"percentile":  Percentile(),
		"mode":        Mode(),
		"product":     Product(),
		"cumSum":      CumSum(),
		"inRange":     InRange(),
		"roundTo":     RoundTo(),
		"normalize":   Normalize(),
		"interpolate": Interpolate(),
		"gcd":         GCD(),
		"lcm":         LCM(),
		"isPrime":     IsPrime(),
		"factorial":   Factorial(),
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

// Product returns the CustomFunc for $product(array).
// Returns the product of all numbers in the array.
func Product() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$product: requires 1 argument")
		}
		nums, err := toFloatSlice("$product", args[0])
		if err != nil {
			return nil, err
		}
		result := 1.0
		for _, n := range nums {
			result *= n
		}
		return result, nil
	}
}

// CumSum returns the CustomFunc for $cumSum(array).
// Returns an array where each element is the cumulative sum up to that index.
func CumSum() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$cumSum: requires 1 argument")
		}
		nums, err := toFloatSlice("$cumSum", args[0])
		if err != nil {
			return nil, err
		}
		result := make([]any, len(nums))
		sum := 0.0
		for i, n := range nums {
			sum += n
			result[i] = sum
		}
		return result, nil
	}
}

// InRange returns the CustomFunc for $inRange(n, min, max).
// Returns true if min <= n <= max.
func InRange() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$inRange: requires 3 arguments")
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$inRange: %w", err)
		}
		min, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$inRange: %w", err)
		}
		max, err := extutil.ToFloat(args[2])
		if err != nil {
			return nil, fmt.Errorf("$inRange: %w", err)
		}
		return n >= min && n <= max, nil
	}
}

// RoundTo returns the CustomFunc for $roundTo(n, places).
// Rounds n to the specified number of decimal places.
func RoundTo() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$roundTo: requires 2 arguments")
		}
		n, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$roundTo: %w", err)
		}
		places, ok := extutil.ToInt(args[1])
		if !ok {
			return nil, fmt.Errorf("$roundTo: places must be an integer")
		}
		factor := math.Pow(10, float64(places))
		return math.Round(n*factor) / factor, nil
	}
}

// Normalize returns the CustomFunc for $normalize(array).
// Returns the array scaled to [0, 1] using min-max normalization.
func Normalize() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$normalize: requires 1 argument")
		}
		nums, err := toFloatSlice("$normalize", args[0])
		if err != nil {
			return nil, err
		}
		if len(nums) == 0 {
			return []any{}, nil
		}
		min_, max_ := nums[0], nums[0]
		for _, n := range nums[1:] {
			if n < min_ {
				min_ = n
			}
			if n > max_ {
				max_ = n
			}
		}
		result := make([]any, len(nums))
		if min_ == max_ {
			for i := range nums {
				result[i] = 0.0
			}
			return result, nil
		}
		for i, n := range nums {
			result[i] = (n - min_) / (max_ - min_)
		}
		return result, nil
	}
}

// Interpolate returns the CustomFunc for $interpolate(a, b, t).
// Returns a + t*(b-a), i.e., linear interpolation between a and b at t.
func Interpolate() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("$interpolate: requires 3 arguments")
		}
		a, err := extutil.ToFloat(args[0])
		if err != nil {
			return nil, fmt.Errorf("$interpolate: %w", err)
		}
		b, err := extutil.ToFloat(args[1])
		if err != nil {
			return nil, fmt.Errorf("$interpolate: %w", err)
		}
		t, err := extutil.ToFloat(args[2])
		if err != nil {
			return nil, fmt.Errorf("$interpolate: %w", err)
		}
		return a + t*(b-a), nil
	}
}

// GCD returns the CustomFunc for $gcd(a, b).
// Returns the greatest common divisor of a and b.
func GCD() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$gcd: requires 2 arguments")
		}
		a, ok1 := extutil.ToInt(args[0])
		b, ok2 := extutil.ToInt(args[1])
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$gcd: arguments must be integers")
		}
		ba := new(big.Int).SetInt64(int64(a))
		bb := new(big.Int).SetInt64(int64(b))
		return float64(new(big.Int).GCD(nil, nil, ba, bb).Int64()), nil
	}
}

// LCM returns the CustomFunc for $lcm(a, b).
// Returns the least common multiple of a and b.
func LCM() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("$lcm: requires 2 arguments")
		}
		a, ok1 := extutil.ToInt(args[0])
		b, ok2 := extutil.ToInt(args[1])
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("$lcm: arguments must be integers")
		}
		if a == 0 || b == 0 {
			return float64(0), nil
		}
		ba := new(big.Int).SetInt64(int64(a))
		bb := new(big.Int).SetInt64(int64(b))
		g := new(big.Int).GCD(nil, nil, ba, bb)
		prod := new(big.Int).Mul(ba, bb)
		if prod.Sign() < 0 {
			prod.Neg(prod)
		}
		return float64(new(big.Int).Div(prod, g).Int64()), nil
	}
}

// IsPrime returns the CustomFunc for $isPrime(n).
// Returns true if n is a prime number (n >= 2).
func IsPrime() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$isPrime: requires 1 argument")
		}
		n, ok := extutil.ToInt(args[0])
		if !ok {
			return nil, fmt.Errorf("$isPrime: argument must be an integer")
		}
		if n < 2 {
			return false, nil
		}
		if n == 2 {
			return true, nil
		}
		if n%2 == 0 {
			return false, nil
		}
		for i := 3; i*i <= n; i += 2 {
			if n%i == 0 {
				return false, nil
			}
		}
		return true, nil
	}
}

// Factorial returns the CustomFunc for $factorial(n).
// Returns n! for non-negative integers up to 20.
func Factorial() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("$factorial: requires 1 argument")
		}
		n, ok := extutil.ToInt(args[0])
		if !ok || n < 0 {
			return nil, fmt.Errorf("$factorial: argument must be a non-negative integer")
		}
		if n > 20 {
			return nil, fmt.Errorf("$factorial: argument too large (max 20)")
		}
		result := 1
		for i := 2; i <= n; i++ {
			result *= i
		}
		return float64(result), nil
	}
}

// toFloatSlice converts an array argument to []float64.
func toFloatSlice(fn string, v any) ([]float64, error) {
	arr, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: argument must be an array", fn)
	}
	result := make([]float64, len(arr))
	for i, x := range arr {
		n, err := extutil.ToFloat(x)
		if err != nil {
			return nil, fmt.Errorf("%s: element %d: %w", fn, i, err)
		}
		result[i] = n
	}
	return result, nil
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
