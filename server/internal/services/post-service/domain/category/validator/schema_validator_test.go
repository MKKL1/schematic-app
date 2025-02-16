package validator

import "testing"

const customSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "afkable": {
      "type": "boolean"
    },
    "mob_type": {
      "type": "string",
      "enum": ["zombie", "skeleton", "spider", "creeper"]
    },
    "spawn_rate": {
      "type": "number",
      "minimum": 100,
      "maximum": 500
    }
  },
  "required": ["afkable", "mob_type", "spawn_rate"],
  "additionalProperties": false
}`

func TestSchemaValidator_ValidDocument(t *testing.T) {
	validator := NewSchemaValidator([]byte(customSchema))
	if validator == nil {
		t.Fatal("Failed to create SchemaValidator")
	}

	// A valid document that meets the schema requirements.
	validDoc := map[string]interface{}{
		"afkable":    false,
		"mob_type":   "spider",
		"spawn_rate": 342,
	}

	err := validator.Validate(validDoc)
	if err != nil {
		t.Fatalf("expected valid document, got error: %v", err)
	}
}

func TestSchemaValidator_InvalidDocument(t *testing.T) {
	validator := NewSchemaValidator([]byte(customSchema))
	if validator == nil {
		t.Fatal("Failed to create SchemaValidator")
	}

	// An invalid document: spawn_rate is out of range and there's an extra field.
	invalidDoc := map[string]interface{}{
		"afkable":    false,
		"mob_type":   "spider",
		"spawn_rate": 50, // Below minimum (should be >= 100)
		"extra":      "should not be here",
	}

	err := validator.Validate(invalidDoc)
	if err == nil {
		t.Fatal("expected invalid document to return an error, got nil")
	} else {
		t.Logf("Validation error as expected: %v", err)
	}
}

func TestSchemaValidator_FieldErrors(t *testing.T) {
	validator := NewSchemaValidator([]byte(customSchema))
	if validator == nil {
		t.Fatal("Failed to create SchemaValidator")
	}

	// An invalid document: spawn_rate is below the minimum and extra field is provided.
	invalidDoc := map[string]interface{}{
		"afkable":    false,
		"mob_type":   "spider",
		"spawn_rate": 50, // Should be >= 100.
		"extra":      "not allowed",
	}

	err := validator.Validate(invalidDoc)
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("Expected a ValidationError, got: %T", err)
	}

	for _, fe := range ve.Errors {
		t.Logf("Field %q error: %s", fe.Field, fe.Message)
	}
}

func BenchmarkSchemaValidator_ValidDocument(b *testing.B) {
	validator := NewSchemaValidator([]byte(customSchema))
	if validator == nil {
		b.Fatal("Failed to create SchemaValidator")
	}

	validDoc := map[string]interface{}{
		"afkable":    false,
		"mob_type":   "spider",
		"spawn_rate": 342,
	}

	err := validator.Validate(validDoc)
	if err != nil {
		b.Fatalf("expected valid document, got error: %v", err)
	}
}
