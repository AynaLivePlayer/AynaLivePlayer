package player

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
)

type PlayerFactory struct {
	EventManager *event.Manager
	Logger       adapter.ILogger
}

func (f *PlayerFactory) NewMPV() adapter.IPlayer {
	return NewMpvPlayer(f.EventManager, f.Logger)
}
