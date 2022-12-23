package main

import (
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/controller/core"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/player"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/textinfo"
	"AynaLivePlayer/plugin/webinfo"
	"AynaLivePlayer/plugin/wylogin"
	"flag"
)

var dev = flag.Bool("dev", false, "generate new translation file")

func createController() controller.IController {
	liveroom := core.NewLiveRoomController()
	lyric := core.NewLyricLoader()
	provider := core.NewProviderController()
	playlist := core.NewPlaylistController(provider)
	plugin := core.NewPluginController()
	mpvPlayer := player.NewMpvPlayer()
	playControl := core.NewPlayerController(mpvPlayer, playlist, lyric, provider)
	ctr := core.NewController(liveroom, playControl, playlist, provider, plugin)
	return ctr
}

func main() {
	flag.Parse()
	logger.Logger.Info("================Program Start================")
	logger.Logger.Infof("================Current Version: %s================", config.Version)
	mainController := createController()
	controller.Instance = mainController
	gui.Initialize()
	plugins := []controller.Plugin{diange.NewDiange(mainController), qiege.NewQiege(mainController),
		textinfo.NewTextInfo(mainController), webinfo.NewWebInfo(mainController),
		wylogin.NewWYLogin()}
	mainController.LoadPlugins(plugins...)
	gui.MainWindow.ShowAndRun()
	mainController.CloseAndSave()
	if *dev {
		i18n.SaveTranslation()
	}
	_ = config.SaveToConfigFile(config.ConfigPath)
	logger.Logger.Info("================Program End================")
}
