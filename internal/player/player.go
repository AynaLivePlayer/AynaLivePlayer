//go:build !mpvOnly && !vlcOnly

package player

import (
	"AynaLivePlayer/internal/player/mpv"
	"AynaLivePlayer/internal/player/vlc"
	"AynaLivePlayer/pkg/config"
)

func SetupPlayer() {
	if config.Experimental.PlayerCore == "vlc" {
		vlc.SetupPlayer()
	} else {
		mpv.SetupPlayer()
	}
}

func StopPlayer() {
	if config.Experimental.PlayerCore == "vlc" {
		//vlc.StopPlayer()
	} else {
		mpv.StopPlayer()
	}
}
