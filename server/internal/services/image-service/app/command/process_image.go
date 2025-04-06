package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/domain/image"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/rs/zerolog"
)

type ProcessUploadedImageCmd struct {
	Hash      string
	ImageType string
}

type ProcessUploadedImage decorator.CommandHandler[ProcessUploadedImageCmd, any]

type processUploadedImage struct {
	repo     image.Repository
	eventBus *cqrs.EventBus
	logger   zerolog.Logger
}

func NewProcessUploadedImageHandler(repo image.Repository, eventBus *cqrs.EventBus, logger zerolog.Logger) ProcessUploadedImage {
	return processUploadedImage{repo, eventBus, logger}
}

func (p processUploadedImage) Handle(ctx context.Context, cmd ProcessUploadedImageCmd) (any, error) {
	//File is already saved permanently, so now we have to generate what image sizes it uses

	//TODO check if already exists
	err := p.repo.CreateImage(ctx, image.CreateParams{
		FileHash:  cmd.Hash,
		ImageType: cmd.ImageType,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
