package liveclient

import (
	"AynaLivePlayer/common/event"
)

const (
	EventStatusChange   event.EventId = "liveclient.status.change"
	EventMessageReceive event.EventId = "liveclient.message.receive"
)

type StatusChangeEvent struct {
	Connected bool
	Client    LiveClient
}
