package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/event"
	"AynaLivePlayer/player"
	"github.com/aynakeya/go-mpv"
)

func handleMpvIdlePlayNext(property *mpv.EventProperty) {
	isIdle := property.Data.(mpv.Node).Value.(bool)
	if isIdle {
		l.Info("mpv went idle, try play next")
		PlayNext()
	}
}

func handlePlaylistAdd(event *event.Event) {
	if MainPlayer.IsIdle() {
		PlayNext()
		return
	}
	if config.Player.SkipPlaylist && CurrentMedia != nil && CurrentMedia.User == player.PlaylistUser {
		PlayNext()
		return
	}
}

func handleLyricUpdate(property *mpv.EventProperty) {
	if property.Data == nil {
		return
	}
	t := property.Data.(mpv.Node).Value.(float64)
	CurrentLyric.Update(t)
}
