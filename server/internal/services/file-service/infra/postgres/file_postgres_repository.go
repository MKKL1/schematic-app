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

func (f FilePostgresRepository) GetTempFile(ctx context.Context, fileHash string) (file.TempFile, error) {
	retrievedFile, err := f.queries.GetFileByHash(ctx, fileHash)
	if err != nil {
		return file.TempFile{}, err
	}

	return toDto(retrievedFile), nil
}

func (f FilePostgresRepository) CreateTempFile(ctx context.Context, params file.CreateTempFileParams) (string, error) {
	key, err := f.queries.CreateFile(ctx, db.CreateFileParams{
		FileHash:    params.FileHash,
		StoreKey:    params.Key,
		FileName:    params.FileName,
		ContentType: params.ContentType,
		FileSize:    params.FileSize,
		ExpiresAt: pgtype.Timestamptz{
			Time:  params.ExpiresAt,
			Valid: true,
		},
	})

	if err != nil {
		return key, err
	}

	return key, nil
}

func (f FilePostgresRepository) GetExpiredFiles(ctx context.Context) ([]file.ExpiredFilesRow, error) {
	files, err := f.queries.ListExpiredFiles(ctx)
	if err != nil {
		return nil, err
	}

	var expiredFiles []file.ExpiredFilesRow
	for i, k := range files {
		expiredFiles[i] = file.ExpiredFilesRow{
			FileHash:  k.FileHash,
			ExpiresAt: k.ExpiresAt.Time,
		}
	}

	return expiredFiles, nil
}

func toDto(model db.TmpFile) file.TempFile {
	return file.TempFile{
		FileHash:    model.FileHash,
		FileName:    model.FileName,
		ContentType: model.ContentType,
		FileSize:    model.FileSize,
		ExpiresAt:   model.ExpiresAt.Time,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
	}
}
