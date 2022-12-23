package controller

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/model"
)

type IPlaylistController interface {
	Size() int
	GetHistory() IPlaylist
	AddToHistory(media *model.Media)
	GetDefault() IPlaylist
	GetCurrent() IPlaylist
	Get(index int) IPlaylist
	Add(pname string, uri string) IPlaylist
	Remove(index int) IPlaylist
	SetDefault(index int) error
	PreparePlaylistByIndex(index int) error
}

type IPlaylist interface {
	Model() *model.Playlist // mutable model (not a copy)
	EventManager() *event.Manager
	Name() string
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
