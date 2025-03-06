package postgres

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/infra/postgres/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type FilePostgresRepository struct {
	queries *db.Queries
}

func NewFilePostgresRepository(queries *db.Queries) *FilePostgresRepository {
	return &FilePostgresRepository{queries}
}

func (f FilePostgresRepository) GetTempFile(ctx context.Context, key string) (file.TempFile, error) {
	retrievedFile, err := f.queries.GetFile(ctx, key)
	if err != nil {
		return file.TempFile{}, err
	}

	return toDto(retrievedFile), nil
}

func (f FilePostgresRepository) CreateTempFile(ctx context.Context, params file.CreateTempFileParams) error {
	err := f.queries.CreateFile(ctx, db.CreateFileParams{
		StoreKey: params.Key,
		FileName: params.FileName,
		ExpiresAt: pgtype.Timestamptz{
			Time:  params.ExpiresAt,
			Valid: true,
		},
	})

	return err
}

func (f FilePostgresRepository) GetExpiredFiles(ctx context.Context) ([]file.ExpiredFilesRow, error) {
	files, err := f.queries.ListExpiredFiles(ctx)
	if err != nil {
		return nil, err
	}

	expiredFiles := make([]file.ExpiredFilesRow, len(files))
	for i, k := range files {
		expiredFiles[i] = file.ExpiredFilesRow{
			Key:       k.StoreKey,
			ExpiresAt: k.ExpiresAt.Time,
		}
	}

	return expiredFiles, nil
}

func (f FilePostgresRepository) DeleteExpiredFilesByKey(ctx context.Context, keys []string) error {
	err := f.queries.DeleteExpiredFiles(ctx, keys)
	if err != nil {
		return err
	}

	return nil
}

func toDto(model db.TmpFile) file.TempFile {
	return file.TempFile{
		FileName:  model.FileName,
		ExpiresAt: model.ExpiresAt.Time,
		CreatedAt: model.CreatedAt.Time,
		UpdatedAt: model.UpdatedAt.Time,
	}
}
