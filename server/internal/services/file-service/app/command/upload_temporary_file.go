package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/minio/minio-go/v7"
	"io"
)

type UploadTempFileParams struct {
	Reader      io.Reader
	FileName    string
	FileSize    int64
	ContentType string
}

type UploadTempFileHandler decorator.CommandHandler[UploadTempFileParams, int64]

type uploadTempFileHandler struct {
	minioClient *minio.Client
}

func NewUploadTempFileHandler(minioClient *minio.Client) UploadTempFileHandler {
	return uploadTempFileHandler{
		minioClient: minioClient,
	}
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileParams) (int64, error) {
	info, err := u.minioClient.PutObject(ctx, "temp-bucket", cmd.FileName, cmd.Reader, -1, minio.PutObjectOptions{ContentType: cmd.ContentType})
	if err != nil {
		return 0, err
	}

	return info.Size, nil
}
