package handlers

import (
	"MyEvents/booking-api/src/daos"
	"MyEvents/booking-api/src/models"
	"MyEvents/booking-api/src/msgqueue"
	"encoding/json"
	"fmt"
	"net/http"
)

type bookingsHandler struct {
	bookingsDao  daos.BookingsDao
	eventEmitter msgqueue.EventEmitter
}

func NewBookingsHandler(bd daos.BookingsDao, emitter msgqueue.EventEmitter) *bookingsHandler {
	return &bookingsHandler{
		bookingsDao:  bd,
		eventEmitter: emitter,
	}
}

func (bh *bookingsHandler) CreateBooking(writer http.ResponseWriter, request *http.Request) {
	bookingRequest := &models.BookingRequest{}
	err := json.NewDecoder(request.Body).Decode(bookingRequest)
	if err != nil {
		writer.WriteHeader(400)
		fmt.Fprintf(writer, "could not decode JSON body: %s", err)
		return
	}

	booking := bookingRequest.MapToBooking()
	if err := booking.Validate(); err != nil {
		writer.WriteHeader(400)
		fmt.Fprintf(writer, "invalid json: %s", err)
		return
	}

	err = bh.bookingsDao.AddBooking(booking)
	if err != nil {
		writer.WriteHeader(500)
		fmt.Fprintf(writer, "error creting booking: %s", err)
		return
	}

	event := &models.EventBookedEvent{
		EventID: bookingRequest.EventID,
		UserID:  bookingRequest.UserID,
	}
	err = bh.eventEmitter.Emit(event)
	if err != nil {
		fmt.Fprintf(writer, "error sending event booked notification: %s", err)
	}

	response := &models.BookingResponse{
		ID:      booking.ID.Hex(),
		UserID:  string(booking.UserID),
		EventID: string(booking.EventID),
		Date:    booking.Date,
		Seats:   booking.Seats,
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(response)
}
