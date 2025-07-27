//go:build vlcOnly

package player

import (
	"AynaLivePlayer/internal/player/vlc"
)

func SetupPlayer() {
	vlc.SetupPlayer()
}

func StopPlayer() {
	vlc.StopPlayer()
}
