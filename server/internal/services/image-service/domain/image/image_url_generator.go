package image

import "context"

type SignedImageUrlGenerator interface {
	GetSignedUrl(ctx context.Context, imageHash string, presetName string) (string, error)
}
