package controller

import (
	"AynaLivePlayer/event"
	"github.com/aynakeya/go-mpv"
)

func handleMpvIdlePlayNext(property *mpv.EventProperty) {
	isIdle := property.Data.(mpv.Node).Value.(bool)
	if isIdle {
		l().Info("mpv went idle, try play next")
		PlayNext()
	}
}

func handlePlaylistAdd(event *event.Event) {
	if MainPlayer.IsIdle() {
		PlayNext()
	}
}

func handleLyricUpdate(property *mpv.EventProperty) {
	if property.Data == nil {
		return
	}
	t := property.Data.(mpv.Node).Value.(float64)
	CurrentLyric.Update(t)
}
