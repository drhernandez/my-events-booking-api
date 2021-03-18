package config

import (
	"MyEvents/boocking-api/src/clients"
	"encoding/json"
	"fmt"
	"os"
)

var (
	DBTypeDefault            = clients.DBTYPE("mongodb")
	DBConnectionDefault      = "mongodb://127.0.0.1"
	RestfulEPDefault         = "localhost:8081"
	SecureRestfulEPDefault   = "localhost:8444"
	AMQPMessageBrokerDefault = "amqp://guest:guest@localhost:5672"
)

type ServiceConfig struct {
	DatabaseType          clients.DBTYPE `json:"database_type"`
	DBConnection          string         `json:"db_connection"`
	RestfulEndpoint       string         `json:"restful_endpoint"`
	SecureRestfulEndpoint string         `json:"secure_restful_endpoint"`
	AMQPMessageBroker     string         `json:"amqp_message_broker"`
}

func ExtractConfiguration(fileName string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		SecureRestfulEPDefault,
		AMQPMessageBrokerDefault,
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values")
	} else {
		err = json.NewDecoder(file).Decode(&conf)
	}

	return conf, err
}
