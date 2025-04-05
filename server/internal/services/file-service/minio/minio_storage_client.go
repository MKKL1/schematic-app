package minio

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/image"
	"github.com/minio/minio-go/v7"
	"io"
)

type StorageClient struct {
	minioClient         *minio.Client
	permanentFileBucket string //files
	temporaryFileBucket string //temp-bucket
	imageBucket         string //images
}

func NewMinioStorageClient(minioClient *minio.Client, permanentFileBucket, temporaryFileBucket, imageBucket string) *StorageClient {
	return &StorageClient{minioClient, permanentFileBucket, temporaryFileBucket, imageBucket}
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

func (m *StorageClient) PutOriginal(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (file.StorageInfo, error) {
	// Ensure bucket exists (optional, depending on setup)
	// m.ensureBucket(ctx, m.originalsBucket)

	uploadInfo, err := m.minioClient.PutObject(ctx, m.imageBucket, key, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return file.StorageInfo{}, fmt.Errorf("%w: %w", image.ErrStorageFailed, err)
	}

	return file.StorageInfo{
		Key:         uploadInfo.Key,
		Size:        uploadInfo.Size,
		ContentType: contentType, // Use provided type, Minio might not return it accurately here
		ETag:        uploadInfo.ETag,
	}, nil
}

func (m *StorageClient) GetOriginal(ctx context.Context, key string) (io.ReadCloser, string, error) {
	obj, err := m.minioClient.GetObject(ctx, m.imageBucket, key, minio.GetObjectOptions{})
	if err != nil {
		// TODO: Map minio errors (e.g., NoSuchKey) to domain errors
		return nil, "", fmt.Errorf("failed to get original object %s: %w", key, err)
	}

	// Stat the object to get ContentType and other metadata reliably
	stat, err := obj.Stat()
	if err != nil {
		obj.Close() // Close object if stat fails
		// TODO: Map minio errors
		return nil, "", fmt.Errorf("failed to stat original object %s: %w", key, err)
	}

	return obj, stat.ContentType, nil
}

func (m *StorageClient) DeleteOriginal(ctx context.Context, key string) error {
	err := m.minioClient.RemoveObject(ctx, m.imageBucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		// TODO: Map minio errors
		return fmt.Errorf("failed to delete original object %s: %w", key, err)
	}
	return nil
}

// Ensure StorageClient implements both interfaces if needed, or use separate clients
var _ file.StorageClient = (*StorageClient)(nil)
var _ file.StorageClient = (*StorageClient)(nil)
