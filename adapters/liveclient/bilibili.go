package liveclient

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"errors"
	"github.com/aynakeya/blivedm"
	"strconv"
	"time"
)

type Bilibili struct {
	client       *blivedm.BLiveWsClient
	eventManager *event.Manager
	roomName     string
	status       bool
	log          adapter.ILogger
}

func BilibiliCtor(id string, em *event.Manager, log adapter.ILogger) (adapter.LiveClient, error) {
	room, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("room id for bilibili should be a integer")
	}
	return NewBilibili(room, em, log), nil
}

func NewBilibili(roomId int, em *event.Manager, log adapter.ILogger) adapter.LiveClient {
	cl := &Bilibili{
		client:       &blivedm.BLiveWsClient{ShortId: roomId, Account: blivedm.DanmuAccount{UID: 0}, HearbeatInterval: 10 * time.Second},
		eventManager: em,
		roomName:     "",
		log:          log,
	}
	cl.client.OnDisconnect = func(client *blivedm.BLiveWsClient) {
		cl.log.Warn("[Bilibili LiveChatSDK] disconnect from websocket connection, maybe try reconnect")
		cl.status = false
		cl.eventManager.CallA(events.LiveRoomStatusChange, events.StatusChangeEvent{Connected: false, Client: cl})
	}
	cl.client.RegHandler(blivedm.CmdDanmaku, cl.handleMsg)
	return cl
}

func (b *Bilibili) ClientName() string {
	return "bilibili"
}

func (b *Bilibili) RoomName() string {
	return b.roomName
}

func (b *Bilibili) Status() bool {
	return b.status
}

func (b *Bilibili) EventManager() *event.Manager {
	return b.eventManager
}

func (b *Bilibili) Connect() bool {
	if b.status {
		return true
	}
	b.log.Info("[Bilibili LiveChatSDK] Trying Connect Danmu Server")
	if b.client.InitRoom() && b.client.ConnectDanmuServer() {
		b.roomName = b.client.RoomInfo.Title
		b.status = true
		b.eventManager.CallA(events.LiveRoomStatusChange, events.StatusChangeEvent{Connected: true, Client: b})
		b.log.Info("[Bilibili LiveChatSDK] Connect Success")
		return true
	}
	b.log.Info("[Bilibili LiveChatSDK] Connect Failed")
	return false
}

func (b *Bilibili) Disconnect() bool {
	b.log.Info("[Bilibili LiveChatSDK] Disconnect from danmu server")
	if b.client == nil {
		return true
	}
	b.client.Disconnect()
	b.eventManager.CallA(events.LiveRoomStatusChange, events.StatusChangeEvent{Connected: false, Client: b})
	return true
}

func (b *Bilibili) handleMsg(context *blivedm.Context) {
	msg, ok := context.ToDanmakuMessage()
	if !ok {
		b.log.Warn("[Bilibili LiveChatSDK] handle message failed, can't convert context to danmu message")
		return
	}
	dmsg := model.DanmuMessage{
		User: model.DanmuUser{
			Uid:      strconv.FormatInt(msg.Uid, 10),
			Username: msg.Uname,
			Medal: model.UserMedal{
				Name:  msg.MedalName,
				Level: int(msg.MedalLevel),
			},
			Admin:     msg.Admin,
			Privilege: int(msg.PrivilegeType),
		},
		Message: msg.Msg,
	}
	b.log.Debug("[Bilibili LiveChatSDK] receive message", dmsg)
	go func() {
		b.eventManager.Call(&event.Event{
			Id:        events.LiveRoomMessageReceive,
			Cancelled: false,
			Data:      &dmsg,
		})
	}()
}
