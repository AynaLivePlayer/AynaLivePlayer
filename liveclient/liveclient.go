package liveclient

import (
	"AynaLivePlayer/common/event"
	"errors"
)

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
	RoomName() string
	Connect() bool
	Disconnect() bool
	Status() bool
	EventManager() *event.Manager
}

type LiveClientCtor func(id string) (LiveClient, error)

var LiveClients map[string]LiveClientCtor = map[string]LiveClientCtor{}

func GetAllClientNames() []string {
	names := make([]string, 0)
	for key, _ := range LiveClients {
		names = append(names, key)
	}
	return names
}

func NewLiveClient(clientName, id string) (LiveClient, error) {
	ctor, ok := LiveClients[clientName]
	if !ok {
		return nil, errors.New("no such client")
	}
	return ctor(id)
}
