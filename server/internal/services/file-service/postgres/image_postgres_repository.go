package postgres

import (
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres/db"
)

type ImagePostgresRepository struct {
	queries *db.Queries
}

func NewImagePostgresRepository(queries *db.Queries) file.ImageRepository { // Return domain interface
	return &ImagePostgresRepository{queries: queries}
}
