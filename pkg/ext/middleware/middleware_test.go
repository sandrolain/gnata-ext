package middleware_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
	"github.com/sandrolain/gnata-ext/pkg/ext/middleware"
)

func evalWith(t *testing.T, funcs map[string]gnata.CustomFunc, expr string, data any) any {
	t.Helper()
	compiled, err := gnata.Compile(expr)
	if err != nil {
		t.Fatalf("compile %q: %v", expr, err)
	}
	env := gnata.NewCustomEnv(funcs)
	result, err := compiled.EvalWithCustomFuncs(context.Background(), data, env)
	if err != nil {
		t.Fatalf("eval %q: %v", expr, err)
	}
	return result
}

func TestWithLogging_DoesNotBreakFunctions(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	wrapped := middleware.WithLogging(ext.AllFuncs(), logger)
	got := evalWith(t, wrapped, `$first([1, 2, 3])`, nil)
	if got != float64(1) {
		t.Fatalf("want 1, got %v", got)
	}
}

func TestWithMemoize_CachesResult(t *testing.T) {
	callCount := 0
	base := map[string]gnata.CustomFunc{
		"counted": func(args []any, _ any) (any, error) {
			callCount++
			return args[0], nil
		},
	}
	wrapped := middleware.WithMemoize(base)

	evalWith(t, wrapped, `$counted(42)`, nil)
	evalWith(t, wrapped, `$counted(42)`, nil)

	if callCount != 1 {
		t.Fatalf("expected 1 call (memoized), got %d", callCount)
	}
}

func TestWithMemoize_ExcludedFunctionNotCached(t *testing.T) {
	callCount := 0
	base := map[string]gnata.CustomFunc{
		"nocache": func(args []any, _ any) (any, error) {
			callCount++
			return callCount, nil
		},
	}
	wrapped := middleware.WithMemoize(base, "nocache")

	evalWith(t, wrapped, `$nocache()`, nil)
	evalWith(t, wrapped, `$nocache()`, nil)

	if callCount != 2 {
		t.Fatalf("expected 2 calls (not cached), got %d", callCount)
	}
}

func TestWithRateLimit_DoesNotBreakFunctions(t *testing.T) {
	wrapped := middleware.WithRateLimit(ext.AllFuncs(), 100) // 100 rps — fast enough for tests
	got := evalWith(t, wrapped, `$sign(-3)`, nil)
	if got != float64(-1) {
		t.Fatalf("want -1, got %v", got)
	}
}

func TestWithRateLimit_ZeroRpsNoop(t *testing.T) {
	// rps=0 should be a no-op and not block
	wrapped := middleware.WithRateLimit(ext.AllFuncs(), 0)
	start := time.Now()
	evalWith(t, wrapped, `$sign(1)`, nil)
	if time.Since(start) > time.Second {
		t.Fatal("zero rps should not block")
	}
}

// ---------------------------------------------------------------------------
// WithTimeout
// ---------------------------------------------------------------------------

func TestWithTimeout_CompletesInTime(t *testing.T) {
	base := map[string]gnata.CustomFunc{
		"fast": func(_ []any, _ any) (any, error) { return "ok", nil },
	}
	wrapped := middleware.WithTimeout(base, time.Second)
	got := evalWith(t, wrapped, `$fast()`, nil)
	if got != "ok" {
		t.Fatalf("want 'ok', got %v", got)
	}
}

func TestWithTimeout_ReturnsErrorOnTimeout(t *testing.T) {
	base := map[string]gnata.CustomFunc{
		"slow": func(_ []any, _ any) (any, error) {
			time.Sleep(200 * time.Millisecond)
			return "ok", nil
		},
	}
	wrapped := middleware.WithTimeout(base, 30*time.Millisecond)
	compiled, _ := gnata.Compile(`$slow()`)
	env := gnata.NewCustomEnv(wrapped)
	_, err := compiled.EvalWithCustomFuncs(context.Background(), nil, env)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

// ---------------------------------------------------------------------------
// WithRetry
// ---------------------------------------------------------------------------

func TestWithRetry_RetriesOnError(t *testing.T) {
	attempts := 0
	base := map[string]gnata.CustomFunc{
		"flaky": func(_ []any, _ any) (any, error) {
			attempts++
			if attempts < 3 {
				return nil, fmt.Errorf("transient")
			}
			return "ok", nil
		},
	}
	wrapped := middleware.WithRetry(base, middleware.MaxAttempts(3), middleware.Backoff(time.Millisecond))
	got := evalWith(t, wrapped, `$flaky()`, nil)
	if got != "ok" {
		t.Fatalf("want 'ok', got %v", got)
	}
	if attempts != 3 {
		t.Fatalf("want 3 attempts, got %d", attempts)
	}
}

func TestWithRetry_RespectsMaxAttempts(t *testing.T) {
	attempts := 0
	base := map[string]gnata.CustomFunc{
		"always": func(_ []any, _ any) (any, error) {
			attempts++
			return nil, fmt.Errorf("always fails")
		},
	}
	wrapped := middleware.WithRetry(base, middleware.MaxAttempts(3), middleware.Backoff(time.Millisecond))
	compiled, _ := gnata.Compile(`$always()`)
	env := gnata.NewCustomEnv(wrapped)
	_, err := compiled.EvalWithCustomFuncs(context.Background(), nil, env)
	if err == nil {
		t.Fatal("expected error after max attempts")
	}
	if attempts != 3 {
		t.Fatalf("want exactly 3 attempts, got %d", attempts)
	}
}

func TestWithRetry_SkipsRetryWhenPredicateFalse(t *testing.T) {
	attempts := 0
	base := map[string]gnata.CustomFunc{
		"nope": func(_ []any, _ any) (any, error) {
			attempts++
			return nil, fmt.Errorf("non-transient")
		},
	}
	wrapped := middleware.WithRetry(base,
		middleware.MaxAttempts(5),
		middleware.Backoff(time.Millisecond),
		middleware.RetryOn(func(error) bool { return false }),
	)
	compiled, _ := gnata.Compile(`$nope()`)
	env := gnata.NewCustomEnv(wrapped)
	compiled.EvalWithCustomFuncs(context.Background(), nil, env) //nolint:errcheck
	if attempts != 1 {
		t.Fatalf("want exactly 1 attempt (no retry), got %d", attempts)
	}
}

// ---------------------------------------------------------------------------
// WithCircuitBreaker
// ---------------------------------------------------------------------------

func TestWithCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	calls := 0
	base := map[string]gnata.CustomFunc{
		"broken": func(_ []any, _ any) (any, error) {
			calls++
			return nil, fmt.Errorf("failure")
		},
	}
	wrapped := middleware.WithCircuitBreaker(base,
		middleware.Threshold(3),
		middleware.ResetAfter(time.Hour),
	)
	compiled, _ := gnata.Compile(`$broken()`)
	env := gnata.NewCustomEnv(wrapped)
	ctx := context.Background()

	// First three calls reach the function and open the circuit.
	for i := 0; i < 3; i++ {
		compiled.EvalWithCustomFuncs(ctx, nil, env) //nolint:errcheck
	}
	if calls != 3 {
		t.Fatalf("want 3 calls before open, got %d", calls)
	}

	// Circuit is now open: next call must fail without invoking the function.
	_, err := compiled.EvalWithCustomFuncs(ctx, nil, env)
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
	if calls != 3 {
		t.Fatalf("function must not be called when circuit is open, calls=%d", calls)
	}
}

func TestWithCircuitBreaker_ResetsAfterDuration(t *testing.T) {
	calls := 0
	base := map[string]gnata.CustomFunc{
		"probe": func(_ []any, _ any) (any, error) {
			calls++
			if calls <= 2 {
				return nil, fmt.Errorf("fail")
			}
			return "ok", nil
		},
	}
	wrapped := middleware.WithCircuitBreaker(base,
		middleware.Threshold(2),
		middleware.ResetAfter(50*time.Millisecond),
	)
	compiled, _ := gnata.Compile(`$probe()`)
	env := gnata.NewCustomEnv(wrapped)
	ctx := context.Background()

	// Trigger open.
	compiled.EvalWithCustomFuncs(ctx, nil, env) //nolint:errcheck
	compiled.EvalWithCustomFuncs(ctx, nil, env) //nolint:errcheck

	// Wait for reset window.
	time.Sleep(100 * time.Millisecond)

	got, err := compiled.EvalWithCustomFuncs(ctx, nil, env)
	if err != nil {
		t.Fatalf("expected success after reset, got %v", err)
	}
	if got != "ok" {
		t.Fatalf("want 'ok', got %v", got)
	}
}

func TestWithCircuitBreaker_ErrCircuitOpenWrapped(t *testing.T) {
	base := map[string]gnata.CustomFunc{
		"fail": func(_ []any, _ any) (any, error) { return nil, fmt.Errorf("err") },
	}
	wrapped := middleware.WithCircuitBreaker(base,
		middleware.Threshold(1),
		middleware.ResetAfter(time.Hour),
	)
	compiled, _ := gnata.Compile(`$fail()`)
	env := gnata.NewCustomEnv(wrapped)
	ctx := context.Background()

	compiled.EvalWithCustomFuncs(ctx, nil, env) //nolint:errcheck — opens circuit

	_, err := compiled.EvalWithCustomFuncs(ctx, nil, env)
	if err == nil {
		t.Fatal("expected error")
	}
	// errors.Is works if gnata preserves the error chain; string fallback otherwise.
	if !errors.Is(err, middleware.ErrCircuitOpen) && !strings.Contains(err.Error(), "circuit breaker open") {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// WithAuditLog
// ---------------------------------------------------------------------------

func TestWithAuditLog_WritesRecord(t *testing.T) {
	base := map[string]gnata.CustomFunc{
		"fn": func(_ []any, _ any) (any, error) { return "result", nil },
	}
	var buf bytes.Buffer
	wrapped := middleware.WithAuditLog(base, &buf)
	evalWith(t, wrapped, `$fn()`, nil)

	line := buf.String()
	if !strings.Contains(line, `"func":"fn"`) {
		t.Fatalf("expected audit record with func name, got: %s", line)
	}
	if !strings.Contains(line, `"time"`) || !strings.Contains(line, `"duration"`) {
		t.Fatalf("expected time and duration fields, got: %s", line)
	}
}

func TestWithAuditLog_RecordsError(t *testing.T) {
	base := map[string]gnata.CustomFunc{
		"oops": func(_ []any, _ any) (any, error) { return nil, fmt.Errorf("boom") },
	}
	var buf bytes.Buffer
	wrapped := middleware.WithAuditLog(base, &buf)
	compiled, _ := gnata.Compile(`$oops()`)
	env := gnata.NewCustomEnv(wrapped)
	compiled.EvalWithCustomFuncs(context.Background(), nil, env) //nolint:errcheck

	line := buf.String()
	if !strings.Contains(line, `"error":"boom"`) {
		t.Fatalf("expected error field in audit record, got: %s", line)
	}
}
