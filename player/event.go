package player

import (
	"AynaLivePlayer/event"
)

const (
	EventPlay              event.EventId = "player.play"
	EventPlaylistPreInsert event.EventId = "playlist.insert.pre"
	EventPlaylistInsert    event.EventId = "playlist.insert.after"
	EventPlaylistUpdate    event.EventId = "playlist.update"
	EventLyricUpdate       event.EventId = "lyric.update"
	EventLyricReload       event.EventId = "lyric.reload"
)

type PlaylistInsertEvent struct {
	Playlist *Playlist
	Index    int
	Media    *Media
}

type PlaylistUpdateEvent struct {
	Playlist *Playlist
}

func newPlaylistUpdateEvent(playlist *Playlist) PlaylistUpdateEvent {
	return PlaylistUpdateEvent{
		Playlist: playlist,
	}
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
