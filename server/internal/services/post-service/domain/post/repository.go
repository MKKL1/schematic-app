package post

import (
	"context"
)

type Repository interface {
	FindById(ctx context.Context, id int64) (Post, error)
	Create(ctx context.Context, model Post) error
	GetCountForTag(ctx context.Context, tag string) (int64, error)
}
