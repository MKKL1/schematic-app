package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/infra/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type CommitFile struct {
	Id string `json:"id"`
}

type EventHandlers struct {
	app app.Application
}

func NewEventHandlers(app app.Application, handler kafka.CqrsHandler) *EventHandlers {
	eh := &EventHandlers{app: app}
	err := handler.CommandProcessor.AddHandlers(
		cqrs.NewCommandHandler("CommitFile", eh.commitFileHandler),
	)
	if err != nil {
		return nil
	}

	return eh
}

func (eh *EventHandlers) commitFileHandler(ctx context.Context, cmd *CommitFile) error {
	_, err := eh.app.Commands.CommitTempFile.Handle(ctx, command.CommitTempParams{Key: cmd.Id})
	return err
}
