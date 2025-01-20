package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	db2 "github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

type UserPostgresRepository struct {
	queries *db2.Queries
}

func NewUserPostgresRepository(queries *db2.Queries) *UserPostgresRepository {
	if queries == nil {
		panic("queries cannot be nil")
	}
	return &UserPostgresRepository{queries}
}

func (ur *UserPostgresRepository) FindById(ctx context.Context, id user.UserID) (db2.User, error) {
	out, err := ur.queries.GetUserById(ctx, int64(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return out, errorDB.ErrNoRows
	}
	return out, err
}

func (ur *UserPostgresRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (db2.User, error) {
	out, err := ur.queries.GetUserByOIDCSub(ctx, pgtype.UUID{Bytes: oidcSub, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return out, errorDB.ErrNoRows
	}
	return out, err
}

func (ur *UserPostgresRepository) FindByName(ctx context.Context, name string) (db2.User, error) {
	out, err := ur.queries.GetUserByName(ctx, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return out, errorDB.ErrNoRows
	}
	return out, err
}

func (ur *UserPostgresRepository) CreateUser(ctx context.Context, user user.User) (int64, error) {
	arr := []db2.CreateUserParams{
		{
			ID:      int64(user.ID),
			Name:    user.Name,
			OidcSub: pgtype.UUID{Bytes: user.OidcSub, Valid: true},
		},
	}

	created, err := ur.queries.CreateUser(ctx, arr)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				switch {
				case strings.Contains(pgErr.ConstraintName, "users_name_unique"):
					return 0, errorDB.NewUniqueConstraintError("Name")
				case strings.Contains(pgErr.ConstraintName, "users_oidc_sub_unique"):
					return 0, errorDB.NewUniqueConstraintError("OidcSub")
				}
			}
		}
		return created, err
	}

	return created, nil
}
