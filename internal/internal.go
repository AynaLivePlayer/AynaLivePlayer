package internal

import (
	"AynaLivePlayer/internal/controller"
	"AynaLivePlayer/internal/liveroom"
	"AynaLivePlayer/internal/player"
	"AynaLivePlayer/internal/playlist"
	"AynaLivePlayer/internal/source"
)

func Initialize() {
	player.SetupMpvPlayer()
	source.Initialize()
	playlist.Initialize()
	controller.Initialize()
	liveroom.Initialize()
}

func Stop() {
	liveroom.StopAndSave()
	playlist.Close()
	player.StopMpvPlayer()
}
