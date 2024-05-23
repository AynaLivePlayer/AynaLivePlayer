package events

const GUISetPlayerWindowOpenCmd = "cmd.gui.player_window.op"

type GUISetPlayerWindowOpenCmdEvent struct {
	SetOpen bool
}
