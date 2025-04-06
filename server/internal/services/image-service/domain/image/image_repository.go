package image

import (
	"context"
)

type CreateParams struct {
	FileHash  string
	ImageType string
}

type Repository interface {
	CreateImage(ctx context.Context, params CreateParams) error
	GetImageTypesForHash(ctx context.Context, hash string) ([]string, error)
}
