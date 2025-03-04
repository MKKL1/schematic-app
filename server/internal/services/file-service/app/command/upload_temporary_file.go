package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type UploadTempFileParams struct {
}

type UploadTempFileHandler decorator.CommandHandler[UploadTempFileParams, int64]

type uploadTempFileHandler struct {
	minioClient *minio.Client
}

func (u uploadTempFileHandler) Handle(ctx context.Context, cmd UploadTempFileParams) (int64, error) {

	info, err := u.minioClient.FPutObject(ctx, "", "", "", minio.PutObjectOptions{ContentType: "contentType"})
	if err != nil {
		log.Fatalln(err)
	}

	u.minioClient.PutObject()
}
