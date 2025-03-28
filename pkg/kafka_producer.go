package pkg

import "github.com/IBM/sarama"

type Producer sarama.SyncProducer

func NewProducer(brokers []string) (Producer, error) {
	// Configure Sarama producer settings
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 5                    // Retry up to 5 times
	config.Producer.Return.Successes = true          // Required for SyncProducer

	// Create a synchronous producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
