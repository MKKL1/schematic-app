package command

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

type UploadTempFileHandler decorator.CommandHandler[UploadTempFileParams, *file.TempFileCreated]

type uploadTempFileHandler struct {
	minioClient *minio.Client
	repo        file.Repository
}

func NewUploadTempFileHandler(minioClient *minio.Client, repo file.Repository) UploadTempFileHandler {
	return uploadTempFileHandler{minioClient, repo}
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileParams) (*file.TempFileCreated, error) {
	objectKey := uuid.New().String()

	// Create an in-memory buffer to hold the file data.
	var buf bytes.Buffer

	// Create a hasher and a MultiWriter that writes to both the buffer and the hash.
	hasher := md5.New()
	mw := io.MultiWriter(&buf, hasher)

	// Read all file data from cmd.Reader into the buffer while computing the hash.
	if _, err := io.Copy(mw, cmd.Reader); err != nil {
		return nil, fmt.Errorf("failed to read input data: %w", err)
	}

	expiresAt := time.Now().Add(time.Hour)

	// Compute the file hash.
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	verifiedObjectKey, err := u.repo.CreateTempFile(ctx, file.CreateTempFileParams{
		FileHash:    fileHash,
		Key:         objectKey,
		FileName:    cmd.FileName,
		ContentType: cmd.ContentType,
		FileSize:    int64(buf.Len()),
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		return nil, err
	}

	//If object key from database and new key are the same, it means that record is unique and file has to be saved
	if verifiedObjectKey == objectKey {
		reader := bytes.NewReader(buf.Bytes())
		info, err := u.minioClient.PutObject(ctx, "temp-bucket", objectKey, reader, -1, minio.PutObjectOptions{ContentType: cmd.ContentType})
		if err != nil {
			return nil, err
		}
		//This approach requires minio to not change the key
		if info.Key != objectKey {
			return nil, fmt.Errorf("key does not match")
		}
	}

	urlExpiry := time.Hour
	// Generate the presigned URL.
	presignedUrl, err := u.minioClient.PresignedGetObject(ctx, "temp-bucket", verifiedObjectKey, urlExpiry, nil)
	if err != nil {
		return nil, err
	}

	return &file.TempFileCreated{
		Key:        verifiedObjectKey,
		Expiration: expiresAt,
		Url:        presignedUrl.String(),
	}, nil
}
