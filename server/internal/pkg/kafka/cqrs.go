package kafka

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"log/slog"
)

type KafkaConfig struct {
	Brokers []string
}

type CqrsHandler struct {
	CommandBus       *cqrs.CommandBus
	CommandProcessor *cqrs.CommandProcessor
	EventBus         *cqrs.EventBus
	EventProcessor   *cqrs.EventProcessor
	router           *message.Router
}

func NewCqrsHandler(config KafkaConfig) CqrsHandler {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	logger := watermill.NewSlogLoggerWithLevelMapping(nil, map[slog.Level]slog.Level{
		slog.LevelInfo: slog.LevelDebug,
	})

	watermillLogger := watermill.NewSlogLoggerWithLevelMapping(
		slog.With("watermill", true),
		map[slog.Level]slog.Level{
			slog.LevelInfo: slog.LevelDebug,
		},
	)

	cqrsMarshaler := cqrs.JSONMarshaler{
		NewUUID:      uuid.New().String,
		GenerateName: cqrs.StructName,
	}

	kafkaMarshaler := kafka.NewWithPartitioningMarshaler(func(topic string, msg *message.Message) (string, error) {
		return msg.Metadata.Get("partition_key"), nil
	})

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   config.Brokers,
			Marshaler: kafkaMarshaler,
		},
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(middleware.Recoverer)
	router.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			slog.Debug("Received message", "metadata", msg.Metadata)
			return h(msg)
		}
	})

	commandBus, err := cqrs.NewCommandBusWithConfig(publisher, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return "commands." + params.CommandName, nil
		},
		Marshaler: cqrsMarshaler,
		Logger:    logger,
	})
	if err != nil {
		panic(err)
	}

	eventBus, err := cqrs.NewEventBusWithConfig(publisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return "events." + params.EventName, nil
		},
		Marshaler: cqrsMarshaler,
		Logger:    logger,
	})
	if err != nil {
		panic(err)
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return "commands." + params.CommandName, nil
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return kafka.NewSubscriber(
					kafka.SubscriberConfig{
						Brokers:       config.Brokers,
						ConsumerGroup: params.HandlerName,
						Unmarshaler:   kafkaMarshaler,
					},
					watermillLogger,
				)
			},
			Marshaler: cqrsMarshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		panic(err)
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return "events." + params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return kafka.NewSubscriber(
					kafka.SubscriberConfig{
						Brokers:       config.Brokers,
						ConsumerGroup: params.HandlerName,
						Unmarshaler:   kafkaMarshaler,
					},
					watermillLogger,
				)
			},
			Marshaler: cqrsMarshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		panic(err)
	}

	//err = commandProcessor.AddHandlers(
	//	cqrs.NewCommandHandler("SubscribeHandler", SubscribeHandler{eventBus}.Handle),
	//	cqrs.NewCommandHandler("UnsubscribeHandler", UnsubscribeHandler{eventBus}.Handle),
	//	cqrs.NewCommandHandler("UpdateEmailHandler", UpdateEmailHandler{eventBus}.Handle),
	//)
	//if err != nil {
	//	panic(err)
	//}
	//
	//activityReadModel := NewActivityTimelineModel()
	//
	//// All messages from this group will have one subscription.
	//// When message arrives, Watermill will match it with the correct handler.
	//err = eventProcessor.AddHandlersGroup(
	//	"ActivityTimelineReadModel",
	//	cqrs.NewGroupEventHandler(activityReadModel.OnSubscribed),
	//	cqrs.NewGroupEventHandler(activityReadModel.OnUnsubscribed),
	//	cqrs.NewGroupEventHandler(activityReadModel.OnEmailUpdated),
	//)
	//if err != nil {
	//	panic(err)
	//}

	return CqrsHandler{
		CommandBus:       commandBus,
		CommandProcessor: commandProcessor,
		EventBus:         eventBus,
		EventProcessor:   eventProcessor,
		router:           router,
	}
}

func (h CqrsHandler) Run(ctx context.Context) {
	go func() {
		err := h.router.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Error running router")
			return
		}
	}()
}
