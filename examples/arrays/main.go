// Package main demonstrates extarray functions with gnata.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extarray"
)

func main() {
	env := gnata.NewCustomEnv(extarray.All())
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

	fmt.Println(eval(`$first([10, 20, 30])`, nil))                   // 10
	fmt.Println(eval(`$last([10, 20, 30])`, nil))                    // 30
	fmt.Println(eval(`$take([1,2,3,4,5], 3)`, nil))                  // [1 2 3]
	fmt.Println(eval(`$skip([1,2,3,4,5], 2)`, nil))                  // [3 4 5]
	fmt.Println(eval(`$slice([1,2,3,4,5], 1, 4)`, nil))              // [2 3 4]
	fmt.Println(eval(`$flatten([[1,2],[3,[4,5]]])`, nil))             // [1 2 3 4 5]
	fmt.Println(eval(`$chunk([1,2,3,4,5], 2)`, nil))                 // [[1 2] [3 4] [5]]
	fmt.Println(eval(`$union([1,2,3], [3,4,5])`, nil))               // [1 2 3 4 5]
	fmt.Println(eval(`$intersection([1,2,3,4], [2,4,6])`, nil))      // [2 4]
	fmt.Println(eval(`$difference([1,2,3,4], [2,4])`, nil))          // [1 3]
	fmt.Println(eval(`$range(1, 6)`, nil))                           // [1 2 3 4 5]
	fmt.Println(eval(`$range(0, 10, 3)`, nil))                       // [0 3 6 9]
	fmt.Println(eval(`$zipLongest([1,2,3], [4,5])`, nil))            // [[1 4] [2 5] [3 <nil>]]
	fmt.Println(eval(`$window([1,2,3,4,5], 3)`, nil))                // [[1 2 3] [2 3 4] [3 4 5]]
}
