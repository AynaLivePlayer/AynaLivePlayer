package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RoomLoggerContainer struct {
}

var RoomLogger = &RoomLoggerContainer{}

func createRoomLogger() fyne.CanvasObject {
	//b := NewAsyncButton("ceshi", func() {
	//	time.Sleep(time.Second * 5)
	//})
	return container.NewVBox(
		widget.NewLabel("广告位招租"),
	)
}
