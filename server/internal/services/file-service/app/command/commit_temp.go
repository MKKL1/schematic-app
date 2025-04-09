package command

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/rs/zerolog"
	"io"
)

type CommitTempCmd struct {
	Key      string
	Type     string
	Metadata map[string]string
}

type CommitTempResult struct {
	Hash string
}

type CommitTempHandler decorator.CommandHandler[CommitTempCmd, CommitTempResult]

type commitTempHandler struct {
	storageClient file.StorageClient
	repo          file.Repository
	eventBus      *cqrs.EventBus
	logger        zerolog.Logger
}

func NewCommitTempHandler(storageClient file.StorageClient, repo file.Repository, eventBus *cqrs.EventBus, logger zerolog.Logger, metrics metrics.Client) CommitTempHandler {
	return decorator.ApplyCommandDecorators[CommitTempCmd, CommitTempResult](
		commitTempHandler{storageClient, repo, eventBus, logger},
		logger,
		metrics,
	)
}

func (h commitTempHandler) Handle(ctx context.Context, cmd CommitTempCmd) (CommitTempResult, error) {
	log := decorator.AddCmdInfo(cmd, h.logger)

	//TODO may need to split into individual commands

	//Ensure that file won't be deleted while it's being processed
	tempFile, err := h.repo.GetAndMarkTempFileProcessing(ctx, cmd.Key)
	if err != nil {
		return CommitTempResult{}, fmt.Errorf("mark file processing %s: %w", cmd.Key, err)
	}

	dstObjName, err := h.computeHash(ctx, tempFile.Key)
	if err != nil {
		return CommitTempResult{}, fmt.Errorf("failed computing hash for %s: %w", tempFile.Key, err)
	}

	//Check file type (image/schematic/unknown/others in future...)
	//Depending on type choose service to perform pre checks for file
	//Image service: verifies that image doesn't contain any inappropriate content
	//Schematic service: verifies that schematic has proper structure
	//Unknown: no operation
	//All of those services take in hash of file and check if they have already checked those files

	//Check if file already exists in bucket, before saving it
	exists, err := h.repo.FileExists(ctx, dstObjName)
	if err != nil {
		return CommitTempResult{}, fmt.Errorf("failed checking existence for hash %s: %w", dstObjName, err)
	}

	var finalFileHash string
	if exists {
		finalFileHash = dstObjName
		log.Trace().Str("hash", cmd.Key).Msg("file already existed")
	} else {
		finalFileHash, err = h.copyToPermanentStorage(ctx, tempFile, dstObjName, log)
		if err != nil {
			return CommitTempResult{}, fmt.Errorf("failed copying %s to %s: %w", tempFile.Key, dstObjName, err)
		}
		log.Trace().Str("hash", cmd.Key).Msg("saved new file to permanent storage")
	}

	err = h.repo.MarkTempFileProcessed(ctx, cmd.Key, finalFileHash)
	if err != nil {
		// Critical error: file is copied but not marked as processed in tmp_file table.
		// This could lead to duplicate processing attempts.
		// For now, logging it. A retry mechanism or dead-letter queue would be better.
		log.Error().Err(err).Msg("CRITICAL: Failed to mark temp file as processed after copy/existence check")
		// Still proceed to publish event as the file *is* available.
	}

	err = h.publishCreatedFileEvent(ctx, file.FileUploaded{
		TempID:   tempFile.Key,
		PermID:   finalFileHash,
		Existed:  exists,
		Type:     cmd.Type,
		Metadata: cmd.Metadata,
	})
	if err != nil {
		// Critical
		log.Error().Err(err).Msg("Failed to publish FileUploaded event")
	}

	//After file is uploaded to bucket, it has to be processed by services again
	//Image service: generate thumbnails
	//Schematic service: check content, generate 3d model, render model, pass renders to image service

	return CommitTempResult{Hash: finalFileHash}, nil
}

// Add logger parameter
func (h commitTempHandler) copyToPermanentStorage(
	ctx context.Context,
	tempFile file.TempFile,
	dstObjName string,
	log zerolog.Logger,
) (string, error) {
	info, err := h.storageClient.CopyTempToPermanent(ctx, tempFile.Key, dstObjName)
	if err != nil {
		log.Error().Err(err).Str("dstObjName", dstObjName).Msg("Minio copy failed")
		return "", fmt.Errorf("storage copy failed: %w", err)
	}

	err = h.repo.CreateFile(ctx, file.CreateFileParams{
		Hash:        info.Key, // Use the key from UploadInfo (should match dstObjName)
		FileSize:    int32(info.Size),
		ContentType: tempFile.ContentType,
	})
	if err != nil {
		log.Warn().Err(err).Str("hash", info.Key).Msg("Failed to create file record in database, attempting to remove orphaned permanent file")
		cleanupErr := h.storageClient.RemovePermObject(ctx, dstObjName)
		if cleanupErr != nil {
			log.Error().Err(cleanupErr).Str("dstObjName", dstObjName).Msg("Failed to clean up orphaned permanent file after DB error")
			// This is a bad state - manual intervention might be needed.
		} else {
			log.Warn().Str("dstObjName", dstObjName).Msg("Cleaned up orphaned permanent file after DB error")
		}
		return "", fmt.Errorf("create file in repo: %w", err)
	}

	log.Trace().Str("hash", info.Key).Int32("size", int32(info.Size)).Msg("File record created in database")
	return info.Key, nil
}

func (h commitTempHandler) computeHash(ctx context.Context, object string) (string, error) {
	obj, err := h.storageClient.GetTempObject(ctx, object)
	if err != nil {
		return "", fmt.Errorf("get temp object %s for hashing: %w", object, err)
	}
	defer obj.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, obj); err != nil {
		return "", fmt.Errorf("read temp object %s for hashing: %w", object, err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (h commitTempHandler) publishCreatedFileEvent(ctx context.Context, event file.FileUploaded) error {
	return h.eventBus.Publish(ctx, event)
}
