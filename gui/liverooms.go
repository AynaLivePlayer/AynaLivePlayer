package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/core/events"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var RoomTab = &struct {
	Rooms         *widget.List
	Index         int
	AddBtn        *widget.Button
	RemoveBtn     *widget.Button
	RoomTitle     *widget.Label
	RoomID        *widget.Label
	Status        *widget.Label
	AutoConnect   *widget.Check
	ConnectBtn    *widget.Button
	DisConnectBtn *widget.Button
}{}

func createRoomSelector() fyne.CanvasObject {
	RoomTab.Rooms = widget.NewList(
		func() int {
			return API.LiveRooms().Size()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(
				API.LiveRooms().Get(id).DisplayName())
		})
	RoomTab.AddBtn = widget.NewButton(i18n.T("gui.room.button.add"), func() {
		clientNameEntry := widget.NewSelect(API.LiveRooms().GetAllClientNames(), nil)
		idEntry := widget.NewEntry()
		dia := dialog.NewCustomConfirm(
			i18n.T("gui.room.add.title"),
			i18n.T("gui.room.add.confirm"),
			i18n.T("gui.room.add.cancel"),
			container.NewVBox(
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel(i18n.T("gui.room.add.client_name")),
					clientNameEntry,
					widget.NewLabel(i18n.T("gui.room.add.id_url")),
					idEntry,
				),
				widget.NewLabel(i18n.T("gui.room.add.prompt")),
			),
			func(b bool) {
				if b && len(clientNameEntry.Selected) > 0 && len(idEntry.Text) > 0 {
					_, err := API.LiveRooms().AddRoom(clientNameEntry.Selected, idEntry.Text)
					showDialogIfError(err)
					RoomTab.Rooms.Refresh()
				}
			},
			MainWindow,
		)
		dia.Resize(fyne.NewSize(512, 256))
		dia.Show()
	})
	RoomTab.RemoveBtn = widget.NewButton(i18n.T("gui.room.button.remove"), func() {
		showDialogIfError(API.LiveRooms().DeleteRoom(PlaylistManager.Index))
		RoomTab.Rooms.Select(0)
		RoomTab.Rooms.Refresh()
	})
	RoomTab.Rooms.OnSelected = func(id widget.ListItemID) {
		rom := API.LiveRooms().Get(PlaylistManager.Index)
		if rom != nil {
			rom.EventManager().Unregister("gui.liveroom.status")
		}
		RoomTab.Index = id
		rom = API.LiveRooms().Get(RoomTab.Index)
		rom.EventManager().RegisterA(events.LiveRoomStatusChange, "gui.liveroom.status", func(event *event.Event) {
			d := event.Data.(events.StatusChangeEvent)
			if d.Connected {
				RoomTab.Status.SetText(i18n.T("gui.room.status.connected"))
			} else {
				RoomTab.Status.SetText(i18n.T("gui.room.status.disconnected"))
			}
			RoomTab.Status.Refresh()
		})
		RoomTab.RoomTitle.SetText(rom.DisplayName())
		RoomTab.RoomID.SetText(rom.Identifier())
		RoomTab.AutoConnect.SetChecked(rom.Model().AutoConnect)
		if API.LiveRooms().GetRoomStatus(RoomTab.Index) {
			RoomTab.Status.SetText(i18n.T("gui.room.status.connected"))
		} else {
			RoomTab.Status.SetText(i18n.T("gui.room.status.disconnected"))
		}
		RoomTab.Status.Refresh()
	}
	return container.NewHBox(
		container.NewBorder(
			nil, container.NewCenter(container.NewHBox(RoomTab.AddBtn, RoomTab.RemoveBtn)),
			nil, nil,
			RoomTab.Rooms,
		),
		widget.NewSeparator(),
	)
}

func createRoomController() fyne.CanvasObject {
	RoomTab.ConnectBtn = widget.NewButton(i18n.T("gui.room.btn.connect"), func() {
		RoomTab.ConnectBtn.Disable()
		go func() {
			_ = API.LiveRooms().Connect(RoomTab.Index)
			RoomTab.ConnectBtn.Enable()
		}()
	})
	RoomTab.DisConnectBtn = widget.NewButton(i18n.T("gui.room.btn.disconnect"), func() {
		_ = API.LiveRooms().Disconnect(RoomTab.Index)
	})
	RoomTab.Status = widget.NewLabel(i18n.T("gui.room.waiting"))
	RoomTab.RoomTitle = widget.NewLabel("")
	RoomTab.RoomID = widget.NewLabel("")
	RoomTab.AutoConnect = widget.NewCheck(i18n.T("gui.room.check.autoconnect"), func(b bool) {
		rom := API.LiveRooms().Get(RoomTab.Index)
		if rom != nil {
			rom.Model().AutoConnect = b
		}
		return
	})
	RoomTab.Rooms.Select(0)
	return container.NewVBox(
		RoomTab.RoomTitle,
		RoomTab.RoomID,
		RoomTab.Status,
		container.NewHBox(widget.NewLabel(i18n.T("gui.room.check.autoconnect")), RoomTab.AutoConnect),
		container.NewHBox(RoomTab.ConnectBtn, RoomTab.DisConnectBtn),
	)
}
