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
	"fmt"
	"github.com/mitchellh/panicwrap"
	"os"
)

func init() {
	exitStatus, _ := panicwrap.BasicWrap(func(s string) {
		logger.Logger.Panic(s)
		os.Exit(1)
		return
	})
	if exitStatus >= 0 {
		os.Exit(exitStatus)
	}
}

var plugins = []controller.Plugin{diange.NewDiange(), qiege.NewQiege(), textinfo.NewTextInfo(), webinfo.NewWebInfo(),
	wylogin.NewWYLogin()}

func main() {
	fmt.Printf("BiliAudioBot Revive %s\n", config.VERSION)
	controller.Initialize()
	controller.LoadPlugins(plugins...)
	defer func() {
		controller.Destroy()
		config.SaveToConfigFile(config.CONFIG_PATH)
		//i18n.SaveTranslation()
	}()
	gui.Initialize()
	gui.MainWindow.ShowAndRun()
	controller.ClosePlugins(plugins...)
}
