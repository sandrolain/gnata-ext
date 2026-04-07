// Package presets provides pre-built gnata custom environments for common use cases.
//
// Each function returns a *evaluator.Environment that can be passed directly to
// Expression.EvalWithCustomFuncs.
//
//	result, _ := expr.EvalWithCustomFuncs(ctx, data, presets.DataEnv())
package presets

import (
	"github.com/recolabs/gnata"

	"github.com/sandrolain/gnata-ext/pkg/ext/extarray"
	"github.com/sandrolain/gnata-ext/pkg/ext/extdatetime"
	"github.com/sandrolain/gnata-ext/pkg/ext/extformat"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
	"github.com/sandrolain/gnata-ext/pkg/ext/extobject"
	"github.com/sandrolain/gnata-ext/pkg/ext/extpath"
	"github.com/sandrolain/gnata-ext/pkg/ext/extstring"
	"github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
)

func merge(sources ...map[string]gnata.CustomFunc) map[string]gnata.CustomFunc {
	out := make(map[string]gnata.CustomFunc)
	for _, m := range sources {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// DataEnv returns functions suited for data transformation:
// extarray + extobject + exttypes + extnumeric + extpath.
func DataEnv() map[string]gnata.CustomFunc {
	return merge(
		extarray.All(),
		extobject.All(),
		exttypes.All(),
		extnumeric.All(),
		extpath.All(),
	)
}

// TextEnv returns functions suited for text processing:
// extstring + extformat + exttypes.
func TextEnv() map[string]gnata.CustomFunc {
	return merge(
		extstring.All(),
		extformat.All(),
		exttypes.All(),
	)
}

// SecureEnv returns functions with uuid removed (non-deterministic).
// Use this when reproducibility and auditability are required.
func SecureEnv() map[string]gnata.CustomFunc {
	funcs := merge(
		extarray.All(),
		extdatetime.All(),
		extformat.All(),
		extnumeric.All(),
		extobject.All(),
		extpath.All(),
		extstring.All(),
		exttypes.All(),
	)
	// Explicitly exclude non-deterministic functions.
	delete(funcs, "uuid")
	return funcs
}

// AnalyticsEnv returns functions suited for analytics workloads:
// extdatetime + extarray + extnumeric.
func AnalyticsEnv() map[string]gnata.CustomFunc {
	return merge(
		extdatetime.All(),
		extarray.All(),
		extnumeric.All(),
	)
}
