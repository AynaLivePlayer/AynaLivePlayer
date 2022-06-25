package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/liveclient"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RoomControllerContainer struct {
	Input         *widget.SelectEntry
	ConnectBtn    *widget.Button
	DisConnectBtn *widget.Button
	Status        *widget.Label
}

var RoomController = &RoomControllerContainer{}

func createRoomController() fyne.CanvasObject {
	RoomController.Input = widget.NewSelectEntry(config.LiveRoom.History)
	RoomController.ConnectBtn = widget.NewButton("Connect", func() {
		RoomController.ConnectBtn.Disable()
		controller.SetDanmuClient(RoomController.Input.Text)
		if controller.LiveClient == nil {
			RoomController.Status.SetText("Set Failed")
			RoomController.ConnectBtn.Enable()
			RoomController.Status.Refresh()
			return
		}
		RoomController.Input.SetOptions(config.LiveRoom.History)
		controller.LiveClient.Handler().RegisterA(liveclient.EventStatusChange, "gui.liveclient.status", func(event *event.Event) {
			d := event.Data.(liveclient.StatusChangeEvent)
			if d.Connected {
				RoomController.Status.SetText("Connected")
			} else {
				RoomController.Status.SetText("Disconnected")
			}
			RoomController.Status.Refresh()
		})
		go func() {
			controller.StartDanmuClient()
			RoomController.ConnectBtn.Enable()
		}()
	})
	RoomController.DisConnectBtn = widget.NewButton("Disconnect", func() {
		controller.ResetDanmuClient()
	})
	RoomController.Status = widget.NewLabel("Waiting")
	return container.NewBorder(
		nil, nil,
		widget.NewLabel("Room ID:"), container.NewHBox(RoomController.ConnectBtn, RoomController.DisConnectBtn),
		container.NewBorder(nil, nil, nil, RoomController.Status, RoomController.Input),
	)
}
