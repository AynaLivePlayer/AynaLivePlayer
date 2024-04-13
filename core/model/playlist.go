package model

import "github.com/AynaLivePlayer/miaosic"

type PlaylistMode int

const (
	PlaylistModeNormal PlaylistMode = iota
	PlaylistModeRandom
	PlaylistModeRepeat
)

type PlaylistID string

const (
	PlaylistIDPlayer  PlaylistID = "player"
	PlaylistIDSystem  PlaylistID = "system"
	PlaylistIDHistory PlaylistID = "history"
)

type PlaylistInfo struct {
	Meta  miaosic.MetaData
	Title string
}

func (p PlaylistInfo) DisplayName() string {
	if p.Title != "" {
		return p.Title
	}
	return p.Meta.ID()
}
