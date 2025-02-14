package postgres

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
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
		return category.Entity{}, err
	}

	return category.Entity{
		Name:             dbCateg.Name,
		ValueDefinitions: dbCateg.MetadataSchema,
	}, nil
}

func (c CategoryPostgresRepository) CreateCategory(ctx context.Context, category category.Entity) (category.Entity, error) {
	//TODO implement me
	panic("implement me")
}
