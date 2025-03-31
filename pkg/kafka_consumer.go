package pkg

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	Client *kafka.Consumer
}

func NewConsumer(bootstrapServers string) (*Consumer, error) {
	// Kafka configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers":                  bootstrapServers,
		"group.id":                           "webhook-notifier",
		"auto.offset.reset":                  "earliest",
		"metadata.max.age.ms":                1000,
		"topic.metadata.refresh.interval.ms": 1000,
	}

	// Create a new consumer instance
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer,
	}, nil
}

func (c *Consumer) Close() {
	err := c.Client.Close()
	Logger.Error(err)
}
