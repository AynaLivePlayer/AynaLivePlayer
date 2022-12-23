package core

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/model"
	"errors"
	"strings"
)

type coreLiveRoom struct {
	model.LiveRoom
	client liveclient.LiveClient
}

func (r *coreLiveRoom) Model() *model.LiveRoom {
	return &r.LiveRoom
}

func (r *coreLiveRoom) Status() bool {
	return r.client.Status()
}

func (r *coreLiveRoom) EventManager() *event.Manager {
	return r.client.EventManager()
}

func (r *coreLiveRoom) init(msgHandler event.HandlerFunc) (err error) {
	if r.client != nil {
		return nil
	}
	r.client, err = liveclient.NewLiveClient(r.ClientName, r.ID)
	if err != nil {
		return
	}
	r.client.EventManager().RegisterA(
		liveclient.EventMessageReceive,
		"controller.danmu.command",
		msgHandler)
	return nil
}

type LiveRoomController struct {
	LiveRoomPath  string
	liveRooms     []*coreLiveRoom
	danmuCommands []controller.DanmuCommandExecutor
}

func NewLiveRoomController() controller.ILiveRoomController {
	lr := &LiveRoomController{
		LiveRoomPath: "liverooms.json",
		liveRooms: []*coreLiveRoom{
			{LiveRoom: model.LiveRoom{
				ClientName:  "bilibili",
				ID:          "9076804",
				AutoConnect: false,
			}},
			{LiveRoom: model.LiveRoom{
				ClientName:  "bilibili",
				ID:          "3819533",
				AutoConnect: false,
			}},
		},
		danmuCommands: make([]controller.DanmuCommandExecutor, 0),
	}
	config.LoadConfig(lr)
	lr.initialize()
	return lr
}

func (lr *LiveRoomController) danmuCommandHandler(event *event.Event) {
	danmu := event.Data.(*liveclient.DanmuMessage)
	args := strings.Split(danmu.Message, " ")
	if len(args[0]) == 0 {
		return
	}
	for _, cmd := range lr.danmuCommands {
		if cmd.Match(args[0]) {
			cmd.Execute(args[0], args[1:], danmu)
		}
	}
}

func (lr *LiveRoomController) initialize() {
	for i := 0; i < len(lr.liveRooms); i++ {
		if lr.liveRooms[i].client == nil {
			_ = lr.liveRooms[i].init(lr.danmuCommandHandler)
		}
	}
	go func() {
		for i := 0; i < len(lr.liveRooms); i++ {
			if lr.liveRooms[i].AutoConnect {
				lr.liveRooms[i].client.Connect()
			}
		}
	}()
}

func (lr *LiveRoomController) Name() string {
	return "LiveRooms"
}

func (lr *LiveRoomController) Size() int {
	return len(lr.liveRooms)
}

func (lr *LiveRoomController) OnLoad() {
	rooms := make([]model.LiveRoom, 0)
	_ = config.LoadJson(lr.LiveRoomPath, &lr.liveRooms)
	if len(rooms) == 0 {
		return
	}
	lr.liveRooms = make([]*coreLiveRoom, len(rooms))
	for i := 0; i < len(rooms); i++ {
		lr.liveRooms[i] = &coreLiveRoom{LiveRoom: rooms[i]}
	}
}

func (lr *LiveRoomController) OnSave() {
	rooms := make([]model.LiveRoom, len(lr.liveRooms))
	for i := 0; i < len(lr.liveRooms); i++ {
		rooms[i] = lr.liveRooms[i].LiveRoom
	}
	_ = config.SaveJson(lr.LiveRoomPath, &rooms)
}

func (lr *LiveRoomController) Get(index int) controller.ILiveRoom {
	if index < 0 || index >= len(lr.liveRooms) {
		return nil
	}
	return lr.liveRooms[index]
}

func (lr *LiveRoomController) GetRoomStatus(index int) bool {
	if index < 0 || index >= len(lr.liveRooms) {
		return false
	}
	return lr.liveRooms[index].client.Status()
}

func (lr *LiveRoomController) Connect(index int) error {
	lg.Infof("[LiveRooms] Try to start LiveRooms.index=%d", index)
	if index < 0 || index >= len(lr.liveRooms) {
		lg.Errorf("[LiveRooms] LiveRooms.index=%d not found", index)
		return errors.New("index out of range")
	}
	lr.liveRooms[index].client.Connect()
	return nil
}

func (lr *LiveRoomController) Disconnect(index int) error {
	lg.Infof("[LiveRooms] Try to Disconnect LiveRooms.index=%d", index)
	if index < 0 || index >= len(lr.liveRooms) {
		lg.Errorf("[LiveRooms] LiveRooms.index=%d not found", index)
		return errors.New("index out of range")
	}
	lr.liveRooms[index].client.Disconnect()
	return nil
}

func (lr *LiveRoomController) AddRoom(clientName, roomId string) (*model.LiveRoom, error) {
	rm := &coreLiveRoom{
		LiveRoom: model.LiveRoom{
			ClientName:  clientName,
			ID:          roomId,
			AutoConnect: false,
		},
	}
	lg.Infof("[LiveRooms] add live room %s", &rm.LiveRoom)
	err := rm.init(lr.danmuCommandHandler)
	if err != nil {
		return nil, err
	}
	lg.Infof("[LiveRooms] %s init failed: %s", &rm.LiveRoom, err)
	if err != nil {
		return nil, err
	}
	lr.liveRooms = append(lr.liveRooms, rm)
	return &rm.LiveRoom, nil
}

func (lr *LiveRoomController) DeleteRoom(index int) error {
	lg.Infof("Try to remove LiveRooms.index=%d", index)
	if index < 0 || index >= len(lr.liveRooms) {
		lg.Warnf("LiveRooms.index=%d not found", index)
		return errors.New("index out of range")
	}
	if len(lr.liveRooms) == 1 {
		return errors.New("can't delete last room")
	}
	_ = lr.liveRooms[index].client.Disconnect()
	lr.liveRooms[index].EventManager().UnregisterAll()
	lr.liveRooms = append(lr.liveRooms[:index], lr.liveRooms[index+1:]...)
	return nil
}

func (lr *LiveRoomController) AddDanmuCommand(executor controller.DanmuCommandExecutor) {
	lr.danmuCommands = append(lr.danmuCommands, executor)
}
