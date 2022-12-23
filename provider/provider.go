package provider

import (
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/model"
	"github.com/sirupsen/logrus"
)

const MODULE_CONTROLLER = "Provider"

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_CONTROLLER)
}

type MediaProvider interface {
	GetName() string
	MatchMedia(keyword string) *model.Media
	GetPlaylist(playlist *model.Meta) ([]*model.Media, error)
	FormatPlaylistUrl(uri string) string
	Search(keyword string) ([]*model.Media, error)
	UpdateMedia(media *model.Media) error
	UpdateMediaUrl(media *model.Media) error
	UpdateMediaLyric(media *model.Media) error
}

var Providers map[string]MediaProvider = make(map[string]MediaProvider)

func GetPlaylist(meta *model.Meta) ([]*model.Media, error) {
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

func MatchMedia(provider string, keyword string) *model.Media {
	if v, ok := Providers[provider]; ok {
		return v.MatchMedia(keyword)
	}
	return nil
}

func Search(provider string, keyword string) ([]*model.Media, error) {
	if v, ok := Providers[provider]; ok {
		return v.Search(keyword)
	}
	return nil, ErrorNoSuchProvider
}

func UpdateMedia(media *model.Media) error {
	if v, ok := Providers[media.Meta.(model.Meta).Name]; ok {
		return v.UpdateMedia(media)
	}
	return ErrorNoSuchProvider
}

func UpdateMediaUrl(media *model.Media) error {
	if v, ok := Providers[media.Meta.(model.Meta).Name]; ok {
		return v.UpdateMediaUrl(media)
	}
	return ErrorNoSuchProvider
}

func UpdateMediaLyric(media *model.Media) error {
	if v, ok := Providers[media.Meta.(model.Meta).Name]; ok {
		return v.UpdateMediaLyric(media)
	}
	return ErrorNoSuchProvider
}
