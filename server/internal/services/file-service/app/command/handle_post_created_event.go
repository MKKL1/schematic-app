package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

//TODO refactor to better fit command pattern

type PostCreatedEvent struct {
	Files []string `json:"files"`
}

type PostCreatedHandler decorator.CommandHandler[PostCreatedEvent, any]

type postCreatedHandler struct {
	commandBus *cqrs.CommandBus
}

func NewPostCreatedHandler(commandBus *cqrs.CommandBus) PostCreatedHandler {
	return postCreatedHandler{commandBus}
}

func (m postCreatedHandler) Handle(ctx context.Context, cmd PostCreatedEvent) (any, error) {
	for _, f := range cmd.Files {
		err := m.commandBus.Send(ctx, file.CommitFile{
			Id: f,
		})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
