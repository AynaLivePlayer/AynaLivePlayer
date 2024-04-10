package adapter

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/event"
)

type PlayerCtor func(ev *event.Manager, log ILogger) IPlayer

type IPlayer interface {
	// Start the player
	Start()
	// Stop the player
	Stop()
	// Play play a media
	Play(media *model.Media) error
	// GetPlaying get playing media
	// if player is idle, return nil
	GetPlaying() *model.Media
	// IsPaused return true if player is paused
	IsPaused() bool
	// Pause pause player
	Pause() error
	// Unpause unpause player
	Unpause() error
	// SetVolume set volume
	SetVolume(volume float64) error
	// IsIdle return true if player is playing anything
	IsIdle() bool
	// Seek to position, if absolute is true, position is absolute time, otherwise position is relative time
	Seek(position float64, absolute bool) error
	// SetWindowHandle set window handle for video output
	SetWindowHandle(handle uintptr) error
	// ObserveProperty observe player property change
	ObserveProperty(property model.PlayerProperty, name string, handler event.HandlerFunc) error
	// GetAudioDeviceList get audio device list
	GetAudioDeviceList() ([]model.AudioDevice, error)
	// SetAudioDevice set audio device
	SetAudioDevice(device string) error
}
