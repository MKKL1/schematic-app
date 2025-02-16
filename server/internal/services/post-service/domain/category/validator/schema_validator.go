package validator

import (
	"fmt"
)
import "github.com/xeipuuv/gojsonschema"

type SchemaValidator struct {
	Schema *gojsonschema.Schema
}

func NewSchemaValidator(schemaRaw []byte) *SchemaValidator {
	schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemaRaw))
	if err != nil {
		return nil
	}
	return &SchemaValidator{
		Schema: schema,
	}
}

func (sv *SchemaValidator) Validate(content interface{}) error {
	documentLoader := gojsonschema.NewGoLoader(content)
	result, err := sv.Schema.Validate(documentLoader)
	if err != nil {
		return fmt.Errorf("error validating document: %w", err)
	}
	if result == nil {
		return fmt.Errorf("result is nil")
	}
	if !result.Valid() {
		var fieldErrors []FieldError
		for _, desc := range result.Errors() {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   desc.Field(),
				Message: desc.Description(),
			})
		}
		return &ValidationError{Errors: fieldErrors}
	}
	return nil
}
