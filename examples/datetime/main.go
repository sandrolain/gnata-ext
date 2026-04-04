// Package main demonstrates extdatetime functions with gnata.
//
// All timestamps are Unix milliseconds (float64), matching JSONata $millis() convention.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extdatetime"
)

func main() {
	env := gnata.NewCustomEnv(extdatetime.All())
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

	// 2024-01-15 12:00:00 UTC as milliseconds
	ts := float64(time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC).UnixMilli())
	data := map[string]any{"ts": ts}

	// dateComponents: extract year/month/day/hour/minute/second/weekday
	fmt.Println("components:", eval(`$dateComponents(ts)`, data))

	// dateAdd: add duration
	fmt.Println("add 7 days:", eval(`$dateAdd(ts, 7, "day")`, data))
	fmt.Println("add 1 year:", eval(`$dateAdd(ts, 1, "year")`, data))

	// dateDiff: difference between two timestamps
	ts2 := float64(time.Date(2024, 1, 25, 12, 0, 0, 0, time.UTC).UnixMilli())
	data2 := map[string]any{"t1": ts, "t2": ts2}
	fmt.Println("diff days:", eval(`$dateDiff(t1, t2, "day")`, data2)) // 10

	// dateStartOf / dateEndOf
	fmt.Println("start of month:", eval(`$dateStartOf(ts, "month")`, data))
	fmt.Println("end of month:  ", eval(`$dateEndOf(ts, "month")`, data))
	fmt.Println("start of day:  ", eval(`$dateStartOf(ts, "day")`, data))
}
