package daos

import (
	"MyEvents/booking-api/src/clients"
	"MyEvents/booking-api/src/models"
	"context"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const EVENTS = "events"

type EventsDao interface {
	AddEvent(ctx context.Context, event *models.Event) ([]byte, error)
	FindEvent(ctx context.Context, id []byte) (*models.Event, error)
	FindEventByName(ctx context.Context, name string) (*models.Event, error)
	FindAllAvailableEvents(ctx context.Context) ([]models.Event, error)
}

type eventsDao struct {
	db *clients.MongoDBClient
}

func NewEventsDao(dbClient *clients.MongoDBClient) EventsDao {
	return &eventsDao{
		db: dbClient,
	}
}

func (dao *eventsDao) AddEvent(ctx context.Context, event *models.Event) ([]byte, error) {
	s := dao.db.GetFreshSession()
	defer s.Close()
	if !event.ID.Valid() {
		event.ID = bson.NewObjectId()
	}

	err := s.DB(clients.DB).C(EVENTS).Insert(event)
	if err != nil {
		log.Printf("[FIL: EventsDao] [METHOD: AddEvent] [ERROR: %s]", err.Error())
		return nil, err
	}

	return []byte(event.ID), nil
}

func (dao *eventsDao) FindEvent(ctx context.Context, id []byte) (*models.Event, error) {
	s := dao.db.GetFreshSession()
	defer s.Close()

	event := &models.Event{}
	err := s.DB(clients.DB).C(EVENTS).FindId(bson.ObjectId(id)).One(event)
	if err != nil {
		log.Printf("[FIL: EventsDao] [METHOD: FindEvent] [ERROR: %s]", err.Error())
		return nil, err
	}

	return event, nil
}

func (dao *eventsDao) FindEventByName(ctx context.Context, name string) (*models.Event, error) {
	s := dao.db.GetFreshSession()
	defer s.Close()

	event := &models.Event{}
	err := s.DB(clients.DB).C(EVENTS).Find(bson.M{"name": name}).One(event)
	if err != nil {
		log.Printf("[FIL: EventsDao] [METHOD: FindEventByName] [ERROR: %s]", err.Error())
		return nil, err
	}

	return event, nil
}

func (dao *eventsDao) FindAllAvailableEvents(ctx context.Context) ([]models.Event, error) {
	s := dao.db.GetFreshSession()
	defer s.Close()

	var events []models.Event
	err := s.DB(clients.DB).C(EVENTS).Find(nil).All(events)
	if err != nil {
		log.Printf("[FIL: EventsDao] [METHOD: FindAllAvailableEvents] [ERROR: %s]", err.Error())
		return nil, err
	}

	return events, nil
}
