package main

import (
	"AynaLivePlayer/adapters"
	"AynaLivePlayer/adapters/player"
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/internal"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/textinfo"
	"AynaLivePlayer/plugin/webinfo"
	"AynaLivePlayer/plugin/wylogin"
	"flag"
)

var dev = flag.Bool("dev", false, "generate new translation file")

type _LogConfig struct {
	config.BaseConfig
	Path           string
	Level          adapter.LogLevel
	RedirectStderr bool
}

func (c *_LogConfig) Name() string {
	return "Log"
}

var Log = &_LogConfig{
	Path:           "./log.txt",
	Level:          adapter.LogLevelInfo,
	RedirectStderr: false, // this should be true if it is in production mode.
}

func createController(log adapter.ILogger) adapter.IControlBridge {
	logbridge := log.WithModule("ControlBridge")
	em := event.MainManager.NewChildManager()
	liveroom := internal.NewLiveRoomController(
		logbridge)
	lyric := internal.NewLyricLoader()
	provider := internal.NewProviderController(logbridge)
	playlist := internal.NewPlaylistController(em, logbridge, provider)
	plugin := internal.NewPluginController(logbridge)
	mpvPlayer := player.NewMpvPlayer(em, logbridge)
	playControl := internal.NewPlayerController(mpvPlayer, playlist, lyric, provider, logbridge)
	ctr := internal.NewController(
		liveroom, playControl, playlist, provider, plugin,
		logbridge,
	)
	return ctr
}

func main() {
	flag.Parse()
	config.LoadFromFile(config.ConfigPath)
	config.LoadConfig(Log)
	i18n.LoadLanguage(config.General.Language)
	log := adapters.Logger.NewLogrus(Log.Path, Log.RedirectStderr)
	log.SetLogLevel(Log.Level)
	log.Info("================Program Start================")
	log.Infof("================Current Version: %s================", config.Version)
	mainController := createController(log)
	gui.API = mainController
	gui.Initialize()
	plugins := []adapter.Plugin{diange.NewDiange(mainController), qiege.NewQiege(mainController),
		textinfo.NewTextInfo(mainController), webinfo.NewWebInfo(mainController),
		wylogin.NewWYLogin(mainController)}
	mainController.LoadPlugins(plugins...)
	gui.MainWindow.ShowAndRun()
	mainController.CloseAndSave()
	if *dev {
		i18n.SaveTranslation()
	}
	_ = config.SaveToConfigFile(config.ConfigPath)
	log.Info("================Program End================")
}

////go:embed all:../webgui/frontend/dist
//var assets embed.FS
//
//func main() {
//	// Create an instance of the app structure
//	app := webgui.NewApp()
//
//	// Create application with options
//	err := wails.Run(&options.App{
//		Title:  "AynaLivePlayer",
//		Width:  1024,
//		Height: 768,
//		AssetServer: &assetserver.Options{
//			Assets: assets,
//		},
//		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
//		OnStartup:        app.Startup,
//		Bind: []interface{}{
//			app,
//		},
//	})
//
//	if err != nil {
//		println("Error:", err.Error())
//	}
//}
