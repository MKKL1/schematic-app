package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

//TODO refactor to better fit command pattern

type FileCommitCommandParams struct {
	Id string
}

type PostCreatedParams struct {
	Files []string `json:"files"`
}

type PostCreatedHandler decorator.CommandHandler[PostCreatedParams, any]

type postCreatedHandler struct {
	eventBus *cqrs.EventBus
}

func NewPostCreatedHandler(eventBus *cqrs.EventBus) PostCreatedHandler {
	return postCreatedHandler{eventBus}
}

func (m postCreatedHandler) Handle(ctx context.Context, cmd PostCreatedParams) (any, error) {
	for _, f := range cmd.Files {
		err := m.eventBus.Publish(ctx, FileCommitCommandParams{Id: f})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
