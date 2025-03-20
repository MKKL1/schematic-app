package file

import (
	"context"
	"io"
)

type UploadInfo struct {
	Key  string
	Size int64
}

type Object interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type RemoveObjectError struct {
	ObjectName string
	Err        error
}

type StorageClient interface {
	CopyTempToPermanent(ctx context.Context, tempObject string, permObject string) (UploadInfo, error)
	RemovePermObject(ctx context.Context, permObject string) error
	RemoveTempObjects(ctx context.Context, tempObject <-chan string) <-chan RemoveObjectError
	GetTempObject(ctx context.Context, tempObject string) (Object, error)
	PutTempObject(ctx context.Context, tempObject string, reader io.Reader, contentType string) (UploadInfo, error)
}
