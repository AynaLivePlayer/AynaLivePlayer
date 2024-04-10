package mpv

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
)

type playerConfig struct {
	Volume float64
}

func (p *playerConfig) Name() string {
	return "Player"
}

func (p *playerConfig) OnLoad() {
	return
}

func (p *playerConfig) OnSave() {
	return
}

var cfg = &playerConfig{
	Volume: 100,
}

func restoreConfig() {
	global.EventManager.CallA(events.PlayerVolumeChangeCmd, events.PlayerVolumeChangeCmdEvent{
		Volume: cfg.Volume,
	})
	global.EventManager.RegisterA(events.PlayerPropertyVolumeUpdate, "player.config.volume", func(evnt *event.Event) {
		data := evnt.Data.(events.PlayerPropertyVolumeUpdateEvent)
		cfg.Volume = data.Volume
	})
}
