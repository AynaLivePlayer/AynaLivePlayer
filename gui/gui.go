package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/resource"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/sirupsen/logrus"
	"os"
)

const MODULE_GUI = "GUI"

type ConfigLayout interface {
	Title() string
	Description() string
	CreatePanel() fyne.CanvasObject
}

var App fyne.App
var MainWindow fyne.Window
var ConfigList = []ConfigLayout{&bascicConfig{}}

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_GUI)
}

func Initialize() {
	l().Info("Initializing GUI")
	os.Setenv("FYNE_FONT", config.GetAssetPath("msyh.ttc"))
	App = app.New()
	MainWindow = App.NewWindow(fmt.Sprintf("%s Ver.%s", config.ProgramName, config.Version))

	tabs := container.NewAppTabs(
		container.NewTabItem(i18n.T("gui.tab.player"),
			newPaddedBoarder(nil, createPlayControllerV2(), nil, nil, createPlaylist()),
		),
		container.NewTabItem(i18n.T("gui.tab.search"),
			newPaddedBoarder(createSearchBar(), nil, nil, nil, createSearchList()),
		),
		container.NewTabItem(i18n.T("gui.tab.room"),
			newPaddedBoarder(nil, nil, createRoomSelector(), nil, createRoomController()),
		),
		container.NewTabItem(i18n.T("gui.tab.playlist"),
			newPaddedBoarder(nil, nil, createPlaylists(), nil, createPlaylistMedias()),
		),
		container.NewTabItem(i18n.T("gui.tab.history"),
			newPaddedBoarder(nil, nil, nil, nil, createHistoryList()),
		),
		container.NewTabItem(i18n.T("gui.tab.config"),
			newPaddedBoarder(nil, nil, nil, nil, createConfigLayout()),
		),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	MainWindow.SetIcon(fyne.NewStaticResource("icon", resource.ProgramIcon))
	MainWindow.SetContent(tabs)
	//MainWindow.Resize(fyne.NewSize(1280, 720))
	MainWindow.Resize(fyne.NewSize(960, 480))
	//MainWindow.SetFixedSize(true)
}

func AddConfigLayout(cfgs ...ConfigLayout) {
	ConfigList = append(ConfigList, cfgs...)
}
