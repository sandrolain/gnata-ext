// Package main demonstrates extobject functions with gnata.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extobject"
)

func main() {
	env := gnata.NewCustomEnv(extobject.All())
	ctx := context.Background()

	eval := func(expr string, data any) any {
		e, err := gnata.Compile(expr)
		if err != nil {
			panic(fmt.Sprintf("compile %q: %v", expr, err))
		}
		result, err := e.EvalWithCustomFuncs(ctx, data, env)
		if err != nil {
			panic(fmt.Sprintf("eval %q: %v", expr, err))
		}
		return result
	}

	obj := map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}

	fmt.Println(eval(`$values({"a":1,"b":2,"c":3})`, nil))           // [1 2 3]
	fmt.Println(eval(`$pairs({"a":1,"b":2})`, nil))                   // [[a 1] [b 2]]
	fmt.Println(eval(`$fromPairs([["x",10],["y",20]])`, nil))         // map[x:10 y:20]
	fmt.Println(eval(`$pick($, ["a","c"])`, obj))                     // map[a:1 c:3]
	fmt.Println(eval(`$omit($, ["b"])`, obj))                         // map[a:1 c:3]
	fmt.Println(eval(`$invert({"a":"x","b":"y"})`, nil))              // map[x:a y:b]
	fmt.Println(eval(`$size({"a":1,"b":2,"c":3})`, nil))              // 3
	fmt.Println(eval(`$rename($, {"a":"alpha"})`, obj))               // map[alpha:1 b:2 c:3]

	// deepMerge
	base := map[string]any{"x": map[string]any{"a": 1.0}, "y": 2.0}
	fmt.Println(eval(`$deepMerge($, {"x":{"b":3},"z":4})`, base))    // x merged, z added
}
