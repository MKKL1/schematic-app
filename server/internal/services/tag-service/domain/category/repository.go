package category

import "context"

type ValueSchema map[string]interface{}

type Entity struct {
	ID               int64
	Name             string
	ValueDefinitions ValueSchema
}

type PostCategoryValue struct {
	PostID     int64
	CategoryID int64
	Values     pq.Jsonb
}

type WithValues struct {
	Entity
	Values pq.Jsonb
}

type Repository interface {
	FindCategoryByID(ctx context.Context, id int64) (Entity, error)
	FindCategoryByName(ctx context.Context, name string) (Entity, error)
	CreateCategory(ctx context.Context, category Entity) (Entity, error)
	CreatePostCategory(ctx context.Context, pcv PostCategoryValue) (PostCategoryValue, error)
	GetPostCategory(ctx context.Context, postID int64, categoryID int64) (PostCategoryValue, error)
	GetPostsByJSONValue(ctx context.Context, categoryID int64, jsonPath string) ([]PostCategoryValue, error)
	ListCategoriesForPost(ctx context.Context, postID int64) ([]WithValues, error)
}
