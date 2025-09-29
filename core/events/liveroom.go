package events

import (
	"AynaLivePlayer/core/model"
	liveroomsdk "github.com/AynaLivePlayer/liveroom-sdk"
)

//const (
//	LiveRoomStatusChange   string = "liveclient.status.change"
//	LiveRoomMessageReceive string = "liveclient.message.receive"
//)
//
//type StatusChangeEvent struct {
//	Connected bool
//	Client    adapter.LiveClient
//}

const LiveRoomAddCmd = "cmd.liveroom.add"

type LiveRoomAddCmdEvent struct {
	Title    string
	Provider string
	RoomKey  string
}

const LiveRoomProviderUpdate = "update.liveroom.provider"

type LiveRoomProviderUpdateEvent struct {
	Providers []model.LiveRoomProviderInfo
}

const LiveRoomRemoveCmd = "cmd.liveroom.remove"

type LiveRoomRemoveCmdEvent struct {
	Identifier string
}

const LiveRoomRoomsUpdate = "update.liveroom.rooms"

type LiveRoomRoomsUpdateEvent struct {
	Rooms []model.LiveRoom
}

const LiveRoomStatusUpdate = "update.liveroom.status"

type LiveRoomStatusUpdateEvent struct {
	Room model.LiveRoom
}

const LiveRoomConfigChangeCmd = "cmd.liveroom.config.change"

type LiveRoomConfigChangeCmdEvent struct {
	Identifier string
	Config     model.LiveRoomConfig
}

const LiveRoomOperationCmd = "cmd.liveroom.operation"

type LiveRoomOperationCmdEvent struct {
	Identifier string
	SetConnect bool // connect or disconnect
}

const LiveRoomOperationFinish = "update.liveroom.operation"

type LiveRoomOperationFinishEvent struct {
}

const LiveRoomMessageReceive = "update.liveroom.message"

type LiveRoomMessageReceiveEvent struct {
	Message *liveroomsdk.Message
}
