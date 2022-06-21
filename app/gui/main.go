package main

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/logger"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Printf("BiliAudioBot Revive %s\n", config.VERSION)
	logger.Logger.SetLevel(logrus.DebugLevel)
	controller.Initialize()
	defer func() {
		controller.Destroy()
		config.SaveToConfigFile(config.CONFIG_PATH)
	}()
	gui.Initialize()
	gui.MainWindow.ShowAndRun()
}
