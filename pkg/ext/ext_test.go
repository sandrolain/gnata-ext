package ext_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext"
)

func TestAllFuncs(t *testing.T) {
	funcs := ext.AllFuncs()
	required := []string{
		// extarray
		"first", "last", "take", "skip", "slice", "flatten", "chunk",
		"union", "intersection", "difference", "symmetricDifference",
		"range", "zipLongest", "window",
		// extcrypto
		"uuid", "hash", "hmac",
		// extdatetime
		"dateAdd", "dateDiff", "dateComponents", "dateStartOf", "dateEndOf",
		// extformat
		"csv", "toCSV", "template",
		// extnumeric
		"log", "sign", "trunc", "clamp",
		"sin", "cos", "tan", "asin", "acos", "atan", "atan2",
		"pi", "e",
		"median", "variance", "stddev", "percentile", "mode",
		// extobject
		"values", "pairs", "fromPairs", "pick", "omit", "deepMerge", "invert", "size", "rename",
		// extstring
		"startsWith", "endsWith", "indexOf", "lastIndexOf",
		"capitalize", "titleCase", "camelCase", "snakeCase", "kebabCase",
		"repeat", "words", "template",
		// exttypes
		"isString", "isNumber", "isBoolean", "isArray", "isObject",
		"isNull", "isUndefined", "isEmpty", "default", "identity",
	}
	for _, name := range required {
		if _, ok := funcs[name]; !ok {
			t.Errorf("AllFuncs(): missing %q", name)
		}
	}
	t.Logf("AllFuncs() returned %d functions", len(funcs))
}

func TestNewEnv(t *testing.T) {
	env := ext.NewEnv()
	if env == nil {
		t.Fatal("NewEnv() returned nil")
	}
}

func TestBuilderWithAllFuncs(t *testing.T) {
	funcs := ext.NewEnvBuilder().WithAllFuncs().Build()
	if len(funcs) == 0 {
		t.Error("WithAllFuncs: expected non-empty")
	}
}

func TestBuilderIndividualPackages(t *testing.T) {
	b := ext.NewEnvBuilder()
	_ = b.WithArrayFuncs()
	_ = b.WithCryptoFuncs()
	_ = b.WithDatetimeFuncs()
	_ = b.WithFormatFuncs()
	_ = b.WithGeoFuncs()
	_ = b.WithJSONFuncs()
	_ = b.WithNetFuncs()
	_ = b.WithObjectFuncs()
	_ = b.WithPathFuncs()
	_ = b.WithTypesFuncs()
	_ = b.WithValidateFuncs()
	funcs := b.Build()
	if _, ok := funcs["first"]; !ok {
		t.Error("builder: missing 'first' from array package")
	}
	if _, ok := funcs["uuid"]; !ok {
		t.Error("builder: missing 'uuid' from crypto package")
	}
}
