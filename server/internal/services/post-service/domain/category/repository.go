package category

import (
	"context"
	"encoding/json"
)

type MetadataSchema json.RawMessage

type Entity struct {
	Name           string
	MetadataSchema MetadataSchema
}

type Repository interface {
	FindCategoryByName(ctx context.Context, name string) (Entity, error)
	CreateCategory(ctx context.Context, category Entity) (Entity, error)
}
