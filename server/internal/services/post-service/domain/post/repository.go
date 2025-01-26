package post

import (
	"context"
)

type Model struct {
	ID          int64
	Description *string
	Owner       int64
	AuthorName  *string
	AuthorID    *int64
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Model, error)
}
