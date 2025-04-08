package command

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type UploadTempFileCmd struct {
	Reader      io.Reader
	FileName    string
	ContentType string
}

type UploadTempFileHandler decorator.CommandHandler[UploadTempFileCmd, file.TempFileCreated]

type uploadTempFileHandler struct {
	storageClient  file.StorageClient
	repo           file.Repository
	eventBus       *cqrs.EventBus
	logger         zerolog.Logger
	expireDuration time.Duration
}

func NewUploadTempFileHandler(storageClient file.StorageClient,
	repo file.Repository,
	eventBus *cqrs.EventBus,
	logger zerolog.Logger,
	metrics metrics.Client,
	expireDuration time.Duration) UploadTempFileHandler {
	return decorator.ApplyCommandDecorators[UploadTempFileCmd, file.TempFileCreated](
		uploadTempFileHandler{storageClient, repo, eventBus, logger, expireDuration},
		logger,
		metrics,
	)
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileCmd) (file.TempFileCreated, error) {
	objectKey := uuid.New().String()
	expiresAt := time.Now().Add(u.expireDuration)

	//TODO transaction

	info, err := u.storageClient.PutTempObject(ctx, objectKey, cmd.Reader, cmd.ContentType)
	if err != nil {
		return file.TempFileCreated{}, fmt.Errorf("store temp file: %w", err)
	}

	err = u.repo.CreateTempFile(ctx, file.CreateTempFileParams{
		Key:       info.Key,
		FileName:  cmd.FileName,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return file.TempFileCreated{}, fmt.Errorf("creating temp file in repo: %w", err)
	}

	err = u.publishFileUploadedEvent(ctx, objectKey, cmd.FileName)
	if err != nil {
		return file.TempFileCreated{}, fmt.Errorf("publishing file uploaded event: %w", err)
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
