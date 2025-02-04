package category

import "encoding/json"

type SchemaProvider interface {
	// GetValidator compiles (or retrieves from cache) the schema from its raw JSON.
	GetValidator(raw json.RawMessage) (ValueSchema, error)
}

type DefaultSchemaProvider struct {
}

func (p DefaultSchemaProvider) GetValidator(raw json.RawMessage) (ValueSchema, error) {
	compiled, err := CompileSchema(raw)
	return compiled, err
}
