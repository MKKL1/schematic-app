package command

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"io"
)

type CommitTempParams struct {
	Key string
}

type CommitTempResult struct {
	Hash string
}

type CommitTempHandler decorator.CommandHandler[CommitTempParams, CommitTempResult]

type commitTempHandler struct {
	minioClient *minio.Client
	repo        file.Repository
	eventBus    *cqrs.EventBus
}

func NewCommitTempHandler(minioClient *minio.Client, repo file.Repository, eventBus *cqrs.EventBus) CommitTempHandler {
	return commitTempHandler{minioClient, repo, eventBus}
}

func (m commitTempHandler) Handle(ctx context.Context, cmd CommitTempParams) (CommitTempResult, error) {
	tempFile, err := m.repo.GetAndMarkTempFileProcessing(ctx, cmd.Key)
	if err != nil {
		return CommitTempResult{}, err
	}

	dstObjName, err := m.computeHash(ctx, tempFile.Key)
	if err != nil {
		return CommitTempResult{}, err
	}

	exists, err := m.repo.FileExists(ctx, dstObjName)
	if err != nil {
		return CommitTempResult{}, err
	}

	var finalFileHash string
	if exists {
		finalFileHash = dstObjName
	} else {
		finalFileHash, err = m.copyToPermanentStorage(ctx, tempFile, dstObjName)
		if err != nil {
			return CommitTempResult{}, err
		}
	}

	err = m.repo.MarkTempFileProcessed(ctx, cmd.Key, finalFileHash)
	if err != nil {
		// Critical error: file is copied but not marked as processed
		//h.publishCleanupMessage(ctx, tempFile.Key, finalFileHash)

		//For now, simply logging it
		//TODO retry
		log.Error().Err(err).Str("key", tempFile.Key).Msg("failed to mark temp file as processed")
		//Since process is already pretty much done, finalize it anyway
	}

	err = m.publishCreatedFileEvent(ctx, FileCreated{
		TempID:  tempFile.Key,
		PermID:  finalFileHash,
		Existed: exists,
	})
	if err != nil {
		//Not sure how to handle it
		log.Error().Err(err).Str("key", tempFile.Key).Msg("failed to publish created file event")
	}

	return CommitTempResult{finalFileHash}, nil
}

func (m commitTempHandler) copyToPermanentStorage(
	ctx context.Context,
	tempFile file.TempFile,
	dstObjName string,
) (string, error) {
	dst := minio.CopyDestOptions{
		Bucket: "files",
		Object: dstObjName,
	}
	src := minio.CopySrcOptions{
		Bucket: "temp-bucket",
		Object: tempFile.Key,
	}

	info, err := m.minioClient.CopyObject(ctx, dst, src)
	if err != nil {
		return "", err
	}

	err = m.repo.CreateFile(ctx, file.CreateFileParams{
		Hash:        info.Key,
		FileSize:    int32(info.Size),
		ContentType: tempFile.ContentType,
	})
	if err != nil {
		// If we fail here, the file is copied but not recorded in the database
		// We should try to remove it from the permanent storage
		cleanupErr := m.minioClient.RemoveObject(ctx, "files", dstObjName, minio.RemoveObjectOptions{})
		if cleanupErr != nil {
			log.Printf("Failed to clean up orphaned file %s: %v", dstObjName, cleanupErr)
		}
		return "", err
	}

	return info.Key, nil
}

func (m commitTempHandler) computeHash(ctx context.Context, object string) (string, error) {
	obj, err := m.minioClient.GetObject(ctx, "temp-bucket", object, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer obj.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, obj); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type FileCreated struct {
	TempID  string `json:"temp_id"`
	PermID  string `json:"perm_id"`
	Existed bool   `json:"existed"`
}

func (m commitTempHandler) publishCreatedFileEvent(ctx context.Context, event FileCreated) error {
	return m.eventBus.Publish(ctx, &event)
}
