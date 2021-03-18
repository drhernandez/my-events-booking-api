package daos

import (
	"MyEvents/boocking-api/src/clients"
	"MyEvents/boocking-api/src/models"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const BOOKINGS = "bookings"

type bookingsDao struct {
	db *clients.MongoDBClient
}

type BookingsDao interface {
	AddBooking(booking *models.Booking) error
}

func NewBookingDao(db *clients.MongoDBClient) BookingsDao {
	return &bookingsDao{
		db: db,
	}
}

func (dao *bookingsDao) AddBooking(booking *models.Booking) error {
	s := dao.db.GetFreshSession()
	defer s.Close()

	if !booking.ID.Valid() {
		booking.ID = bson.NewObjectId()
	}

	err := s.DB(clients.DB).C(BOOKINGS).Insert(booking)
	if err != nil {
		log.Printf("[FIL: BookingsDao] [METHOD: AddBooking] [ERROR: %s]", err.Error())
	}
	return err
}