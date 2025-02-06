package post

import (
	"context"
)

type Entity struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Entity, error)
	Create(ctx context.Context, model Entity) error
}
