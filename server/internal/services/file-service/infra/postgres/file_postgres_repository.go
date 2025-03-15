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
	retrievedFile, err := f.queries.GetTempFile(ctx, key)
	if err != nil {
		return file.TempFile{}, err
	}

	return toDto(retrievedFile), nil
}

func (f FilePostgresRepository) CreateTempFile(ctx context.Context, params file.CreateTempFileParams) error {
	err := f.queries.CreateTempFile(ctx, db.CreateTempFileParams{
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

func (f FilePostgresRepository) DeleteTmpFilesByKey(ctx context.Context, keys []string) error {
	err := f.queries.DeleteTmpFiles(ctx, keys)
	if err != nil {
		return err
	}

	return nil
}

func (f FilePostgresRepository) GetFileByHash(ctx context.Context, hash string) (file.PermFile, error) {
	//TODO implement me
	panic("implement me")
}

func (f FilePostgresRepository) FileExists(ctx context.Context, hash string) (bool, error) {
	exists, err := f.queries.FileExistsByHash(ctx, hash)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (f FilePostgresRepository) CreateFile(ctx context.Context, params file.CreateFileParams) error {
	err := f.queries.CreateFile(ctx, db.CreateFileParams{
		Hash:        params.Hash,
		FileSize:    params.FileSize,
		ContentType: params.ContentType,
	})
	if err != nil {
		return err
	}
	return nil
}

func (f FilePostgresRepository) GetAndMarkTempFileProcessing(ctx context.Context, key string) (file.TempFile, error) {
	tmpFile, err := f.queries.GetAndMarkTempFileProcessing(ctx, key)
	if err != nil {
		return file.TempFile{}, err
	}

	return toDto(tmpFile), nil
}

func (f FilePostgresRepository) MarkTempFileFailed(ctx context.Context, key string, reason string) error {
	err := f.queries.MarkTempFileFailed(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (f FilePostgresRepository) MarkTempFileProcessed(ctx context.Context, key string, finalHash string) error {
	err := f.queries.MarkTempFileProcessed(ctx, db.MarkTempFileProcessedParams{
		StoreKey:  key,
		FinalHash: &finalHash,
	})
	if err != nil {
		return err
	}
	return nil
}

func toDto(model db.TmpFile) file.TempFile {
	return file.TempFile{
		Key:         model.StoreKey,
		FileName:    model.FileName,
		ContentType: model.ContentType,
		Status:      model.Status,
		ErrorReason: model.ErrorReason,
		FinalHash:   model.FinalHash,
		ExpiresAt:   model.ExpiresAt.Time,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
	}
}
