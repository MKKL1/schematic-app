package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
)

type PostPostgresRepository struct {
	queries db.Queries
}

func NewPostPostgresRepository(queries db.Queries) *PostPostgresRepository {
	return &PostPostgresRepository{queries}
}

func (p PostPostgresRepository) FindById(ctx context.Context, id int64) (post.Entity, error) {
	row, err := p.queries.GetPost(ctx, id)
	if err != nil {
		return post.Entity{}, err
	}

	var categoryVars []post.CategoryVarsEntity
	if row.CategoryVars != nil {
		data, ok := row.CategoryVars.([]byte)
		if !ok {
			return post.Entity{}, fmt.Errorf("invalid type for CategoryVars")
		}
		if err := json.Unmarshal(data, &categoryVars); err != nil {
			return post.Entity{}, fmt.Errorf("failed to unmarshal CategoryVars: %w", err)
		}
	}

	var tags []string
	if row.Tags != nil {
		tagArray, ok := row.Tags.([]interface{})
		if !ok {
			return post.Entity{}, fmt.Errorf("invalid type for Tags")
		}
		for _, tag := range tagArray {
			if str, ok := tag.(string); ok {
				tags = append(tags, str)
			} else {
				return post.Entity{}, fmt.Errorf("invalid tag type: expected string, got %T", tag)
			}
		}
	}

	postEntity := post.Entity{
		ID:           row.ID,
		Name:         row.Name,
		Description:  row.Description,
		Owner:        row.Owner,
		AuthorID:     row.AuthorID,
		CategoryVars: categoryVars,
		Tags:         tags,
	}

	return postEntity, nil
}

func (p PostPostgresRepository) Create(ctx context.Context, model post.Entity) error {
	//TODO implement me
	panic("implement me")
}

func (p PostPostgresRepository) GetCountForTag(ctx context.Context, tag string) (int64, error) {
	//TODO implement me
	panic("implement me")
}
