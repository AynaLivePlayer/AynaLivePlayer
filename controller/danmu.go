package controller

import (
	"AynaLivePlayer/event"
	"AynaLivePlayer/liveclient"
)

var DanmuHandlers []DanmuHandler

type DanmuHandler interface {
	Execute(anmu *liveclient.DanmuMessage)
}

func AddDanmuHandler(handlers ...DanmuHandler) {
	DanmuHandlers = append(DanmuHandlers, handlers...)
}

func danmuHandler(event *event.Event) {
	danmu := event.Data.(*liveclient.DanmuMessage)
	for _, cmd := range DanmuHandlers {
		cmd.Execute(danmu)
	}
}
