package player

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/model"
)

var lg = logger.Logger.WithField("Module", "Player")

type IPlayer interface {
	Start()
	Stop()
	Play(media *model.Media) error
	IsPaused() bool
	Pause() error
	Unpause() error
	SetVolume(volume float64) error
	IsIdle() bool
	Seek(position float64, absolute bool) error
	ObserveProperty(property model.PlayerProperty, name string, handler event.HandlerFunc) error
	GetAudioDeviceList() ([]model.AudioDevice, error)
	SetAudioDevice(device string) error
}
