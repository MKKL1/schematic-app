package image

import "context"

type SignedImageUrlGenerator interface {
	GetSignedUrl(ctx context.Context, imageHash string, preset Preset) (string, error)
	GetSignedUrlBulk(ctx context.Context, imageHash string, presets []Preset) ([]string, error)
}
