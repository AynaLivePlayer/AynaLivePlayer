package model

type PlaylistMode int

const (
	PlaylistModeNormal PlaylistMode = iota
	PlaylistModeRandom
	PlaylistModeRepeat
)

type PlaylistID string

const (
	PlaylistIDPlayer    PlaylistID = "player"
	PlaylistIDSystem    PlaylistID = "system"
	PlaylistIDHistory   PlaylistID = "history"
	PlaylistIDPlaylists PlaylistID = "playlists"
)
