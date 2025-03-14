package command

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/infra/kafka"
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
	pub         *kafka.KafkaPublisher
}

func NewUploadTempFileHandler(minioClient *minio.Client, repo file.Repository, pub *kafka.KafkaPublisher) UploadTempFileHandler {
	return uploadTempFileHandler{minioClient, repo, pub}
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

	err = publishFileUploadedEvent(u.pub, objectKey, cmd.FileName)
	if err != nil {
		return file.TempFileCreated{}, err
	}

	return file.TempFileCreated{
		Key:        info.Key,
		Expiration: expiresAt,
	}, nil
}

type FileUploadedEvent struct {
	FileID string `json:"file_id"`
	Path   string `json:"path"`
}

func publishFileUploadedEvent(publisher *kafka.KafkaPublisher, fileID, path string) error {
	event := FileUploadedEvent{FileID: fileID, Path: path}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return publisher.Publish("file.temp.created", payload)
}
