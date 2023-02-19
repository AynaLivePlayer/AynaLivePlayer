package adapter

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/model"
)

// IControlBridge is the interface for all controller and
// all system use cases.
type IControlBridge interface {
	LiveRooms() ILiveRoomController
	PlayControl() IPlayController
	Playlists() IPlaylistController
	Provider() IProviderController
	Plugin() IPluginController
	LoadPlugins(plugins ...Plugin)
	Logger() ILogger
	CloseAndSave()
}

type ILiveRoomController interface {
	Size() int
	Get(index int) ILiveRoom
	GetRoomStatus(index int) bool
	Connect(index int) error
	Disconnect(index int) error
	AddRoom(clientName, roomId string) (*model.LiveRoom, error)
	DeleteRoom(index int) error
	AddDanmuCommand(executor LiveRoomExecutor)
	GetAllClientNames() []string
}

type IProviderController interface {
	GetPriority() []string
	PrepareMedia(media *model.Media) error
	MediaMatch(keyword string) *model.Media
	Search(keyword string) ([]*model.Media, error)
	SearchWithProvider(keyword string, provider string) ([]*model.Media, error)
	PreparePlaylist(playlist IPlaylist) error
}

type IPluginController interface {
	LoadPlugin(plugin Plugin)
	LoadPlugins(plugins ...Plugin)
	ClosePlugins()
}

type IPlaylistController interface {
	Size() int
	GetHistory() IPlaylist
	AddToHistory(media *model.Media)
	GetDefault() IPlaylist
	GetCurrent() IPlaylist
	Get(index int) IPlaylist
	Add(pname string, uri string) (IPlaylist, error)
	Remove(index int) (IPlaylist, error)
	SetDefault(index int) error
	PreparePlaylistByIndex(index int) error
}

type IPlayControlConfig struct {
	SkipPlaylist     bool
	AutoNextWhenFail bool
}

type IPlayController interface {
	EventManager() *event.Manager
	GetPlaying() *model.Media
	GetPlayer() IPlayer
	PlayNext()
	Play(media *model.Media) error
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
	Config() *IPlayControlConfig
}

type ILyricLoader interface {
	EventManager() *event.Manager
	Get() *model.Lyric
	Reload(lyric string)
	Update(time float64)
}
