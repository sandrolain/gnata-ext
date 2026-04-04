// Package main demonstrates using all gnata-ext extension functions together
// via ext.AllFuncs() and gnata.NewCustomEnv.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
)

func main() {
	env := gnata.NewCustomEnv(ext.AllFuncs())
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

	fmt.Println("=== String functions ===")
	fmt.Println(eval(`$camelCase("hello world foo")`, nil))     // helloWorldFoo
	fmt.Println(eval(`$snakeCase("helloWorldFoo")`, nil))       // hello_world_foo
	fmt.Println(eval(`$titleCase("hello world")`, nil))         // Hello World
	fmt.Println(eval(`$repeat("ab", 3)`, nil))                  // ababab

	fmt.Println("\n=== Numeric functions ===")
	fmt.Println(eval(`$clamp(15, 0, 10)`, nil))                // 10
	fmt.Println(eval(`$sign(-42)`, nil))                       // -1
	fmt.Println(eval(`$median([1,2,3,4,5])`, nil))             // 3

	fmt.Println("\n=== Array functions ===")
	fmt.Println(eval(`$flatten([[1,2],[3,[4,5]]])`, nil))      // [1 2 3 4 5]
	fmt.Println(eval(`$chunk([1,2,3,4,5], 2)`, nil))           // [[1 2] [3 4] [5]]

	fmt.Println("\n=== Object functions ===")
	data := map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}
	fmt.Println(eval(`$pick($, ["a", "c"])`, data))            // map[a:1 c:3]
	fmt.Println(eval(`$omit($, ["b"])`, data))                 // map[a:1 c:3]

	fmt.Println("\n=== Type functions ===")
	fmt.Println(eval(`$isArray([1,2,3])`, nil))                // true
	fmt.Println(eval(`$default(null, "fallback")`, nil))       // fallback

	fmt.Println("\n=== Crypto functions ===")
	fmt.Println(eval(`$uuid()`, nil))                          // e.g. 550e8400-e29b-41d4-a716-446655440000
	fmt.Println(eval(`$hash("sha256", "hello")`, nil))         // hex hash
}
