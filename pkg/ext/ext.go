// Package ext provides a single entry point to all gnata-ext extension functions.
//
// Use AllFuncs to get a merged map of every extension function, or NewEnv to
// create a *gnata.CustomEnv ready to pass to gnata expressions.
//
//	env := ext.NewEnv()
//	expr, _ := gnata.Compile(`$uuid()`)
//	result, _ := expr.EvalWithCustomFuncs(ctx, nil, env)
package ext

import (
	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extarray"
	"github.com/sandrolain/gnata-ext/pkg/ext/extcolor"
	"github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
	"github.com/sandrolain/gnata-ext/pkg/ext/extdatetime"
	"github.com/sandrolain/gnata-ext/pkg/ext/extdiff"
	"github.com/sandrolain/gnata-ext/pkg/ext/extformat"
	"github.com/sandrolain/gnata-ext/pkg/ext/extgeo"
	"github.com/sandrolain/gnata-ext/pkg/ext/extjson"
	"github.com/sandrolain/gnata-ext/pkg/ext/extlogic"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnet"
	"github.com/sandrolain/gnata-ext/pkg/ext/extregex"
	"github.com/sandrolain/gnata-ext/pkg/ext/exturi"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
	"github.com/sandrolain/gnata-ext/pkg/ext/extobject"
	"github.com/sandrolain/gnata-ext/pkg/ext/extpath"
	"github.com/sandrolain/gnata-ext/pkg/ext/extschema"
	"github.com/sandrolain/gnata-ext/pkg/ext/extsemver"
	"github.com/sandrolain/gnata-ext/pkg/ext/extstring"
	"github.com/sandrolain/gnata-ext/pkg/ext/exttext"
	"github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
	"github.com/sandrolain/gnata-ext/pkg/ext/extvalidate"
)

// AllFuncs returns a merged map of all extension functions from every sub-package.
// Keys are function names (without the leading $).
func AllFuncs() map[string]gnata.CustomFunc {
	sources := []map[string]gnata.CustomFunc{
		extarray.All(),
		extcolor.All(),
		extcrypto.All(),
		extdatetime.All(),
		extdiff.All(),
		extformat.All(),
		extgeo.All(),
		extjson.All(),
		extlogic.All(),
		extnet.All(),
		extregex.All(),
		exturi.All(),
		extnumeric.All(),
		extobject.All(),
		extpath.All(),
		extschema.All(),
		extsemver.All(),
		extstring.All(),
		exttext.All(),
		exttypes.All(),
		extvalidate.All(),
	}
	merged := make(map[string]gnata.CustomFunc)
	for _, m := range sources {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

// NewEnv creates a gnata custom environment containing all extension functions.
// The returned value is *evaluator.Environment (gnata internal type) and can be
// passed directly to Expression.EvalWithCustomFuncs.
func NewEnv() any {
	return gnata.NewCustomEnv(AllFuncs())
}
