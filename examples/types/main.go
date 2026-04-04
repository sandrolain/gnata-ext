// Package main demonstrates exttypes functions with gnata.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/exttypes"
)

func main() {
	env := gnata.NewCustomEnv(exttypes.All())
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

	// Type inspection
	fmt.Println(eval(`$isString("hello")`, nil))          // true
	fmt.Println(eval(`$isString(42)`, nil))               // false
	fmt.Println(eval(`$isNumber(3.14)`, nil))             // true
	fmt.Println(eval(`$isBoolean(true)`, nil))            // true
	fmt.Println(eval(`$isArray([1,2,3])`, nil))           // true
	fmt.Println(eval(`$isObject({"a":1})`, nil))          // true
	fmt.Println(eval(`$isNull(null)`, nil))               // true

	// isEmpty checks
	fmt.Println(eval(`$isEmpty("")`, nil))                // true
	fmt.Println(eval(`$isEmpty([])`, nil))                // true
	fmt.Println(eval(`$isEmpty({})`, nil))                // true
	fmt.Println(eval(`$isEmpty("hello")`, nil))           // false

	// $default(v, fallback) — returns fallback when v is null
	fmt.Println(eval(`$default(null, "N/A")`, nil))       // N/A
	fmt.Println(eval(`$default("real", "N/A")`, nil))     // real

	// $identity passes the value through unchanged
	fmt.Println(eval(`$identity(42)`, nil))               // 42
	fmt.Println(eval(`$identity([1,2,3])`, nil))          // [1 2 3]

	// Use with data from the input document
	data := map[string]any{"score": nil, "label": ""}
	fmt.Println(eval(`$default(score, 0)`, data))         // 0 (score is null)
	fmt.Println(eval(`$isEmpty(label)`, data))            // true
}
