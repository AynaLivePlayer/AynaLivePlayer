//go:build mpvOnly

package player

import (
	"AynaLivePlayer/internal/player/mpv"
)

func SetupPlayer() {
	mpv.SetupPlayer()
}

func StopPlayer() {
	mpv.StopPlayer()
}
