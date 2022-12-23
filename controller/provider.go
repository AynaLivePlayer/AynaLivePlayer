package controller

import (
	"AynaLivePlayer/model"
)

var PlaylistUser = &model.User{Name: "Playlists"}
var SystemUser = &model.User{Name: "System"}
var HistoryUser = &model.User{Name: "History"}

type IProviderController interface {
	GetPriority() []string
	PrepareMedia(media *model.Media) error
	MediaMatch(keyword string) *model.Media
	Search(keyword string) ([]*model.Media, error)
	SearchWithProvider(keyword string, provider string) ([]*model.Media, error)
	PreparePlaylist(playlist IPlaylist) error
}

func ApplyUser(medias []*model.Media, user interface{}) {
	for _, m := range medias {
		m.User = user
	}
}

func ToSpMedia(media *model.Media, user *model.User) *model.Media {
	media = media.Copy()
	media.User = user
	return media
}
