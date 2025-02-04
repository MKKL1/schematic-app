package category

import (
	"context"
	"encoding/json"
)

type Entity struct {
	ID               int64
	Name             string
	ValueDefinitions json.RawMessage
}

type PostCategoryValue struct {
	PostID     int64
	CategoryID int64
	Values     json.RawMessage
}

type WithValues struct {
	Entity
	Values json.RawMessage
}

type Repository interface {
	FindCategoryByID(ctx context.Context, id int64) (Entity, error)
	CreateCategory(ctx context.Context, category Entity) (Entity, error)
	CreatePostCategory(ctx context.Context, pcv PostCategoryValue) (PostCategoryValue, error)
}
