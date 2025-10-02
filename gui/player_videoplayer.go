package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gutil"
	"fyne.io/fyne/v2"
)

func setupPlayerWindow() {
	playerWindow = App.NewWindow("CorePlayerPreview")
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
		logger.Infof("video output window handle: %d", playerWindowHandle)
		if playerWindowHandle != 0 {
			_ = global.EventBus.PublishToChannel(eventChannel, events.PlayerVideoPlayerSetWindowHandleCmd,
				events.PlayerVideoPlayerSetWindowHandleCmdEvent{Handle: playerWindowHandle})
		}
	}
}
