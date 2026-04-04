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
