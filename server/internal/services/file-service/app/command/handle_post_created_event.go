package command

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/rs/zerolog"
)

//TODO it's more of a role for mapper service

// Event payload structure
type PostCreated struct {
	Files []PostCreatedFile `json:"files"`
}

type PostCreatedFile struct {
	TempId   string            `json:"tempId"`
	Type     string            `json:"type"` //Image,schematic,unknown
	Metadata map[string]string `json:"metadata"`
}

type PostCreatedHandler decorator.CommandHandler[PostCreated, any] // Remains a command handler internally

type postCreatedHandler struct {
	commandBus *cqrs.CommandBus
	logger     zerolog.Logger
}

func NewPostCreatedHandler(commandBus *cqrs.CommandBus, logger zerolog.Logger) PostCreatedHandler {
	return postCreatedHandler{commandBus: commandBus, logger: logger}
}

func (h postCreatedHandler) Handle(ctx context.Context, event PostCreated) (any, error) {
	for _, f := range event.Files {
		err := h.commandBus.Send(ctx, file.CommitFile{
			Id:       f.TempId,
			Type:     f.Type,
			Metadata: f.Metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send CommitFile command for %s: %w", f.TempId, err)
		}
	}
	return nil, nil
}
