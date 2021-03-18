package clients

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

const (
	DB              = "bookings-db"
	//MONGODB  DBTYPE = "mongodb"
	//DYNAMODB DBTYPE = "dynamodb"
)

type DBTYPE string

type MongoDBClient struct {
	session *mgo.Session
}

func NewMongoDBClient(connection string) *MongoDBClient {
	fmt.Println("Connecting to database")
	s, err := mgo.Dial(connection)
	if err != nil {
		panic(err)
	}

	return &MongoDBClient{
		session: s,
	}
}

func (client *MongoDBClient) GetFreshSession() *mgo.Session {
	return client.session.Copy()
}
