package mappers

import (
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	db "github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres/db"
	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgtype"
)

// Maps domain Metadata to DB parameters for saving
func ImageMetadataToCreateParams(metadata file.Metadata) db.CreateImageMetadataParams {
	return db.CreateImageMetadataParams{
		ImageID:     metadata.ID.Int64(),
		StorageKey:  metadata.StorageKey,
		ContentType: metadata.ContentType,
		Width:       metadata.Width,
		Height:      metadata.Height,
		FileSize:    metadata.FileSize,
		CreatedAt:   pgtype.Timestamptz{Time: metadata.CreatedAt, Valid: !metadata.CreatedAt.IsZero()},
	}
}

// Maps DB model to domain Metadata for retrieval
func ImageModelToDomain(model db.ImageMetadatum) file.Metadata { // sqlc might singularize table name
	return file.Metadata{
		ID:          snowflake.ID(model.ImageID),
		StorageKey:  model.StorageKey,
		ContentType: model.ContentType,
		Width:       model.Width,
		Height:      model.Height,
		FileSize:    model.FileSize,
		CreatedAt:   model.CreatedAt.Time,
	}
}
