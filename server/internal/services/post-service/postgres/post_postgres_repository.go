package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/mappers"
	db2 "github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres/db"
	"github.com/jackc/pgx/v5"
)

type PostPostgresRepository struct {
	queries *db2.Queries
}

func NewPostPostgresRepository(queries *db2.Queries) *PostPostgresRepository {
	return &PostPostgresRepository{queries}
}

func (p PostPostgresRepository) FindById(ctx context.Context, id int64) (post.Post, error) {
	row, err := p.queries.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return post.Post{}, errorDB.ErrNoRows
		}
		return post.Post{}, err
	}

	return mappers.GetPostRowToApp(row)
}

func (p PostPostgresRepository) Create(ctx context.Context, params post.CreatePostParams) error {
	dbParams, err := mappers.CreatePostParamsToQuery(params)
	if err != nil {
		return err
	}

	err = p.queries.CreatePost(ctx, dbParams)
	if err != nil {
		return err
	}
	return nil
}

func (p PostPostgresRepository) GetCountForTag(ctx context.Context, tag string) (int64, error) {
	//TODO implement me
	panic("implement me")
}
