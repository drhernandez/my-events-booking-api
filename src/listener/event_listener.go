package listener

import (
	"MyEvents/boocking-api/src/daos"
	"MyEvents/boocking-api/src/models"
	"MyEvents/boocking-api/src/msgqueue"
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
		ep.eventDao.AddEvent(context, &models.Event{ID: bson.ObjectId(e.ID)})
	case *models.LocationCreatedEvent:
		log.Printf("location %s created: %s", e.ID, e)
		ep.locationsDao.AddLocation(context, &models.Location{ID: bson.ObjectId(e.ID)})
	default:
		log.Printf("unknown event: %t", e)
	}
}
