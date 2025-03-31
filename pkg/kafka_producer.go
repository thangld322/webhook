package pkg

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	Client *kafka.Producer
}

func NewProducer(bootstrapServers string) (*Producer, error) {
	// Create a new producer using default configuration
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":                  bootstrapServers,
		"metadata.max.age.ms":                1000,
		"topic.metadata.refresh.interval.ms": 1000,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer,
	}, nil
}

func (p *Producer) Produce(topic string, event []byte) error {
	return p.Client.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          event,
	}, nil)
}

func (p *Producer) Close() {
	p.Client.Flush(15000)
	p.Client.Close()
}
