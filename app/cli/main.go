package main

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/logger"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Printf("BiliAudioBot Revive %s\n", config.VERSION)
	logger.Logger.SetLevel(logrus.DebugLevel)
	fmt.Println("Please enter room id")
	var roomid string

	// Taking input from user
	fmt.Scanln(&roomid)
	controller.Initialize()
	controller.SetDanmuClient(roomid)
	ch := make(chan int)
	<-ch
}
