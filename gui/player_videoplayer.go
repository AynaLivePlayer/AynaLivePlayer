package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/xfyne"
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
	playerWindow.Show()
	if playerWindowHandle == 0 {
		playerWindowHandle = xfyne.GetWindowHandle(playerWindow)
		logger.Infof("video output window handle: %d", playerWindowHandle)
		if playerWindowHandle != 0 {
			global.EventManager.CallA(events.PlayerVideoPlayerSetWindowHandleCmd,
				events.PlayerVideoPlayerSetWindowHandleCmdEvent{Handle: playerWindowHandle})
		}
	}
}
