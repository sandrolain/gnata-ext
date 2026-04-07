# FuncCatalog

Discover and introspect registered functions at runtime.

---

## Basic Usage

`Catalog()` returns metadata for every registered function; useful for documentation, tooling, or runtime introspection:

```go
import (
    "fmt"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

// List all functions
for _, f := range ext.Catalog() {
    fmt.Printf("$%s %s — %s\n", f.Name, f.Signature, f.Description)
}
```

---

## Group by Package

```go
// Group by package
byPkg := ext.CatalogByPackage()
for _, meta := range byPkg["extarray"] {
    fmt.Printf("  $%s — %s\n", meta.Name, meta.Description)
}

// List all packages
for pkg, funcs := range byPkg {
    fmt.Printf("%s (%d funcs)\n", pkg, len(funcs))
}
```

---

## FuncMeta Fields

Each catalog entry carries these fields:

```go
type FuncMeta struct {
    Name        string // Function name (e.g. "uuid")
    Package     string // Package (e.g. "extcrypto")
    Signature   string // Function signature (e.g. "uuid() -> string")
    Description string // Human-readable description
}
```

---

## Common Patterns

### Documentation generation

```go
import (
    "os"
    "text/template"
    "github.com/sandrolain/gnata-ext/pkg/ext"
)

tmpl := template.Must(template.New("funcs").Parse(`
{{- range .}}
- **${{.Name}}** ({{.Package}})
  {{.Description}}
{{- end}}
`))

tmpl.Execute(os.Stdout, ext.Catalog())
```

### Tooling and IDE support

```go
catalog := ext.CatalogByPackage()

// Provide autocomplete suggestions
suggestions := catalog["extstring"]
for _, f := range suggestions {
    // Feed to autocomplete engine
}
```

### Schema validation

```go
registered := make(map[string]bool)
for _, f := range ext.Catalog() {
    registered[f.Name] = true
}

// Check if a function is registered
if !registered["hash"] {
    // Handle missing function
}
```
