package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
)

type FileCreated struct {
	TempId  string `json:"temp_id"`
	PermId  string `json:"perm_id"`
	Existed bool   `json:"existed"`
}

type EventHandlers struct {
	app app.Application
}

func NewEventHandlers(app app.Application, handler kafka.CqrsHandler) *EventHandlers {
	eh := &EventHandlers{app: app}
	err := handler.EventProcessor.AddHandlers(
		cqrs.NewEventHandler("UpdateAttachedFiles", eh.handleFileCreated),
	)
	if err != nil {
		return nil
	}

	return eh
}

func (eh EventHandlers) handleFileCreated(ctx context.Context, cmd *FileCreated) error {
	parse, err := uuid.Parse(cmd.TempId)
	if err != nil {
		return err
	}
	_, err = eh.app.Commands.UpdateFileHash.Handle(ctx, command.UpdateFileHashCommand{Ids: []command.FileHashTempId{
		{
			TempId: parse,
			Hash:   cmd.PermId,
		},
	}})

	return err
}
