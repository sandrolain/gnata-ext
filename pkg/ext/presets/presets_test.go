package presets_test

import (
	"context"
	"testing"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/presets"
)

func eval(t *testing.T, funcs map[string]gnata.CustomFunc, expr string, data any) any {
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

func TestDataEnv(t *testing.T) {
	// array function
	got := eval(t, presets.DataEnv(), `$first([5, 10])`, nil)
	if got != float64(5) {
		t.Fatalf("DataEnv: want 5, got %v", got)
	}
	// numeric function
	got2 := eval(t, presets.DataEnv(), `$sign(-1)`, nil)
	if got2 != float64(-1) {
		t.Fatalf("DataEnv numeric: want -1, got %v", got2)
	}
}

func TestTextEnv(t *testing.T) {
	got := eval(t, presets.TextEnv(), `$capitalize("hello")`, nil)
	if got != "Hello" {
		t.Fatalf("TextEnv: want Hello, got %v", got)
	}
}

func TestSecureEnv_NoUUID(t *testing.T) {
	compiled, err := gnata.Compile(`$uuid()`)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	env := gnata.NewCustomEnv(presets.SecureEnv())
	_, err = compiled.EvalWithCustomFuncs(context.Background(), nil, env)
	if err == nil {
		t.Fatal("SecureEnv: $uuid() should not be available")
	}
}

func TestSecureEnv_HasOtherFuncs(t *testing.T) {
	got := eval(t, presets.SecureEnv(), `$first([1, 2, 3])`, nil)
	if got != float64(1) {
		t.Fatalf("SecureEnv: want 1, got %v", got)
	}
}

func TestAnalyticsEnv(t *testing.T) {
	got := eval(t, presets.AnalyticsEnv(), `$median([1, 2, 3, 4, 5])`, nil)
	if got != float64(3) {
		t.Fatalf("AnalyticsEnv: want 3, got %v", got)
	}
}
