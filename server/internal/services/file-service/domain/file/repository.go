package file

import (
	"context"
	"time"
)

type TempFile struct {
	Key       string
	FileName  string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTempFileParams struct {
	Key       string
	FileName  string
	ExpiresAt time.Time
}

type ExpiredFilesRow struct {
	Key       string
	ExpiresAt time.Time
}

type Repository interface {
	GetTempFile(ctx context.Context, key string) (TempFile, error)
	CreateTempFile(ctx context.Context, params CreateTempFileParams) error
	GetExpiredFiles(ctx context.Context) ([]ExpiredFilesRow, error)
	DeleteExpiredFilesByKey(ctx context.Context, keys []string) error
}
