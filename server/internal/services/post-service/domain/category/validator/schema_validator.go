package validator

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
)

type Validator interface {
	// Validate returns an error if the value does not conform.
	Validate(value interface{}) (interface{}, error)
}

// ValidatorBuilder creates a Validator from a property definition node.
type ValidatorBuilder func(key string, node *ast.Node) (Validator, error)

var validatorBuilders = map[string]ValidatorBuilder{
	"boolean": NewBooleanValidator,
	"enum":    NewEnumValidator,
	"range":   NewRangeValidator,
}

func CompileSchema(schema []byte) (MetadataSchemaValidator, error) {
	// Parse the raw JSON schema into an AST node.
	root, err := sonic.Get(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	// Interpret the root as a JSON object mapping property names to definitions.
	props, err := root.MapUseNode()
	if err != nil {
		return nil, fmt.Errorf("schema must be a JSON object: %w", err)
	}

	compiled := make(MetadataSchemaValidator)
	for key, node := range props {
		// Extract the "type" field.
		typeNode := node.Get("type")
		if !typeNode.Exists() {
			return nil, fmt.Errorf("missing required field 'type' for key %s", key)
		}
		t, err := typeNode.String()
		if err != nil {
			return nil, fmt.Errorf("invalid 'type' for key %s: %w", key, err)
		}

		// Look up the builder function.
		builder, ok := validatorBuilders[t]
		if !ok {
			return nil, fmt.Errorf("unknown type %q for key %s", t, key)
		}
		validator, err := builder(key, &node)
		if err != nil {
			return nil, err
		}
		compiled[key] = validator
	}
	return compiled, nil
}

// MetadataSchemaValidator maps a property name to its Validator.
type MetadataSchemaValidator map[string]Validator

// ValidateData uses the compiled schema to check the provided data.
func (vs MetadataSchemaValidator) ValidateData(data map[string]interface{}) (map[string]interface{}, error) {
	cleanedData := make(map[string]interface{})
	for key, validator := range vs {
		value, exists := data[key]
		if !exists {
			return nil, fmt.Errorf("missing key: %s", key)
		}
		cleanValue, err := validator.Validate(value)
		if err != nil {
			return nil, fmt.Errorf("validation failed for key %s: %w", key, err)
		}
		cleanedData[key] = cleanValue
	}
	return cleanedData, nil
}

func (vs MetadataSchemaValidator) Map() map[string]Validator {
	return vs
}
