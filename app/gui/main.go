package main

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/textinfo"
	"AynaLivePlayer/plugin/webinfo"
	"AynaLivePlayer/plugin/wylogin"
	"flag"
)

var dev = flag.Bool("dev", false, "generate new translation file")

var plugins = []controller.Plugin{diange.NewDiange(), qiege.NewQiege(), textinfo.NewTextInfo(), webinfo.NewWebInfo(),
	wylogin.NewWYLogin()}

func main() {
	flag.Parse()
	logger.Logger.Info("================Program Start================")
	logger.Logger.Infof("================Current Version: %s================", config.Version)
	controller.Initialize()
	gui.Initialize()
	controller.LoadPlugins(plugins...)
	gui.MainWindow.ShowAndRun()
	controller.ClosePlugins(plugins...)
	controller.Destroy()
	controller.CloseAndSave()
	if *dev {
		i18n.SaveTranslation()
	}
	logger.Logger.Info("================Program End================")
}
