// Package exttesting provides a lightweight test helper for expressions that
// use gnata-ext custom functions.
//
// It wraps a set of gnata custom functions with optional overrides (frozen time,
// deterministic UUID) so that tests are reproducible.
//
//	env := exttesting.New(
//	    ext.AllFuncs(),
//	    exttesting.WithDeterministicUUID("00000000-0000-0000-0000-000000000001"),
//	    exttesting.WithFrozenTime(1705319400000),
//	)
//	env.AssertEqual(t, `$uuid()`, nil, "00000000-0000-0000-0000-000000000001")
package exttesting

import (
	"context"
	"reflect"
	"testing"

	"github.com/recolabs/gnata"
)

// Option is a functional option for New.
type Option func(*config)

type config struct {
	frozenTime *int64
	fixedUUID  string
	extraFuncs map[string]gnata.CustomFunc
}

// WithFrozenTime replaces the "millis" and "now" functions so that they return
// ts (Unix milliseconds) instead of the current time.
func WithFrozenTime(ts int64) Option {
	return func(c *config) { c.frozenTime = &ts }
}

// WithDeterministicUUID replaces the "uuid" function so that it always returns
// fixed instead of a random UUID.
func WithDeterministicUUID(fixed string) Option {
	return func(c *config) { c.fixedUUID = fixed }
}

// WithExtraFuncs adds or overrides functions in the environment.
func WithExtraFuncs(funcs map[string]gnata.CustomFunc) Option {
	return func(c *config) { c.extraFuncs = funcs }
}

// Env is a test helper for evaluating gnata expressions.
type Env struct {
	funcs map[string]gnata.CustomFunc
}

// New builds a test Env from funcs applying the given options.
func New(funcs map[string]gnata.CustomFunc, opts ...Option) *Env {
	cfg := &config{}
	for _, o := range opts {
		o(cfg)
	}

	merged := make(map[string]gnata.CustomFunc, len(funcs))
	for k, v := range funcs {
		merged[k] = v
	}
	for k, v := range cfg.extraFuncs {
		merged[k] = v
	}
	if cfg.frozenTime != nil {
		ts := float64(*cfg.frozenTime)
		frozen := func(_ []any, _ any) (any, error) { return ts, nil }
		merged["millis"] = frozen
		merged["now"] = frozen
	}
	if cfg.fixedUUID != "" {
		fixed := cfg.fixedUUID
		merged["uuid"] = func(_ []any, _ any) (any, error) { return fixed, nil }
	}

	return &Env{funcs: merged}
}

// Eval compiles and evaluates expr against data. Calls t.Fatal on any error.
func (e *Env) Eval(t testing.TB, expr string, data any) any {
	t.Helper()
	compiled, err := gnata.Compile(expr)
	if err != nil {
		t.Fatalf("exttesting.Eval: compile %q: %v", expr, err)
	}
	env := gnata.NewCustomEnv(e.funcs)
	result, err := compiled.EvalWithCustomFuncs(context.Background(), data, env)
	if err != nil {
		t.Fatalf("exttesting.Eval: eval %q: %v", expr, err)
	}
	return result
}

// AssertEqual evaluates expr and compares the result to want using reflect.DeepEqual.
func (e *Env) AssertEqual(t testing.TB, expr string, data any, want any) {
	t.Helper()
	got := e.Eval(t, expr, data)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expr %q\n\tgot  %v (%T)\n\twant %v (%T)", expr, got, got, want, want)
	}
}

// AssertError evaluates expr and calls t.Fatal if no error is returned.
func (e *Env) AssertError(t testing.TB, expr string, data any) {
	t.Helper()
	compiled, err := gnata.Compile(expr)
	if err != nil {
		return // compile-time error counts
	}
	env := gnata.NewCustomEnv(e.funcs)
	_, err = compiled.EvalWithCustomFuncs(context.Background(), data, env)
	if err == nil {
		t.Errorf("expr %q: expected an error but got none", expr)
	}
}
