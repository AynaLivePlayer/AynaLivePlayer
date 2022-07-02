package main

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/textinfo"
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

func main() {
	fmt.Printf("BiliAudioBot Revive %s\n", config.VERSION)
	//logger.Logger.SetLevel(logrus.DebugLevel)
	controller.Initialize()
	controller.LoadPlugins(diange.NewDiange(), qiege.NewQiege(), textinfo.NewTextInfo())
	defer func() {
		controller.Destroy()
		config.SaveToConfigFile(config.CONFIG_PATH)
		//i18n.SaveTranslation()
	}()
	gui.Initialize()
	gui.MainWindow.ShowAndRun()
}
