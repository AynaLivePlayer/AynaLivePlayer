package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/event"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
	"AynaLivePlayer/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

const MODULE_CONTROLLER = "Controller"

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_CONTROLLER)
}

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
	}
}

func Destroy() {
	MainPlayer.Stop()
}

func SetDanmuClient(roomId string) {
	ResetDanmuClient()
	l().Infof("setting live client for %s", roomId)
	room, err := strconv.Atoi(roomId)
	if err != nil {
		l().Warn("parse room id error", err)
		return
	}
	if !util.StringSliceContains(config.LiveRoom.History, roomId) {
		config.LiveRoom.History = append(config.LiveRoom.History, roomId)
	}
	LiveClient = liveclient.NewBilibili(room)
	LiveClient.Handler().Register(&event.EventHandler{
		EventId: liveclient.EventMessageReceive,
		Name:    "controller.commandexecutor",
		Handler: danmuCommandHandler,
	})
	LiveClient.Handler().RegisterA(
		liveclient.EventMessageReceive,
		"controller.danmu.handler",
		danmuHandler)
	l().Infof("setting live client for %s success", roomId)
}

func StartDanmuClient() {
	LiveClient.Connect()
}

func ResetDanmuClient() {
	if LiveClient != nil {
		l().Infof("disconnect from current live client %s", LiveClient.ClientName())
		LiveClient.Disconnect()
		LiveClient.Handler().UnregisterAll()
		LiveClient = nil
	}
}

func AddPlaylist(pname string, uri string) *player.Playlist {
	l().Infof("try add playlist %s with provider %s", uri, pname)
	id, err := provider.FormatPlaylistUrl(pname, uri)
	if err != nil || id == "" {
		l().Warnf("fail to format %s playlist id for %s", uri, pname)
		return nil
	}
	p := player.NewPlaylist(fmt.Sprintf("%s-%s", pname, id), player.PlaylistConfig{})
	p.Meta = provider.Meta{
		Name: pname,
		Id:   id,
	}
	PlaylistManager = append(PlaylistManager, p)
	config.Player.Playlists = append(config.Player.Playlists, id)
	config.Player.PlaylistsProvider = append(config.Player.PlaylistsProvider, pname)
	return p
}

func RemovePlaylist(index int) {
	l().Infof("Try to remove playlist.index=%d", index)
	if index < 0 || index >= len(PlaylistManager) {
		l().Warnf("playlist.index=%d not found", index)
		return
	}
	if index == config.Player.PlaylistIndex {
		l().Info("Delete current system playlist, reset system playlist to index = 0")
		SetSystemPlaylist(0)
	}
	if index < config.Player.PlaylistIndex {
		l().Debugf("Delete playlist before system playlist (index=%d), reduce system playlist index by 1", config.Player.PlaylistIndex)
		config.Player.PlaylistIndex = config.Player.PlaylistIndex - 1
	}
	PlaylistManager = append(PlaylistManager[:index], PlaylistManager[index+1:]...)
	config.Player.Playlists = append(config.Player.Playlists[:index], config.Player.Playlists[index+1:]...)
	config.Player.PlaylistsProvider = append(config.Player.PlaylistsProvider[:index], config.Player.PlaylistsProvider[index+1:]...)
}

func SetSystemPlaylist(index int) {
	l().Infof("try set system playlist to playlist.id=%d", index)
	if index < 0 || index >= len(PlaylistManager) {
		l().Warn("playlist.index=%d not found", index)
		return
	}
	err := PreparePlaylist(PlaylistManager[index])
	if err != nil {
		return
	}
	medias := PlaylistManager[index].Playlist
	config.Player.PlaylistIndex = index
	ApplyUser(medias, player.PlaylistUser)
	SystemPlaylist.Replace(medias)
}

func PreparePlaylistByIndex(index int) {
	l().Infof("try prepare playlist.id=%d", index)
	if index < 0 || index >= len(PlaylistManager) {
		l().Warn("playlist.id=%d not found", index)
		return
	}
	err := PreparePlaylist(PlaylistManager[index])
	if err != nil {
		return
	}
}
