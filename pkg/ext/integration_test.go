package ext_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
	"github.com/sandrolain/gnata-ext/pkg/ext/extarray"
	"github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
	"github.com/sandrolain/gnata-ext/pkg/ext/extobject"
	"github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
)

// evalExpr compiles and evaluates a gnata expression with the given custom functions.
func evalExpr(t *testing.T, funcs map[string]gnata.CustomFunc, expr string, data any) any {
	t.Helper()
	e, err := gnata.Compile(expr)
	if err != nil {
		t.Fatalf("Compile(%q): %v", expr, err)
	}
	env := gnata.NewCustomEnv(funcs)
	ctx := context.Background()
	result, err := e.EvalWithCustomFuncs(ctx, data, env)
	if err != nil {
		t.Fatalf("EvalWithCustomFuncs(%q): %v", expr, err)
	}
	return result
}

// TestIntegration_AllFuncs verifies that AllFuncs can be used end-to-end with gnata.
func TestIntegration_AllFuncs(t *testing.T) {
	funcs := ext.AllFuncs()
	env := gnata.NewCustomEnv(funcs)
	ctx := context.Background()

	expr, err := gnata.Compile(`$camelCase("hello_world_foo")`)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	result, err := expr.EvalWithCustomFuncs(ctx, nil, env)
	if err != nil {
		t.Fatalf("EvalWithCustomFuncs: %v", err)
	}
	if result != "helloWorldFoo" {
		t.Errorf("AllFuncs $camelCase: want %q, got %v", "helloWorldFoo", result)
	}
}

// TestIntegration_SinglePackage verifies that a single sub-package can be wired up in isolation.
func TestIntegration_SinglePackage(t *testing.T) {
	env := gnata.NewCustomEnv(extcrypto.All())
	ctx := context.Background()

	expr, err := gnata.Compile(`$uuid()`)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	result, err := expr.EvalWithCustomFuncs(ctx, nil, env)
	if err != nil {
		t.Fatalf("EvalWithCustomFuncs: %v", err)
	}
	s, ok := result.(string)
	if !ok || len(s) != 36 {
		t.Errorf("$uuid: want 36-char string, got %v (%T)", result, result)
	}
}

// TestIntegration_SelectivePackages verifies that multiple sub-packages can be merged manually.
func TestIntegration_SelectivePackages(t *testing.T) {
	merged := make(map[string]gnata.CustomFunc)
	for k, v := range extarray.All() {
		merged[k] = v
	}
	for k, v := range exttypes.All() {
		merged[k] = v
	}
	for k, v := range extobject.All() {
		merged[k] = v
	}

	env := gnata.NewCustomEnv(merged)
	ctx := context.Background()

	expr, err := gnata.Compile(`$isArray([1,2,3])`)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	result, err := expr.EvalWithCustomFuncs(ctx, nil, env)
	if err != nil {
		t.Fatalf("EvalWithCustomFuncs: %v", err)
	}
	if result != true {
		t.Errorf("$isArray: want true, got %v", result)
	}
}

// TestIntegration_StreamEvaluator verifies integration with gnata StreamEvaluator.
func TestIntegration_StreamEvaluator(t *testing.T) {
	expr, err := gnata.Compile(`$snakeCase(name)`)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	se := gnata.NewStreamEvaluator(
		[]*gnata.Expression{expr},
		gnata.WithCustomFunctions(ext.AllFuncs()),
	)

	data, _ := json.Marshal(map[string]string{"name": "helloWorldFoo"})
	ctx := context.Background()
	result, err := se.EvalOne(ctx, data, "test-schema", 0)
	if err != nil {
		t.Fatalf("EvalOne: %v", err)
	}
	if result != "hello_world_foo" {
		t.Errorf("StreamEvaluator $snakeCase: want %q, got %v", "hello_world_foo", result)
	}
}

// TestIntegration_Expressions is a table-driven end-to-end test covering each package.
func TestIntegration_Expressions(t *testing.T) {
	funcs := ext.AllFuncs()
	env := gnata.NewCustomEnv(funcs)
	ctx := context.Background()

	tests := []struct {
		name string
		expr string
		data any
		want any
	}{
		// extstring
		{name: "camelCase", expr: `$camelCase("hello_world_foo")`, want: "helloWorldFoo"},
		{name: "snakeCase", expr: `$snakeCase("helloWorldFoo")`, want: "hello_world_foo"},
		{name: "kebabCase", expr: `$kebabCase("hello world")`, want: "hello-world"},
		{name: "capitalize", expr: `$capitalize("hello WORLD")`, want: "Hello world"},
		{name: "titleCase", expr: `$titleCase("hello world")`, want: "Hello World"},
		{name: "startsWith true", expr: `$startsWith("gnata-ext", "gnata")`, want: true},
		{name: "endsWith true", expr: `$endsWith("gnata-ext", "ext")`, want: true},
		{name: "repeat", expr: `$repeat("ab", 3)`, want: "ababab"},
		// extnumeric
		{name: "clamp high", expr: `$clamp(15, 1, 10)`, want: float64(10)},
		{name: "clamp low", expr: `$clamp(-5, 1, 10)`, want: float64(1)},
		{name: "trunc", expr: `$trunc(3.9)`, want: float64(3)},
		{name: "sign negative", expr: `$sign(-5)`, want: float64(-1)},
		{name: "sign positive", expr: `$sign(5)`, want: float64(1)},
		{name: "sign zero", expr: `$sign(0)`, want: float64(0)},
		{name: "median odd", expr: `$median([1,2,3,4,5])`, want: float64(3)},
		{name: "median even", expr: `$median([1,2,3,4])`, want: float64(2.5)},
		// exttypes
		{name: "isArray true", expr: `$isArray([1,2,3])`, want: true},
		{name: "isArray false", expr: `$isArray("str")`, want: false},
		{name: "isString true", expr: `$isString("hello")`, want: true},
		{name: "isNumber true", expr: `$isNumber(42)`, want: true},
		{name: "isEmpty empty string", expr: `$isEmpty("")`, want: true},
		{name: "default null", expr: `$default(null, "fallback")`, want: "fallback"},
		{name: "identity", expr: `$identity(42)`, want: float64(42)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, err := gnata.Compile(tc.expr)
			if err != nil {
				t.Fatalf("Compile(%q): %v", tc.expr, err)
			}
			got, err := e.EvalWithCustomFuncs(ctx, tc.data, env)
			if err != nil {
				t.Fatalf("EvalWithCustomFuncs(%q): %v", tc.expr, err)
			}
			if got != tc.want {
				t.Errorf("expr %q: want %v (%T), got %v (%T)", tc.expr, tc.want, tc.want, got, got)
			}
		})
	}
}

// TestIntegration_DataDriven tests expressions that reference input data fields.
func TestIntegration_DataDriven(t *testing.T) {
	funcs := ext.AllFuncs()
	env := gnata.NewCustomEnv(funcs)
	ctx := context.Background()

	t.Run("pick from object", func(t *testing.T) {
		data := map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}
		e, _ := gnata.Compile(`$pick($, ["a", "c"])`)
		result, err := e.EvalWithCustomFuncs(ctx, data, env)
		if err != nil {
			t.Fatalf("EvalWithCustomFuncs: %v", err)
		}
		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("$pick: expected map, got %T: %v", result, result)
		}
		if _, ok := m["a"]; !ok {
			t.Error("$pick: key 'a' missing")
		}
		if _, ok := m["c"]; !ok {
			t.Error("$pick: key 'c' missing")
		}
		if _, ok := m["b"]; ok {
			t.Error("$pick: key 'b' should be absent")
		}
	})

	t.Run("hash sha256", func(t *testing.T) {
		e, _ := gnata.Compile(`$hash("sha256", "hello")`)
		result, err := e.EvalWithCustomFuncs(ctx, nil, env)
		if err != nil {
			t.Fatalf("EvalWithCustomFuncs: %v", err)
		}
		s, ok := result.(string)
		if !ok || len(s) == 0 {
			t.Errorf("$hash: unexpected result %v", result)
		}
	})

	t.Run("flatten nested array", func(t *testing.T) {
		e, _ := gnata.Compile(`$flatten([[1,2],[3,4]])`)
		result, err := e.EvalWithCustomFuncs(ctx, nil, env)
		if err != nil {
			t.Fatalf("EvalWithCustomFuncs: %v", err)
		}
		arr, ok := result.([]any)
		if !ok || len(arr) != 4 {
			t.Errorf("$flatten: expected 4-element array, got %v", result)
		}
	})

	t.Run("clamp on field data", func(t *testing.T) {
		data := map[string]any{"value": float64(15)}
		e, _ := gnata.Compile(`$clamp(value, 0, 10)`)
		result, err := e.EvalWithCustomFuncs(ctx, data, env)
		if err != nil {
			t.Fatalf("EvalWithCustomFuncs: %v", err)
		}
		if result != float64(10) {
			t.Errorf("$clamp on field: want 10, got %v", result)
		}
	})
}

// TestIntegration_DirectFuncCall shows that sub-package functions can be called as
// plain Go functions without going through gnata expression evaluation.
func TestIntegration_DirectFuncCall(t *testing.T) {
	fn := extnumeric.Median()
	result, err := fn([]any{[]any{float64(1), float64(2), float64(3), float64(4), float64(5)}}, nil)
	if err != nil {
		t.Fatalf("Median: %v", err)
	}
	if result != float64(3) {
		t.Errorf("Median direct call: want 3, got %v", result)
	}
}
