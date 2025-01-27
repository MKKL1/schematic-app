package user

import (
	"context"
	"github.com/google/uuid"
)

type Entity struct {
	ID      int64
	Name    string
	OidcSub uuid.UUID
}

func EntityToDTO(entity Entity) User {
	return User{
		ID:      ID(entity.ID),
		Name:    entity.Name,
		OidcSub: entity.OidcSub,
	}
}

type Repository interface {
	FindById(ctx context.Context, id int64) (Entity, error)
	FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (Entity, error)
	FindByName(ctx context.Context, name string) (Entity, error)
	CreateUser(ctx context.Context, user Entity) (int64, error)
}
