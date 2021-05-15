package amqp

import (
	"MyEvents/booking-api/src/msgqueue"
	"encoding/json"
	"github.com/streadway/amqp"
)

type amqpEventEmitter struct {
	connection *amqp.Connection
	exchange   string
}

func NewAMQPEventEmitter(conn *amqp.Connection, exchange string) (msgqueue.EventEmitter, error) {
	emitter := &amqpEventEmitter{
		connection: conn,
		exchange: exchange,
	}

	err := emitter.setup()
	if err != nil {
		return nil, err
	}

	return emitter, nil
}

func (emitter *amqpEventEmitter) setup() error {
	channel, err := emitter.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return channel.ExchangeDeclare(emitter.exchange, "topic", true, false, false, false, nil)
}

func (emitter *amqpEventEmitter) Emit(event msgqueue.Event) error {
	jsonDoc, err := json.Marshal(event)
	if err != nil {
		return err
	}

	channel, err := emitter.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.EventName()},
		Body:        jsonDoc,
		ContentType: "application/json",
	}

	return channel.Publish(emitter.exchange, event.EventName(), false, false, msg)
}
