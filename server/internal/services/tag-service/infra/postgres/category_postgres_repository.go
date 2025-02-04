package postgres

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/infra/postgres/db"
)

type CategoryPostgresRepository struct {
	queries *db.Queries
}

func NewCategoryPostgresRepository(queries *db.Queries) *CategoryPostgresRepository {
	return &CategoryPostgresRepository{queries: queries}
}

func (c CategoryPostgresRepository) FindCategoryByID(ctx context.Context, id int64) (category.Entity, error) {
	dbCategory, err := c.queries.GetCategoryByID(ctx, id)
	if err != nil {
		return category.Entity{}, err
	}

	return category.Entity{
		ID:               dbCategory.ID,
		Name:             dbCategory.Name,
		ValueDefinitions: dbCategory.ValueDefinitions,
	}, nil
}

func (c CategoryPostgresRepository) CreateCategory(ctx context.Context, category category.Entity) (category.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (c CategoryPostgresRepository) CreatePostCategory(ctx context.Context, pcv category.PostCategoryValue) (category.PostCategoryValue, error) {
	//TODO implement me
	panic("implement me")
}
