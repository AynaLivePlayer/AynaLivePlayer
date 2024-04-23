package events

import (
	"AynaLivePlayer/core/model"
)

const PlayerVolumeChangeCmd = "cmd.player.op.change_volume"

type PlayerVolumeChangeCmdEvent struct {
	Volume float64 // Volume from 0-100
}

const PlayerPlayCmd = "cmd.player.op.play"

type PlayerPlayCmdEvent struct {
	Media model.Media
}

const PlayerSeekCmd = "cmd.player.op.seek"

type PlayerSeekCmdEvent struct {
	Position float64
	// Absolute is the seek mode.
	// if absolute = true : position is the time in second
	// if absolute = false: position is in percentage eg 0.1 0.2
	Absolute bool
}

const PlayerToggleCmd = "cmd.player.op.toggle"

type PlayerToggleCmdEvent struct {
}

const PlayerPlayNextCmd = "cmd.player.op.next"

type PlayerPlayNextCmdEvent struct {
}
