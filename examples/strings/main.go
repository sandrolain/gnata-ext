// Package main demonstrates extstring functions with gnata.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extstring"
)

func main() {
	env := gnata.NewCustomEnv(extstring.All())
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

	fmt.Println(eval(`$startsWith("gnata-ext", "gnata")`, nil))   // true
	fmt.Println(eval(`$endsWith("gnata-ext", "ext")`, nil))       // true
	fmt.Println(eval(`$indexOf("hello world", "world")`, nil))    // 6
	fmt.Println(eval(`$lastIndexOf("abcabc", "b")`, nil))         // 4
	fmt.Println(eval(`$capitalize("hello WORLD")`, nil))          // Hello world
	fmt.Println(eval(`$titleCase("hello world foo")`, nil))       // Hello World Foo
	fmt.Println(eval(`$camelCase("hello_world_foo")`, nil))       // helloWorldFoo
	fmt.Println(eval(`$snakeCase("helloWorldFoo")`, nil))         // hello_world_foo
	fmt.Println(eval(`$kebabCase("Hello World Foo")`, nil))       // hello-world-foo
	fmt.Println(eval(`$repeat("ab", 4)`, nil))                    // abababab
	fmt.Println(eval(`$words("hello world foo")`, nil))           // [hello world foo]
	fmt.Println(eval(`$template("Hi {{name}}!", {"name": "World"})`, nil)) // Hi World!
}
