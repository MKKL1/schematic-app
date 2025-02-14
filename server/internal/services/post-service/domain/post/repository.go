package post

import (
	"context"
)

type CreatePostParams struct {
	ID          int64
	Name        string
	Description *string
	AuthorID    *int64
	Owner       int64
	Categories  []CreatePostCategoryParams
	Tags        []string
}

type CreatePostCategoryParams struct {
	Name     string
	Metadata map[string]interface{}
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Post, error)
	Create(ctx context.Context, params CreatePostParams) error
	GetCountForTag(ctx context.Context, tag string) (int64, error)
}
