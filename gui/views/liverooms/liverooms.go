package liverooms

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sync"
)

func CreateView() fyne.CanvasObject {
	view := container.NewBorder(nil, nil, createRoomSelector(), nil, createRoomController())
	registerRoomHandlers()
	return view
}

var providers []model.LiveRoomProviderInfo = make([]model.LiveRoomProviderInfo, 0)
var rooms []model.LiveRoom = make([]model.LiveRoom, 0)

var lock sync.RWMutex

var currentRoomView = &struct {
	roomTitle     *widget.Label
	roomID        *widget.Label
	status        *widget.Label
	autoConnect   *widget.Check
	connectBtn    *widget.Button
	disConnectBtn *widget.Button
}{}

var currentIndex int = 0

func getCurrentRoom() (model.LiveRoom, bool) {
	lock.RLock()
	if currentIndex >= len(rooms) {
		lock.RUnlock()
		return model.LiveRoom{}, false
	}
	room := rooms[currentIndex]
	lock.RUnlock()
	return room, true
}

func renderCurrentRoom() {
	room, ok := getCurrentRoom()
	if !ok {
		currentRoomView.roomTitle.SetText("")
		currentRoomView.roomID.SetText("")
		currentRoomView.autoConnect.SetChecked(false)
		currentRoomView.status.SetText(i18n.T("gui.room.waiting"))
		return
	}
	currentRoomView.roomTitle.SetText(room.DisplayName())
	currentRoomView.roomID.SetText(room.LiveRoom.Identifier())
	currentRoomView.autoConnect.SetChecked(room.Config.AutoConnect)
	if room.Status {
		currentRoomView.status.SetText(i18n.T("gui.room.status.connected"))
	} else {
		currentRoomView.status.SetText(i18n.T("gui.room.status.disconnected"))
	}
}

func registerRoomHandlers() {
	global.EventBus.Subscribe(gctx.EventChannel,
		events.LiveRoomProviderUpdate,
		"gui.liveroom.provider_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			providers = event.Data.(events.LiveRoomProviderUpdateEvent).Providers
			//RoomTab.Rooms.Refresh()
		}))
	global.EventBus.Subscribe(gctx.EventChannel,
		events.UpdateLiveRoomRooms,
		"gui.liveroom.rooms_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			gctx.Logger.Infof("Update rooms")
			lock.Lock()
			rooms = event.Data.(events.UpdateLiveRoomRoomsData).Rooms
			lock.Unlock()
			renderRoomList()
			renderCurrentRoom()
		}))
	global.EventBus.Subscribe(gctx.EventChannel,
		events.UpdateLiveRoomStatus,
		"gui.liveroom.room_status_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			room := event.Data.(events.UpdateLiveRoomStatusData).Room
			lock.Lock()
			index := -1
			for i := 0; i < len(rooms); i++ {
				if rooms[i].LiveRoom.Identifier() == room.LiveRoom.Identifier() {
					index = i
					break
				}
			}
			if index == -1 {
				lock.Unlock()
				return
			}
			rooms[index] = room
			lock.Unlock()
			if index == currentIndex {
				renderCurrentRoom()
			}
		}))

}

func createRoomController() fyne.CanvasObject {
	currentRoomView.connectBtn = widget.NewButton(i18n.T("gui.room.btn.connect"), func() {
		room, ok := getCurrentRoom()
		if !ok {
			return
		}
		gctx.Logger.Infof("Connect to room %s", room.LiveRoom.Identifier())
		currentRoomView.connectBtn.Disable()
		go func() {
			resp, err := global.EventBus.Call(events.CmdLiveRoomOperation, events.ReplyLiveRoomOperation, events.CmdLiveRoomOperationData{
				Identifier: room.LiveRoom.Identifier(),
				SetConnect: true,
			})
			if err != nil {
				gctx.Logger.Errorf("failed to connect to room %s", room.LiveRoom.Identifier())
				gutil.RunInFyneThread(currentRoomView.connectBtn.Enable)
				return
			}
			if resp.Data.(events.ReplyLiveRoomOperationData).Err != nil {
				err = resp.Data.(events.ReplyLiveRoomOperationData).Err
			}
			if err != nil {
				// todo: show error
			}
			gutil.RunInFyneThread(currentRoomView.connectBtn.Enable)
		}()
	})
	currentRoomView.disConnectBtn = widget.NewButton(i18n.T("gui.room.btn.disconnect"), func() {
		room, ok := getCurrentRoom()
		if !ok {
			return
		}
		gctx.Logger.Infof("disconnect to room %s", room.LiveRoom.Identifier())
		currentRoomView.disConnectBtn.Disable()
		go func() {
			resp, err := global.EventBus.Call(events.CmdLiveRoomOperation, events.ReplyLiveRoomOperation, events.CmdLiveRoomOperationData{
				Identifier: room.LiveRoom.Identifier(),
				SetConnect: false,
			})
			if err != nil {
				gctx.Logger.Errorf("failed to disconnect to room %s", room.LiveRoom.Identifier())
				gutil.RunInFyneThread(currentRoomView.disConnectBtn.Enable)
				return
			}
			if resp.Data.(events.ReplyLiveRoomOperationData).Err != nil {
				err = resp.Data.(events.ReplyLiveRoomOperationData).Err
			}
			if err != nil {
				// todo: show error
			}
			gutil.RunInFyneThread(currentRoomView.disConnectBtn.Enable)
		}()
	})

	currentRoomView.status = widget.NewLabel(i18n.T("gui.room.waiting"))
	currentRoomView.roomTitle = widget.NewLabel("")
	currentRoomView.roomID = widget.NewLabel("")

	currentRoomView.autoConnect = widget.NewCheck(i18n.T("gui.room.check.autoconnect"), func(b bool) {
		room, ok := getCurrentRoom()
		if !ok {
			return
		}
		gctx.Logger.Infof("Change room %s autoconnect to %v", room.LiveRoom.Identifier(), b)
		_ = global.EventBus.PublishToChannel(gctx.EventChannel,
			events.CmdLiveRoomConfigChange,
			events.CmdLiveRoomConfigChangeData{
				Identifier: room.LiveRoom.Identifier(),
				Config: model.LiveRoomConfig{
					AutoConnect: b,
				},
			})
		return
	})
	return container.NewVBox(
		currentRoomView.roomTitle,
		currentRoomView.roomID,
		currentRoomView.status,
		container.NewHBox(widget.NewLabel(i18n.T("gui.room.check.autoconnect")), currentRoomView.autoConnect),
		container.NewHBox(currentRoomView.connectBtn, currentRoomView.disConnectBtn),
	)
}
