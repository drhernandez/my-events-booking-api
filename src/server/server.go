package server

import (
	"MyEvents/booking-api/src/clients"
	"MyEvents/booking-api/src/config"
	"MyEvents/booking-api/src/daos"
	"MyEvents/booking-api/src/handlers"
	"MyEvents/booking-api/src/listener"
	msgqueue_amqp "MyEvents/booking-api/src/msgqueue/amqp"
	handlers2 "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"net/http"
)

func NewServer(config config.ServiceConfig) *mux.Router {

	dbClient := clients.NewMongoDBClient(config.DBConnection)
	eventsDao := daos.NewEventsDao(dbClient)
	locationsDao := daos.NewLocationsDao(dbClient)
	bookingsDao := daos.NewBookingDao(dbClient)

	conn, err := amqp.Dial(config.AMQPMessageBroker)
	if err != nil {
		panic(err)
	}
	eventsEmitter, err := msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
	if err != nil {
		panic(err)
	}
	eventsListener, err := msgqueue_amqp.NewAMQPEventListener(conn, "events")
	if err != nil {
		panic(err)
	}

	//kafkaClient := clients.NewKafkaClient(config.KafkaMessageBrokers)
	//eventsEmitter := kafka.NewKafkaEventEmitter(kafkaClient, "events")
	//eventsListener := kafka.NewKafkaEventListener(kafkaClient, "events")

	//Async instances
	eventProcessor := listener.NewEventProcessor(eventsListener, eventsDao, locationsDao)
	go eventProcessor.ProcessEvent()

	healthHandler := handlers.HealthHandler{}
	bookingsHandler := handlers.NewBookingsHandler(bookingsDao, eventsEmitter)

	r := mux.NewRouter()
	r.Use(handlers2.CORS())
	r.Use(setApplicationJsonContentType)
	r.Methods(http.MethodGet).Path("/ping").HandlerFunc(healthHandler.PingHandler)
	r.Methods(http.MethodPost).Path("/bookings").HandlerFunc(bookingsHandler.CreateBooking)

	return r
}
