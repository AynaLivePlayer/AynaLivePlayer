package adapter

import "AynaLivePlayer/core/model"

type MediaProviderConfig map[string]string
type MediaProviderCtor func(config MediaProviderConfig) MediaProvider

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
