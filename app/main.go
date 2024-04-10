package main

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/internal"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	loggerRepo "AynaLivePlayer/pkg/logger/repository"
	"flag"
)

var dev = flag.Bool("dev", false, "dev")

type _LogConfig struct {
	config.BaseConfig
	Path           string
	Level          logger.LogLevel
	RedirectStderr bool
	MaxSize        int64
}

func (c *_LogConfig) Name() string {
	return "Log"
}

var Log = &_LogConfig{
	Path:           "./log.txt",
	Level:          logger.LogLevelInfo,
	RedirectStderr: false, // this should be true if it is in production mode.
	MaxSize:        5,
}

//func createController(log adapter.ILogger) adapter.IControlBridge {
//	logbridge := log.WithModule("ControlBridge")
//	em := event.MainManager.NewChildManager()
//	liveroom := internal.NewLiveRoomController(
//		logbridge)
//	lyric := internal.NewLyricLoader()
//	provider := internal.NewProviderController(logbridge)
//	playlist := internal.NewPlaylistController(em, logbridge, provider)
//	plugin := internal.NewPluginController(logbridge)
//	mpvPlayer := player.NewMpvPlayer(em, logbridge)
//	playControl := internal.NewPlayerController(mpvPlayer, playlist, lyric, provider, logbridge)
//	ctr := internal.NewController(
//		liveroom, playControl, playlist, provider, plugin,
//		logbridge,
//	)
//	return ctr
//}

func setupGlobal() {
	global.EventManager = event.NewManger(128, 16)
	global.Logger = loggerRepo.NewZapColoredLogger()
	global.Logger.SetLogLevel(Log.Level)
}

func main() {
	flag.Parse()
	config.LoadFromFile(config.ConfigPath)
	config.LoadConfig(Log)
	i18n.LoadLanguage(config.General.Language)
	setupGlobal()
	global.Logger.Info("================Program Start================")
	global.Logger.Infof("================Current Version: %s================", model.Version(config.Version))
	//mainController := createController(log)
	internal.Initialize()
	gui.Initialize()
	global.EventManager.Start()
	//plugins := []adapter.Plugin{diange.NewDiange(mainController), qiege.NewQiege(mainController),
	//	textinfo.NewTextInfo(mainController), webinfo.NewWebInfo(mainController),
	//	wylogin.NewWYLogin(mainController)}
	//mainController.LoadPlugins(plugins...)
	gui.MainWindow.ShowAndRun()
	internal.Stop()
	global.EventManager.Stop()
	if *dev {
		i18n.SaveTranslation()
	}
	_ = config.SaveToConfigFile(config.ConfigPath)
	global.Logger.Info("================Program End================")
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
