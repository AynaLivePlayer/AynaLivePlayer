package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/resource"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os"

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
	if config.General.CustomFonts != "" {
		_ = os.Setenv("FYNE_FONT", config.GetAssetPath(config.General.CustomFonts))
	}
	App = app.NewWithID(config.ProgramName)
	//App.Settings().SetTheme(&myTheme{})
	MainWindow = App.NewWindow(getAppTitle())

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
	MainWindow.Resize(fyne.NewSize(config.General.Width, config.General.Height))

	// todo: fix, window were created even if not show. this block gui from closing
	// i can't create sub window before the main window shows.
	// setupPlayerWindow()

	// register error
	global.EventBus.Subscribe("",
		events.ErrorUpdate, "gui.show_error", gutil.ThreadSafeHandler(func(e *eventbus.Event) {
			err := e.Data.(events.ErrorUpdateEvent).Error
			logger.Warnf("gui received error event: %v, %v", err, err == nil)
			if err == nil {
				return
			}
			dialog.ShowError(err, MainWindow)
		}))

	checkUpdate()
	MainWindow.SetFixedSize(config.General.FixedSize)
	if config.General.ShowSystemTray {
		setupSysTray()
	} else {
		MainWindow.SetCloseIntercept(
			func() {
				// todo: save twice i don;t care
				_ = config.SaveToConfigFile(config.ConfigPath)
				MainWindow.Close()
			})
	}
	MainWindow.SetOnClosed(func() {
		logger.Infof("GUI closing")
		if playerWindow != nil {
			logger.Infof("player window closing")
			playerWindow.Close()
		}
	})
}
