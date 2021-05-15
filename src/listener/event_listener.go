package listener

import (
	"MyEvents/booking-api/src/daos"
	"MyEvents/booking-api/src/models"
	"MyEvents/booking-api/src/msgqueue"
	"context"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type EventProcessor interface {
	ProcessEvent() error
}

type eventProcessor struct {
	eventListener msgqueue.EventListener
	eventDao      daos.EventsDao
	locationsDao  daos.LocationsDao
}

func NewEventProcessor(listener msgqueue.EventListener, eventsDao daos.EventsDao, locationsDao daos.LocationsDao) EventProcessor {
	return &eventProcessor{
		eventListener: listener,
		eventDao: eventsDao,
		locationsDao: locationsDao,
	}
}

func (ep *eventProcessor) ProcessEvent() error {
	log.Println("Listening to events...")
	eventsChan, errorsChan, err := ep.eventListener.Listen("event.created")
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-eventsChan:
			ep.handleEvent(context.Background(), event)
		case err := <-errorsChan:
			log.Printf("received error while processing msg: %s", err)
		}
	}
}

func (ep *eventProcessor) handleEvent(context context.Context, event msgqueue.Event) {
	switch e := event.(type) {
	case *models.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)
		newEvent := &models.Event{
			ID: bson.ObjectIdHex(e.ID),
			Name: e.Name,
			StartDate: e.Start.Unix(),
			EndDate: e.End.Unix(),
		}
		ep.eventDao.AddEvent(context, newEvent)
	case *models.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
		ep.locationsDao.AddLocation(context, &models.Location{ID: bson.ObjectId(e.ID)})
	default:
		log.Printf("unknown event: %t", e)
	}
}
