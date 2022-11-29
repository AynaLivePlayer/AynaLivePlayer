package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
	"fmt"
)

var UserPlaylist *player.Playlist
var History *player.Playlist
var HistoryUser *player.User
var SystemPlaylist *player.Playlist
var PlaylistManager []*player.Playlist

func AddToHistory(media *player.Media) {
	l.Tracef("add media %s (%s) to history", media.Title, media.Artist)
	media = media.Copy()
	// reset url for future use
	media.Url = ""
	if History.Size() >= 1024 {
		History.Replace([]*player.Media{})
	}
	History.Push(media)
	return
}

func ToHistoryMedia(media *player.Media) *player.Media {
	media = media.Copy()
	media.User = HistoryUser
	return media
}

func ToSystemMedia(media *player.Media) *player.Media {
	media = media.Copy()
	media.User = player.SystemUser
	return media
}

func AddPlaylist(pname string, uri string) *player.Playlist {
	l.Infof("try add playlist %s with provider %s", uri, pname)
	id, err := provider.FormatPlaylistUrl(pname, uri)
	if err != nil || id == "" {
		l.Warnf("fail to format %s playlist id for %s", uri, pname)
		return nil
	}
	p := player.NewPlaylist(fmt.Sprintf("%s-%s", pname, id), player.PlaylistConfig{})
	p.Meta = provider.Meta{
		Name: pname,
		Id:   id,
	}
	PlaylistManager = append(PlaylistManager, p)
	config.Player.Playlists = append(config.Player.Playlists, &config.PlayerPlaylist{
		ID:       uri,
		Provider: pname,
	})
	return p
}

func RemovePlaylist(index int) {
	l.Infof("Try to remove playlist.index=%d", index)
	if index < 0 || index >= len(PlaylistManager) {
		l.Warnf("playlist.index=%d not found", index)
		return
	}
	if index == config.Player.PlaylistIndex {
		l.Info("Delete current system playlist, reset system playlist to index = 0")
		SetSystemPlaylist(0)
	}
	if index < config.Player.PlaylistIndex {
		l.Debugf("Delete playlist before system playlist (index=%d), reduce system playlist index by 1", config.Player.PlaylistIndex)
		config.Player.PlaylistIndex = config.Player.PlaylistIndex - 1
	}
	PlaylistManager = append(PlaylistManager[:index], PlaylistManager[index+1:]...)
	config.Player.Playlists = append(config.Player.Playlists[:index], config.Player.Playlists[index+1:]...)
}

func SetSystemPlaylist(index int) {
	l.Infof("try set system playlist to playlist.id=%d", index)
	if index < 0 || index >= len(PlaylistManager) {
		l.Warn("playlist.index=%d not found", index)
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
	l.Infof("try prepare playlist.id=%d", index)
	if index < 0 || index >= len(PlaylistManager) {
		l.Warn("playlist.id=%d not found", index)
		return
	}
	err := PreparePlaylist(PlaylistManager[index])
	if err != nil {
		return
	}
}
