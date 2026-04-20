package extschema_test

import (
	"testing"

	"github.com/sandrolain/gnata-ext/pkg/ext/extschema"
)

// --- shared helpers ---

var stringSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"name": map[string]any{"type": "string"},
		"age":  map[string]any{"type": "number"},
	},
	"required": []any{"name"},
}

var validData = map[string]any{
	"name": "Alice",
	"age":  float64(30),
}

var invalidData = map[string]any{
	"age": float64(30),
	// missing required "name"
}

// ---------- ValidateSchema ----------

func TestValidateSchema_Valid(t *testing.T) {
	fn := extschema.ValidateSchema()
	got, err := fn([]any{validData, stringSchema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result := got.(map[string]any)
	if result["valid"] != true {
		t.Errorf("expected valid=true, errors=%v", result["errors"])
	}
	errs := result["errors"].([]any)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateSchema_Invalid(t *testing.T) {
	fn := extschema.ValidateSchema()
	got, err := fn([]any{invalidData, stringSchema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result := got.(map[string]any)
	if result["valid"] != false {
		t.Error("expected valid=false")
	}
	errs := result["errors"].([]any)
	if len(errs) == 0 {
		t.Error("expected at least one error")
	}
}

func TestValidateSchema_NoArgs(t *testing.T) {
	fn := extschema.ValidateSchema()
	_, err := fn([]any{validData}, nil)
	if err == nil {
		t.Error("expected error for missing schema arg")
	}
}

func TestValidateSchema_SchemaNonObject(t *testing.T) {
	fn := extschema.ValidateSchema()
	_, err := fn([]any{validData, "not-an-object"}, nil)
	if err == nil {
		t.Error("expected error for non-object schema")
	}
}

// ---------- IsValid ----------

func TestIsValid_True(t *testing.T) {
	fn := extschema.IsValid()
	got, err := fn([]any{validData, stringSchema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != true {
		t.Error("expected true")
	}
}

func TestIsValid_False(t *testing.T) {
	fn := extschema.IsValid()
	got, err := fn([]any{invalidData, stringSchema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != false {
		t.Error("expected false")
	}
}

func TestIsValid_NoArgs(t *testing.T) {
	fn := extschema.IsValid()
	_, err := fn([]any{validData}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestIsValid_SchemaNonObject(t *testing.T) {
	fn := extschema.IsValid()
	_, err := fn([]any{validData, 42}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

// ---------- SchemaErrors ----------

func TestSchemaErrors_Valid(t *testing.T) {
	fn := extschema.SchemaErrors()
	got, err := fn([]any{validData, stringSchema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	errs := got.([]any)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestSchemaErrors_Invalid(t *testing.T) {
	fn := extschema.SchemaErrors()
	got, err := fn([]any{invalidData, stringSchema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	errs := got.([]any)
	if len(errs) == 0 {
		t.Error("expected at least one error message")
	}
	for _, e := range errs {
		if _, ok := e.(string); !ok {
			t.Errorf("expected string error, got %T", e)
		}
	}
}

func TestSchemaErrors_NoArgs(t *testing.T) {
	fn := extschema.SchemaErrors()
	_, err := fn([]any{validData}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestSchemaErrors_SchemaNonObject(t *testing.T) {
	fn := extschema.SchemaErrors()
	_, err := fn([]any{validData, []any{}}, nil)
	if err == nil {
		t.Error("expected error for non-object schema")
	}
}

func TestSchemaErrors_WrongTypeConstraint(t *testing.T) {
	// Value "Alice" where number expected.
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"age": map[string]any{"type": "number"},
		},
	}
	data := map[string]any{"age": "not-a-number"}
	fn := extschema.SchemaErrors()
	got, err := fn([]any{data, schema}, nil)
	if err != nil {
		t.Fatal(err)
	}
	errs := got.([]any)
	if len(errs) == 0 {
		t.Error("expected type error")
	}
}

// ---------- All ----------

func TestAll_Keys(t *testing.T) {
	m := extschema.All()
	expected := []string{"validateSchema", "isValid", "schemaErrors"}
	for _, k := range expected {
		if _, ok := m[k]; !ok {
			t.Errorf("All(): missing key %q", k)
		}
	}
	if len(m) != len(expected) {
		t.Errorf("All(): expected %d keys, got %d", len(expected), len(m))
	}
}
