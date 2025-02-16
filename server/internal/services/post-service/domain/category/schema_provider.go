package category

import "github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category/validator"

type SchemaProvider interface {
	// GetValidator compiles (or retrieves from cache) the schema from its raw JSON.
	GetValidator(raw MetadataSchema) (validator.MetadataSchemaValidator, error)
}

type DefaultSchemaProvider struct {
}

func (p DefaultSchemaProvider) GetValidator(raw MetadataSchema) (validator.MetadataSchemaValidator, error) {
	compiled, err := validator.CompileSchema(raw)
	return compiled, err
}
