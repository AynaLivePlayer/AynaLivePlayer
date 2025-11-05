package events

import (
	"AynaLivePlayer/core/model"
	liveroomsdk "github.com/AynaLivePlayer/liveroom-sdk"
)

const CmdLiveRoomAdd = "cmd.liveroom.add"

type CmdLiveRoomAddData struct {
	Title    string
	Provider string
	RoomKey  string
}

const CmdLiveRoomRemove = "cmd.liveroom.remove"

type CmdLiveRoomRemoveData struct {
	Identifier string
}

const CmdLiveRoomConfigChange = "cmd.liveroom.config.change"

type CmdLiveRoomConfigChangeData struct {
	Identifier string
	Config     model.LiveRoomConfig
}

const LiveRoomProviderUpdate = "update.liveroom.provider"

type LiveRoomProviderUpdateEvent struct {
	Providers []model.LiveRoomProviderInfo
}

const UpdateLiveRoomRooms = "update.liveroom.rooms"

type UpdateLiveRoomRoomsData struct {
	Rooms []model.LiveRoom
}

const UpdateLiveRoomStatus = "update.liveroom.status"

type UpdateLiveRoomStatusData struct {
	Room model.LiveRoom
}

const CmdLiveRoomOperation = "cmd.liveroom.operation"

type CmdLiveRoomOperationData struct {
	Identifier string
	SetConnect bool // connect or disconnect
}

const ReplyLiveRoomOperation = "reply.liveroom.operation"

type ReplyLiveRoomOperationData struct {
	Err error
}

const LiveRoomMessageReceive = "update.liveroom.message"

type LiveRoomMessageReceiveEvent struct {
	Message *liveroomsdk.Message
}
