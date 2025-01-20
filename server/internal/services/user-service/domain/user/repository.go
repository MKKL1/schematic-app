package user

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/google/uuid"
)

type Repository interface {
	FindById(ctx context.Context, id UserID) (db.User, error)
	FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (db.User, error)
	FindByName(ctx context.Context, name string) (db.User, error)
	CreateUser(ctx context.Context, user User) (int64, error)
}
