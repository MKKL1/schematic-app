package postgres

import (
	"context"
	domain "github.com/MKKL1/schematic-app/server/internal/services/tag-service/domain/tag"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/infra/postgres/db"
)

type TagPostgresRepository struct {
	queries *db.Queries
}

func NewTagPostgresRepository(queries *db.Queries) *TagPostgresRepository {
	if queries == nil {
		panic("queries cannot be nil")
	}
	return &TagPostgresRepository{queries}
}

func (r TagPostgresRepository) GetForPost(ctx context.Context, id int64) ([]domain.Tag, error) {
	tags, err := r.queries.GetTagsForPost(ctx, id)
	return domain.ToTag(tags), err
}

func (r TagPostgresRepository) GetCountForTag(ctx context.Context, tag domain.Tag) (int64, error) {
	return r.queries.CountPostsForTag(ctx, string(tag))
}
