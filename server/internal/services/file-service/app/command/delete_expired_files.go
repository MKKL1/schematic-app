package command

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/rs/zerolog"
	"strings"
)

type DeleteExpiredFilesCmd struct {
}

type DeleteExpiredFilesHandler decorator.CommandHandler[DeleteExpiredFilesCmd, any]

type deleteExpiredFilesHandler struct {
	storageClient file.StorageClient
	repo          file.Repository
	logger        zerolog.Logger
}

func NewDeleteExpiredFilesHandler(storageClient file.StorageClient, repo file.Repository, logger zerolog.Logger, metrics metrics.Client) DeleteExpiredFilesHandler {
	return decorator.ApplyCommandDecorators[DeleteExpiredFilesCmd, any](
		deleteExpiredFilesHandler{storageClient, repo, logger},
		logger,
		metrics,
	)
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NoSuchKey")
}

func (d deleteExpiredFilesHandler) Handle(ctx context.Context, cmd DeleteExpiredFilesCmd) (any, error) {
	logger := decorator.AddCmdInfo(cmd, d.logger)
	files, err := d.repo.GetExpiredFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting expired files: %w", err)
	}

	if len(files) == 0 {
		return nil, nil
	}

	// Create a map to track removal status.
	// We'll assume that all files are successfully removed unless we see an error.
	successfulRemovals := make(map[string]bool, len(files))
	for _, f := range files {
		successfulRemovals[f.Key] = true
	}

	objectsCh := make(chan string)
	go func() {
		defer close(objectsCh)
		for _, f := range files {
			objectsCh <- f.Key
		}
	}()

	errCh := d.storageClient.RemoveTempObjects(ctx, objectsCh)
	for removeErr := range errCh {
		if removeErr.Err != nil {
			//TODO abstract

			// Check if the error indicates the object is already removed.
			if isNotFoundError(removeErr.Err) {
				logger.Warn().Str("object_name", removeErr.ObjectName).Msg("Object not found, treating as already removed")
			} else {
				// Otherwise, mark the file as not removed.
				delete(successfulRemovals, removeErr.ObjectName)
				logger.Error().Err(removeErr.Err).Str("object_name", removeErr.ObjectName).Msg("Failed to remove object, continuing anyway")
			}
		}
	}

	// Collect file keys that were successfully removed.
	var keysToDelete []string
	for key := range successfulRemovals {
		keysToDelete = append(keysToDelete, key)
	}

	// Remove only those records from the database whose files were successfully removed.
	if len(keysToDelete) > 0 {
		if err := d.repo.DeleteTmpFilesByKey(ctx, keysToDelete); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
