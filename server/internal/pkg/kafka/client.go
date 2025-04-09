package kafka

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/zerowater"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

type CqrsHandler struct {
	CommandBus       *cqrs.CommandBus
	CommandProcessor *cqrs.CommandProcessor
	EventBus         *cqrs.EventBus
	EventProcessor   *cqrs.EventProcessor
	router           *message.Router
	logger           zerolog.Logger
}

type KafkaConfig struct {
	Brokers []string `koanf:"brokers"`
}

func NewCqrsHandler(cfg KafkaConfig, zerologger zerolog.Logger) CqrsHandler {
	compLogger := zerologger.With().Str("component", "cqrs-handler").Logger()
	logger := zerowater.NewZerologLoggerAdapter(compLogger)
	watermillLogger := zerowater.NewZerologLoggerAdapterMapped(zerologger.With().Str("component", "windmill").Logger())

	cqrsMarshaler := cqrs.JSONMarshaler{
		NewUUID:      uuid.New().String,
		GenerateName: cqrs.StructName,
	}

	kafkaMarshaler := kafka.NewWithPartitioningMarshaler(func(topic string, msg *message.Message) (string, error) {
		return msg.Metadata.Get("partition_key"), nil
	})

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   cfg.Brokers,
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

	retryMiddleware := middleware.Retry{
		MaxRetries:          5,
		InitialInterval:     500 * time.Millisecond,
		Multiplier:          3.0,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      60 * time.Second,
		RandomizationFactor: 0.5,
		Logger:              logger,
	}

	poisMiddleware, err := middleware.PoisonQueue(publisher, "failed.events")
	if err != nil {
		panic(err)
	}
	router.AddMiddleware(
		middleware.Recoverer,
		poisMiddleware,
		retryMiddleware.Middleware,
		middleware.CorrelationID,
	)

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
						Brokers:       cfg.Brokers,
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
						Brokers:       cfg.Brokers,
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

	return CqrsHandler{
		CommandBus:       commandBus,
		CommandProcessor: commandProcessor,
		EventBus:         eventBus,
		EventProcessor:   eventProcessor,
		router:           router,
		logger:           compLogger,
	}
}

func (h *CqrsHandler) Run(ctx context.Context) {
	err := h.router.Run(ctx)
	if err != nil {
		h.logger.Fatal().Err(err).Msg("Error running router")
		return
	}
}

func (h *CqrsHandler) Close(ctx context.Context) error {
	h.logger.Info().Msg("Closing Kafka connections...")
	err := h.router.Close()
	if err != nil {
		return fmt.Errorf("closing Kafka router: %w", err)
	}

	h.logger.Info().Msg("Kafka connections closed.")
	return nil
}
