package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres/db"
	"github.com/jackc/pgx/v5"
)

type PostPostgresRepository struct {
	queries *db.Queries
}

func NewPostPostgresRepository(queries *db.Queries) *PostPostgresRepository {
	if queries == nil {
		panic("queries cannot be nil")
	}
	return &PostPostgresRepository{queries}
}

func (p *PostPostgresRepository) FindById(ctx context.Context, id int64) (post.Model, error) {
	out, err := p.queries.GetPostById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return post.Model{}, errorDB.ErrNoRows
	} else if err != nil {
		return post.Model{}, err
	}
	return toModel(out)
}

func toModel(dbPost db.Post) (post.Model, error) {
	var desc *string = nil
	if dbPost.Desc.Valid {
		desc = &dbPost.Desc.String
	}

	var authorName *string = nil
	if dbPost.AuthorUnknown.Valid {
		authorName = &dbPost.AuthorUnknown.String
	}
	var authorId *int64 = nil
	if dbPost.AuthorKnown.Valid {
		authorId = &dbPost.AuthorKnown.Int64
	}

	return post.Model{
		ID:          dbPost.ID,
		Description: desc,
		Owner:       dbPost.Owner,
		AuthorName:  authorName,
		AuthorID:    authorId,
	}, nil
}
