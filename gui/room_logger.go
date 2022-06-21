package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type RoomLoggerContainer struct {
}

var RoomLogger = &RoomLoggerContainer{}

func createRoomLogger() fyne.CanvasObject {
	return widget.NewLabel("广告位招租")
}
