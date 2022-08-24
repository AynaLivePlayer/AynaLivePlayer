package main

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/textinfo"
	"AynaLivePlayer/plugin/webinfo"
	"AynaLivePlayer/plugin/wylogin"
)

var plugins = []controller.Plugin{diange.NewDiange(), qiege.NewQiege(), textinfo.NewTextInfo(), webinfo.NewWebInfo(),
	wylogin.NewWYLogin()}

func main() {
	logger.Logger.Info("================Program Start================")
	logger.Logger.Infof("================Current Version: %s================", config.Version)
	controller.Initialize()
	controller.LoadPlugins(plugins...)
	gui.Initialize()
	gui.MainWindow.ShowAndRun()
	controller.ClosePlugins(plugins...)
	controller.Destroy()
	_ = config.SaveToConfigFile(config.ConfigPath)
	//i18n.SaveTranslation()
	logger.Logger.Info("================Program End================")
}
