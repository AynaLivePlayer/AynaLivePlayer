package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/sirupsen/logrus"
	"os"
)

const MODULE_GUI = "GUI"

var App fyne.App
var MainWindow fyne.Window

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_GUI)
}

func Initialize() {
	os.Setenv("FYNE_FONT", config.GetAssetPath("msyh.ttc"))
	App = app.New()
	MainWindow = App.NewWindow("AynaLivePlayer")

	tabs := container.NewAppTabs(
		container.NewTabItem("Player",
			newPaddedBoarder(nil, createPlayController(), nil, nil, createPlaylist()),
		),
		container.NewTabItem("Search",
			newPaddedBoarder(createSearchBar(), nil, nil, nil, createSearchList()),
		),
		container.NewTabItem("Room",
			newPaddedBoarder(createRoomController(), nil, nil, nil, createRoomLogger()),
		),
		container.NewTabItem("Playlist",
			newPaddedBoarder(nil, nil, createPlaylists(), nil, createPlaylistMedias()),
		),
	)

	tabs.SetTabLocation(container.TabLocationTop)

	MainWindow.SetContent(tabs)
	//MainWindow.Resize(fyne.NewSize(1280, 720))
	MainWindow.Resize(fyne.NewSize(960, 480))
	//MainWindow.SetFixedSize(true)
}
