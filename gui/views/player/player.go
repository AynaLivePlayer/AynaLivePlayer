package player

import (
	"AynaLivePlayer/gui/gctx"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateView() fyne.CanvasObject {
	setupLyricViewer()
	registerHandlers()
	gctx.Context.OnMainWindowClosing(func() {
		if playerWindow != nil {
			gctx.Logger.Infof("closing player window")
			go playerWindow.Close()
		}
	})
	return container.NewBorder(nil, createPlayControllerV2(), nil, nil, createPlaylist())
}
