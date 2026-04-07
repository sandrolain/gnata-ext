// gnata-ext-cli evaluates JSONata expressions with all gnata-ext functions,
// lists available functions, and describes individual functions.
//
// Usage:
//
//	gnata-ext-cli eval '<expr>' [--data '<json>'] [--data-file <file>]
//	gnata-ext-cli list [--package <pkg>]
//	gnata-ext-cli describe <funcName>
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "eval":
		runEval(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	case "describe":
		runDescribe(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `gnata-ext-cli — evaluate JSONata expressions with gnata-ext functions

Usage:
  gnata-ext-cli eval '<expr>' [--data '<json>'] [--data-file <file>]
  gnata-ext-cli list [--package <pkg>]
  gnata-ext-cli describe <funcName>`)
}

// --- eval subcommand ---

func runEval(args []string) {
	fs := flag.NewFlagSet("eval", flag.ExitOnError)
	dataStr := fs.String("data", "", "JSON data string")
	dataFile := fs.String("data-file", "", "path to JSON data file")
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "eval: missing expression argument")
		os.Exit(1)
	}
	expr := fs.Arg(0)

	var data any
	switch {
	case *dataFile != "":
		raw, err := os.ReadFile(*dataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "eval: read data file: %v\n", err)
			os.Exit(1)
		}
		if err := json.Unmarshal(raw, &data); err != nil {
			fmt.Fprintf(os.Stderr, "eval: parse data file: %v\n", err)
			os.Exit(1)
		}
	case *dataStr != "":
		if err := json.Unmarshal([]byte(*dataStr), &data); err != nil {
			fmt.Fprintf(os.Stderr, "eval: parse --data: %v\n", err)
			os.Exit(1)
		}
	}

	compiled, err := gnata.Compile(expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval: compile: %v\n", err)
		os.Exit(1)
	}
	env := gnata.NewCustomEnv(ext.AllFuncs())
	result, err := compiled.EvalWithCustomFuncs(context.Background(), data, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval: %v\n", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval: marshal result: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}

// --- list subcommand ---

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	pkg := fs.String("package", "", "filter by package name")
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	catalog := ext.Catalog()
	if *pkg != "" {
		filtered := catalog[:0]
		for _, f := range catalog {
			if strings.EqualFold(f.Package, *pkg) {
				filtered = append(filtered, f)
			}
		}
		catalog = filtered
	}

	sort.Slice(catalog, func(i, j int) bool {
		if catalog[i].Package != catalog[j].Package {
			return catalog[i].Package < catalog[j].Package
		}
		return catalog[i].Name < catalog[j].Name
	})

	fmt.Printf("%-22s %-14s %s\n", "FUNCTION", "PACKAGE", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 80))
	for _, f := range catalog {
		fmt.Printf("$%-21s %-14s %s\n", f.Name, f.Package, f.Description)
	}
}

// --- describe subcommand ---

func runDescribe(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "describe: missing function name")
		os.Exit(1)
	}
	name := strings.TrimPrefix(args[0], "$")

	for _, f := range ext.Catalog() {
		if f.Name == name {
			fmt.Printf("Function:    $%s\n", f.Name)
			fmt.Printf("Package:     %s\n", f.Package)
			fmt.Printf("Signature:   $%s%s\n", f.Name, f.Signature)
			fmt.Printf("Description: %s\n", f.Description)
			return
		}
	}
	fmt.Fprintf(os.Stderr, "describe: unknown function %q\n", name)
	os.Exit(1)
}
