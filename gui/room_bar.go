package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/i18n"
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
	RoomController.ConnectBtn = widget.NewButton(i18n.T("gui.room.btn.connect"), func() {
		RoomController.ConnectBtn.Disable()
		controller.SetDanmuClient(RoomController.Input.Text)
		if controller.LiveClient == nil {
			RoomController.Status.SetText(i18n.T("gui.room.status.failed"))
			RoomController.ConnectBtn.Enable()
			RoomController.Status.Refresh()
			return
		}
		RoomController.Input.SetOptions(config.LiveRoom.History)
		controller.LiveClient.Handler().RegisterA(liveclient.EventStatusChange, "gui.liveclient.status", func(event *event.Event) {
			d := event.Data.(liveclient.StatusChangeEvent)
			if d.Connected {
				RoomController.Status.SetText(i18n.T("gui.room.status.connected"))
			} else {
				RoomController.Status.SetText(i18n.T("gui.room.status.disconnected"))
			}
			RoomController.Status.Refresh()
		})
		go func() {
			controller.StartDanmuClient()
			RoomController.ConnectBtn.Enable()
		}()
	})
	RoomController.DisConnectBtn = widget.NewButton(i18n.T("gui.room.btn.disconnect"), func() {
		controller.ResetDanmuClient()
	})
	RoomController.Status = widget.NewLabel(i18n.T("gui.room.waiting"))
	return container.NewBorder(
		nil, nil,
		widget.NewLabel(i18n.T("gui.room.id")), container.NewHBox(RoomController.ConnectBtn, RoomController.DisConnectBtn),
		container.NewBorder(nil, nil, nil, RoomController.Status, RoomController.Input),
	)
}
