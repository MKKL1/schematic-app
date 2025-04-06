package postgres

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/domain/image"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/postgres/db"
)

type ImagePostgresRepository struct {
	queries *db.Queries
}

func NewImagePostgresRepository(queries *db.Queries) image.Repository { // Return domain interface
	return &ImagePostgresRepository{queries: queries}
}

func (i ImagePostgresRepository) CreateImage(ctx context.Context, params image.CreateParams) error {
	err := i.queries.CreateImage(ctx, db.CreateImageParams{
		FileHash:  params.FileHash,
		ImageType: params.ImageType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (i ImagePostgresRepository) GetImageTypesForHash(ctx context.Context, hash string) ([]string, error) {
	types, err := i.queries.GetImageTypesForHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return types, nil
}
