package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

type UserPostgresRepository struct {
	queries *db.Queries
}

func NewUserPostgresRepository(queries *db.Queries) *UserPostgresRepository {
	if queries == nil {
		panic("queries cannot be nil")
	}
	return &UserPostgresRepository{queries}
}

func (ur *UserPostgresRepository) FindById(ctx context.Context, id int64) (user.Model, error) {
	out, err := ur.queries.GetUserById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Model{}, errorDB.ErrNoRows
	} else if err != nil {
		return user.Model{}, err
	}
	return toModel(out)
}

func (ur *UserPostgresRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (user.Model, error) {
	out, err := ur.queries.GetUserByOIDCSub(ctx, pgtype.UUID{Bytes: oidcSub, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Model{}, errorDB.ErrNoRows
	} else if err != nil {
		return user.Model{}, err
	}
	return toModel(out)
}

func (ur *UserPostgresRepository) FindByName(ctx context.Context, name string) (user.Model, error) {
	out, err := ur.queries.GetUserByName(ctx, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Model{}, errorDB.ErrNoRows
	} else if err != nil {
		return user.Model{}, err
	}
	return toModel(out)
}

func (ur *UserPostgresRepository) CreateUser(ctx context.Context, user user.Model) (int64, error) {
	arr := []db.CreateUserParams{
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

func toModel(dbUser db.User) (user.Model, error) {
	sub, err := uuid.FromBytes(dbUser.OidcSub.Bytes[:])
	return user.Model{
		ID:      dbUser.ID,
		Name:    dbUser.Name,
		OidcSub: sub,
	}, err
}
