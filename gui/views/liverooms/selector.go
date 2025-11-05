package liverooms

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var roomSelectorView = &struct {
	rooms     *widget.List
	addBtn    *widget.Button
	removeBtn *widget.Button
}{}

func renderRoomList() {
	lock.Lock()
	roomSelectorView.rooms.Refresh()
	lock.Unlock()
}

func createRoomSelector() fyne.CanvasObject {
	roomSelectorView.rooms = widget.NewList(
		func() int {
			return len(rooms)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(
				rooms[id].DisplayName())
		})
	roomSelectorView.addBtn = widget.NewButton(i18n.T("gui.room.button.add"), func() {
		providerNames := make([]string, len(providers))
		for i := 0; i < len(providers); i++ {
			providerNames[i] = providers[i].Name
		}
		descriptionLabel := widget.NewLabel(i18n.T("gui.room.add.prompt"))
		clientNameEntry := widget.NewSelect(providerNames, func(s string) {
			for i := 0; i < len(providers); i++ {
				if providers[i].Name == s {
					descriptionLabel.SetText(i18n.T(providers[i].Description))
					break
				}
				descriptionLabel.SetText("")
			}
		})
		idEntry := widget.NewEntry()
		nameEntry := widget.NewEntry()
		dia := dialog.NewCustomConfirm(
			i18n.T("gui.room.add.title"),
			i18n.T("gui.room.add.confirm"),
			i18n.T("gui.room.add.cancel"),
			container.NewVBox(
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel(i18n.T("gui.room.add.name")),
					nameEntry,
					widget.NewLabel(i18n.T("gui.room.add.client_name")),
					clientNameEntry,
					widget.NewLabel(i18n.T("gui.room.add.id_url")),
					idEntry,
				),
				descriptionLabel,
			),
			func(b bool) {
				if b && len(clientNameEntry.Selected) > 0 && len(idEntry.Text) > 0 {
					gctx.Logger.Infof("Add room %s %s", clientNameEntry.Selected, idEntry.Text)
					_ = global.EventBus.PublishToChannel(gctx.EventChannel,
						events.CmdLiveRoomAdd,
						events.CmdLiveRoomAddData{
							Title:    nameEntry.Text,
							Provider: clientNameEntry.Selected,
							RoomKey:  idEntry.Text,
						})
				}
			},
			gctx.Context.Window,
		)
		dia.Resize(fyne.NewSize(512, 256))
		dia.Show()
	})
	roomSelectorView.removeBtn = widget.NewButton(i18n.T("gui.room.button.remove"), func() {
		room, ok := getCurrentRoom()
		if !ok {
			return
		}
		_ = global.EventBus.PublishToChannel(gctx.EventChannel,
			events.CmdLiveRoomRemove,
			events.CmdLiveRoomRemoveData{
				Identifier: room.LiveRoom.Identifier(),
			})
	})
	roomSelectorView.rooms.OnSelected = func(id widget.ListItemID) {
		if id >= len(rooms) {
			return
		}
		currentIndex = id
		room, ok := getCurrentRoom()
		if !ok {
			return
		}
		gctx.Logger.Infof("Select room %s", room.LiveRoom.Identifier())
		renderCurrentRoom()
	}
	return container.NewHBox(
		container.NewBorder(
			nil, container.NewCenter(container.NewHBox(roomSelectorView.addBtn, roomSelectorView.removeBtn)),
			nil, nil,
			roomSelectorView.rooms,
		),
		widget.NewSeparator(),
	)
}
