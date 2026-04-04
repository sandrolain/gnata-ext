// Package main demonstrates extformat functions with gnata.
package main

import (
	"context"
	"fmt"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext/extformat"
)

func main() {
	env := gnata.NewCustomEnv(extformat.All())
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

	// $csv(text) — parse CSV text into array of objects (first row = headers)
	csvText := "name,age,city\nAlice,30,NYC\nBob,25,LA"
	result := eval(`$csv(text)`, map[string]any{"text": csvText})
	fmt.Println("parsed CSV:", result)

	// $toCSV(array) — serialise array of objects to CSV text
	records := []any{
		map[string]any{"name": "Alice", "age": 30.0},
		map[string]any{"name": "Bob", "age": 25.0},
	}
	fmt.Println("to CSV:\n" + fmt.Sprint(eval(`$toCSV(rows)`, map[string]any{"rows": records})))

	// $template(str, vars) — {{key}} placeholder replacement
	fmt.Println(eval(`$template("Hello, {{name}}! You are {{age}} years old.", {"name":"Alice","age":30})`, nil))
	fmt.Println(eval(`$template("Server: {{host}}:{{port}}", {"host":"localhost","port":8080})`, nil))
}
