package wshub

import (
	"encoding/json"
)

type EventData struct {
	EventID string
	Data    interface{}
}

type EventDataReceived struct {
	EventID string
	Data    json.RawMessage
}

var eventCache []*EventData

func init() {
	eventCache = make([]*EventData, 0)
}
