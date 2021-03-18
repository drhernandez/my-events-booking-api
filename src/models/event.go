package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Event struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}

type Location struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty"`
	Capacity int    `json:"capacity"`
}

type EventCreatedEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LocationID string    `json:"location_id"`
	Start      time.Time `json:"start_time"`
	End        time.Time `json:"end_time"`
}

func (event *EventCreatedEvent) EventName() string {
	return "event.created"
}

type LocationCreatedEvent struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Country string `json:"country"`
	Halls   []Hall `json:"halls"`
}

func (event *LocationCreatedEvent) EventName() string {
	return "location.created"
}
