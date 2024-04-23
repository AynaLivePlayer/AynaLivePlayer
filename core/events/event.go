package events

//const (
//	EventPlay              event.EventId = "player.play"
//	EventPlayed            event.EventId = "player.played"
//	EventPlaylistPreInsert event.EventId = "playlist.insert.pre"
//	EventPlaylistInsert    event.EventId = "playlist.insert.after"
//	EventPlaylistUpdate    event.EventId = "playlist.update"
//	EventLyricUpdate       event.EventId = "lyric.update"
//	EventLyricReload       event.EventId = "lyric.reload"
//)

const ErrorUpdate = "update.error"

type ErrorUpdateEvent struct {
	Error error
}

//
//func EventPlayerPropertyUpdate(property model.PlayerProperty) event.EventId {
//	return event.EventId("player.property.update." + string(property))
//}
//
//type PlaylistInsertEvent struct {
//	Playlist *model.Playlist
//	Index    int
//	Media    *model.Media
//}
//
//type PlaylistUpdateEvent struct {
//	Playlist *model.Playlist // Playlist is a copy of the playlist
//}
//
//type PlayEvent struct {
//	Media *model.Media
//}
//
//type LyricUpdateEvent struct {
//	Lyrics *model.Lyric
//	Time   float64
//	Lyric  *model.LyricContext
//}
//
//type LyricReloadEvent struct {
//	Lyrics *model.Lyric
//}
//
//type PlayerPropertyUpdateEvent struct {
//	Property model.PlayerProperty
//	Value    model.PlayerPropertyValue
//}
//
//type LiveRoomStatusUpdateEvent struct {
//	RoomTitle string
//	Status    bool
//}
