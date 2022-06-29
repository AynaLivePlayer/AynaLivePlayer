package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
)

func PlayNext() {
	l().Info("try to play next possible media")
	if UserPlaylist.Size() == 0 && SystemPlaylist.Size() == 0 {
		return
	}
	var media *player.Media
	if UserPlaylist.Size() != 0 {
		media = UserPlaylist.Pop()
	} else if SystemPlaylist.Size() != 0 {
		media = SystemPlaylist.Next()
	}
	Play(media)
}

func Play(media *player.Media) {
	l().Info("prepare media")
	err := PrepareMedia(media)
	if err != nil {
		l().Warn("prepare media failed. try play next")
		PlayNext()
		return
	}
	CurrentMedia = media
	if err := MainPlayer.Play(media); err != nil {
		l().Warn("play failed", err)
	}
	CurrentLyric.Reload(media.Lyric)
	// reset
	media.Url = ""
}

func Add(keyword string, user interface{}) {
	medias, err := Search(keyword)
	if err != nil {
		l().Warnf("search for %s, got error %s", keyword, err)
		return
	}
	if len(medias) == 0 {
		l().Info("search for %s, got no result", keyword)
		return
	}
	media := medias[0]
	media.User = user
	l().Infof("add media %s (%s)", media.Title, media.Artist)
	UserPlaylist.Insert(-1, media)
}

func AddWithProvider(keyword string, pname string, user interface{}) {
	medias, err := provider.Search(pname, keyword)
	if err != nil {
		l().Warnf("search for %s, got error %s", keyword, err)
		return
	}
	if len(medias) == 0 {
		l().Info("search for %s, got no result", keyword)
	}
	media := medias[0]
	media.User = user
	l().Info("add media %s (%s)", media.Title, media.Artist)
	UserPlaylist.Insert(-1, media)
}

func Seek(position float64, absolute bool) {
	if err := MainPlayer.Seek(position, absolute); err != nil {
		l().Warnf("seek to position %f (%t) failed, %s", position, absolute, err)
	}
}

func Toggle() (b bool) {
	var err error
	if MainPlayer.IsPaused() {
		err = MainPlayer.Unpause()
		b = false
	} else {
		err = MainPlayer.Pause()
		b = true
	}
	if err != nil {
		l().Warn("toggle failed", err)
	}
	return
}

func SetVolume(volume float64) {
	if MainPlayer.SetVolume(volume) != nil {
		l().Warnf("set mpv volume to %f failed", volume)
		return
	}
	config.Player.Volume = volume
}

func Destroy() {
	MainPlayer.Stop()
}

func GetAudioDevices() []player.AudioDevice {
	dl, err := MainPlayer.GetAudioDeviceList()
	if err != nil {
		return make([]player.AudioDevice, 0)
	}
	return dl
}

func SetAudioDevice(device string) {
	l().Infof("set audio device to %s", device)
	if err := MainPlayer.SetAudioDevice(device); err != nil {
		l().Warnf("set mpv audio device to %s failed, %s", device, err)
		MainPlayer.SetAudioDevice("auto")
		config.Player.AudioDevice = "auto"
		return
	}
	config.Player.AudioDevice = device
}
