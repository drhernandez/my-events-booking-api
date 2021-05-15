package daos

import (
	"MyEvents/booking-api/src/clients"
	"MyEvents/booking-api/src/models"
	"context"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const LOCATIONS = "locations"

type LocationsDao interface {
	AddLocation(ctx context.Context, location *models.Location) ([]byte, error)
}

type locationsDao struct {
	db *clients.MongoDBClient
}

func NewLocationsDao(db *clients.MongoDBClient) LocationsDao {
	return &locationsDao{
		db: db,
	}
}

func (dao *locationsDao) AddLocation(ctx context.Context, location *models.Location) ([]byte, error) {
	s := dao.db.GetFreshSession()
	defer s.Close()
	if !location.ID.Valid() {
		location.ID = bson.NewObjectId()
	}

	err := s.DB(clients.DB).C(LOCATIONS).Insert(location)
	if err != nil {
		log.Printf("[FIL: LocationsDao] [METHOD: AddLocation] [ERROR: %s]", err.Error())
		return nil, err
	}

	return []byte(location.ID), nil
}
