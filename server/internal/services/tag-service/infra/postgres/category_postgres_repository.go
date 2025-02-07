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

func (c CategoryPostgresRepository) FindCategoryByID(ctx context.Context, name string) (category.Entity, error) {
	dbCategory, err := c.queries.GetCategoryByName(ctx, name)
	if err != nil {
		return category.Entity{}, err
	}

	return category.Entity{
		Name:             dbCategory.Name,
		ValueDefinitions: dbCategory.ValueDefinitions,
	}, nil
}

func (c CategoryPostgresRepository) CreateCategory(ctx context.Context, category category.Entity) (category.Entity, error) {
	//TODO implement me
	panic("implement me")
}

func (c CategoryPostgresRepository) CreatePostCategory(ctx context.Context, pcv category.PostCategoryVars) error {
	err := c.queries.CreatePostCategory(ctx, db.CreatePostCategoryParams{
		PostID:   pcv.PostID,
		Category: pcv.Category,
		Values:   pcv.Values,
	})
	return err
}

func (c CategoryPostgresRepository) FindCategVarsByPostID(ctx context.Context, postID int64) ([]category.PostCategoryVars, error) {
	dbRows, err := c.queries.GetCategVarsForPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	pcvs := make([]category.PostCategoryVars, len(dbRows))
	for i, row := range dbRows {
		pcvs[i] = category.PostCategoryVars{
			PostID:   postID,
			Category: row.Category,
			Values:   row.Values,
		}
	}

	return pcvs, nil
}
