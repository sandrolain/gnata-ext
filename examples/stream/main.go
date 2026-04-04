// Package main demonstrates using gnata StreamEvaluator with gnata-ext functions.
//
// StreamEvaluator is designed for high-throughput evaluation of the same
// set of compiled expressions against many inputs.
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
)

func main() {
	// Compile expressions once.
	exprName, _ := gnata.Compile(`$camelCase(name)`)
	exprScore, _ := gnata.Compile(`$clamp(score, 0, 100)`)
	exprTags, _ := gnata.Compile(`$flatten(tags)`)

	// Build a StreamEvaluator with all extension functions.
	se := gnata.NewStreamEvaluator(
		[]*gnata.Expression{exprName, exprScore, exprTags},
		gnata.WithCustomFunctions(ext.AllFuncs()),
	)

	records := []map[string]any{
		{"name": "hello world", "score": 42.0, "tags": []any{[]any{"go", "ext"}, []any{"jsonata"}}},
		{"name": "foo bar baz", "score": 110.0, "tags": []any{[]any{"data"}}},
		{"name": "open source", "score": -5.0, "tags": []any{[]any{"oss", "tools"}}},
	}

	ctx := context.Background()
	for _, rec := range records {
		raw, _ := json.Marshal(rec)

		name, _ := se.EvalOne(ctx, raw, "record", 0)   // expression index 0
		score, _ := se.EvalOne(ctx, raw, "record", 1)  // expression index 1
		tags, _ := se.EvalOne(ctx, raw, "record", 2)   // expression index 2

		fmt.Printf("name=%-20v  score=%-6v  tags=%v\n", name, score, tags)
	}

	// Use EvalMany to evaluate all expressions in a single pass.
	fmt.Println("\n--- EvalMany ---")
	for _, rec := range records {
		raw, _ := json.Marshal(rec)
		results, _ := se.EvalMany(ctx, raw, "record", []int{0, 1, 2})
		fmt.Println(results)
	}
}
