package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
)

type EventListener struct {
	sub message.Subscriber
	App app.Application
}

func NewEventListener(ctx context.Context, sub message.Subscriber, app app.Application) *EventListener {

	//TODO bad way to register listeners
	ob := &EventListener{sub: sub, App: app}
	err := ob.registerPostCreatedEvent(ctx)
	if err != nil {
		return nil
	}

	err = ob.registerFileCommitCommandEvent(ctx)
	if err != nil {
		return nil
	}

	return ob
}

func (el EventListener) registerPostCreatedEvent(ctx context.Context) error {
	msgChannel, err := el.sub.Subscribe(ctx, "post.created")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgChannel {
			var params command.PostCreatedParams
			err := sonic.Unmarshal(msg.Payload, &params)
			if err != nil {
				log.Error().Err(err).Msg("failed to unmarshal post created event")
				continue
			}

			if params.Files == nil {
				log.Error().Msg("post created event has no files (nil pointer)")
			}

			_, err = el.App.Commands.PostCreatedHandler.Handle(ctx, params)
			if err != nil {
				//Retry
				log.Error().Err(err).Msg("could not handle post created event")
				continue
			}
			msg.Ack()
		}
	}()

	return nil
}

func (el EventListener) registerFileCommitCommandEvent(ctx context.Context) error {
	msgChannel, err := el.sub.Subscribe(ctx, "file.cmd.commit")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgChannel {
			var params command.FileCommitCommandParams
			err = sonic.Unmarshal(msg.Payload, &params)
			if err != nil {
				log.Error().Err(err).Msg("failed to unmarshal file commit command event")
				continue
			}

			_, err = el.App.Commands.CommitTempFile.Handle(ctx, command.CommitTempParams{Key: params.Id})
			if err != nil {
				//Retry
				log.Error().Err(err).Msg("could not handle file commit command event")
				continue
			}
			msg.Ack()
		}
	}()

	return nil
}
