package player

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/pkg/eventbus"
)

func registerHandlers() {
	global.EventBus.Subscribe(gctx.EventChannel, events.GUISetPlayerWindowOpenCmd, "gui.player.videoplayer.handleopen", func(event *eventbus.Event) {
		data := event.Data.(events.GUISetPlayerWindowOpenCmdEvent)
		if data.SetOpen {
			playerWindow.Close()
		} else {
			showPlayerWindow()
		}
	})
}
