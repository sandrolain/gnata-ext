// Package main demonstrates extnumeric functions with gnata.
package main

import (
	"context"
	"fmt"
	"math"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extnumeric"
)

func main() {
	env := gnata.NewCustomEnv(extnumeric.All())
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

	fmt.Println(eval(`$clamp(15, 0, 10)`, nil))           // 10
	fmt.Println(eval(`$clamp(-5, 0, 10)`, nil))           // 0
	fmt.Println(eval(`$sign(-42)`, nil))                  // -1
	fmt.Println(eval(`$sign(0)`, nil))                    // 0
	fmt.Println(eval(`$trunc(3.9)`, nil))                 // 3
	fmt.Println(eval(`$log(100, 10)`, nil))               // 2
	fmt.Println(eval(`$pi()`, nil))                       // 3.141592653589793
	fmt.Println(eval(`$e()`, nil))                        // 2.718281828459045

	fmt.Println("\n=== Statistics ===")
	fmt.Println(eval(`$median([1,2,3,4,5])`, nil))        // 3
	fmt.Println(eval(`$variance([2,4,4,4,5,5,7,9])`, nil)) // 4
	fmt.Println(eval(`$stddev([2,4,4,4,5,5,7,9])`, nil))  // 2
	fmt.Println(eval(`$percentile([1,2,3,4,5], 50)`, nil)) // 3
	fmt.Println(eval(`$mode([1,2,2,3,3,3,4])`, nil))      // 3

	fmt.Println("\n=== Trigonometry ===")
	fmt.Printf("sin(π/2) = %v\n", eval(`$sin($pi() / 2)`, nil)) // 1
	_ = math.Pi // satisfy import
}
