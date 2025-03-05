package command

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	FileSize    int64
	ContentType string
}

type UploadTempFileHandler decorator.CommandHandler[UploadTempFileParams, *file.TempFileCreated]

type uploadTempFileHandler struct {
	minioClient *minio.Client
	repo        file.Repository
}

func NewUploadTempFileHandler(minioClient *minio.Client, repo file.Repository) UploadTempFileHandler {
	return uploadTempFileHandler{minioClient, repo}
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileParams) (*file.TempFileCreated, error) {
	hasher := sha256.New()
	teeReader := io.TeeReader(cmd.Reader, hasher)
	_, err := u.repo.CreateTempFile(ctx, file.CreateTempFileParams{
		FileHash:    hex.EncodeToString(hasher.Sum(nil)),
		FileName:    cmd.FileName,
		ContentType: cmd.ContentType,
		FileSize:    cmd.FileSize,
		ExpiresAt:   time.Now().Add(time.Hour),
	})
	if err != nil {
		return nil, err
	}

	info, err := u.minioClient.PutObject(ctx, "temp-bucket", uuid.New().String(), teeReader, -1, minio.PutObjectOptions{ContentType: cmd.ContentType})
	if err != nil {
		return nil, err
	}

	urlExpiry := time.Hour
	// Generate the presigned URL.
	presignedUrl, err := u.minioClient.PresignedGetObject(ctx, "temp-bucket", info.Key, urlExpiry, nil)
	if err != nil {
		return nil, err
	}

	return &file.TempFileCreated{
		Key:        info.Key,
		Expiration: urlExpiry,
		Url:        presignedUrl.String(),
	}, nil
}
