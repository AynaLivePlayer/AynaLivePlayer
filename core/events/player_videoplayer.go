package events

import "AynaLivePlayer/core/model"

const PlayerVideoPlayerSetWindowHandleCmd = "cmd.player.videoplayer.set_window_handle"

type PlayerVideoPlayerSetWindowHandleCmdEvent struct {
	Handle uintptr
}

const PlayerSetAudioDeviceCmd = "cmd.player.set_audio_device"

type PlayerSetAudioDeviceCmdEvent struct {
	Device string
}

const PlayerAudioDeviceUpdate = "update.player.audio_device"

type PlayerAudioDeviceUpdateEvent struct {
	Current string
	Devices []model.AudioDevice
}
