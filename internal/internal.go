package internal

import (
	"AynaLivePlayer/internal/controller"
	"AynaLivePlayer/internal/liveroom"
	"AynaLivePlayer/internal/player"
	"AynaLivePlayer/internal/playlist"
)

func Initialize() {
	player.SetupMpvPlayer()
	playlist.Initialize()
	controller.Initialize()
	liveroom.Initialize()
}

func Stop() {
	liveroom.StopAndSave()
	player.StopMpvPlayer()
}
