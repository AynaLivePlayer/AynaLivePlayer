package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
)

func registerHandlers() {
	global.EventManager.RegisterA(events.GUISetPlayerWindowOpenCmd, "gui.player.videoplayer.handleopen", func(event *event.Event) {
		data := event.Data.(events.GUISetPlayerWindowOpenCmdEvent)
		if data.SetOpen {
			playerWindow.Close()
		} else {
			showPlayerWindow()
		}
	})
}
