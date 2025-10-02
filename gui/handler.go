package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
)

func registerHandlers() {
	global.EventBus.Subscribe(eventChannel,  events.GUISetPlayerWindowOpenCmd, "gui.player.videoplayer.handleopen", func(event *eventbus.Event) {
		data := event.Data.(events.GUISetPlayerWindowOpenCmdEvent)
		if data.SetOpen {
			playerWindow.Close()
		} else {
			showPlayerWindow()
		}
	})
}
