package adapter

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/model"
)

type IPlaylist interface {
	Identifier() string     // must unique for each playlist
	Model() *model.Playlist // mutable model (not a copy)
	EventManager() *event.Manager
	DisplayName() string
	Size() int
	Get(index int) *model.Media
	Pop() *model.Media
	Replace(medias []*model.Media)
	Push(media *model.Media)
	Insert(index int, media *model.Media)
	Delete(index int) *model.Media
	Move(src int, dst int)
	Next() *model.Media
}
