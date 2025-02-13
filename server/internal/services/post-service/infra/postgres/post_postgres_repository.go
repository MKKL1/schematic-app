package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
)

type PostPostgresRepository struct {
	queries *db.Queries
}

func NewPostPostgresRepository(queries *db.Queries) *PostPostgresRepository {
	return &PostPostgresRepository{queries}
}

func (p PostPostgresRepository) FindById(ctx context.Context, id int64) (post.Post, error) {
	row, err := p.queries.GetPost(ctx, id)
	if err != nil {
		return post.Post{}, err
	}

	var categoryVars post.CategoryVars
	if row.CategoryVars != nil {
		s, ok := row.CategoryVars.(string)
		if !ok {
			return post.Post{}, errors.New("invalid type for category var")
		}
		categoryVars = []byte(s)
	}

	var tags []string
	if row.Tags != nil {
		tagArray, ok := row.Tags.([]interface{})
		if !ok {
			return post.Post{}, fmt.Errorf("invalid type for Tags")
		}
		for _, tag := range tagArray {
			if str, ok := tag.(string); ok {
				tags = append(tags, str)
			} else {
				return post.Post{}, fmt.Errorf("invalid tag type: expected string, got %T", tag)
			}
		}
	}

	postEntity := post.Post{
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

func (p PostPostgresRepository) Create(ctx context.Context, model post.Post) error {
	//TODO implement me
	panic("implement me")
}

func (p PostPostgresRepository) GetCountForTag(ctx context.Context, tag string) (int64, error) {
	//TODO implement me
	panic("implement me")
}
