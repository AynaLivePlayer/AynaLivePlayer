package wshub

import (
	"AynaLivePlayer/pkg/event"
	"encoding/json"
)

type EventData struct {
	EventID event.EventId
	Data    interface{}
}

type EventDataReceived struct {
	EventID event.EventId
	Data    json.RawMessage
}
