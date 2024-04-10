package playlist

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
)

var PlayerPlaylist *playlist = nil
var HistoryPlaylist *playlist = nil
var SystemPlaylist *playlist = nil
var PlaylistsPlaylist *playlist = nil

type playlistConfig struct {
	SystemPlaylistMode model.PlaylistMode
}

func (p *playlistConfig) Name() string {
	return "playlist"
}

func (p *playlistConfig) OnLoad() {
	return
}

func (p *playlistConfig) OnSave() {
	return
}

var cfg = &playlistConfig{}

func Initialize() {
	PlayerPlaylist = newPlaylist(model.PlaylistIDPlayer)
	SystemPlaylist = newPlaylist(model.PlaylistIDSystem)
	HistoryPlaylist = newPlaylist(model.PlaylistIDHistory)
	PlaylistsPlaylist = newPlaylist(model.PlaylistIDPlaylists)
	config.LoadConfig(cfg)

	global.EventManager.RegisterA(events.PlaylistModeChangeCmd(model.PlaylistIDSystem), "internal.playlist.system_init", func(event *event.Event) {
		cfg.SystemPlaylistMode = event.Data.(events.PlaylistModeChangeUpdateEvent).Mode
	})

	global.EventManager.CallA(events.PlaylistModeChangeUpdate(model.PlaylistIDSystem), events.PlaylistModeChangeUpdateEvent{
		Mode: cfg.SystemPlaylistMode,
	})
}
