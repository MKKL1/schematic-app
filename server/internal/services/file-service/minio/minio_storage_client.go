package minio

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/minio/minio-go/v7"
	"io"
)

type StorageClient struct {
	minioClient         *minio.Client
	permanentFileBucket string //files
	temporaryFileBucket string //temp-bucket
}

func NewMinioStorageClient(minioClient *minio.Client, permanentFileBucket string, temporaryFileBucket string) *StorageClient {
	return &StorageClient{minioClient, permanentFileBucket, temporaryFileBucket}
}

// --- File Domain Methods ---

func (m StorageClient) CopyTempToPermanent(ctx context.Context, tempObject string, permObject string) (file.UploadInfo, error) {
	dst := minio.CopyDestOptions{
		Bucket: m.permanentFileBucket,
		Object: permObject,
	}
	src := minio.CopySrcOptions{
		Bucket: m.temporaryFileBucket,
		Object: tempObject,
	}

	info, err := m.minioClient.CopyObject(ctx, dst, src)
	if err != nil {
		return file.UploadInfo{}, err
	}

	return toDomainObject(info), nil
}

func (m StorageClient) RemovePermObject(ctx context.Context, permObject string) error {
	cleanupErr := m.minioClient.RemoveObject(ctx, m.permanentFileBucket, permObject, minio.RemoveObjectOptions{})
	return cleanupErr
}

func (m StorageClient) RemoveTempObjects(ctx context.Context, tempObject <-chan string) <-chan file.RemoveObjectError {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for o := range tempObject {
			objectsCh <- minio.ObjectInfo{Key: o}
		}
	}()

	domainErrCh := make(chan file.RemoveObjectError)
	errCh := m.minioClient.RemoveObjects(ctx, m.temporaryFileBucket, objectsCh, minio.RemoveObjectsOptions{})
	for e := range errCh {
		domainErrCh <- file.RemoveObjectError{
			ObjectName: e.ObjectName,
			Err:        e.Err,
		}
	}

	return domainErrCh
}

func (m StorageClient) GetTempObject(ctx context.Context, tempObject string) (file.Object, error) {
	obj, err := m.minioClient.GetObject(ctx, m.temporaryFileBucket, tempObject, minio.GetObjectOptions{})
	return obj, err
}

func (m StorageClient) PutTempObject(ctx context.Context, tempObject string, reader io.Reader, contentType string) (file.UploadInfo, error) {
	info, err := m.minioClient.PutObject(ctx, m.temporaryFileBucket, tempObject, reader, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return file.UploadInfo{}, err
	}

	return toDomainObject(info), nil
}

func toDomainObject(info minio.UploadInfo) file.UploadInfo {
	return file.UploadInfo{
		Key:  info.Key,
		Size: info.Size,
	}
}

// Ensure StorageClient implements both interfaces if needed, or use separate clients
var _ file.StorageClient = (*StorageClient)(nil)
