package events

import (
	"AynaLivePlayer/core/model"
)

const PlaylistManagerSetSystemCmd = "cmd.playlist.manager.set.system"

type PlaylistManagerSetSystemCmdEvent struct {
	PlaylistID string
}

const PlaylistManagerSystemUpdate = "update.playlist.manager.system"

type PlaylistManagerSystemUpdateEvent struct {
	Info model.PlaylistInfo
}

const PlaylistManagerRefreshCurrentCmd = "cmd.playlist.manager.refresh.current"

type PlaylistManagerRefreshCurrentCmdEvent struct {
	PlaylistID string
}

const PlaylistManagerGetCurrentCmd = "cmd.playlist.manager.get.current"

type PlaylistManagerGetCurrentCmdEvent struct {
	PlaylistID string
}

const PlaylistManagerCurrentUpdate = "update.playlist.manager.current"

type PlaylistManagerCurrentUpdateEvent struct {
	Medias []model.Media
}

const PlaylistManagerInfoUpdate = "update.playlist.manager.info"

type PlaylistManagerInfoUpdateEvent struct {
	Playlists []model.PlaylistInfo
}

const PlaylistManagerAddPlaylistCmd = "cmd.playlist.manager.add"

type PlaylistManagerAddPlaylistCmdEvent struct {
	Provider string
	URL      string
}

const PlaylistManagerRemovePlaylistCmd = "cmd.playlist.manager.remove"

type PlaylistManagerRemovePlaylistCmdEvent struct {
	PlaylistID string
}
