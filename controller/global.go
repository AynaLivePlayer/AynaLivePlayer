package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
	"fmt"
)

var MainPlayer *player.Player
var UserPlaylist *player.Playlist
var SystemPlaylist *player.Playlist
var LiveClient liveclient.LiveClient
var PlaylistManager []*player.Playlist
var CurrentLyric *player.Lyric
var CurrentMedia *player.Media

func Initialize() {

	MainPlayer = player.NewPlayer()
	SetAudioDevice(config.Player.AudioDevice)
	SetVolume(config.Player.Volume)
	UserPlaylist = player.NewPlaylist("user", player.PlaylistConfig{RandomNext: false})
	SystemPlaylist = player.NewPlaylist("system", player.PlaylistConfig{RandomNext: config.Player.PlaylistRandom})
	PlaylistManager = make([]*player.Playlist, 0)
	CurrentLyric = player.NewLyric("")
	loadPlaylists()

	MainPlayer.ObserveProperty("idle-active", handleMpvIdlePlayNext)
	UserPlaylist.Handler.RegisterA(player.EventPlaylistInsert, "controller.playnextwhenadd", handlePlaylistAdd)
	MainPlayer.ObserveProperty("time-pos", handleLyricUpdate)
	MainPlayer.Start()

}

func loadPlaylists() {
	l().Info("Loading playlists ", config.Player.Playlists, config.Player.PlaylistsProvider)
	if len(config.Player.Playlists) != len(config.Player.Playlists) {
		l().Warn("playlist id and provider does not have same length")
		return
	}
	for i := 0; i < len(config.Player.Playlists); i++ {
		pname := config.Player.PlaylistsProvider[i]
		id := config.Player.Playlists[i]
		p := player.NewPlaylist(fmt.Sprintf("%s-%s", pname, id), player.PlaylistConfig{})
		p.Meta = provider.Meta{
			Name: pname,
			Id:   id,
		}
		PlaylistManager = append(PlaylistManager, p)
	}
	if config.Player.PlaylistIndex < 0 || config.Player.PlaylistIndex >= len(config.Player.Playlists) {
		l().Warn("playlist index did not find")
		return
	}
	go func() {
		c := config.Player.PlaylistIndex
		err := PreparePlaylist(PlaylistManager[c])
		if err != nil {
			return
		}
		SetSystemPlaylist(c)
	}()
}
