package events

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
)

const (
	LiveRoomStatusChange   event.EventId = "liveclient.status.change"
	LiveRoomMessageReceive event.EventId = "liveclient.message.receive"
)

type StatusChangeEvent struct {
	Connected bool
	Client    adapter.LiveClient
}
