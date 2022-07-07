package provider

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/player"
	"os"
)

type _LocalPlaylist struct {
	Name   string
	Medias []*player.Media
}

type Local struct {
	Playlists []_LocalPlaylist
}

var LocalAPI *Local

func init() {
	LocalAPI = _newLocal()
	//Providers[LocalAPI.GetName()] = LocalAPI
}

func _newLocal() *Local {
	l := &Local{Playlists: make([]_LocalPlaylist, 0)}
	if err := os.MkdirAll(config.Provider.LocalDir, 0755); err != nil {
		return l
	}

	return l
}

func (l *Local) GetName() string {
	return "local"
}
func (l *Local) MatchMedia(keyword string) *player.Media {
	//TODO implement me
	panic("implement me")
}

func (l *Local) UpdateMediaLyric(media *player.Media) error {
	//TODO implement me
	panic("implement me")
}

func (l *Local) FormatPlaylistUrl(uri string) string {
	return ""
}

func (l *Local) GetPlaylist(playlist string) ([]*player.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Local) Search(keyword string) ([]*player.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Local) UpdateMedia(media *player.Media) error {
	//TODO implement me
	panic("implement me")
}

func (l *Local) UpdateMediaUrl(media *player.Media) error {
	//TODO implement me
	panic("implement me")
}
