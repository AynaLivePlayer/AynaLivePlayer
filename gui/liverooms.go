package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sync"
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
	providers     []model.LiveRoomProviderInfo
	rooms         []model.LiveRoom
	lock          sync.RWMutex
}{}

func createRoomSelector() fyne.CanvasObject {
	RoomTab.Rooms = widget.NewList(
		func() int {
			return len(RoomTab.rooms)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(
				RoomTab.rooms[id].DisplayName())
		})
	RoomTab.AddBtn = widget.NewButton(i18n.T("gui.room.button.add"), func() {
		providerNames := make([]string, len(RoomTab.providers))
		for i := 0; i < len(RoomTab.providers); i++ {
			providerNames[i] = RoomTab.providers[i].Name
		}
		descriptionLabel := widget.NewLabel(i18n.T("gui.room.add.prompt"))
		clientNameEntry := widget.NewSelect(providerNames, func(s string) {
			for i := 0; i < len(RoomTab.providers); i++ {
				if RoomTab.providers[i].Name == s {
					descriptionLabel.SetText(i18n.T(RoomTab.providers[i].Description))
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
					logger.Infof("Add room %s %s", clientNameEntry.Selected, idEntry.Text)
					_ = global.EventBus.PublishToChannel(eventChannel,
						events.LiveRoomAddCmd,
						events.LiveRoomAddCmdEvent{
							Title:    nameEntry.Text,
							Provider: clientNameEntry.Selected,
							RoomKey:  idEntry.Text,
						})
				}
			},
			MainWindow,
		)
		dia.Resize(fyne.NewSize(512, 256))
		dia.Show()
	})
	RoomTab.RemoveBtn = widget.NewButton(i18n.T("gui.room.button.remove"), func() {
		if len(RoomTab.rooms) == 0 {
			return
		}
		_ = global.EventBus.PublishToChannel(eventChannel,
			events.LiveRoomRemoveCmd,
			events.LiveRoomRemoveCmdEvent{
				Identifier: RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier(),
			})
	})
	RoomTab.Rooms.OnSelected = func(id widget.ListItemID) {
		if id >= len(RoomTab.rooms) {
			return
		}
		logger.Infof("Select room %s", RoomTab.rooms[id].LiveRoom.Identifier())
		RoomTab.Index = id
		room := RoomTab.rooms[RoomTab.Index]
		RoomTab.RoomTitle.SetText(room.DisplayName())
		RoomTab.RoomID.SetText(room.LiveRoom.Identifier())
		RoomTab.AutoConnect.SetChecked(room.Config.AutoConnect)
		if room.Status {
			RoomTab.Status.SetText(i18n.T("gui.room.status.connected"))
		} else {
			RoomTab.Status.SetText(i18n.T("gui.room.status.disconnected"))
		}
		RoomTab.Status.Refresh()
	}
	registerRoomHandlers()
	return container.NewHBox(
		container.NewBorder(
			nil, container.NewCenter(container.NewHBox(RoomTab.AddBtn, RoomTab.RemoveBtn)),
			nil, nil,
			RoomTab.Rooms,
		),
		widget.NewSeparator(),
	)
}

func registerRoomHandlers() {
	global.EventBus.Subscribe(eventChannel, 
		events.LiveRoomProviderUpdate,
		"gui.liveroom.provider_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			RoomTab.providers = event.Data.(events.LiveRoomProviderUpdateEvent).Providers
			RoomTab.Rooms.Refresh()
		}))
	global.EventBus.Subscribe(eventChannel, 
		events.LiveRoomRoomsUpdate,
		"gui.liveroom.rooms_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			logger.Infof("Update rooms")
			data := event.Data.(events.LiveRoomRoomsUpdateEvent)
			RoomTab.lock.Lock()
			RoomTab.rooms = data.Rooms
			RoomTab.Rooms.Select(0)
			RoomTab.Rooms.Refresh()
			RoomTab.lock.Unlock()
		}))
	global.EventBus.Subscribe(eventChannel, 
		events.LiveRoomStatusUpdate,
		"gui.liveroom.room_status_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			room := event.Data.(events.LiveRoomStatusUpdateEvent).Room
			index := -1
			for i := 0; i < len(RoomTab.rooms); i++ {
				if RoomTab.rooms[i].LiveRoom.Identifier() == room.LiveRoom.Identifier() {
					index = i
					break
				}
			}
			if index == -1 {
				return
			}
			RoomTab.rooms[index] = room
			// add lock to avoid race condition
			RoomTab.lock.Lock()
			RoomTab.Rooms.Refresh()
			RoomTab.lock.Unlock()
			if index == RoomTab.Index {
				RoomTab.RoomTitle.SetText(room.DisplayName())
				RoomTab.RoomID.SetText(room.LiveRoom.Identifier())
				RoomTab.AutoConnect.SetChecked(room.Config.AutoConnect)
				if room.Status {
					RoomTab.Status.SetText(i18n.T("gui.room.status.connected"))
				} else {
					RoomTab.Status.SetText(i18n.T("gui.room.status.disconnected"))
				}
				RoomTab.Status.Refresh()
			}
		}))

}

func createRoomController() fyne.CanvasObject {
	RoomTab.ConnectBtn = widget.NewButton(i18n.T("gui.room.btn.connect"), func() {
		if RoomTab.Index >= len(RoomTab.rooms) {
			return
		}
		RoomTab.ConnectBtn.Disable()
		logger.Infof("Connect to room %s", RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier())
		_ = global.EventBus.PublishToChannel(eventChannel,
			events.LiveRoomOperationCmd,
			events.LiveRoomOperationCmdEvent{
				Identifier: RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier(),
				SetConnect: true,
			})
	})
	RoomTab.DisConnectBtn = widget.NewButton(i18n.T("gui.room.btn.disconnect"), func() {
		if RoomTab.Index >= len(RoomTab.rooms) {
			return
		}
		RoomTab.DisConnectBtn.Disable()
		logger.Infof("Disconnect from room %s", RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier())
		_ = global.EventBus.PublishToChannel(eventChannel,
			events.LiveRoomOperationCmd,
			events.LiveRoomOperationCmdEvent{
				Identifier: RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier(),
				SetConnect: false,
			})
	})
	global.EventBus.Subscribe(eventChannel, 
		events.LiveRoomOperationFinish,
		"gui.liveroom.operation_finish",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			RoomTab.ConnectBtn.Enable()
			RoomTab.DisConnectBtn.Enable()
		}))
	RoomTab.Status = widget.NewLabel(i18n.T("gui.room.waiting"))
	RoomTab.RoomTitle = widget.NewLabel("")
	RoomTab.RoomID = widget.NewLabel("")
	RoomTab.AutoConnect = widget.NewCheck(i18n.T("gui.room.check.autoconnect"), func(b bool) {
		if RoomTab.Index >= len(RoomTab.rooms) {
			return
		}
		logger.Infof("Change room %s autoconnect to %v", RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier(), b)
		_ = global.EventBus.PublishToChannel(eventChannel,
			events.LiveRoomConfigChangeCmd,
			events.LiveRoomConfigChangeCmdEvent{
				Identifier: RoomTab.rooms[RoomTab.Index].LiveRoom.Identifier(),
				Config: model.LiveRoomConfig{
					AutoConnect: b,
				},
			})
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
