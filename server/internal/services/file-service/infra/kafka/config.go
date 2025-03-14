package kafka

import "github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"

func NewKafkaSubscriber(brokers []string) (*kafka.Subscriber, error) {
	config := kafka.SubscriberConfig{
		Brokers:               brokers,
		Unmarshaler:           kafka.DefaultMarshaler{},
		OverwriteSaramaConfig: kafka.DefaultSaramaSubscriberConfig(),
	}
	return kafka.NewSubscriber(config, nil)
}

func NewKafkaPublisher(brokers []string) (*kafka.Publisher, error) {
	config := kafka.PublisherConfig{
		Brokers:   brokers,
		Marshaler: kafka.DefaultMarshaler{},
	}
	return kafka.NewPublisher(config, nil)
}
