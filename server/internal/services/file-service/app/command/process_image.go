package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"

	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/rs/zerolog"
)

type UploadImageParams struct {
}

type UploadImageResult struct {
}

type UploadImageHandler decorator.CommandHandler[UploadImageParams, UploadImageResult]

type uploadImageHandler struct {
	repo     file.ImageRepository
	eventBus *cqrs.EventBus
	logger   zerolog.Logger
}

func (h uploadImageHandler) Handle(ctx context.Context, cmd UploadImageParams) (UploadImageResult, error) {
	//File is already saved permanently, so now we have to generate what image sizes it uses

}
