package player

import (
	"AynaLivePlayer/internal/player/mpv"
	"AynaLivePlayer/internal/player/vlc"
	"AynaLivePlayer/pkg/config"
)

func SetupMpvPlayer() {
	if config.Experimental.PlayerCore == "vlc" {
		vlc.SetupPlayer()
	} else {
		mpv.SetupPlayer()
	}

}

func StopMpvPlayer() {
	if config.Experimental.PlayerCore == "vlc" {
		vlc.StopPlayer()
	} else {
		mpv.StopPlayer()
	}
}
