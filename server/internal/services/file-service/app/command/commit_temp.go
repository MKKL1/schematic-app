package command

import (
	"context"
	"errors"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/minio/minio-go/v7"
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
}

func NewCommitTempHandler(minioClient *minio.Client, repo file.Repository) CommitTempHandler {
	return commitTempHandler{minioClient, repo}
}

func (m commitTempHandler) Handle(ctx context.Context, cmd CommitTempParams) (CommitTempResult, error) {
	tempFile, err := m.repo.GetTempFile(ctx, cmd.Key)
	if err != nil {
		return CommitTempResult{}, err
	}

	if tempFile.Key == "" {
		return CommitTempResult{}, errors.New("key is empty")
	}

	attributes, err := m.minioClient.GetObjectAttributes(ctx, "temp-bucket", tempFile.Key, minio.ObjectAttributesOptions{})
	if err != nil {
		return CommitTempResult{}, err
	}

	dstObjName := attributes.Checksum.ChecksumCRC32C
	exists, err := m.repo.FileExists(ctx, dstObjName)
	if err != nil {
		return CommitTempResult{}, err
	}

	var finalFileHash string
	if exists {
		finalFileHash = dstObjName
	} else {
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
			return CommitTempResult{}, err
		}

		err = m.repo.CreateFile(ctx, file.CreateFileParams{
			Hash:        info.Key,
			FileSize:    int32(info.Size),
			ContentType: tempFile.ContentType,
		})
		if err != nil {
			return CommitTempResult{}, err
		}
		finalFileHash = info.Key
	}

	err = m.minioClient.RemoveObject(ctx, "temp-bucket", tempFile.Key, minio.RemoveObjectOptions{})
	if err != nil {
		return CommitTempResult{}, err
	}

	err = m.repo.DeleteTmpFilesByKey(ctx, []string{tempFile.Key})
	if err != nil {
		return CommitTempResult{}, err
	}

	return CommitTempResult{finalFileHash}, nil
}
