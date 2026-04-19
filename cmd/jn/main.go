// jn — JSONata processor with gnata-ext extension functions.
//
// Usage (inspired by jn):
//
//	jn [flags] [expr] [file...]
//	jn list [--package <pkg>]
//	jn describe <funcName>
//	jn version
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/mattn/go-colorable"
	json "github.com/neilotoole/jsoncolor"
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
	rawOutput0  bool // output raw strings with NUL terminator instead of newline
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
	colorOutput bool // -C: force color even when not a terminal
	monoOutput  bool // -M: disable color output
	unbuffered  bool // flush output after each JSON object
}

func buildRoot() *cobra.Command {
	var opts evalOpts

	root := &cobra.Command{
		Use:   "jn [flags] [expr] [file...]",
		Short: "JSONata processor with extended functions",
		Long: `jn — JSONata expression processor (inspired by jn)

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

  # colored output (auto-detected when stdout is a terminal)
  echo '{"x":1}' | jn -C '$'

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
	f.BoolVar(&opts.rawOutput0, "raw-output0", false, "like -r but use NUL as separator instead of newline")
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
	f.BoolVarP(&opts.colorOutput, "color-output", "C", false, "force colorized output even when not writing to a terminal")
	f.BoolVarP(&opts.monoOutput, "monochrome-output", "M", false, "disable colorized output")
	f.BoolVar(&opts.unbuffered, "unbuffered", false, "flush output after each JSON object is printed")

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
	p := newPrinter(opts)

	// Null-input mode: evaluate once with nil
	if opts.nullInput {
		result, err := compiled.EvalWithCustomFuncs(ctx, nil, env)
		if err != nil {
			return fmt.Errorf("eval: %w", err)
		}
		if err := p.print(result); err != nil {
			return err
		}
		if err := p.flush(); err != nil {
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

	if err := p.flush(); err != nil {
		return err
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
	opts      evalOpts
	lastVal   any
	hasValue  bool
	w         io.Writer     // effective output writer (may be buffered)
	bw        *bufio.Writer // non-nil when using buffered output
	colorMode bool          // whether ANSI color is active
	colors    *json.Colors  // color config (nil = no color)
	indent    string        // indentation string (empty when compact)
}

// newPrinter constructs a printer configured from opts.
// Color is auto-detected from the terminal unless overridden by -C/-M or NO_COLOR.
func newPrinter(opts evalOpts) *printer {
	p := &printer{opts: opts}

	// Resolve indent string
	if !opts.compact {
		if opts.tabIndent {
			p.indent = "\t"
		} else {
			p.indent = strings.Repeat(" ", clampIndent(opts.indentN))
		}
	}

	// Resolve color mode:
	//   NO_COLOR env (non-empty) disables color unless -C is explicit
	//   -M disables color
	//   -C forces color
	//   otherwise auto-detect from terminal
	noColorEnv := os.Getenv("NO_COLOR") != ""
	switch {
	case opts.monoOutput:
		p.colorMode = false
	case opts.colorOutput:
		p.colorMode = true
	case noColorEnv:
		p.colorMode = false
	default:
		p.colorMode = json.IsColorTerminal(os.Stdout)
	}

	// Build the output writer
	var baseWriter io.Writer
	if p.colorMode {
		baseWriter = colorable.NewColorable(os.Stdout)
	} else {
		baseWriter = os.Stdout
	}
	p.bw = bufio.NewWriter(baseWriter)
	p.w = p.bw

	// Build color config (respects JN_COLORS env variable)
	if p.colorMode {
		p.colors = resolveColors()
	}

	return p
}

// resolveColors returns a Colors config, optionally customised via JN_COLORS.
// JN_COLORS format (colon-separated, 8 fields): null:false:true:numbers:strings:arrays:objects:keys
// Each field is a terminal escape code like "1;31".
func resolveColors() *json.Colors {
	clrs := json.DefaultColors()
	jnColors := os.Getenv("JN_COLORS")
	if jnColors == "" {
		return clrs
	}
	parts := strings.Split(jnColors, ":")
	applyColor := func(idx int, target *json.Color) {
		if idx < len(parts) && parts[idx] != "" {
			*target = json.Color("\x1b[" + parts[idx] + "m")
		}
	}
	applyColor(0, &clrs.Null)   // null
	applyColor(1, &clrs.Bool)   // false (true uses same field)
	applyColor(3, &clrs.Number) // numbers
	applyColor(4, &clrs.String) // strings
	applyColor(5, &clrs.Punc)   // arrays  (brackets)
	applyColor(7, &clrs.Key)    // object keys
	return clrs
}

func (p *printer) print(v any) error {
	p.lastVal = v
	p.hasValue = true

	// Raw output: emit the string value without JSON encoding
	if p.opts.rawOutput || p.opts.rawOutput0 {
		if s, ok := v.(string); ok {
			if _, err := fmt.Fprint(p.w, s); err != nil {
				return err
			}
			switch {
			case p.opts.rawOutput0:
				if _, err := p.w.Write([]byte{0}); err != nil {
					return err
				}
			case !p.opts.joinOutput:
				if _, err := p.w.Write([]byte{'\n'}); err != nil {
					return err
				}
			}
			return p.flushIfUnbuffered()
		}
		// Non-string values fall through to normal JSON encoding
	}

	// Encode to JSON (with optional color)
	out, err := p.encodeJSON(v)
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	if _, err := p.w.Write(out); err != nil {
		return err
	}
	if !p.opts.joinOutput {
		if _, err := p.w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	return p.flushIfUnbuffered()
}

// encodeJSON encodes v to JSON bytes, applying color and indentation as configured.
// The trailing newline that json.Encoder.Encode appends is stripped.
func (p *printer) encodeJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if p.colorMode && p.colors != nil {
		enc.SetColors(p.colors)
	}
	if !p.opts.compact && p.indent != "" {
		enc.SetIndent("", p.indent)
	}
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	// json.Encoder.Encode always appends '\n'; strip it so callers control line endings
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

// flushIfUnbuffered flushes the buffered writer immediately when --unbuffered is set.
func (p *printer) flushIfUnbuffered() error {
	if p.opts.unbuffered && p.bw != nil {
		return p.bw.Flush()
	}
	return nil
}

// flush flushes any buffered output. Must be called after all printing is done.
func (p *printer) flush() error {
	if p.bw != nil {
		return p.bw.Flush()
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
