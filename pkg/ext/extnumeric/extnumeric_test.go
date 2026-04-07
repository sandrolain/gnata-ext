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
