package systray

import (
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/resource"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func SetupSysTray() {
	if desk, ok := gctx.Context.App.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem(i18n.T("gui.tray.btn.show"), func() {
				gctx.Context.Window.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(resource.ImageIcon)
	}
	gctx.Context.Window.SetCloseIntercept(func() {
		gctx.Context.Window.Hide()
	})
}
