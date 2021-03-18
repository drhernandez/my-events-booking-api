package handlers

import "net/http"

type HealthHandler struct {}

func (h HealthHandler) PingHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("pong"))
}