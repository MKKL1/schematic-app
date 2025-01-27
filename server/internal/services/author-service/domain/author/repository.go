package author

import (
	"context"
)

type Entity struct {
	ID       int64
	Name     string
	UserID   int64
	Metadata map[string]string
}

type Repository interface {
	FindByID(ctx context.Context, id int64) (Entity, error)
	FindByName(ctx context.Context, name string) (Entity, error)
	Create(ctx context.Context, author Entity) error
}
