package user

import (
	"context"
	"github.com/google/uuid"
)

//TODO it probably shouldn't be in domain

type Model struct {
	ID      int64
	Name    string
	OidcSub uuid.UUID
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Model, error)
	FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (Model, error)
	FindByName(ctx context.Context, name string) (Model, error)
	CreateUser(ctx context.Context, user Model) (int64, error)
}
