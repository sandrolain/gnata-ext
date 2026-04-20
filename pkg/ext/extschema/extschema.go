// Package extschema provides JSON Schema validation functions for gnata.
//
// Functions:
//
//   - $validateSchema(data, schema) – validate data and return {valid, errors}
//   - $isValid(data, schema)        – true if data is valid according to schema
//   - $schemaErrors(data, schema)   – array of validation error strings (empty if valid)
package extschema

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/recolabs/gnata"
)

// All returns a map of all extschema functions.
func All() map[string]gnata.CustomFunc {
	return map[string]gnata.CustomFunc{
		"validateSchema": ValidateSchema(),
		"isValid":        IsValid(),
		"schemaErrors":   SchemaErrors(),
	}
}

// ValidateSchema returns the CustomFunc for $validateSchema(data, schema).
// Returns {valid: bool, errors: [string]}.
func ValidateSchema() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		data, schema, err := parseArgs(args, "$validateSchema")
		if err != nil {
			return nil, err
		}
		errs := validateAgainstSchema(data, schema)
		return map[string]any{
			"valid":  len(errs) == 0,
			"errors": errs,
		}, nil
	}
}

// IsValid returns the CustomFunc for $isValid(data, schema).
// Returns true if data is valid according to the JSON Schema.
func IsValid() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		data, schema, err := parseArgs(args, "$isValid")
		if err != nil {
			return nil, err
		}
		errs := validateAgainstSchema(data, schema)
		return len(errs) == 0, nil
	}
}

// SchemaErrors returns the CustomFunc for $schemaErrors(data, schema).
// Returns an array of validation error strings. Empty array means valid.
func SchemaErrors() gnata.CustomFunc {
	return func(args []any, _ any) (any, error) {
		data, schema, err := parseArgs(args, "$schemaErrors")
		if err != nil {
			return nil, err
		}
		return validateAgainstSchema(data, schema), nil
	}
}

// --- helpers ---

// parseArgs validates that both arguments are present and returns them.
func parseArgs(args []any, name string) (data any, schema map[string]any, err error) {
	if len(args) < 2 {
		return nil, nil, fmt.Errorf("%s: requires 2 arguments (data, schema)", name)
	}
	schemaObj, ok := args[1].(map[string]any)
	if !ok {
		return nil, nil, fmt.Errorf("%s: schema must be an object", name)
	}
	return args[0], schemaObj, nil
}

// validateAgainstSchema compiles the schema object and validates the data.
// Returns a []any of error message strings so it can be embedded in gnata results.
func validateAgainstSchema(data any, schemaObj map[string]any) []any {
	// Round-trip through JSON so the jsonschema library can handle it.
	schemaBytes, err := json.Marshal(schemaObj)
	if err != nil {
		return []any{fmt.Sprintf("schema marshal error: %s", err)}
	}

	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	const schemaURI = "schema.json"
	if err := compiler.AddResource(schemaURI, strings.NewReader(string(schemaBytes))); err != nil {
		return []any{fmt.Sprintf("schema compile error: %s", err)}
	}
	compiled, err := compiler.Compile(schemaURI)
	if err != nil {
		return []any{fmt.Sprintf("schema compile error: %s", err)}
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return []any{fmt.Sprintf("data marshal error: %s", err)}
	}

	var unmarshalled any
	if err := json.Unmarshal(dataBytes, &unmarshalled); err != nil {
		return []any{fmt.Sprintf("data unmarshal error: %s", err)}
	}

	if err := compiled.Validate(unmarshalled); err != nil {
		var ve *jsonschema.ValidationError
		if ok := isValidationError(err, &ve); ok {
			msgs := collectErrors(ve)
			out := make([]any, len(msgs))
			for i, m := range msgs {
				out[i] = m
			}
			return out
		}
		return []any{err.Error()}
	}
	return []any{}
}

func isValidationError(err error, target **jsonschema.ValidationError) bool {
	if ve, ok := err.(*jsonschema.ValidationError); ok {
		*target = ve
		return true
	}
	return false
}

func collectErrors(ve *jsonschema.ValidationError) []string {
	var msgs []string
	if ve.Message != "" {
		msgs = append(msgs, ve.InstanceLocation+": "+ve.Message)
	}
	for _, cause := range ve.Causes {
		msgs = append(msgs, collectErrors(cause)...)
	}
	return msgs
}
