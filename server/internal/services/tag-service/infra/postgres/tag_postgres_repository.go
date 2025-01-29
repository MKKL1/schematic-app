package postgres

import (
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
