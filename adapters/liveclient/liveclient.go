package liveclient

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
	"errors"
)

type LiveClientFactory struct {
	LiveClients  map[string]adapter.LiveClientCtor
	EventManager *event.Manager
	Logger       adapter.ILogger
}

func (f *LiveClientFactory) GetAllClientNames() []string {
	names := make([]string, 0)
	for key, _ := range f.LiveClients {
		names = append(names, key)
	}
	return names
}

func (f *LiveClientFactory) NewLiveClient(clientName, id string) (adapter.LiveClient, error) {
	ctor, ok := f.LiveClients[clientName]
	if !ok {
		return nil, errors.New("no such client")
	}
	return ctor(id, f.EventManager.NewChildManager(), f.Logger)
}
