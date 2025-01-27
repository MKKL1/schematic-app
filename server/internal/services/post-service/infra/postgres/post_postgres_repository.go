package postgres

import (
	"context"
	"errors"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	db2 "github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostPostgresRepository struct {
	queries *db2.Queries
}

func NewPostPostgresRepository(queries *db2.Queries) *PostPostgresRepository {
	if queries == nil {
		panic("queries cannot be nil")
	}
	return &PostPostgresRepository{queries}
}

func (p *PostPostgresRepository) FindById(ctx context.Context, id int64) (post.Entity, error) {
	out, err := p.queries.GetPostById(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return post.Entity{}, errorDB.ErrNoRows
	} else if err != nil {
		return post.Entity{}, err
	}
	return toModel(out)
}

func (p *PostPostgresRepository) Create(ctx context.Context, model post.Entity) error {
	params := []db2.CreatePostParams{
		{
			ID:            model.ID,
			Name:          model.Name,
			Desc:          toText(model.Description),
			Owner:         model.Owner,
			AuthorKnown:   toInt(model.AuthorID),
			AuthorUnknown: toText(model.AuthorName),
		},
	}

	_, err := p.queries.CreatePost(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func toText(str *string) pgtype.Text {
	if str == nil {
		return pgtype.Text{
			Valid: false,
		}
	}
	return pgtype.Text{
		String: *str,
		Valid:  true,
	}
}

func toInt(val *int64) pgtype.Int8 {
	if val == nil {
		return pgtype.Int8{
			Valid: false,
		}
	}
	return pgtype.Int8{
		Int64: *val,
		Valid: true,
	}
}

func toModel(dbPost db2.Post) (post.Entity, error) {
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

	return post.Entity{
		ID:          dbPost.ID,
		Name:        dbPost.Name,
		Description: desc,
		Owner:       dbPost.Owner,
		AuthorName:  authorName,
		AuthorID:    authorId,
	}, nil
}
