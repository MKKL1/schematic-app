package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/rs/zerolog/log"
	"strings"
)

type DeleteExpiredFilesParams struct {
}

type DeleteExpiredFilesHandler decorator.CommandHandler[DeleteExpiredFilesParams, any]

type deleteExpiredFilesHandler struct {
	storageClient file.StorageClient
	repo          file.Repository
}

func NewDeleteExpiredFilesHandler(storageClient file.StorageClient, repo file.Repository) DeleteExpiredFilesHandler {
	return deleteExpiredFilesHandler{storageClient, repo}
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NoSuchKey")
}

func (d deleteExpiredFilesHandler) Handle(ctx context.Context, cmd DeleteExpiredFilesParams) (any, error) {
	files, err := d.repo.GetExpiredFiles(ctx)
	if err != nil {
		return nil, err
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
				log.Printf("Object %s not found, treating as already removed", removeErr.ObjectName)
			} else {
				// Otherwise, mark the file as not removed.
				delete(successfulRemovals, removeErr.ObjectName)
				log.Printf("Failed to remove object '%s': %v", removeErr.ObjectName, removeErr.Err)
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
