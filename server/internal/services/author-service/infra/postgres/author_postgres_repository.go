package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/error/db"
	"github.com/MKKL1/schematic-app/server/internal/services/author-service/domain/author"
	"github.com/MKKL1/schematic-app/server/internal/services/author-service/infra/postgres/db"
)

type AuthorPostgresRepository struct {
	queries *db.Queries
}

func NewAuthorPostgresRepository(queries *db.Queries) *AuthorPostgresRepository {
	if queries == nil {
		panic("queries cannot be nil")
	}
	return &AuthorPostgresRepository{queries}
}

func (a AuthorPostgresRepository) FindByID(ctx context.Context, id int64) (author.Entity, error) {
	dbAuthor, err := a.queries.GetAuthorByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return author.Entity{}, errorDB.ErrNoRows
		}
		return author.Entity{}, fmt.Errorf("failed to get author by ID: %w", err)
	}

	return convertDBAuthor(dbAuthor)
}

func (a AuthorPostgresRepository) FindByName(ctx context.Context, name string) (author.Entity, error) {
	dbAuthor, err := a.queries.GetAuthorByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return author.Entity{}, errorDB.ErrNoRows
		}
		return author.Entity{}, fmt.Errorf("failed to get author by name: %w", err)
	}

	return convertDBAuthor(dbAuthor)
}

func (a AuthorPostgresRepository) Create(ctx context.Context, authorEntity author.Entity) error {
	metadataJSON, err := json.Marshal(authorEntity.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	params := []db.CreateAuthorParams{
		{
			Name:     authorEntity.Name,
			UserID:   authorEntity.UserID,
			Metadata: metadataJSON,
		},
	}

	_, err = a.queries.CreateAuthor(ctx, params)
	if err != nil {
		//TODO handle conflict
		return fmt.Errorf("failed to create author: %w", err)
	}

	return nil
}

func convertDBAuthor(dbAuthor db.Author) (author.Entity, error) {
	var metadata map[string]string
	if len(dbAuthor.Metadata) > 0 {
		if err := json.Unmarshal(dbAuthor.Metadata, &metadata); err != nil {
			return author.Entity{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return author.Entity{
		ID:       dbAuthor.ID,
		Name:     dbAuthor.Name,
		UserID:   dbAuthor.UserID,
		Metadata: metadata,
	}, nil
}
