package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	db2 "github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/postgres/db"
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

func (ur *UserPostgresRepository) FindById(ctx context.Context, id int64) (user.Entity, error) {
	out, err := ur.queries.GetUserById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Entity{}, errorDB.ErrNoRows
	} else if err != nil {
		return user.Entity{}, err
	}
	return toModel(out)
}

func (ur *UserPostgresRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (user.Entity, error) {
	out, err := ur.queries.GetUserByOIDCSub(ctx, pgtype.UUID{Bytes: oidcSub, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Entity{}, errorDB.ErrNoRows
	} else if err != nil {
		return user.Entity{}, err
	}
	return toModel(out)
}

func (ur *UserPostgresRepository) FindByName(ctx context.Context, name string) (user.Entity, error) {
	out, err := ur.queries.GetUserByName(ctx, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Entity{}, errorDB.ErrNoRows
	} else if err != nil {
		return user.Entity{}, err
	}
	return toModel(out)
}

func (ur *UserPostgresRepository) CreateUser(ctx context.Context, user user.Entity) (int64, error) {
	arr := []db2.CreateUserParams{
		{
			ID:      user.ID,
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

func toModel(dbUser db2.User) (user.Entity, error) {
	sub, err := uuid.FromBytes(dbUser.OidcSub.Bytes[:])
	return user.Entity{
		ID:      dbUser.ID,
		Name:    dbUser.Name,
		OidcSub: sub,
	}, err
}
