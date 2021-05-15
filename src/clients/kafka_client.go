package clients

import (
	"github.com/Shopify/sarama"
)

func NewKafkaClient(kafkaMessageBrokers []string) sarama.Client {
	kafkaConfig := sarama.NewConfig()
	client, err := sarama.NewClient(kafkaMessageBrokers, kafkaConfig)
	if err != nil {
		panic(err)
	}

	return client
}
