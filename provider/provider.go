package provider

import (
	"AynaLivePlayer/logger"
	"AynaLivePlayer/player"
	"github.com/sirupsen/logrus"
)

const MODULE_CONTROLLER = "Provider"

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_CONTROLLER)
}

type Meta struct {
	Name string
	Id   string
}

type MediaProvider interface {
	GetName() string
	MatchMedia(keyword string) *player.Media
	GetPlaylist(playlist Meta) ([]*player.Media, error)
	FormatPlaylistUrl(uri string) string
	Search(keyword string) ([]*player.Media, error)
	UpdateMedia(media *player.Media) error
	UpdateMediaUrl(media *player.Media) error
	UpdateMediaLyric(media *player.Media) error
}

var Providers map[string]MediaProvider = make(map[string]MediaProvider)

func GetPlaylist(meta Meta) ([]*player.Media, error) {
	if v, ok := Providers[meta.Name]; ok {
		return v.GetPlaylist(meta)
	}
	return nil, ErrorNoSuchProvider
}

func FormatPlaylistUrl(pname, uri string) (string, error) {
	if v, ok := Providers[pname]; ok {
		return v.FormatPlaylistUrl(uri), nil
	}
	return "", ErrorNoSuchProvider
}

func Search(provider string, keyword string) ([]*player.Media, error) {
	if v, ok := Providers[provider]; ok {
		return v.Search(keyword)
	}
	return nil, ErrorNoSuchProvider
}

func UpdateMedia(media *player.Media) error {
	if v, ok := Providers[media.Meta.(Meta).Name]; ok {
		return v.UpdateMedia(media)
	}
	return ErrorNoSuchProvider
}

func UpdateMediaUrl(media *player.Media) error {
	if v, ok := Providers[media.Meta.(Meta).Name]; ok {
		return v.UpdateMediaUrl(media)
	}
	return ErrorNoSuchProvider
}

func UpdateMediaLyric(media *player.Media) error {
	if v, ok := Providers[media.Meta.(Meta).Name]; ok {
		return v.UpdateMediaLyric(media)
	}
	return ErrorNoSuchProvider
}
