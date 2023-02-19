package adapter

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/model"
)

type LiveClientCtor func(id string, ev *event.Manager, log ILogger) (LiveClient, error)

type LiveClient interface {
	ClientName() string
	RoomName() string
	Connect() bool
	Disconnect() bool
	Status() bool
	EventManager() *event.Manager
}

type ILiveRoom interface {
	Client() LiveClient
	Model() *model.LiveRoom // should return mutable model (not a copy)
	Identifier() string
	DisplayName() string
	Status() bool
	EventManager() *event.Manager
}

type LiveRoomExecutor interface {
	Match(command string) bool
	Execute(command string, args []string, danmu *model.DanmuMessage)
}
