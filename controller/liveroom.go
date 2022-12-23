package controller

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/model"
)

type DanmuCommandExecutor interface {
	Match(command string) bool
	Execute(command string, args []string, danmu *liveclient.DanmuMessage)
}

type ILiveRoomController interface {
	Size() int
	Get(index int) ILiveRoom
	GetRoomStatus(index int) bool
	Connect(index int) error
	Disconnect(index int) error
	AddRoom(clientName, roomId string) (*model.LiveRoom, error)
	DeleteRoom(index int) error
	AddDanmuCommand(executor DanmuCommandExecutor)
}

type ILiveRoom interface {
	Model() *model.LiveRoom // should return mutable model (not a copy)
	Title() string          // should be same as Model().Title
	Status() bool
	EventManager() *event.Manager
}
