package kafka

type message struct {
	EventName string      `json:"event_name"`
	Payload   interface{} `json:"payload"`
}
