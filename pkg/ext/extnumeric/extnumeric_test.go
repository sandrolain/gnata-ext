package extnumeric_test

import (
	"math"
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
)

func call(fn func() interface { /* placeholder */
}, args ...any) (any, error) {
	return nil, nil
}

func callFunc(f interface{}, args []any) (any, error) {
	type cf interface {
		// gnata.CustomFunc compatible call
	}
	_ = cf(nil)
	return nil, nil
}

func invoke(f func([]any, any) (any, error), args ...any) (any, error) {
	return f(args, nil)
}

func TestLog(t *testing.T) {
	f := extnumeric.Log()
	got, err := invoke(f, 1.0)
	if err != nil || math.Abs(got.(float64)) > 1e-9 {
		t.Errorf("log(1): got %v, %v", got, err)
	}
	got, err = invoke(f, 100.0, 10.0)
	if err != nil || math.Abs(got.(float64)-2.0) > 1e-9 {
		t.Errorf("log(100, 10): got %v, %v", got, err)
	}
	_, err = invoke(f, -1.0)
	if err == nil {
		t.Error("log(-1): expected error")
	}
}

func TestSign(t *testing.T) {
	f := extnumeric.Sign()
	cases := []struct {
		in   float64
		want float64
	}{
		{-5, -1}, {0, 0}, {3, 1},
	}
	for _, c := range cases {
		got, err := invoke(f, c.in)
		if err != nil || got.(float64) != c.want {
			t.Errorf("sign(%v): got %v, %v; want %v", c.in, got, err, c.want)
		}
	}
}

func TestTrunc(t *testing.T) {
	f := extnumeric.Trunc()
	got, err := invoke(f, 3.9)
	if err != nil || got.(float64) != 3.0 {
		t.Errorf("trunc(3.9): got %v, %v", got, err)
	}
	got, err = invoke(f, -3.9)
	if err != nil || got.(float64) != -3.0 {
		t.Errorf("trunc(-3.9): got %v, %v", got, err)
	}
}

func TestClamp(t *testing.T) {
	f := extnumeric.Clamp()
	got, err := invoke(f, 5.0, 1.0, 10.0)
	if err != nil || got.(float64) != 5.0 {
		t.Errorf("clamp(5,1,10): got %v, %v", got, err)
	}
	got, err = invoke(f, -5.0, 1.0, 10.0)
	if err != nil || got.(float64) != 1.0 {
		t.Errorf("clamp(-5,1,10): got %v, %v", got, err)
	}
	got, err = invoke(f, 15.0, 1.0, 10.0)
	if err != nil || got.(float64) != 10.0 {
		t.Errorf("clamp(15,1,10): got %v, %v", got, err)
	}
}

func TestTrig(t *testing.T) {
	all := extnumeric.All()
	sin := all["sin"]
	got, err := sin([]any{0.0}, nil)
	if err != nil || math.Abs(got.(float64)) > 1e-9 {
		t.Errorf("sin(0): got %v, %v", got, err)
	}
	cos := all["cos"]
	got, err = cos([]any{0.0}, nil)
	if err != nil || math.Abs(got.(float64)-1.0) > 1e-9 {
		t.Errorf("cos(0): got %v, %v", got, err)
	}
}

func TestPiE(t *testing.T) {
	all := extnumeric.All()
	pi := all["pi"]
	got, err := pi([]any{}, nil)
	if err != nil || math.Abs(got.(float64)-math.Pi) > 1e-9 {
		t.Errorf("pi(): got %v, %v", got, err)
	}
	e := all["e"]
	got, err = e([]any{}, nil)
	if err != nil || math.Abs(got.(float64)-math.E) > 1e-9 {
		t.Errorf("e(): got %v, %v", got, err)
	}
}

func TestMedian(t *testing.T) {
	f := extnumeric.Median()
	got, err := invoke(f, []any{1.0, 2.0, 3.0})
	if err != nil || got.(float64) != 2.0 {
		t.Errorf("median([1,2,3]): got %v, %v", got, err)
	}
	got, err = invoke(f, []any{1.0, 2.0, 3.0, 4.0})
	if err != nil || got.(float64) != 2.5 {
		t.Errorf("median([1,2,3,4]): got %v, %v", got, err)
	}
}

func TestVarianceStddev(t *testing.T) {
	f := extnumeric.Variance()
	got, err := invoke(f, []any{2.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0})
	if err != nil || math.Abs(got.(float64)-4.0) > 1e-9 {
		t.Errorf("variance: got %v, %v; want 4.0", got, err)
	}
	fs := extnumeric.Stddev()
	got, err = invoke(fs, []any{2.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0})
	if err != nil || math.Abs(got.(float64)-2.0) > 1e-9 {
		t.Errorf("stddev: got %v, %v; want 2.0", got, err)
	}
}

func TestPercentile(t *testing.T) {
	f := extnumeric.Percentile()
	got, err := invoke(f, []any{1.0, 2.0, 3.0, 4.0, 5.0}, 50.0)
	if err != nil || got.(float64) != 3.0 {
		t.Errorf("percentile([1..5], 50): got %v, %v", got, err)
	}
}

func TestMode(t *testing.T) {
	f := extnumeric.Mode()
	got, err := invoke(f, []any{1.0, 2.0, 2.0, 3.0})
	if err != nil || got.(float64) != 2.0 {
		t.Errorf("mode([1,2,2,3]): got %v, %v", got, err)
	}
}

// --- Additional coverage tests ---

func TestAtan2(t *testing.T) {
	f := extnumeric.Atan2()
	got, err := invoke(f, 1.0, 1.0)
	if err != nil {
		t.Fatalf("atan2: unexpected error: %v", err)
	}
	if got.(float64) <= 0 {
		t.Errorf("atan2(1,1): expected positive, got %v", got)
	}
}

func TestAtan2Errors(t *testing.T) {
	f := extnumeric.Atan2()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad", 1.0)
	if err == nil {
		t.Error("expected error for non-numeric y")
	}
	_, err = invoke(f, 1.0, "bad")
	if err == nil {
		t.Error("expected error for non-numeric x")
	}
}

func TestLogErrors(t *testing.T) {
	f := extnumeric.Log()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-numeric")
	}
	_, err = invoke(f, -1.0)
	if err == nil {
		t.Error("expected error for non-positive")
	}
	// two-arg form — bad base
	_, err = invoke(f, 4.0, "bad")
	if err == nil {
		t.Error("expected error for non-numeric base")
	}
	// base <= 0
	_, err = invoke(f, 4.0, -1.0)
	if err == nil {
		t.Error("expected error for non-positive base")
	}
	// base == 1
	_, err = invoke(f, 4.0, 1.0)
	if err == nil {
		t.Error("expected error for base=1")
	}
}

func TestLogBase(t *testing.T) {
	f := extnumeric.Log()
	// log_2(8) = 3
	got, err := invoke(f, 8.0, 2.0)
	if err != nil {
		t.Fatalf("log base 2: %v", err)
	}
	if got.(float64) < 2.9 || got.(float64) > 3.1 {
		t.Errorf("log_2(8): expected ~3, got %v", got)
	}
}

func TestSignErrors(t *testing.T) {
	f := extnumeric.Sign()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-numeric")
	}
}

func TestTruncErrors(t *testing.T) {
	f := extnumeric.Trunc()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-numeric")
	}
}

func TestClampErrors(t *testing.T) {
	f := extnumeric.Clamp()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad", 0.0, 10.0)
	if err == nil {
		t.Error("expected error for non-numeric n")
	}
	_, err = invoke(f, 5.0, "bad", 10.0)
	if err == nil {
		t.Error("expected error for non-numeric min")
	}
	_, err = invoke(f, 5.0, 0.0, "bad")
	if err == nil {
		t.Error("expected error for non-numeric max")
	}
}

func TestMedianErrors(t *testing.T) {
	f := extnumeric.Median()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestVarianceErrors(t *testing.T) {
	f := extnumeric.Variance()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestStddevErrors(t *testing.T) {
	f := extnumeric.Stddev()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestPercentileErrors(t *testing.T) {
	f := extnumeric.Percentile()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad", 50.0)
	if err == nil {
		t.Error("expected error for non-array")
	}
	_, err = invoke(f, []any{1.0, 2.0}, "bad")
	if err == nil {
		t.Error("expected error for non-numeric p")
	}
	_, err = invoke(f, []any{1.0, 2.0}, -1.0)
	if err == nil {
		t.Error("expected error for p < 0")
	}
	_, err = invoke(f, []any{1.0, 2.0}, 101.0)
	if err == nil {
		t.Error("expected error for p > 100")
	}
}

func TestPercentileEdgeCases(t *testing.T) {
	f := extnumeric.Percentile()
	// p=0 → min, p=100 → max
	got, err := invoke(f, []any{3.0, 1.0, 2.0}, 0.0)
	if err != nil || got.(float64) != 1.0 {
		t.Errorf("p=0: got %v, %v", got, err)
	}
	got, err = invoke(f, []any{3.0, 1.0, 2.0}, 100.0)
	if err != nil || got.(float64) != 3.0 {
		t.Errorf("p=100: got %v, %v", got, err)
	}
}

func TestModeErrors(t *testing.T) {
	f := extnumeric.Mode()
	_, err := f([]any{}, nil)
	if err == nil {
		t.Error("expected error for 0 args")
	}
	_, err = invoke(f, "bad")
	if err == nil {
		t.Error("expected error for non-array")
	}
}

func TestMathFunc1Errors(t *testing.T) {
	// sin covers mathFunc1 -- test error paths via All() map
	all := extnumeric.All()
	sinFn := all["sin"]
	_, err := sinFn([]any{}, nil)
	if err == nil {
		t.Error("sin: expected error for 0 args")
	}
	_, err = sinFn([]any{"bad"}, nil)
	if err == nil {
		t.Error("sin: expected error for non-numeric")
	}
}

func TestProduct(t *testing.T) {
	fn := extnumeric.Product()
	cases := []struct {
		arr  []any
		want float64
	}{
		{[]any{float64(2), float64(3), float64(4)}, 24},
		{[]any{float64(1)}, 1},
		{[]any{}, 1},
	}
	for _, c := range cases {
		got, err := fn([]any{c.arr}, nil)
		if err != nil {
			t.Errorf("product %v: unexpected error: %v", c.arr, err)
		}
		if got != c.want {
			t.Errorf("product %v: got %v, want %v", c.arr, got, c.want)
		}
	}
}

func TestProductErrors(t *testing.T) {
	fn := extnumeric.Product()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("product: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array"}, nil); err == nil {
		t.Error("product: expected error for non-array")
	}
	if _, err := fn([]any{[]any{"bad"}}, nil); err == nil {
		t.Error("product: expected error for non-numeric element")
	}
}

func TestCumSum(t *testing.T) {
	fn := extnumeric.CumSum()
	got, err := fn([]any{[]any{float64(1), float64(2), float64(3)}}, nil)
	if err != nil {
		t.Errorf("cumSum: unexpected error: %v", err)
	}
	want := []any{float64(1), float64(3), float64(6)}
	arr := got.([]any)
	for i, w := range want {
		if arr[i] != w {
			t.Errorf("cumSum[%d]: got %v, want %v", i, arr[i], w)
		}
	}
}

func TestCumSumErrors(t *testing.T) {
	fn := extnumeric.CumSum()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("cumSum: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array"}, nil); err == nil {
		t.Error("cumSum: expected error for non-array")
	}
}

func TestInRange(t *testing.T) {
	fn := extnumeric.InRange()
	cases := []struct {
		args []any
		want bool
	}{
		{[]any{float64(5), float64(1), float64(10)}, true},
		{[]any{float64(0), float64(1), float64(10)}, false},
		{[]any{float64(10), float64(1), float64(10)}, true},
		{[]any{float64(11), float64(1), float64(10)}, false},
	}
	for _, c := range cases {
		got, err := fn(c.args, nil)
		if err != nil {
			t.Errorf("inRange %v: unexpected error: %v", c.args, err)
		}
		if got != c.want {
			t.Errorf("inRange %v: got %v, want %v", c.args, got, c.want)
		}
	}
}

func TestInRangeErrors(t *testing.T) {
	fn := extnumeric.InRange()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("inRange: expected error for 0 args")
	}
	if _, err := fn([]any{"bad", float64(1), float64(10)}, nil); err == nil {
		t.Error("inRange: expected error for non-numeric n")
	}
	if _, err := fn([]any{float64(5), "bad", float64(10)}, nil); err == nil {
		t.Error("inRange: expected error for non-numeric min")
	}
	if _, err := fn([]any{float64(5), float64(1), "bad"}, nil); err == nil {
		t.Error("inRange: expected error for non-numeric max")
	}
}

func TestRoundTo(t *testing.T) {
	fn := extnumeric.RoundTo()
	cases := []struct {
		n, places float64
		want      float64
	}{
		{3.14159, 2, 3.14},
		{2.555, 2, 2.56},
		{100.0, 0, 100.0},
	}
	for _, c := range cases {
		got, err := fn([]any{c.n, c.places}, nil)
		if err != nil {
			t.Errorf("roundTo(%v,%v): unexpected error: %v", c.n, c.places, err)
		}
		if got != c.want {
			t.Errorf("roundTo(%v,%v): got %v, want %v", c.n, c.places, got, c.want)
		}
	}
}

func TestRoundToErrors(t *testing.T) {
	fn := extnumeric.RoundTo()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("roundTo: expected error for 0 args")
	}
	if _, err := fn([]any{"bad", float64(2)}, nil); err == nil {
		t.Error("roundTo: expected error for non-numeric n")
	}
	if _, err := fn([]any{float64(1.5), "bad"}, nil); err == nil {
		t.Error("roundTo: expected error for non-integer places")
	}
}

func TestNormalize(t *testing.T) {
	fn := extnumeric.Normalize()
	got, err := fn([]any{[]any{float64(0), float64(5), float64(10)}}, nil)
	if err != nil {
		t.Errorf("normalize: unexpected error: %v", err)
	}
	arr := got.([]any)
	want := []float64{0, 0.5, 1}
	for i, w := range want {
		if arr[i] != w {
			t.Errorf("normalize[%d]: got %v, want %v", i, arr[i], w)
		}
	}
	// all same value → all 0
	got2, _ := fn([]any{[]any{float64(5), float64(5)}}, nil)
	arr2 := got2.([]any)
	for _, v := range arr2 {
		if v != 0.0 {
			t.Errorf("normalize same values: got %v, want 0", v)
		}
	}
	// empty array
	got3, _ := fn([]any{[]any{}}, nil)
	arr3 := got3.([]any)
	if len(arr3) != 0 {
		t.Errorf("normalize empty: got %v", arr3)
	}
}

func TestNormalizeErrors(t *testing.T) {
	fn := extnumeric.Normalize()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("normalize: expected error for 0 args")
	}
	if _, err := fn([]any{"not-array"}, nil); err == nil {
		t.Error("normalize: expected error for non-array")
	}
}

func TestInterpolate(t *testing.T) {
	fn := extnumeric.Interpolate()
	cases := []struct {
		a, b, t, want float64
	}{
		{0, 10, 0.5, 5},
		{0, 10, 0, 0},
		{0, 10, 1, 10},
	}
	for _, c := range cases {
		got, err := fn([]any{c.a, c.b, c.t}, nil)
		if err != nil {
			t.Errorf("interpolate: unexpected error: %v", err)
		}
		if got != c.want {
			t.Errorf("interpolate(%v,%v,%v): got %v, want %v", c.a, c.b, c.t, got, c.want)
		}
	}
}

func TestInterpolateErrors(t *testing.T) {
	fn := extnumeric.Interpolate()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("interpolate: expected error for 0 args")
	}
	if _, err := fn([]any{"bad", float64(10), float64(0.5)}, nil); err == nil {
		t.Error("interpolate: expected error for non-numeric a")
	}
	if _, err := fn([]any{float64(0), "bad", float64(0.5)}, nil); err == nil {
		t.Error("interpolate: expected error for non-numeric b")
	}
	if _, err := fn([]any{float64(0), float64(10), "bad"}, nil); err == nil {
		t.Error("interpolate: expected error for non-numeric t")
	}
}

func TestGCD(t *testing.T) {
	fn := extnumeric.GCD()
	cases := []struct {
		a, b int
		want float64
	}{
		{12, 8, 4},
		{7, 3, 1},
		{15, 25, 5},
	}
	for _, c := range cases {
		got, err := fn([]any{float64(c.a), float64(c.b)}, nil)
		if err != nil {
			t.Errorf("gcd(%v,%v): unexpected error: %v", c.a, c.b, err)
		}
		if got != c.want {
			t.Errorf("gcd(%v,%v): got %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestGCDErrors(t *testing.T) {
	fn := extnumeric.GCD()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("gcd: expected error for 0 args")
	}
	if _, err := fn([]any{"bad", float64(3)}, nil); err == nil {
		t.Error("gcd: expected error for non-integer a")
	}
	if _, err := fn([]any{float64(3), "bad"}, nil); err == nil {
		t.Error("gcd: expected error for non-integer b")
	}
}

func TestLCM(t *testing.T) {
	fn := extnumeric.LCM()
	cases := []struct {
		a, b int
		want float64
	}{
		{4, 6, 12},
		{3, 7, 21},
		{0, 5, 0},
	}
	for _, c := range cases {
		got, err := fn([]any{float64(c.a), float64(c.b)}, nil)
		if err != nil {
			t.Errorf("lcm(%v,%v): unexpected error: %v", c.a, c.b, err)
		}
		if got != c.want {
			t.Errorf("lcm(%v,%v): got %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestLCMErrors(t *testing.T) {
	fn := extnumeric.LCM()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("lcm: expected error for 0 args")
	}
	if _, err := fn([]any{"bad", float64(6)}, nil); err == nil {
		t.Error("lcm: expected error for non-integer a")
	}
	if _, err := fn([]any{float64(4), "bad"}, nil); err == nil {
		t.Error("lcm: expected error for non-integer b")
	}
}

func TestIsPrime(t *testing.T) {
	fn := extnumeric.IsPrime()
	cases := []struct {
		n    int
		want bool
	}{
		{2, true}, {3, true}, {4, false}, {7, true},
		{1, false}, {0, false}, {-1, false}, {17, true},
	}
	for _, c := range cases {
		got, err := fn([]any{float64(c.n)}, nil)
		if err != nil {
			t.Errorf("isPrime(%v): unexpected error: %v", c.n, err)
		}
		if got != c.want {
			t.Errorf("isPrime(%v): got %v, want %v", c.n, got, c.want)
		}
	}
}

func TestIsPrimeErrors(t *testing.T) {
	fn := extnumeric.IsPrime()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("isPrime: expected error for 0 args")
	}
	if _, err := fn([]any{"bad"}, nil); err == nil {
		t.Error("isPrime: expected error for non-integer")
	}
}

func TestFactorial(t *testing.T) {
	fn := extnumeric.Factorial()
	cases := []struct {
		n    int
		want float64
	}{
		{0, 1}, {1, 1}, {5, 120}, {10, 3628800}, {20, 2432902008176640000},
	}
	for _, c := range cases {
		got, err := fn([]any{float64(c.n)}, nil)
		if err != nil {
			t.Errorf("factorial(%v): unexpected error: %v", c.n, err)
		}
		if got != c.want {
			t.Errorf("factorial(%v): got %v, want %v", c.n, got, c.want)
		}
	}
}

func TestFactorialErrors(t *testing.T) {
	fn := extnumeric.Factorial()
	if _, err := fn([]any{}, nil); err == nil {
		t.Error("factorial: expected error for 0 args")
	}
	if _, err := fn([]any{"bad"}, nil); err == nil {
		t.Error("factorial: expected error for non-integer")
	}
	if _, err := fn([]any{float64(-1)}, nil); err == nil {
		t.Error("factorial: expected error for negative")
	}
	if _, err := fn([]any{float64(21)}, nil); err == nil {
		t.Error("factorial: expected error for too large")
	}
}
