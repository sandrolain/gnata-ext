package ext

import (
	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extarray"
	"github.com/sandrolain/gnata-ext/pkg/ext/extcrypto"
	"github.com/sandrolain/gnata-ext/pkg/ext/extdatetime"
	"github.com/sandrolain/gnata-ext/pkg/ext/extformat"
	"github.com/sandrolain/gnata-ext/pkg/ext/extgeo"
	"github.com/sandrolain/gnata-ext/pkg/ext/extjson"
	"github.com/sandrolain/gnata-ext/pkg/ext/extlogic"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnet"
	"github.com/sandrolain/gnata-ext/pkg/ext/exturi"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
	"github.com/sandrolain/gnata-ext/pkg/ext/extobject"
	"github.com/sandrolain/gnata-ext/pkg/ext/extpath"
	"github.com/sandrolain/gnata-ext/pkg/ext/extstring"
	"github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
	"github.com/sandrolain/gnata-ext/pkg/ext/extvalidate"
)

// EnvBuilder constructs a gnata custom environment from selected extension packages.
//
// Start with NewEnvBuilder(), chain With* methods to add packages or individual
// functions, optionally call Without to exclude specific names, then call Build
// to obtain a *evaluator.Environment ready for EvalWithCustomFuncs.
//
//	env := ext.NewEnvBuilder().
//	    WithArrayFuncs().
//	    WithNumericFuncs().
//	    Without("range").
//	    Build()
type EnvBuilder struct {
	funcs map[string]gnata.CustomFunc
}

// NewEnvBuilder returns an empty EnvBuilder.
func NewEnvBuilder() *EnvBuilder {
	return &EnvBuilder{funcs: make(map[string]gnata.CustomFunc)}
}

func (b *EnvBuilder) merge(pkg map[string]gnata.CustomFunc) *EnvBuilder {
	for k, v := range pkg {
		b.funcs[k] = v
	}
	return b
}

// WithAllFuncs adds every extension function (equivalent to AllFuncs()).
func (b *EnvBuilder) WithAllFuncs() *EnvBuilder      { return b.merge(AllFuncs()) }
func (b *EnvBuilder) WithArrayFuncs() *EnvBuilder    { return b.merge(extarray.All()) }
func (b *EnvBuilder) WithCryptoFuncs() *EnvBuilder   { return b.merge(extcrypto.All()) }
func (b *EnvBuilder) WithDatetimeFuncs() *EnvBuilder { return b.merge(extdatetime.All()) }
func (b *EnvBuilder) WithFormatFuncs() *EnvBuilder   { return b.merge(extformat.All()) }
func (b *EnvBuilder) WithGeoFuncs() *EnvBuilder      { return b.merge(extgeo.All()) }
func (b *EnvBuilder) WithJSONFuncs() *EnvBuilder     { return b.merge(extjson.All()) }
func (b *EnvBuilder) WithLogicFuncs() *EnvBuilder    { return b.merge(extlogic.All()) }
func (b *EnvBuilder) WithNetFuncs() *EnvBuilder      { return b.merge(extnet.All()) }
func (b *EnvBuilder) WithURIFuncs() *EnvBuilder      { return b.merge(exturi.All()) }
func (b *EnvBuilder) WithNumericFuncs() *EnvBuilder  { return b.merge(extnumeric.All()) }
func (b *EnvBuilder) WithObjectFuncs() *EnvBuilder   { return b.merge(extobject.All()) }
func (b *EnvBuilder) WithPathFuncs() *EnvBuilder     { return b.merge(extpath.All()) }
func (b *EnvBuilder) WithStringFuncs() *EnvBuilder   { return b.merge(extstring.All()) }
func (b *EnvBuilder) WithTypesFuncs() *EnvBuilder    { return b.merge(exttypes.All()) }
func (b *EnvBuilder) WithValidateFuncs() *EnvBuilder { return b.merge(extvalidate.All()) }

// WithPackage merges an arbitrary map of functions into the builder.
func (b *EnvBuilder) WithPackage(pkg map[string]gnata.CustomFunc) *EnvBuilder {
	return b.merge(pkg)
}

// WithFunc adds a single named function to the builder.
func (b *EnvBuilder) WithFunc(name string, fn gnata.CustomFunc) *EnvBuilder {
	b.funcs[name] = fn
	return b
}

// Without removes the named functions from the builder.
func (b *EnvBuilder) Without(names ...string) *EnvBuilder {
	for _, n := range names {
		delete(b.funcs, n)
	}
	return b
}

// Funcs returns a copy of the current function map.
func (b *EnvBuilder) Funcs() map[string]gnata.CustomFunc {
	out := make(map[string]gnata.CustomFunc, len(b.funcs))
	for k, v := range b.funcs {
		out[k] = v
	}
	return out
}

// Build returns a copy of the accumulated function map.
// Pass the result to gnata.NewCustomEnv to create a *evaluator.Environment:
//
//	env := gnata.NewCustomEnv(b.Build())
//	expr.EvalWithCustomFuncs(ctx, data, env)
func (b *EnvBuilder) Build() map[string]gnata.CustomFunc {
	return b.Funcs()
}
