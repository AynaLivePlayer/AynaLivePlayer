package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/event"
	"AynaLivePlayer/liveclient"
	"errors"
	"fmt"
)

var LiveRoomManager = &LiveRooms{
	LiveRoomPath: "liverooms.json",
	LiveRooms: []*LiveRoom{
		{
			ClientName:  "bilibili",
			ID:          "9076804",
			AutoConnect: false,
		},
		{
			ClientName:  "bilibili",
			ID:          "3819533",
			AutoConnect: false,
		},
	},
}

type LiveRooms struct {
	LiveRoomPath string
	LiveRooms    []*LiveRoom `ini:"-"`
}

func (lr *LiveRooms) Name() string {
	return "LiveRoom"
}

func (lr *LiveRooms) Size() int {
	return len(lr.LiveRooms)
}

func (lr *LiveRooms) OnLoad() {
	_ = config.LoadJson(lr.LiveRoomPath, &lr.LiveRooms)
}

func (lr *LiveRooms) OnSave() {
	_ = config.SaveJson(lr.LiveRoomPath, &lr.LiveRooms)
}

func (lr *LiveRooms) InitializeRooms() {
	for i := 0; i < len(lr.LiveRooms); i++ {
		if lr.LiveRooms[i].client == nil {
			lr.LiveRooms[i].Init()
		}
	}
	go func() {
		for i := 0; i < len(lr.LiveRooms); i++ {
			if lr.LiveRooms[i].AutoConnect {
				go lr.LiveRooms[i].Connect()
			}
		}
	}()
}

func (lr *LiveRooms) GetRoom(index int) *LiveRoom {
	if index < 0 || index >= len(lr.LiveRooms) {
		return nil
	}
	return lr.LiveRooms[index]
}

func (lr *LiveRooms) AddRoom(clientName, roomId string) (*LiveRoom, error) {
	l.Infof("add live client (%s) for %s", clientName, roomId)
	rm := &LiveRoom{
		ClientName:  clientName,
		ID:          roomId,
		AutoConnect: false,
	}
	err := rm.Init()
	l.Infof("live client (%s) %s init failed", clientName, roomId)
	if err != nil {
		return nil, err
	}
	lr.LiveRooms = append(lr.LiveRooms, rm)
	return rm, nil
}

func (lr *LiveRooms) ConnectRoom(index int) error {
	l.Infof("Try to start LiveRoom.index=%d", index)
	if index < 0 || index >= len(lr.LiveRooms) {
		l.Warnf("LiveRoom.index=%d not found", index)
		return errors.New("index out of range")
	}
	lr.LiveRooms[index].client.Connect()
	return nil
}

func (lr *LiveRooms) DisconnectRoom(index int) error {
	l.Infof("Try to Disconnect LiveRoom.index=%d", index)
	if index < 0 || index >= len(lr.LiveRooms) {
		l.Warnf("LiveRoom.index=%d not found", index)
		return errors.New("index out of range")
	}
	lr.LiveRooms[index].client.Disconnect()
	return nil
}

func (lr *LiveRooms) DeleteRoom(index int) error {
	l.Infof("Try to remove LiveRoom.index=%d", index)
	if index < 0 || index >= len(lr.LiveRooms) {
		l.Warnf("LiveRoom.index=%d not found", index)
		return errors.New("index out of range")
	}
	if len(lr.LiveRooms) == 1 {
		return errors.New("can't delete last room")
	}
	lr.LiveRooms[index].client.Handler().UnregisterAll()
	_ = lr.LiveRooms[index].Disconnect()
	lr.LiveRooms = append(lr.LiveRooms[:index], lr.LiveRooms[index+1:]...)
	return nil
}

type LiveRoom struct {
	ClientName  string
	ID          string
	AutoConnect bool
	client      liveclient.LiveClient
}

func (r *LiveRoom) Init() (err error) {
	if r.client != nil {
		return nil
	}
	r.client, err = liveclient.NewLiveClient(r.ClientName, r.ID)
	if err != nil {
		return
	}
	r.client.Handler().Register(&event.EventHandler{
		EventId: liveclient.EventMessageReceive,
		Name:    "controller.commandexecutor",
		Handler: danmuCommandHandler,
	})
	r.client.Handler().RegisterA(
		liveclient.EventMessageReceive,
		"controller.danmu.handler",
		danmuHandler)
	return nil
}

func (r *LiveRoom) Connect() error {
	if r.client == nil {
		return errors.New("client hasn't initialized yet")
	}
	r.client.Connect()
	return nil
}

func (r *LiveRoom) Disconnect() error {
	if r.client == nil {
		return errors.New("client hasn't initialized yet")
	}
	r.client.Disconnect()
	return nil
}

func (r *LiveRoom) Title() string {
	return fmt.Sprintf("%s-%s", r.ClientName, r.ID)
}

func (r *LiveRoom) Client() liveclient.LiveClient {
	return r.client
}
