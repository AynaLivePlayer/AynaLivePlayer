package controller

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/model"
	"AynaLivePlayer/player"
)

type IPlayController interface {
	EventManager() *event.Manager
	GetPlaying() *model.Media
	GetPlayer() player.IPlayer
	PlayNext()
	Play(media *model.Media)
	Add(keyword string, user interface{})
	AddWithProvider(keyword string, provider string, user interface{})
	Seek(position float64, absolute bool)
	Toggle() bool
	SetVolume(volume float64)
	Destroy()
	GetCurrentAudioDevice() string
	GetAudioDevices() []model.AudioDevice
	SetAudioDevice(device string)
	GetLyric() ILyricLoader
	GetSkipPlaylist() bool
	SetSkipPlaylist(b bool)
}

type ILyricLoader interface {
	EventManager() *event.Manager
	Get() *model.Lyric
	Reload(lyric string)
	Update(time float64)
}
