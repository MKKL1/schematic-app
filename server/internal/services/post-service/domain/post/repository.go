package post

import (
	"context"
)

type Model struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorName  *string
	AuthorID    *int64
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Model, error)
	Create(ctx context.Context, model Model) error
}
