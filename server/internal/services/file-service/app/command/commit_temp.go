package command

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/rs/zerolog" // Use zerolog
	"io"
)

type CommitTempParams struct {
	Key      string
	Type     string
	Metadata map[string]string
}

type CommitTempResult struct {
	Hash string
}

type CommitTempHandler decorator.CommandHandler[CommitTempParams, CommitTempResult]

type commitTempHandler struct {
	storageClient file.StorageClient
	repo          file.Repository
	eventBus      *cqrs.EventBus
	logger        zerolog.Logger
}

func NewCommitTempHandler(storageClient file.StorageClient, repo file.Repository, eventBus *cqrs.EventBus, logger zerolog.Logger) CommitTempHandler {
	return commitTempHandler{storageClient, repo, eventBus, logger}
}

func (h commitTempHandler) Handle(ctx context.Context, cmd CommitTempParams) (CommitTempResult, error) {
	log := h.logger.With().Str("handler", "CommitTempHandler").Str("tempKey", cmd.Key).Logger()
	log.Info().Msg("Committing temporary file")

	tempFile, err := h.repo.GetAndMarkTempFileProcessing(ctx, cmd.Key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get and mark temp file processing")
		return CommitTempResult{}, fmt.Errorf("failed to get/mark temp file %s: %w", cmd.Key, err)
	}

	dstObjName, err := h.computeHash(ctx, tempFile.Key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to compute hash")
		_ = h.repo.MarkTempFileFailed(ctx, tempFile.Key, "failed to compute hash") // Mark failed
		return CommitTempResult{}, fmt.Errorf("failed computing hash for %s: %w", tempFile.Key, err)
	}

	exists, err := h.repo.FileExists(ctx, dstObjName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if file exists")
		_ = h.repo.MarkTempFileFailed(ctx, tempFile.Key, "failed to check file existence") // Mark failed
		return CommitTempResult{}, fmt.Errorf("failed checking existence for hash %s: %w", dstObjName, err)
	}

	var finalFileHash string
	if exists {
		finalFileHash = dstObjName
	} else {
		finalFileHash, err = h.copyToPermanentStorage(ctx, tempFile, dstObjName, log) // Pass logger
		if err != nil {
			log.Error().Err(err).Msg("Failed to copy file to permanent storage")
			_ = h.repo.MarkTempFileFailed(ctx, tempFile.Key, "failed to copy to permanent storage") // Mark failed
			return CommitTempResult{}, fmt.Errorf("failed copying %s to %s: %w", tempFile.Key, dstObjName, err)
		}
	}

	err = h.repo.MarkTempFileProcessed(ctx, cmd.Key, finalFileHash)
	if err != nil {
		// Critical error: file is copied but not marked as processed in tmp_file table.
		// This could lead to duplicate processing attempts.
		// For now, logging it. A retry mechanism or dead-letter queue would be better.
		log.Error().Err(err).Msg("CRITICAL: Failed to mark temp file as processed after copy/existence check")
		// Still proceed to publish event as the file *is* available.
	}

	err = h.publishCreatedFileEvent(ctx, file.FileCreated{
		TempID:   tempFile.Key,
		PermID:   finalFileHash,
		Existed:  exists,
		Type:     cmd.Type,
		Metadata: cmd.Metadata,
	})
	if err != nil {
		// Log non-critical error: Event publishing failed, main operation succeeded.
		log.Error().Err(err).Msg("Failed to publish FileCreated event")
	}

	return CommitTempResult{Hash: finalFileHash}, nil
}

// Add logger parameter
func (h commitTempHandler) copyToPermanentStorage(
	ctx context.Context,
	tempFile file.TempFile,
	dstObjName string,
	log zerolog.Logger, // Use logger
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
		log.Error().Err(err).Str("hash", info.Key).Msg("Failed to create file record in database")
		// If DB insert fails, the file is copied but not recorded. Attempt cleanup.
		cleanupErr := h.storageClient.RemovePermObject(ctx, dstObjName)
		if cleanupErr != nil {
			log.Error().Err(cleanupErr).Str("dstObjName", dstObjName).Msg("Failed to clean up orphaned permanent file after DB error")
			// This is a bad state - manual intervention might be needed.
		} else {
			log.Warn().Str("dstObjName", dstObjName).Msg("Cleaned up orphaned permanent file after DB error")
		}
		return "", fmt.Errorf("db create failed after copy: %w", err)
	}

	log.Info().Str("hash", info.Key).Int32("size", int32(info.Size)).Msg("File record created in database")
	return info.Key, nil
}

func (h commitTempHandler) computeHash(ctx context.Context, object string) (string, error) {
	obj, err := h.storageClient.GetTempObject(ctx, object)
	if err != nil {
		return "", fmt.Errorf("failed to get temp object %s for hashing: %w", object, err)
	}
	defer obj.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, obj); err != nil {
		return "", fmt.Errorf("failed to read temp object %s for hashing: %w", object, err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (h commitTempHandler) publishCreatedFileEvent(ctx context.Context, event file.FileCreated) error {
	return h.eventBus.Publish(ctx, event) // Publish event object directly
}
