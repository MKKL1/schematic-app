package ports

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/command"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ImageUploaded struct {
	TempID   string            `json:"temp_id"`
	PermID   string            `json:"perm_id"`
	Existed  bool              `json:"existed"`
	Type     string            `json:"type"` //Should be "image"
	Metadata map[string]string `json:"metadata"`
}

type EventHandlers struct {
	app app.Application
}

func NewEventHandlers(app app.Application, handler kafka.CqrsHandler) *EventHandlers {
	eh := &EventHandlers{app: app}
	err := handler.EventProcessor.AddHandlers(
		cqrs.NewEventHandler("ProcessUploadedImage", eh.handleFileUploaded),
	)
	if err != nil {
		return nil
	}

	return eh
}

func (eh EventHandlers) handleFileUploaded(ctx context.Context, cmd *ImageUploaded) error {
	if cmd.Type != "image" {
		return fmt.Errorf("invalid event ImageUploaded file type ('%s' should be 'image')", cmd.TempID)
	}

	imageType, found := cmd.Metadata["image_type"]
	if !found {
		return fmt.Errorf("image type not found in metadata")
	}

	_, err := eh.app.Commands.ProcessUploadedImage.Handle(ctx, command.ProcessUploadedImageCmd{
		Hash:      cmd.PermID,
		ImageType: imageType,
	})

	return err
}
