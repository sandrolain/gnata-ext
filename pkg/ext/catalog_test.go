package ext_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext"
)

func TestCatalog_NotEmpty(t *testing.T) {
	c := ext.Catalog()
	if len(c) == 0 {
		t.Fatal("Catalog() must not be empty")
	}
}

func TestCatalog_Sorted(t *testing.T) {
	c := ext.Catalog()
	for i := 1; i < len(c); i++ {
		if c[i].Name < c[i-1].Name {
			t.Fatalf("Catalog() not sorted: %q before %q", c[i-1].Name, c[i].Name)
		}
	}
}

func TestCatalog_AllHavePackage(t *testing.T) {
	for _, f := range ext.Catalog() {
		if f.Package == "" {
			t.Errorf("function %q has empty Package", f.Name)
		}
		if f.Signature == "" {
			t.Errorf("function %q has empty Signature", f.Name)
		}
		if f.Description == "" {
			t.Errorf("function %q has empty Description", f.Name)
		}
	}
}

func TestCatalogByPackage(t *testing.T) {
	m := ext.CatalogByPackage()
	expectedPkgs := []string{
		"extarray", "extcrypto", "extdatetime", "extformat",
		"extgeo", "extjson", "extnet", "extnumeric",
		"extobject", "extpath", "extstring", "exttypes", "extvalidate",
	}
	for _, pkg := range expectedPkgs {
		if len(m[pkg]) == 0 {
			t.Errorf("CatalogByPackage() missing entries for %q", pkg)
		}
	}
}

func TestCatalog_KnownFunctions(t *testing.T) {
	byName := make(map[string]ext.FuncMeta)
	for _, f := range ext.Catalog() {
		byName[f.Name] = f
	}
	known := []string{"first", "uuid", "dateAdd", "csv", "log", "values", "startsWith", "isString", "get", "isEmail", "jsonParse", "haversine", "ipVersion"}
	for _, name := range known {
		if _, ok := byName[name]; !ok {
			t.Errorf("Catalog() missing function %q", name)
		}
	}
}
