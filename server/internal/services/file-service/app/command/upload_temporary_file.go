package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"io"
	"time"
)

type UploadTempFileParams struct {
	Reader      io.Reader
	FileName    string
	ContentType string
}

type UploadTempFileHandler decorator.CommandHandler[UploadTempFileParams, file.TempFileCreated]

type uploadTempFileHandler struct {
	minioClient *minio.Client
	repo        file.Repository
}

func NewUploadTempFileHandler(minioClient *minio.Client, repo file.Repository) UploadTempFileHandler {
	return uploadTempFileHandler{minioClient, repo}
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileParams) (file.TempFileCreated, error) {
	objectKey := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour)

	info, err := u.minioClient.PutObject(ctx, "temp-bucket", objectKey, cmd.Reader, -1, minio.PutObjectOptions{ContentType: cmd.ContentType})
	if err != nil {
		return file.TempFileCreated{}, err
	}

	err = u.repo.CreateTempFile(ctx, file.CreateTempFileParams{
		Key:       info.Key,
		FileName:  cmd.FileName,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return file.TempFileCreated{}, err
	}

	return file.TempFileCreated{
		Key:        info.Key,
		Expiration: expiresAt,
	}, nil
}
