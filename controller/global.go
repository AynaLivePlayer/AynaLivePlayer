package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
	"fmt"
)

var MainPlayer *player.Player
var LiveClient liveclient.LiveClient
var CurrentLyric *player.Lyric
var CurrentMedia *player.Media

func Initialize() {
	MainPlayer = player.NewPlayer()

	SetAudioDevice(config.Player.AudioDevice)
	SetVolume(config.Player.Volume)

	UserPlaylist = player.NewPlaylist("user", player.PlaylistConfig{RandomNext: config.Player.UserPlaylistRandom})
	SystemPlaylist = player.NewPlaylist("system", player.PlaylistConfig{RandomNext: config.Player.PlaylistRandom})
	PlaylistManager = make([]*player.Playlist, 0)
	History = player.NewPlaylist("history", player.PlaylistConfig{RandomNext: false})
	HistoryUser = &player.User{Name: "History"}
	loadPlaylists()

	config.LoadConfig(LiveRoomManager)
	LiveRoomManager.InitializeRooms()

	CurrentLyric = player.NewLyric("")

	MainPlayer.ObserveProperty("idle-active", handleMpvIdlePlayNext)
	UserPlaylist.Handler.RegisterA(player.EventPlaylistInsert, "controller.playnextwhenadd", handlePlaylistAdd)
	MainPlayer.ObserveProperty("time-pos", handleLyricUpdate)
	MainPlayer.Start()

}

func CloseAndSave() {
	// set value to global config
	config.Player.PlaylistRandom = SystemPlaylist.Config.RandomNext
	config.Player.UserPlaylistRandom = UserPlaylist.Config.RandomNext
	_ = config.SaveToConfigFile(config.ConfigPath)
}

func loadPlaylists() {
	l.Info("Loading playlists ", config.Player.Playlists)
	if len(config.Player.Playlists) != len(config.Player.Playlists) {
		l.Warn("playlist id and provider does not have same length")
		return
	}
	for i := 0; i < len(config.Player.Playlists); i++ {
		pc := config.Player.Playlists[i]
		p := player.NewPlaylist(fmt.Sprintf("%s-%s", pc.Provider, pc.ID), player.PlaylistConfig{})
		p.Meta = provider.Meta{
			Name: pc.Provider,
			Id:   pc.ID,
		}
		PlaylistManager = append(PlaylistManager, p)
	}
	if config.Player.PlaylistIndex < 0 || config.Player.PlaylistIndex >= len(config.Player.Playlists) {
		config.Player.PlaylistIndex = 0
		l.Warn("playlist index did not find")
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
