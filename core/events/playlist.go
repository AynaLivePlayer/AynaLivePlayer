package events

import (
	"AynaLivePlayer/core/model"
)

func PlaylistDetailUpdate(id model.PlaylistID) string {
	return string("update.playlist.detail." + id)
}

type PlaylistDetailUpdateEvent struct {
	Medias []model.Media
}

func PlaylistMoveCmd(id model.PlaylistID) string {
	return string("cmd.playlist.move." + id)
}

type PlaylistMoveCmdEvent struct {
	From int
	To   int
}

func PlaylistSetIndexCmd(id model.PlaylistID) string {
	return string("cmd.playlist.setindex." + id)
}

type PlaylistSetIndexCmdEvent struct {
	Index int
}

func PlaylistDeleteCmd(id model.PlaylistID) string {
	return string("cmd.playlist.delete." + id)
}

type PlaylistDeleteCmdEvent struct {
	Index int
}

func PlaylistInsertCmd(id model.PlaylistID) string {
	return string("cmd.playlist.insert." + id)
}

type PlaylistInsertCmdEvent struct {
	Position int // position to insert, -1 means last one
	Media    model.Media
}

func PlaylistInsertUpdate(id model.PlaylistID) string {
	return string("update.playlist.insert." + id)
}

type PlaylistInsertUpdateEvent struct {
	Position int // position to insert, -1 means last one
	Media    model.Media
}

func PlaylistNextCmd(id model.PlaylistID) string {
	return string("cmd.playlist.next." + id)
}

type PlaylistNextCmdEvent struct {
	Remove bool // remove the media after next
}

func PlaylistNextUpdate(id model.PlaylistID) string {
	return string("update.playlist.next." + id)
}

type PlaylistNextUpdateEvent struct {
	Media model.Media
}

func PlaylistModeChangeCmd(id model.PlaylistID) string {
	return string("cmd.playlist.mode." + id)
}

type PlaylistModeChangeCmdEvent struct {
	Mode model.PlaylistMode
}

func PlaylistModeChangeUpdate(id model.PlaylistID) string {
	return string("update.playlist.mode." + id)
}

type PlaylistModeChangeUpdateEvent struct {
	Mode model.PlaylistMode
}
