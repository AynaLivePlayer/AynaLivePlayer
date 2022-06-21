package liveclient

import (
	"AynaLivePlayer/event"
)

const (
	EventStatusChange   event.EventId = "liveclient.status.change"
	EventMessageReceive event.EventId = "liveclient.message.receive"
)

type StatusChangeEvent struct {
	Connected bool
	Client    LiveClient
}
