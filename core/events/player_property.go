package events

import "AynaLivePlayer/core/model"

const PlayerPlayingUpdate = "update.player.playing"

type PlayerPlayingUpdateEvent struct {
	Media   model.Media
	Removed bool // if no media is playing, removed is true
}

const PlayerPropertyPauseUpdate = "update.player.property.pause"

type PlayerPropertyPauseUpdateEvent struct {
	Paused bool
}

const PlayerPropertyPercentPosUpdate = "update.player.property.percent_pos"

type PlayerPropertyPercentPosUpdateEvent struct {
	PercentPos float64
}

const PlayerPropertyIdleActiveUpdate = "update.player.property.idle_active"

type PlayerPropertyIdleActiveUpdateEvent struct {
	IsIdle bool
}

const PlayerPropertyTimePosUpdate = "update.player.property.time_pos"

type PlayerPropertyTimePosUpdateEvent struct {
	TimePos float64 // Time in seconds
}

const PlayerPropertyDurationUpdate = "update.player.property.duration"

type PlayerPropertyDurationUpdateEvent struct {
	Duration float64 // Duration in seconds
}

const PlayerPropertyVolumeUpdate = "update.player.property.volume"

type PlayerPropertyVolumeUpdateEvent struct {
	Volume float64 // Volume from 0-100
}
