package kafka

import (
	"MyEvents/booking-api/src/msgqueue"
	"encoding/json"
	"github.com/Shopify/sarama"
)

type kafkaEventEmitter struct {
	topic    string
	producer sarama.SyncProducer
}

func NewKafkaEventEmitter(client sarama.Client, topic string) msgqueue.EventEmitter {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return &kafkaEventEmitter{
		producer: producer,
		topic: topic,
	}
}

func (emitter *kafkaEventEmitter) Emit(event msgqueue.Event) error {
	message := &message{
		EventName: event.EventName(),
		Payload:   event,
	}
	jsonBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: emitter.topic,
		Value: sarama.ByteEncoder(jsonBody),
	}

	_, _, err = emitter.producer.SendMessage(msg)
	return err
}
