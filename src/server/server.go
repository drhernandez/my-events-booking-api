package server

import (
	"MyEvents/boocking-api/src/clients"
	"MyEvents/boocking-api/src/config"
	"MyEvents/boocking-api/src/daos"
	"MyEvents/boocking-api/src/handlers"
	"MyEvents/boocking-api/src/listener"
	msgqueue_amqp "MyEvents/boocking-api/src/msgqueue/amqp"
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
	eventListener, err := msgqueue_amqp.NewAMQPEventListener(conn, "events")
	if err != nil {
		panic(err)
	}

	//Async instances
	eventProcessor := listener.NewEventProcessor(eventListener, eventsDao, locationsDao)
	go eventProcessor.ProcessEvent()

	healthHandler := handlers.HealthHandler{}
	bookingsHandler := handlers.NewBookingsHandler(bookingsDao, eventsEmitter)

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/ping").HandlerFunc(healthHandler.PingHandler)
	r.Methods(http.MethodPost).Path("/bookings").HandlerFunc(bookingsHandler.CreateBooking)

	return r
}
