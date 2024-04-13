package gui

import (
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/resource"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func setupSysTray() {
	if desk, ok := App.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem(i18n.T("gui.tray.btn.show"), func() {
				MainWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(resource.ImageIcon)
	}
	MainWindow.SetCloseIntercept(func() {
		MainWindow.Hide()
	})
}
