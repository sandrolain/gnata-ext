package middleware_test

import (
	"context"
	"log/slog"
	"os"
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
