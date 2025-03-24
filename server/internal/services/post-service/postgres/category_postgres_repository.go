package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres/db"
	"github.com/jackc/pgx/v5"
)

type CategoryPostgresRepository struct {
	queries *db.Queries
}

func NewCategoryPostgresRepository(queries *db.Queries) *CategoryPostgresRepository {
	return &CategoryPostgresRepository{queries}
}

func (c CategoryPostgresRepository) FindCategoryByName(ctx context.Context, name string) (category.Entity, error) {
	dbCateg, err := c.queries.GetCategory(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return category.Entity{}, errorDB.ErrNoRows
		}
		return category.Entity{}, err
	}

	return category.Entity{
		Name:           dbCateg.Name,
		MetadataSchema: dbCateg.MetadataSchema,
	}, nil
}

func (c CategoryPostgresRepository) CreateCategory(ctx context.Context, category category.Entity) (category.Entity, error) {
	//TODO implement me
	panic("implement me")
}
