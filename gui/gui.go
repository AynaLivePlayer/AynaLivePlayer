package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	configView "AynaLivePlayer/gui/views/config"
	"AynaLivePlayer/gui/views/history"
	"AynaLivePlayer/gui/views/liverooms"
	"AynaLivePlayer/gui/views/player"
	"AynaLivePlayer/gui/views/playlists"
	"AynaLivePlayer/gui/views/search"
	"AynaLivePlayer/gui/views/systray"
	"AynaLivePlayer/gui/views/updater"
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

var logger _logger.ILogger = nil

func black_magic() {
	widget.RichTextStyleStrong.TextStyle.Bold = false
}

func Initialize() {
	logger = global.Logger.WithPrefix("GUI")
	gctx.Logger = logger
	black_magic()
	logger.Info("Initializing GUI")

	if config.General.CustomFonts != "" {
		_ = os.Setenv("FYNE_FONT", config.GetAssetPath(config.General.CustomFonts))
	}

	mainApp := app.NewWithID(config.ProgramName)
	MainWindow := mainApp.NewWindow(getAppTitle())

	gctx.Context = gctx.NewGuiContext(mainApp, MainWindow)
	gctx.Context.Init()

	gctx.Context.OnMainWindowClosing(func() {
		_ = config.SaveToConfigFile(config.ConfigPath)
		logger.Infof("config saved to %s", config.ConfigPath)
	})

	updater.CreateUpdaterPopUp()

	tabs := container.NewAppTabs(
		container.NewTabItem(i18n.T("gui.tab.player"), player.CreateView()),
		container.NewTabItem(i18n.T("gui.tab.search"), search.CreateView()),
		container.NewTabItem(i18n.T("gui.tab.room"), liverooms.CreateView()),
		container.NewTabItem(i18n.T("gui.tab.playlist"), playlists.CreateView()),
		container.NewTabItem(i18n.T("gui.tab.history"), history.CreateView()),
		container.NewTabItem(i18n.T("gui.tab.config"), configView.CreateView()),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	MainWindow.SetIcon(resource.ImageIcon)
	MainWindow.SetContent(tabs)
	MainWindow.Resize(fyne.NewSize(config.General.Width, config.General.Height))
	MainWindow.SetFixedSize(config.General.FixedSize)

	// todo: fix, window were created even if not show. this block gui from closing
	// i can't create sub window before the main window shows.
	// setupPlayerWindow()

	// register error
	global.EventBus.Subscribe(gctx.EventChannel,
		events.ErrorUpdate, "gui.show_error", gutil.ThreadSafeHandler(func(e *eventbus.Event) {
			err := e.Data.(events.ErrorUpdateEvent).Error
			logger.Warnf("gui received error event: %v, %v", err, err == nil)
			if err == nil {
				return
			}
			dialog.ShowError(err, MainWindow)
		}))

	if config.General.ShowSystemTray {
		systray.SetupSysTray()
	}
}
