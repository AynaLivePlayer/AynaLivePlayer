package events

const PlayerVideoPlayerSetWindowHandleCmd = "cmd.player.videoplayer.set_window_handle"

type PlayerVideoPlayerSetWindowHandleCmdEvent struct {
	Handle uintptr
}
