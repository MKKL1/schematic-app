package category

import (
	"context"
	"encoding/json"
)

type Entity struct {
	Name             string
	ValueDefinitions json.RawMessage
}

type PostCategoryVars struct {
	PostID   int64
	Category string
	Values   json.RawMessage
}

type Repository interface {
	FindCategoryByID(ctx context.Context, name string) (Entity, error)
	FindCategVarsByPostID(ctx context.Context, postID int64) ([]PostCategoryVars, error)
	CreateCategory(ctx context.Context, category Entity) (Entity, error)
	CreatePostCategory(ctx context.Context, pcv PostCategoryVars) error
}
