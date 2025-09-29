package vlc

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
)

type playerConfig struct {
	Volume            float64
	AudioDevice       string
	DisplayMusicCover bool
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
	Volume:            100,
	DisplayMusicCover: true,
}

func restoreConfig() {
	_ = global.EventBus.Publish(events.PlayerVolumeChangeCmd, events.PlayerVolumeChangeCmdEvent{
		Volume: cfg.Volume,
	})
	global.EventBus.Subscribe("", events.PlayerPropertyVolumeUpdate, "player.config.volume", func(evnt *eventbus.Event) {
		data := evnt.Data.(events.PlayerPropertyVolumeUpdateEvent)
		if data.Volume < 0 {
			return
		}
		cfg.Volume = data.Volume
	})
	_ = global.EventBus.Publish(events.PlayerSetAudioDeviceCmd, events.PlayerSetAudioDeviceCmdEvent{
		Device: cfg.AudioDevice,
	})
}
