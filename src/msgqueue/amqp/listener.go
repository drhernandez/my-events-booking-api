package amqp

import (
	"MyEvents/boocking-api/src/models"
	"MyEvents/boocking-api/src/msgqueue"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
}

func (listener *amqpEventListener) setup() error {
	channel, err := listener.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare("events", true, false, false, false, nil)
	return err
}

func NewAMQPEventListener(connection *amqp.Connection, queue string) (msgqueue.EventListener, error) {
	listener := &amqpEventListener{
		connection: connection,
		queue: queue,
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (listener *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := listener.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()

	for _, eventName := range eventNames {
		err := channel.QueueBind("events", eventName, "events", false, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	msgsChan, err := channel.Consume("events", "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgsChan {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("msg did not contain x-event-name header")
				msg.Nack(false, false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("event name is not a string but %t", rawEventName)
				msg.Nack(false, false)
				continue
			}

			var event msgqueue.Event
			switch eventName {
				case "event.created":
					event = &models.EventCreatedEvent{}
				default:
					errors <- fmt.Errorf("event type %s is unknown", eventName)
					msg.Nack(false, false)
					continue
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errors <- err
				continue
			}

			events <- event
		}
	}()

	return events, errors, nil
}