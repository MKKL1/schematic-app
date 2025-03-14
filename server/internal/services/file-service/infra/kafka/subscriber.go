package kafka

import (
	"context"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type KafkaSubscriber struct {
	subscriber *kafka.Subscriber
}

func NewKafkaSubscriberInstance(brokers []string) (*KafkaSubscriber, error) {
	sub, err := NewKafkaSubscriber(brokers)
	if err != nil {
		return nil, err
	}
	return &KafkaSubscriber{subscriber: sub}, nil
}

func (s *KafkaSubscriber) Subscribe(topic string, handler func(msg *message.Message)) error {
	messages, err := s.subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			handler(msg)
			msg.Ack()
		}
	}()

	return nil
}
