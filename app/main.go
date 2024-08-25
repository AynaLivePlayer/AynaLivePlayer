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
	"time"
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

func setupGlobal() {
	global.EventManager = event.NewManger(128, 16)
	global.Logger = loggerRepo.NewZapColoredLogger(Log.Path, !*dev)
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
	internal.Initialize()
	gui.Initialize()
	go func() {
		// temporary fix for gui not render correctly.
		// wait until gui rendered then start event dispatching
		time.Sleep(1 * time.Second)
		global.EventManager.Start()
	}()
	gui.MainWindow.ShowAndRun()
	global.Logger.Info("closing internal server")
	internal.Stop()
	global.Logger.Infof("closing event manager")
	global.EventManager.Stop()
	if *dev {
		i18n.SaveTranslation()
	}
	_ = config.SaveToConfigFile(config.ConfigPath)
	global.Logger.Info("================Program End================")
}
