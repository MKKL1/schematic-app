package command

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/decorator"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bytedance/sonic"
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
	pub message.Publisher
}

func NewPostCreatedHandler(pub message.Publisher) PostCreatedHandler {
	return postCreatedHandler{pub}
}

func (m postCreatedHandler) Handle(ctx context.Context, cmd PostCreatedParams) (any, error) {
	for _, f := range cmd.Files {
		payload, err := sonic.Marshal(FileCommitCommandParams{Id: f})
		if err != nil {
			//TODO we don't have to fail entire event handle process
			return nil, err
		}
		msg := message.NewMessage(watermill.NewUUID(), payload)
		err = m.pub.Publish("file.cmd.commit", msg)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
