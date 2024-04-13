package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/resource"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	_logger "AynaLivePlayer/pkg/logger"
)

var App fyne.App
var MainWindow fyne.Window
var playerWindow fyne.Window
var playerWindowHandle uintptr

var logger _logger.ILogger = nil

func black_magic() {
	widget.RichTextStyleStrong.TextStyle.Bold = false
}

func Initialize() {
	logger = global.Logger.WithPrefix("GUI")
	black_magic()
	logger.Info("Initializing GUI")
	//os.Setenv("FYNE_FONT", config.GetAssetPath("msyh.ttc"))
	App = app.NewWithID(config.ProgramName)
	App.Settings().SetTheme(&myTheme{})
	MainWindow = App.NewWindow(fmt.Sprintf("%s Ver %s", config.ProgramName, model.Version(config.Version)))

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
		//container.NewTabItem(i18n.T("gui.tab.config"),
		//	createConfigLayout(),
		//),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	MainWindow.SetIcon(resource.ImageIcon)
	MainWindow.SetContent(tabs)
	//MainWindow.Resize(fyne.NewSize(1280, 720))
	MainWindow.Resize(fyne.NewSize(config.General.Width, config.General.Height))

	setupPlayerWindow()

	// register error
	global.EventManager.RegisterA(
		events.ErrorUpdate, "gui.show_error", func(e *event.Event) {
			err := e.Data.(events.ErrorUpdateEvent).Error
			logger.Warnf("gui received error event: %v, %v", err, err == nil)
			if err == nil {
				return
			}
			dialog.ShowError(err, MainWindow)
		})

	MainWindow.SetFixedSize(true)
	if config.General.ShowSystemTray {
		setupSysTray()
	}
	//if config2.General.AutoCheckUpdate {
	//	go checkUpdate()
	//}
}

//
//func checkUpdate() {
//	l().Info("checking updates...")
//	err := API.App().CheckUpdate()
//	if err != nil {
//		showDialogIfError(err)
//		l().Warnf("check update failed", err)
//		return
//	}
//	l().Infof("latest version: v%s", API.App().LatestVersion().Version)
//	if API.App().LatestVersion().Version > API.App().Version().Version {
//		l().Info("new update available")
//		dialog.ShowCustom(
//			i18n.T("gui.update.new_version"),
//			"OK",
//			widget.NewRichTextFromMarkdown(API.App().LatestVersion().Info),
//			MainWindow)
//	}
//}
