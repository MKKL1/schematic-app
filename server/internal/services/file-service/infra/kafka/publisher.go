package kafka

import (
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type KafkaPublisher struct {
	Publisher *kafka.Publisher
}

func NewKafkaPublisherInstance(brokers []string) (*KafkaPublisher, error) {
	pub, err := NewKafkaPublisher(brokers)
	if err != nil {
		return nil, err
	}
	return &KafkaPublisher{Publisher: pub}, nil
}

func (p *KafkaPublisher) Publish(topic string, payload []byte) error {
	msg := message.NewMessage(uuid.New().String(), payload)
	return p.Publisher.Publish(topic, message.Messages{msg}...)
}
