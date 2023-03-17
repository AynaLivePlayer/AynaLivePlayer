package gui

import (
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/resource"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var API adapter.IControlBridge
var App fyne.App
var MainWindow fyne.Window
var playerWindow fyne.Window
var playerWindowHandle uintptr

func l() adapter.ILogger {
	return API.Logger().WithModule("GUI")
}

func black_magic() {
	widget.RichTextStyleStrong.TextStyle.Bold = false
}

func Initialize() {
	black_magic()
	l().Info("Initializing GUI")
	//os.Setenv("FYNE_FONT", config.GetAssetPath("msyh.ttc"))
	App = app.New()
	App.Settings().SetTheme(&myTheme{})
	MainWindow = App.NewWindow(fmt.Sprintf("%s Ver.%s", config.ProgramName, model.Version(config.Version)))

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

	playerWindow = App.NewWindow("CorePlayerPreview")
	playerWindow.Resize(fyne.NewSize(480, 240))
	playerWindow.SetCloseIntercept(func() {
		playerWindow.Hide()
	})
	MainWindow.SetOnClosed(func() {
		playerWindow.Close()
	})
	playerWindow.Hide()

	//MainWindow.SetFixedSize(true)
	if config.General.AutoCheckUpdate {
		go checkUpdate()
	}
}

func showPlayerWindow() {
	playerWindow.Show()
	if playerWindowHandle == 0 {
		playerWindowHandle = util.GetWindowHandle("CorePlayerPreview")
		l().Infof("video output window handle: %d", playerWindowHandle)
		if playerWindowHandle != 0 {
			_ = API.PlayControl().GetPlayer().SetWindowHandle(playerWindowHandle)
		}
	}
}

func addShortCut() {
	key := &desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}
	MainWindow.Canvas().AddShortcut(key, func(shortcut fyne.Shortcut) {
		l().Info("Shortcut pressed: Ctrl+Shift+Right")
		API.PlayControl().PlayNext()
	})
}

func checkUpdate() {
	l().Info("checking updates...")
	err := API.App().CheckUpdate()
	if err != nil {
		showDialogIfError(err)
		l().Warnf("check update failed", err)
		return
	}
	l().Infof("latest version: v%s", API.App().LatestVersion().Version)
	if API.App().LatestVersion().Version > API.App().Version().Version {
		l().Info("new update available")
		dialog.ShowCustom(
			i18n.T("gui.update.new_version"),
			"OK",
			widget.NewRichTextFromMarkdown(API.App().LatestVersion().Info),
			MainWindow)
	}
}
