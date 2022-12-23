package model

import (
	"AynaLivePlayer/common/event"
)

const (
	EventPlay              event.EventId = "player.play"
	EventPlaylistPreInsert event.EventId = "playlist.insert.pre"
	EventPlaylistInsert    event.EventId = "playlist.insert.after"
	EventPlaylistUpdate    event.EventId = "playlist.update"
	EventLyricUpdate       event.EventId = "lyric.update"
	EventLyricReload       event.EventId = "lyric.reload"
)

func EventPlayerPropertyUpdate(property PlayerProperty) event.EventId {
	return event.EventId("player.property.update." + string(property))
}

type PlaylistInsertEvent struct {
	Playlist *Playlist
	Index    int
	Media    *Media
}

type PlaylistUpdateEvent struct {
	Playlist *Playlist // Playlist is a copy of the playlist
}

type PlayEvent struct {
	Media *Media
}

type LyricUpdateEvent struct {
	Lyrics *Lyric
	Time   float64
	Lyric  *LyricLine
}

type LyricReloadEvent struct {
	Lyrics *Lyric
}

type PlayerPropertyUpdateEvent struct {
	Property PlayerProperty
	Value    PlayerPropertyValue
}

type LiveRoomStatusUpdateEvent struct {
	RoomTitle string
	Status    bool
}
