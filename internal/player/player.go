package player

import "AynaLivePlayer/internal/player/mpv"

func SetupMpvPlayer() {
	mpv.SetupPlayer()
}

func StopMpvPlayer() {
	mpv.StopPlayer()
}
