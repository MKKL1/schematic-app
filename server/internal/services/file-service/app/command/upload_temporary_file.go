package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
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
	storageClient file.StorageClient
	repo          file.Repository
	eventBus      *cqrs.EventBus
}

func NewUploadTempFileHandler(storageClient file.StorageClient, repo file.Repository, eventBus *cqrs.EventBus) UploadTempFileHandler {
	return uploadTempFileHandler{storageClient, repo, eventBus}
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileParams) (file.TempFileCreated, error) {
	objectKey := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour)

	info, err := u.storageClient.PutTempObject(ctx, objectKey, cmd.Reader, cmd.ContentType)
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

	err = u.publishFileUploadedEvent(ctx, objectKey, cmd.FileName)
	if err != nil {
		return file.TempFileCreated{}, err
	}

	return file.TempFileCreated{
		Key:        info.Key,
		Expiration: expiresAt,
	}, nil
}

func (u uploadTempFileHandler) publishFileUploadedEvent(ctx context.Context, fileID string, path string) error {
	event := file.TmpFileUploaded{FileID: fileID, Path: path}
	return u.eventBus.Publish(ctx, event)
}
