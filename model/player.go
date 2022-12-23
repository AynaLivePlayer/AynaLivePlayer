package model

type AudioDevice struct {
	Name        string
	Description string
}

type PlayerPropertyValue any
type PlayerProperty string

const (
	PlayerPropIdleActive PlayerProperty = "idle-active"
	PlayerPropTimePos    PlayerProperty = "time-pos"
	PlayerPropDuration   PlayerProperty = "duration"
	PlayerPropPercentPos PlayerProperty = "percent-pos"
	PlayerPropPause      PlayerProperty = "pause"
	PlayerPropVolume     PlayerProperty = "volume"
)
