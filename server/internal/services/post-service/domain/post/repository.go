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
	Files       []CreatePostFileParams
}

type CreatePostCategoryParams struct {
	Name     string
	Metadata map[string]interface{}
}

type CreatePostFileParams struct {
	Name   string `json:"name"` //json for database query, may need to move it to infra
	TempId string `json:"temp_id"`
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Post, error)
	Create(ctx context.Context, params CreatePostParams) error
	GetCountForTag(ctx context.Context, tag string) (int64, error)
}
