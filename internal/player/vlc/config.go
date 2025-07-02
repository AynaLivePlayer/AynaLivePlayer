package vlc

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
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
	global.EventManager.CallA(events.PlayerVolumeChangeCmd, events.PlayerVolumeChangeCmdEvent{
		Volume: cfg.Volume,
	})
	global.EventManager.RegisterA(events.PlayerPropertyVolumeUpdate, "player.config.volume", func(evnt *event.Event) {
		data := evnt.Data.(events.PlayerPropertyVolumeUpdateEvent)
		if data.Volume < 0 {
			return
		}
		cfg.Volume = data.Volume
	})
	global.EventManager.CallA(events.PlayerSetAudioDeviceCmd, events.PlayerSetAudioDeviceCmdEvent{
		Device: cfg.AudioDevice,
	})
}
