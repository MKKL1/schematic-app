package file

import (
	"context"
	"time"
)

type TempFile struct {
	Key         string
	FileName    string
	ContentType string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PermFile struct {
	Hash        string
	FileSize    int32
	ContentType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTempFileParams struct {
	Key         string
	FileName    string
	ContentType string
	ExpiresAt   time.Time
}

type CreateFileParams struct {
	Hash        string
	FileSize    int32
	ContentType string
}

type ExpiredFilesRow struct {
	Key       string
	ExpiresAt time.Time
}

type Repository interface {
	GetTempFile(ctx context.Context, key string) (TempFile, error)
	CreateTempFile(ctx context.Context, params CreateTempFileParams) error
	GetExpiredFiles(ctx context.Context) ([]ExpiredFilesRow, error)
	DeleteTmpFilesByKey(ctx context.Context, keys []string) error
	GetFileByHash(ctx context.Context, hash string) (PermFile, error) //TODO separate
	FileExists(ctx context.Context, hash string) (bool, error)
	CreateFile(ctx context.Context, params CreateFileParams) error
}
