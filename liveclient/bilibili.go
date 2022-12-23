package liveclient

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/logger"
	"errors"
	"github.com/aynakeya/blivedm"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Bilibili struct {
	client       *blivedm.BLiveWsClient
	eventManager *event.Manager
	roomName     string
	status       bool
}

func init() {
	LiveClients["bilibili"] = func(id string) (LiveClient, error) {
		room, err := strconv.Atoi(id)
		if err != nil {
			return nil, errors.New("room id for bilibili should be a integer")
		}
		return NewBilibili(room), nil
	}
}

func NewBilibili(roomId int) LiveClient {
	cl := &Bilibili{
		client:       &blivedm.BLiveWsClient{ShortId: roomId, Account: blivedm.DanmuAccount{UID: 0}, HearbeatInterval: 10 * time.Second},
		eventManager: event.MainManager.NewChildManager(),
		roomName:     "Unknown",
	}
	cl.client.OnDisconnect = func(client *blivedm.BLiveWsClient) {
		cl.l().Warn("disconnect from websocket connection, maybe try reconnect")
		cl.status = false
		cl.eventManager.CallA(EventStatusChange, StatusChangeEvent{Connected: false, Client: cl})
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
	b.l().Info("Trying Connect Danmu Server")
	if b.client.InitRoom() && b.client.ConnectDanmuServer() {
		b.roomName = b.client.RoomInfo.Title
		b.status = true
		b.eventManager.CallA(EventStatusChange, StatusChangeEvent{Connected: true, Client: b})
		b.l().Info("Connect Success")
		return true
	}
	b.l().Info("Connect Failed")
	return false
}

func (b *Bilibili) Disconnect() bool {
	b.l().Info("Disconnect from danmu server")
	if b.client == nil {
		return true
	}
	b.client.Disconnect()
	b.eventManager.CallA(EventStatusChange, StatusChangeEvent{Connected: false, Client: b})
	return true
}

func (b *Bilibili) l() *logrus.Entry {
	return logger.Logger.WithFields(logrus.Fields{
		"Module":     MODULE_NAME,
		"ClientName": b.ClientName(),
		"RoomId":     b.client.ShortId,
	})
}

func (b *Bilibili) handleMsg(context *blivedm.Context) {
	msg, ok := context.ToDanmakuMessage()
	if !ok {
		b.l().Warn("handle message failed, can't convert context to danmu message")
		return
	}
	dmsg := DanmuMessage{
		User: DanmuUser{
			Uid:      strconv.FormatInt(msg.Uid, 10),
			Username: msg.Uname,
			Medal: UserMedal{
				Name:  msg.MedalName,
				Level: int(msg.MedalLevel),
			},
			Admin:     msg.Admin,
			Privilege: int(msg.PrivilegeType),
		},
		Message: msg.Msg,
	}
	b.l().Debug("receive message", dmsg)
	go func() {
		b.eventManager.Call(&event.Event{
			Id:        EventMessageReceive,
			Cancelled: false,
			Data:      &dmsg,
		})
	}()
}
