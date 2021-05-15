package kafka

import (
	"MyEvents/booking-api/src/models"
	"MyEvents/booking-api/src/msgqueue"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"log"
)

type kafkaEventListener struct {
	consumer   sarama.Consumer
	partitions []int32
	topic      string
}

func NewKafkaEventListener(client sarama.Client, topic string) msgqueue.EventListener {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		panic(err)
	}

	return &kafkaEventListener{
		consumer:   consumer,
		topic: topic,
	}
}

func (listener *kafkaEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	var err error
	events := make(chan msgqueue.Event)
	errors := make(chan error)

	partitions := listener.partitions
	if len(partitions) == 0 {
		partitions, err = listener.consumer.Partitions(listener.topic)
		if err != nil {
			return nil, nil, err
		}
	}
	log.Printf("Topic %s has partitions: %v", listener.topic, partitions)

	for _, partition := range partitions {
		con, err := listener.consumer.ConsumePartition(listener.topic, partition, 0)
		if err != nil {
			return nil, nil, err
		}

		go func() {
			for msg := range con.Messages() {
				body := &message{}
				err := json.Unmarshal(msg.Value, body)
				if err != nil {
					errors <- fmt.Errorf("could not JSON-decode message: %s", err)
					continue
				}

				var event msgqueue.Event
				switch body.EventName {
				case "event.created":
					event = &models.EventCreatedEvent{}
				default:
					errors <- fmt.Errorf("unknown event name: %s", body.EventName)
					continue
				}

				cfg := &mapstructure.DecoderConfig{
					Result: event,
					TagName: "json",
				}

				decoder, err := mapstructure.NewDecoder(cfg)
				if err != nil {
					errors <- fmt.Errorf("could not create new decoder: %s", err)
				}

				err = decoder.Decode(body.Payload)
				if err != nil {
					errors <- fmt.Errorf("could not decode event %s: %s", event.EventName(), err)
				}

				events <- event
			}
		}()
	}

	return events, errors, nil
}
