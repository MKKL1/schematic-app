package category

type SchemaProvider interface {
	// GetValidator compiles (or retrieves from cache) the schema from its raw JSON.
	GetValidator(raw MetadataSchema) (ValueSchemaValidator, error)
}

type DefaultSchemaProvider struct {
}

func (p DefaultSchemaProvider) GetValidator(raw MetadataSchema) (ValueSchemaValidator, error) {
	compiled, err := CompileSchema(raw)
	return compiled, err
}
