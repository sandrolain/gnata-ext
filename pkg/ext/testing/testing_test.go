package exttesting_test

import (
	"testing"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
	exttesting "github.com/sandrolain/gnata-ext/pkg/ext/testing"
)

func TestEnv_Eval_Basic(t *testing.T) {
	env := exttesting.New(ext.AllFuncs())
	got := env.Eval(t, `$first([10, 20, 30])`, nil)
	if got != float64(10) {
		t.Fatalf("want 10, got %v", got)
	}
}

func TestEnv_AssertEqual(t *testing.T) {
	env := exttesting.New(ext.AllFuncs())
	env.AssertEqual(t, `$capitalize("hello")`, nil, "Hello")
}

func TestEnv_WithFrozenTime(t *testing.T) {
	const frozen int64 = 1705319400000
	env := exttesting.New(ext.AllFuncs(), exttesting.WithFrozenTime(frozen))
	env.AssertEqual(t, `$millis()`, nil, float64(frozen))
}

func TestEnv_WithDeterministicUUID(t *testing.T) {
	const fixed = "00000000-0000-0000-0000-000000000001"
	env := exttesting.New(ext.AllFuncs(), exttesting.WithDeterministicUUID(fixed))
	env.AssertEqual(t, `$uuid()`, nil, fixed)
}

func TestEnv_AssertError(t *testing.T) {
	env := exttesting.New(ext.AllFuncs())
	// $hash with an unknown algorithm should produce an error.
	env.AssertError(t, `$hash("unknown-algo", "data")`, nil)
}

func TestEnv_WithExtraFuncs(t *testing.T) {
	extra := map[string]gnata.CustomFunc{
		"triple": func(args []any, _ any) (any, error) {
			return args[0].(float64) * 3, nil
		},
	}
	env := exttesting.New(ext.AllFuncs(), exttesting.WithExtraFuncs(extra))
	env.AssertEqual(t, `$triple(4)`, nil, float64(12))
}
