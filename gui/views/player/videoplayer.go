package player

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	"fyne.io/fyne/v2"
)

var playerWindow fyne.Window
var playerWindowHandle uintptr

func setupPlayerWindow() {
	playerWindow = gctx.Context.App.NewWindow("CorePlayerPreview")
	playerWindow.Resize(fyne.NewSize(480, 240))
	playerWindow.SetCloseIntercept(func() {
		playerWindow.Hide()
	})
	playerWindow.Hide()
}

func showPlayerWindow() {
	if playerWindow == nil {
		setupPlayerWindow()
	}
	playerWindow.Show()
	if playerWindowHandle == 0 {
		playerWindowHandle = gutil.GetWindowHandle(playerWindow)
		gctx.Logger.Infof("video output window handle: %d", playerWindowHandle)
		if playerWindowHandle != 0 {
			_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerVideoPlayerSetWindowHandleCmd,
				events.PlayerVideoPlayerSetWindowHandleCmdEvent{Handle: playerWindowHandle})
		}
	}
}
