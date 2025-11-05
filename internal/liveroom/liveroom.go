package liveroom

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/logger"
	"errors"
	liveroomsdk "github.com/AynaLivePlayer/liveroom-sdk"
	"github.com/AynaLivePlayer/liveroom-sdk/provider/openblive"
	"github.com/AynaLivePlayer/liveroom-sdk/provider/webdm"
	"sort"
)

type liveroom struct {
	room  liveroomsdk.ILiveRoom
	model model.LiveRoom
}

var liveRooms = map[string]*liveroom{}
var log logger.ILogger

func Initialize() {
	log = global.Logger.WithPrefix("LiveRoom")
	config.LoadConfig(cfg)

	liveroomsdk.RegisterProvider(openblive.NewOpenBLiveClientProvider(cfg.ApiServer, 1661006726438))
	// ignore web danmu client
	liveroomsdk.RegisterProvider(webdm.NewWebDanmuClientProvider(cfg.ApiServer))

	liveRooms = make(map[string]*liveroom, 0)

	callEvents()
	registerHandlers()
}

func StopAndSave() {
	log.Infof("Stop and save live rooms")
	for _, r := range liveRooms {
		log.Infof("Disconnect room %s: %v", r.room.Config().Identifier(), r.room.Disconnect())
	}
	liveroomConfigs := make([]model.LiveRoom, 0)
	for _, r := range liveRooms {
		liveroomConfigs = append(liveroomConfigs, r.model)
	}
	cfg.liveRooms = liveroomConfigs
}

func addLiveRoom(roomModel model.LiveRoom) {
	log.Info("Add live room")
	room, err := liveroomsdk.CreateLiveRoom(roomModel.LiveRoom)
	// handle failed to create liveroom
	if err != nil {
		log.Errorf("Create live room failed: %s", err)
		_ = global.EventBus.Publish(
			events.ErrorUpdate, events.ErrorUpdateEvent{
				Error: err,
			})
		return
	}
	if _, ok := liveRooms[room.Config().Identifier()]; ok {
		log.Errorf("fail to add, room %s already exists", room.Config().Identifier())
		_ = global.EventBus.Publish(
			events.ErrorUpdate, events.ErrorUpdateEvent{
				Error: errors.New("room already exists"),
			})
		return
	}
	if err != nil {
		log.Errorf("Create live room failed: %s", err)
		_ = global.EventBus.Publish(
			events.ErrorUpdate, events.ErrorUpdateEvent{
				Error: err,
			})
		return
	}
	lr := &liveroom{
		room:  room,
		model: roomModel,
	}
	liveRooms[room.Config().Identifier()] = lr

	room.OnStatusChange(func(connected bool) {
		log.Infof("room %s status change to %t", room.Config().Identifier(), connected)
		lr.model.Status = connected
		sendRoomStatusUpdateEvent(room.Config().Identifier())
	})
	room.OnMessage(func(message *liveroomsdk.Message) {
		log.Debugf("room %s receive message: %s", room.Config().Identifier(), message.Message)
		_ = global.EventBus.Publish(
			events.LiveRoomMessageReceive,
			events.LiveRoomMessageReceiveEvent{
				Message: message,
			})
	})

	log.Infof("success add live room %s", room.Config().Identifier())
	sendRoomsUpdateEvent()
}

func registerHandlers() {
	global.EventBus.Subscribe("",
		events.CmdLiveRoomAdd, "internal.liveroom.add", func(event *eventbus.Event) {
			data := event.Data.(events.CmdLiveRoomAddData)
			addLiveRoom(model.LiveRoom{
				LiveRoom: liveroomsdk.LiveRoom{
					Provider: data.Provider,
					Room:     data.RoomKey,
				},
				Config: model.LiveRoomConfig{
					AutoConnect: false,
				},
				Title:  data.Title,
				Status: false,
			})
		})

	global.EventBus.Subscribe("",
		events.CmdLiveRoomRemove, "internal.liveroom.remove", func(event *eventbus.Event) {
			data := event.Data.(events.CmdLiveRoomRemoveData)
			room, ok := liveRooms[data.Identifier]
			if !ok {
				log.Errorf("remove room failed, room %s not found", data.Identifier)
				return

			}
			_ = room.room.Disconnect()
			room.room.OnStatusChange(nil)
			room.room.OnMessage(nil)
			delete(liveRooms, data.Identifier)
			log.Infof("success remove live room %s", data.Identifier)
			sendRoomsUpdateEvent()
		})

	global.EventBus.Subscribe("",
		events.CmdLiveRoomConfigChange, "internal.liveroom.config.change", func(event *eventbus.Event) {
			data := event.Data.(events.CmdLiveRoomConfigChangeData)
			if room, ok := liveRooms[data.Identifier]; ok {
				room.model.Config = data.Config
				sendRoomStatusUpdateEvent(data.Identifier)
			}
		})

	global.EventBus.Subscribe("",
		events.CmdLiveRoomOperation, "internal.liveroom.operation", func(event *eventbus.Event) {
			data := event.Data.(events.CmdLiveRoomOperationData)
			log.Infof("Live room operation SetConnect %v", data.SetConnect)
			room, ok := liveRooms[data.Identifier]
			if !ok {
				log.Errorf("Room %s not found", data.Identifier)
				return
			}
			var err error
			if data.SetConnect {
				err = room.room.Connect()
			} else {
				err = room.room.Disconnect()
			}
			if err != nil {
				log.Errorf("Room %s operation failed: %s", data.Identifier, err)
				_ = global.EventBus.Publish(
					events.ErrorUpdate, events.ErrorUpdateEvent{
						Error: err,
					})
			}
			_ = global.EventBus.Reply(event,
				events.ReplyLiveRoomOperation, events.ReplyLiveRoomOperationData{})
			sendRoomStatusUpdateEvent(data.Identifier)
		})
}

func sendRoomStatusUpdateEvent(roomId string) {
	room, ok := liveRooms[roomId]
	if !ok {
		log.Errorf("send room status update event failed, room %s not found", roomId)
		return
	}
	log.Infof("send room status update event, room %s", roomId)
	_ = global.EventBus.Publish(
		events.UpdateLiveRoomStatus,
		events.UpdateLiveRoomStatusData{
			Room: room.model,
		})
}

func sendRoomsUpdateEvent() {
	rooms := make([]model.LiveRoom, 0)
	for _, r := range liveRooms {
		rooms = append(rooms, r.model)
	}
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].LiveRoom.Identifier() < rooms[j].LiveRoom.Identifier()
	})
	_ = global.EventBus.Publish(
		events.UpdateLiveRoomRooms,
		events.UpdateLiveRoomRoomsData{
			Rooms: rooms,
		})
}

func callEvents() {
	providers := liveroomsdk.ListAvailableProviders()
	providerInfo := make([]model.LiveRoomProviderInfo, 0)
	for _, p := range providers {
		provider, _ := liveroomsdk.GetProvider(p)
		providerInfo = append(providerInfo, model.LiveRoomProviderInfo{
			Name:        provider.GetName(),
			Description: provider.GetDescription(),
		})
	}
	for _, roomCfg := range cfg.liveRooms {
		addLiveRoom(roomCfg)
	}
	_ = global.EventBus.Publish(
		events.LiveRoomProviderUpdate,
		events.LiveRoomProviderUpdateEvent{
			Providers: providerInfo,
		})
	sendRoomsUpdateEvent()
	for _, r := range liveRooms {
		if r.model.Config.AutoConnect {
			_ = global.EventBus.Publish(
				events.CmdLiveRoomOperation,
				events.CmdLiveRoomOperationData{
					Identifier: r.room.Config().Identifier(),
					SetConnect: true,
				})
		}
	}
}
