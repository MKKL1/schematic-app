package file

import (
	"context"
	"time"
)

type TempFile struct {
	FileHash    string
	Key         string
	FileName    string
	ContentType string
	FileSize    int64
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTempFileParams struct {
	FileHash    string
	Key         string
	FileName    string
	ContentType string
	FileSize    int64
	ExpiresAt   time.Time
}

type ExpiredFilesRow struct {
	FileHash  string
	Key       string
	ExpiresAt time.Time
}

type Repository interface {
	GetTempFile(ctx context.Context, fileHash string) (TempFile, error)
	CreateTempFile(ctx context.Context, params CreateTempFileParams) (string, error)
	GetExpiredFiles(ctx context.Context) ([]ExpiredFilesRow, error)
}
