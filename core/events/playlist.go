package events

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/event"
)

func PlaylistDetailUpdate(id model.PlaylistID) event.EventId {
	return event.EventId("update.playlist.detail." + id)
}

type PlaylistDetailUpdateEvent struct {
	Medias []model.Media
}

func PlaylistMoveCmd(id model.PlaylistID) event.EventId {
	return event.EventId("cmd.playlist.move." + id)
}

type PlaylistMoveCmdEvent struct {
	From int
	To   int
}

func PlaylistSetIndexCmd(id model.PlaylistID) event.EventId {
	return event.EventId("cmd.playlist.setindex." + id)
}

type PlaylistSetIndexCmdEvent struct {
	Index int
}

func PlaylistDeleteCmd(id model.PlaylistID) event.EventId {
	return event.EventId("cmd.playlist.delete." + id)
}

type PlaylistDeleteCmdEvent struct {
	Index int
}

func PlaylistInsertCmd(id model.PlaylistID) event.EventId {
	return event.EventId("cmd.playlist.insert." + id)
}

type PlaylistInsertCmdEvent struct {
	Position int // position to insert, -1 means last one
	Media    model.Media
}

func PlaylistInsertUpdate(id model.PlaylistID) event.EventId {
	return event.EventId("update.playlist.insert." + id)
}

type PlaylistInsertUpdateEvent struct {
	Position int // position to insert, -1 means last one
	Media    model.Media
}

func PlaylistNextCmd(id model.PlaylistID) event.EventId {
	return event.EventId("cmd.playlist.next." + id)
}

type PlaylistNextCmdEvent struct {
	Remove bool // remove the media after next
}

func PlaylistNextUpdate(id model.PlaylistID) event.EventId {
	return event.EventId("update.playlist.next." + id)
}

type PlaylistNextUpdateEvent struct {
	Media model.Media
}

func PlaylistModeChangeCmd(id model.PlaylistID) event.EventId {
	return event.EventId("cmd.playlist.mode." + id)
}

type PlaylistModeChangeCmdEvent struct {
	Mode model.PlaylistMode
}

func PlaylistModeChangeUpdate(id model.PlaylistID) event.EventId {
	return event.EventId("update.playlist.mode." + id)
}

type PlaylistModeChangeUpdateEvent struct {
	Mode model.PlaylistMode
}
