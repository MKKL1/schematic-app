package file

import (
	"context"
	"time"
)

type TempFile struct {
	FileHash    string
	FileName    string
	ContentType string
	FileSize    int64
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTempFileParams struct {
	FileHash    string
	FileName    string
	ContentType string
	FileSize    int64
	ExpiresAt   time.Time
}

type ExpiredFilesRow struct {
	FileHash  string
	ExpiresAt time.Time
}

type Repository interface {
	GetTempFile(ctx context.Context, fileHash string) (TempFile, error)
	CreateTempFile(ctx context.Context, params CreateTempFileParams) (TempFile, error)
	GetExpiredFiles(ctx context.Context) ([]ExpiredFilesRow, error)
}
