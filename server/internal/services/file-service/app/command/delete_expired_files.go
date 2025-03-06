package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"strings"
)

type DeleteExpiredFilesParams struct {
}

type DeleteExpiredFilesHandler decorator.CommandHandler[DeleteExpiredFilesParams, any]

type deleteExpiredFilesHandler struct {
	minioClient *minio.Client
	repo        file.Repository
}

func NewDeleteExpiredFilesHandler(minioClient *minio.Client, repo file.Repository) DeleteExpiredFilesHandler {
	return deleteExpiredFilesHandler{minioClient, repo}
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

	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, f := range files {
			objectsCh <- minio.ObjectInfo{Key: f.Key}
		}
	}()

	errCh := d.minioClient.RemoveObjects(ctx, "temp-bucket", objectsCh, minio.RemoveObjectsOptions{})
	for removeErr := range errCh {
		if removeErr.Err != nil {
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
		if err := d.repo.DeleteExpiredFilesByKey(ctx, keysToDelete); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
