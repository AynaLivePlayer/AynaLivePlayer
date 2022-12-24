package gui

import (
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/resource"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

const MODULE_GUI = "GUI"

var App fyne.App
var MainWindow fyne.Window

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_GUI)
}

func black_magic() {
	widget.RichTextStyleStrong.TextStyle.Bold = false
}

func Initialize() {
	//black_magic()
	l().Info("Initializing GUI")
	//os.Setenv("FYNE_FONT", config.GetAssetPath("msyh.ttc"))
	App = app.New()
	App.Settings().SetTheme(&myTheme{})
	MainWindow = App.NewWindow(fmt.Sprintf("%s Ver.%s", config.ProgramName, config.Version))

	tabs := container.NewAppTabs(
		container.NewTabItem(i18n.T("gui.tab.player"),
			container.NewBorder(nil, createPlayControllerV2(), nil, nil, createPlaylist()),
		),
		container.NewTabItem(i18n.T("gui.tab.search"),
			container.NewBorder(createSearchBar(), nil, nil, nil, createSearchList()),
		),
		container.NewTabItem(i18n.T("gui.tab.room"),
			container.NewBorder(nil, nil, createRoomSelector(), nil, createRoomController()),
		),
		container.NewTabItem(i18n.T("gui.tab.playlist"),
			container.NewBorder(nil, nil, createPlaylists(), nil, createPlaylistMedias()),
		),
		container.NewTabItem(i18n.T("gui.tab.history"),
			container.NewBorder(nil, nil, nil, nil, createHistoryList()),
		),
		container.NewTabItem(i18n.T("gui.tab.config"),
			createConfigLayout(),
		),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	MainWindow.SetIcon(resource.ImageIcon)
	MainWindow.SetContent(tabs)
	//MainWindow.Resize(fyne.NewSize(1280, 720))
	MainWindow.Resize(fyne.NewSize(960, 480))
	//MainWindow.SetFixedSize(true)
}

func addShortCut() {
	key := &desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}
	MainWindow.Canvas().AddShortcut(key, func(shortcut fyne.Shortcut) {
		l().Info("Shortcut pressed: Ctrl+Shift+Right")
		controller.Instance.PlayControl().PlayNext()
	})
}
