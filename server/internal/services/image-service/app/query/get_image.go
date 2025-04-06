package query

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/domain/image"
	"github.com/rs/zerolog"
)

type GetImageParams struct {
	ImageID string
}
type GetImageResult struct {
	Sizes []GetImageResultSize
}

type GetImageResultSize struct {
	URL    string
	Preset image.Preset
}

type GetImageSizes decorator.QueryHandler[GetImageParams, GetImageResult]

type getImageHandler struct {
	repo   image.Repository
	urlGen image.SignedImageUrlGenerator
	logger zerolog.Logger
}

func NewGetImageHandler(repo image.Repository, urlGen image.SignedImageUrlGenerator, logger zerolog.Logger /*, cfg *SomeConfig */) GetImageSizes {
	return &getImageHandler{repo, urlGen, logger}
}

func (h getImageHandler) Handle(ctx context.Context, query GetImageParams) (GetImageResult, error) {
	imgTypes, err := h.repo.GetImageTypesForHash(ctx, query.ImageID)
	if err != nil {
		return GetImageResult{}, fmt.Errorf("getting image types for hash %s: %w", query.ImageID, err)
	}

	var allPresets []image.Preset
	for _, imgType := range imgTypes {
		presets, err := image.GeneratePreset(imgType)
		if err != nil {
			return GetImageResult{}, fmt.Errorf("generating presets for %s: %w", imgType, err)
		}
		allPresets = append(allPresets, presets...)
	}

	sizes := make([]GetImageResultSize, len(allPresets))
	urls, err := h.urlGen.GetSignedUrlBulk(ctx, query.ImageID, allPresets)
	if err != nil {
		return GetImageResult{}, fmt.Errorf("getting signed urls for %s: %w", query.ImageID, err)
	}

	for i, url := range urls {
		sizes[i] = GetImageResultSize{
			URL:    url,
			Preset: allPresets[i],
		}
	}

	return GetImageResult{Sizes: sizes}, nil
}
