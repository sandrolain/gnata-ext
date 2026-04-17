// jn — JSONata processor with gnata-ext extension functions.
//
// Usage (inspired by jq):
//
//	jn [flags] [expr] [file...]
//	jn list [--package <pkg>]
//	jn describe <funcName>
//	jn version
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/recolabs/gnata"
	"github.com/sandrolain/gnata-ext/pkg/ext"
	"github.com/spf13/cobra"
)

// Version is injected by goreleaser via -ldflags "-X main.Version=<tag>".
var Version = "dev"

func main() {
	if err := buildRoot().Execute(); err != nil {
		os.Exit(1)
	}
}

// ─── root (eval) command ───────────────────────────────────────────────────

type evalOpts struct {
	dataStr     string
	dataFile    string
	compact     bool
	rawOutput   bool
	rawInput    bool
	nullInput   bool
	exitStatus  bool
	fromFile    string
	tabIndent   bool
	indentN     int
	sortKeys    bool
	joinOutput  bool
	slurp       bool
	argVars     []string
	argjsonVars []string
}

func buildRoot() *cobra.Command {
	var opts evalOpts

	root := &cobra.Command{
		Use:   "jn [flags] [expr] [file...]",
		Short: "JSONata processor with extended functions",
		Long: `jn — JSONata expression processor (inspired by jq)

Reads JSON from files or stdin, evaluates a JSONata expression with
gnata-ext extension functions, and writes results to stdout.

EXAMPLES
  # identity (pretty-print input)
  echo '{"x":1}' | jn '$'

  # field access
  echo '{"name":"Alice"}' | jn '$.name'

  # use an extension function
  echo '"hello world"' | jn '$camelCase($)'

  # null input (no data needed)
  jn -n '$uuid()'

  # from file
  jn '$.users.$count($)' data.json

  # compact output
  echo '[1,2,3]' | jn -c '$reverse($)'

  # raw string output
  echo '{"msg":"hello"}' | jn -r '$.msg'

  # slurp all JSON values into array
  cat lines.json | jn -s '$count($)'

  # list all extension functions
  jn list

  # describe a function
  jn describe haversine`,
		SilenceUsage: true,
		Args:         cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEval(opts, args)
		},
	}

	f := root.Flags()
	f.StringVar(&opts.dataStr, "data", "", "inline JSON input data")
	f.StringVar(&opts.dataFile, "data-file", "", "path to JSON input file")
	f.BoolVarP(&opts.compact, "compact", "c", false, "compact output (no indentation)")
	f.BoolVarP(&opts.rawOutput, "raw-output", "r", false, "output raw strings, not JSON-encoded")
	f.BoolVarP(&opts.rawInput, "raw-input", "R", false, "read each input line as a raw string")
	f.BoolVarP(&opts.nullInput, "null-input", "n", false, "use null as input (no data required)")
	f.BoolVarP(&opts.exitStatus, "exit-status", "e", false, "exit 5 if last output is null or false")
	f.StringVarP(&opts.fromFile, "from-file", "f", "", "read JSONata expression from file")
	f.BoolVar(&opts.tabIndent, "tab", false, "use tab indentation")
	f.IntVar(&opts.indentN, "indent", 2, "indentation width in spaces (0–7)")
	f.BoolVarP(&opts.sortKeys, "sort-keys", "S", false, "sort object keys (no-op: Go already sorts)")
	f.BoolVarP(&opts.joinOutput, "join-output", "j", false, "no trailing newline after each output")
	f.BoolVarP(&opts.slurp, "slurp", "s", false, "slurp all inputs into an array before eval")
	f.StringArrayVar(&opts.argVars, "arg", nil, "bind string variable: --arg name=value")
	f.StringArrayVar(&opts.argjsonVars, "argjson", nil, "bind JSON variable: --argjson name=json")

	root.AddCommand(buildListCmd())
	root.AddCommand(buildDescribeCmd())
	root.AddCommand(buildVersionCmd())

	return root
}

// ─── eval logic ───────────────────────────────────────────────────────────

func runEval(opts evalOpts, posArgs []string) error {
	// Resolve expression (default: return input as-is)
	expr := "$"
	fileArgs := posArgs

	switch {
	case opts.fromFile != "":
		raw, err := os.ReadFile(opts.fromFile)
		if err != nil {
			return fmt.Errorf("read expression file: %w", err)
		}
		expr = strings.TrimSpace(string(raw))
	case len(posArgs) > 0:
		expr = posArgs[0]
		fileArgs = posArgs[1:]
	}

	// Prepend variable bindings when --arg / --argjson are used
	var err error
	expr, err = injectVars(expr, opts.argVars, opts.argjsonVars)
	if err != nil {
		return err
	}

	compiled, err := gnata.Compile(expr)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}
	env := gnata.NewCustomEnv(ext.AllFuncs())
	ctx := context.Background()
	p := &printer{opts: opts}

	// Null-input mode: evaluate once with nil
	if opts.nullInput {
		result, err := compiled.EvalWithCustomFuncs(ctx, nil, env)
		if err != nil {
			return fmt.Errorf("eval: %w", err)
		}
		if err := p.print(result); err != nil {
			return err
		}
		return p.checkExit()
	}

	// Build input sources
	sources := buildSources(opts, fileArgs)

	if opts.slurp {
		var all []any
		for _, src := range sources {
			if err := src.eachValue(opts.rawInput, func(v any) error {
				all = append(all, v)
				return nil
			}); err != nil {
				return err
			}
		}
		result, err := compiled.EvalWithCustomFuncs(ctx, all, env)
		if err != nil {
			return fmt.Errorf("eval: %w", err)
		}
		if err := p.print(result); err != nil {
			return err
		}
	} else {
		for _, src := range sources {
			if err := src.eachValue(opts.rawInput, func(v any) error {
				result, err := compiled.EvalWithCustomFuncs(ctx, v, env)
				if err != nil {
					return fmt.Errorf("eval: %w", err)
				}
				return p.print(result)
			}); err != nil {
				return err
			}
		}
	}

	return p.checkExit()
}

// ─── input abstraction ────────────────────────────────────────────────────

type dataSource struct {
	open func() (io.ReadCloser, error)
}

func buildSources(opts evalOpts, fileArgs []string) []dataSource {
	switch {
	case opts.dataStr != "":
		return []dataSource{stringSource(opts.dataStr)}
	case opts.dataFile != "":
		return []dataSource{fileSource(opts.dataFile)}
	case len(fileArgs) > 0:
		srcs := make([]dataSource, len(fileArgs))
		for i, p := range fileArgs {
			srcs[i] = fileSource(p)
		}
		return srcs
	default:
		return []dataSource{stdinSource()}
	}
}

// eachValue iterates over every JSON value (or raw line) in the source,
// calling fn for each one.
func (ds dataSource) eachValue(rawInput bool, fn func(any) error) error {
	rc, err := ds.open()
	if err != nil {
		return err
	}
	defer rc.Close()

	if rawInput {
		sc := bufio.NewScanner(rc)
		for sc.Scan() {
			if err := fn(sc.Text()); err != nil {
				return err
			}
		}
		return sc.Err()
	}

	dec := json.NewDecoder(rc)
	for {
		var v any
		if err := dec.Decode(&v); err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("decode JSON: %w", err)
		}
		if err := fn(v); err != nil {
			return err
		}
	}
}

func stringSource(s string) dataSource {
	return dataSource{open: func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(s)), nil
	}}
}

func fileSource(path string) dataSource {
	return dataSource{open: func() (io.ReadCloser, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", path, err)
		}
		return f, nil
	}}
}

func stdinSource() dataSource {
	return dataSource{open: func() (io.ReadCloser, error) {
		return io.NopCloser(os.Stdin), nil
	}}
}

// ─── output formatting ────────────────────────────────────────────────────

type printer struct {
	opts     evalOpts
	lastVal  any
	hasValue bool
}

func (p *printer) print(v any) error {
	p.lastVal = v
	p.hasValue = true

	if p.opts.rawOutput {
		if s, ok := v.(string); ok {
			fmt.Print(s)
			if !p.opts.joinOutput {
				fmt.Println()
			}
			return nil
		}
	}

	var (
		out []byte
		err error
	)
	if p.opts.compact {
		out, err = json.Marshal(v)
	} else {
		ind := strings.Repeat(" ", clampIndent(p.opts.indentN))
		if p.opts.tabIndent {
			ind = "\t"
		}
		out, err = json.MarshalIndent(v, "", ind)
	}
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	fmt.Print(string(out))
	if !p.opts.joinOutput {
		fmt.Println()
	}
	return nil
}

func (p *printer) checkExit() error {
	if !p.opts.exitStatus {
		return nil
	}
	if !p.hasValue || !isTruthy(p.lastVal) {
		os.Exit(5)
	}
	return nil
}

func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func clampIndent(n int) int {
	if n < 0 {
		return 0
	}
	if n > 7 {
		return 7
	}
	return n
}

// ─── variable injection ───────────────────────────────────────────────────

// injectVars wraps expr with JSONata variable bindings produced from
// --arg name=value and --argjson name=json flags.
func injectVars(expr string, argVars, argjsonVars []string) (string, error) {
	var bindings []string

	for _, kv := range argVars {
		name, val, ok := strings.Cut(kv, "=")
		if !ok || name == "" {
			return "", fmt.Errorf("--arg: expected name=value, got %q", kv)
		}
		encoded, _ := json.Marshal(val)
		bindings = append(bindings, fmt.Sprintf("$%s := %s", name, string(encoded)))
	}

	for _, kv := range argjsonVars {
		name, val, ok := strings.Cut(kv, "=")
		if !ok || name == "" {
			return "", fmt.Errorf("--argjson: expected name=json, got %q", kv)
		}
		if !json.Valid([]byte(val)) {
			return "", fmt.Errorf("--argjson %s: invalid JSON: %q", name, val)
		}
		bindings = append(bindings, fmt.Sprintf("$%s := %s", name, val))
	}

	if len(bindings) == 0 {
		return expr, nil
	}
	return fmt.Sprintf("(%s; %s)", strings.Join(bindings, "; "), expr), nil
}

// ─── list subcommand ──────────────────────────────────────────────────────

func buildListCmd() *cobra.Command {
	var pkg string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available extension functions",
		Long: `List all gnata-ext functions, grouped by package.

Use --package / -p to filter to a specific package (e.g. extarray, extnumeric).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			catalog := ext.Catalog()
			if pkg != "" {
				var filtered []ext.FuncMeta
				for _, f := range catalog {
					if strings.EqualFold(f.Package, pkg) {
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
			return nil
		},
	}

	cmd.Flags().StringVarP(&pkg, "package", "p", "", "filter by package name (e.g. extarray)")
	return cmd
}

// ─── describe subcommand ──────────────────────────────────────────────────

func buildDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe <funcName>",
		Short: "Show signature and description for an extension function",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			name := strings.TrimPrefix(args[0], "$")
			for _, f := range ext.Catalog() {
				if f.Name == name {
					fmt.Printf("Function:    $%s\n", f.Name)
					fmt.Printf("Package:     %s\n", f.Package)
					fmt.Printf("Signature:   $%s%s\n", f.Name, f.Signature)
					fmt.Printf("Description: %s\n", f.Description)
					return nil
				}
			}
			return fmt.Errorf("unknown function %q", name)
		},
	}
}

// ─── version subcommand ───────────────────────────────────────────────────

func buildVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print jn version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("jn version %s\n", Version)
		},
	}
}
