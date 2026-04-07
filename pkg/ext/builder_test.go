package ext_test

import (
	"context"
	"testing"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
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

func TestEnvBuilder_WithAllFuncs(t *testing.T) {
	funcs := ext.NewEnvBuilder().WithAllFuncs().Build()
	got := eval(t, funcs, `$first([10, 20, 30])`, nil)
	if got != float64(10) {
		t.Fatalf("want 10, got %v", got)
	}
}

func TestEnvBuilder_SelectivePackages(t *testing.T) {
	funcs := ext.NewEnvBuilder().WithStringFuncs().WithNumericFuncs().Build()

	// string function works
	got := eval(t, funcs, `$capitalize("hello")`, nil)
	if got != "Hello" {
		t.Fatalf("want Hello, got %v", got)
	}

	// numeric function works
	got2 := eval(t, funcs, `$sign(-5)`, nil)
	if got2 != float64(-1) {
		t.Fatalf("want -1, got %v", got2)
	}
}

func TestEnvBuilder_WithFunc(t *testing.T) {
	funcs := ext.NewEnvBuilder().
		WithFunc("double", func(args []any, _ any) (any, error) {
			n := args[0].(float64)
			return n * 2, nil
		}).Build()
	got := eval(t, funcs, `$double(7)`, nil)
	if got != float64(14) {
		t.Fatalf("want 14, got %v", got)
	}
}

func TestEnvBuilder_Without(t *testing.T) {
	b := ext.NewEnvBuilder().WithCryptoFuncs().Without("uuid")
	funcs := b.Funcs()
	if _, ok := funcs["uuid"]; ok {
		t.Fatal("uuid should have been removed")
	}
	if _, ok := funcs["hash"]; !ok {
		t.Fatal("hash should still be present")
	}
}

func TestEnvBuilder_WithPackage(t *testing.T) {
	pkg := map[string]gnata.CustomFunc{
		"triple": func(args []any, _ any) (any, error) {
			return args[0].(float64) * 3, nil
		},
	}
	env := ext.NewEnvBuilder().WithPackage(pkg).Build()
	got := eval(t, env, `$triple(4)`, nil)
	if got != float64(12) {
		t.Fatalf("want 12, got %v", got)
	}
}

func TestEnvBuilder_Funcs_IsACopy(t *testing.T) {
	b := ext.NewEnvBuilder().WithArrayFuncs()
	f1 := b.Funcs()
	f2 := b.Funcs()
	// Mutating f1 must not affect f2
	delete(f1, "first")
	if _, ok := f2["first"]; !ok {
		t.Fatal("Funcs() should return an independent copy")
	}
}
