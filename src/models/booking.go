package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Booking struct {
	ID      bson.ObjectId `bson:"_id"`
	UserID  []byte
	EventID []byte
	Date    time.Time
	Seats   int
}

func (booking *Booking) Validate() error {
	if booking.Seats <= 0 {
		return fmt.Errorf("invalid number of seats: %d", booking.Seats)
	}

	return nil
}

type BookingRequest struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Seats   int    `json:"seats"`
}

func (req *BookingRequest) MapToBooking() *Booking {
	return &Booking{
		UserID:  []byte(req.UserID),
		EventID: []byte(req.EventID),
		Date:    time.Now().UTC(),
		Seats:   req.Seats,
	}
}

type BookingResponse struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	EventID string    `json:"event_id"`
	Date    time.Time `json:"date_created"`
	Seats   int       `json:"seats"`
}

type EventBookedEvent struct {
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
}

func (event *EventBookedEvent) EventName() string {
	return "event.booked"
}
