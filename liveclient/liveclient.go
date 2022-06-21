package liveclient

import "AynaLivePlayer/event"

const MODULE_NAME = "LiveClient"

type UserMedal struct {
	Name  string
	Level int
}

type DanmuUser struct {
	Uid       string
	Username  string
	Medal     UserMedal
	Admin     bool
	Privilege int
}

type DanmuMessage struct {
	User    DanmuUser
	Message string
}

type LiveClient interface {
	ClientName() string
	Connect() bool
	Disconnect() bool
	Handler() *event.Handler
}
