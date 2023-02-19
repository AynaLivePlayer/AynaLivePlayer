package internal

import (
	"AynaLivePlayer/adapters"
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"errors"
	"strings"
)

type liveRoomImpl struct {
	model.LiveRoom
	client adapter.LiveClient
}

func (r *liveRoomImpl) Status() bool {
	return r.client.Status()
}

func (r *liveRoomImpl) EventManager() *event.Manager {
	return r.client.EventManager()
}

func (r *liveRoomImpl) Model() *model.LiveRoom {
	return &r.LiveRoom
}

func (r *liveRoomImpl) Client() adapter.LiveClient {
	return r.client
}

func (r *liveRoomImpl) DisplayName() string {
	// todo need to be fixed
	r.LiveRoom.Title = r.client.RoomName()
	if r.LiveRoom.Title != "" {
		return r.LiveRoom.Title
	}
	return r.LiveRoom.Identifier()
}

func (r *liveRoomImpl) init(msgHandler event.HandlerFunc) (err error) {
	if r.client != nil {
		return nil
	}
	r.client, err = adapters.LiveClient.NewLiveClient(r.ClientName, r.ID)
	if err != nil {
		return
	}
	r.LiveRoom.Title = r.client.RoomName()
	r.client.EventManager().RegisterA(
		events.LiveRoomMessageReceive,
		"adapter.danmu.command",
		msgHandler)
	return nil
}

type LiveRoomController struct {
	LiveRoomPath  string
	liveRooms     []*liveRoomImpl
	danmuCommands []adapter.LiveRoomExecutor
	log           adapter.ILogger
}

func (lr *LiveRoomController) GetAllClientNames() []string {
	return adapters.LiveClient.GetAllClientNames()
}

func NewLiveRoomController(
	log adapter.ILogger,
) adapter.ILiveRoomController {
	lr := &LiveRoomController{
		LiveRoomPath: "liverooms.json",
		liveRooms: []*liveRoomImpl{
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
		danmuCommands: make([]adapter.LiveRoomExecutor, 0),
		log:           log,
	}
	config.LoadConfig(lr)
	lr.initialize()
	return lr
}

func (lr *LiveRoomController) danmuCommandHandler(event *event.Event) {
	danmu := event.Data.(*model.DanmuMessage)
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
	lr.liveRooms = make([]*liveRoomImpl, len(rooms))
	for i := 0; i < len(rooms); i++ {
		lr.liveRooms[i] = &liveRoomImpl{LiveRoom: rooms[i]}
	}
}

func (lr *LiveRoomController) OnSave() {
	rooms := make([]model.LiveRoom, len(lr.liveRooms))
	for i := 0; i < len(lr.liveRooms); i++ {
		rooms[i] = lr.liveRooms[i].LiveRoom
	}
	_ = config.SaveJson(lr.LiveRoomPath, &rooms)
}

func (lr *LiveRoomController) Get(index int) adapter.ILiveRoom {
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
	lr.log.Infof("[LiveRooms] Try to start LiveRooms.index=%d", index)
	if index < 0 || index >= len(lr.liveRooms) {
		lr.log.Errorf("[LiveRooms] LiveRooms.index=%d not found", index)
		return errors.New("index out of range")
	}
	lr.liveRooms[index].client.Connect()
	return nil
}

func (lr *LiveRoomController) Disconnect(index int) error {
	lr.log.Infof("[LiveRooms] Try to Disconnect LiveRooms.index=%d", index)
	if index < 0 || index >= len(lr.liveRooms) {
		lr.log.Errorf("[LiveRooms] LiveRooms.index=%d not found", index)
		return errors.New("index out of range")
	}
	lr.liveRooms[index].client.Disconnect()
	return nil
}

func (lr *LiveRoomController) AddRoom(clientName, roomId string) (*model.LiveRoom, error) {
	rm := &liveRoomImpl{
		LiveRoom: model.LiveRoom{
			ClientName:  clientName,
			ID:          roomId,
			AutoConnect: false,
		},
	}
	lr.log.Infof("[LiveRooms] add live room %s", &rm.LiveRoom)
	err := rm.init(lr.danmuCommandHandler)
	if err != nil {
		return nil, err
	}
	lr.log.Infof("[LiveRooms] %s init failed: %s", &rm.LiveRoom, err)
	if err != nil {
		return nil, err
	}
	lr.liveRooms = append(lr.liveRooms, rm)
	return &rm.LiveRoom, nil
}

func (lr *LiveRoomController) DeleteRoom(index int) error {
	lr.log.Infof("Try to remove LiveRooms.index=%d", index)
	if index < 0 || index >= len(lr.liveRooms) {
		lr.log.Warnf("LiveRooms.index=%d not found", index)
		return errors.New("index out of range")
	}
	if len(lr.liveRooms) == 1 {
		return errors.New("can't delete last room")
	}
	_ = lr.liveRooms[index].client.Disconnect()
	lr.liveRooms[index].Client().EventManager().UnregisterAll()
	lr.liveRooms = append(lr.liveRooms[:index], lr.liveRooms[index+1:]...)
	return nil
}

func (lr *LiveRoomController) AddDanmuCommand(executor adapter.LiveRoomExecutor) {
	lr.danmuCommands = append(lr.danmuCommands, executor)
}
