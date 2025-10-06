package events

//const (
//	EventPlay              string = "player.play"
//	EventPlayed            string = "player.played"
//	EventPlaylistPreInsert string = "playlist.insert.pre"
//	EventPlaylistInsert    string = "playlist.insert.after"
//	EventPlaylistUpdate    string = "playlist.update"
//	EventLyricUpdate       string = "lyric.update"
//	EventLyricReload       string = "lyric.reload"
//)

const ErrorUpdate = "update.error"

type ErrorUpdateEvent struct {
	Error error
}

//
//func EventPlayerPropertyUpdate(property model.PlayerProperty) string {
//	return string("player.property.update." + string(property))
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
//type UpdateLiveRoomStatusData struct {
//	RoomTitle string
//	Status    bool
//}
