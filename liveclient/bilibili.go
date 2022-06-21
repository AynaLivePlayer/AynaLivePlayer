package liveclient

import (
	"AynaLivePlayer/event"
	"AynaLivePlayer/logger"
	"github.com/aynakeya/blivedm"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Bilibili struct {
	client   *blivedm.BLiveWsClient
	handlers *event.Handler
}

func NewBilibili(roomId int) LiveClient {
	cl := &Bilibili{
		client:   &blivedm.BLiveWsClient{ShortId: roomId, Account: blivedm.DanmuAccount{UID: 0}, HearbeatInterval: 10 * time.Second},
		handlers: event.NewHandler(),
	}
	cl.client.OnDisconnect = func(client *blivedm.BLiveWsClient) {
		cl.l().Warn("disconnect from websocket connection, maybe try reconnect")
		cl.Handler().CallA(EventStatusChange, StatusChangeEvent{Connected: false, Client: cl})
	}
	cl.client.RegHandler(blivedm.CmdDanmaku, cl.handleMsg)
	return cl
}

func (b *Bilibili) ClientName() string {
	return "bilibili"
}

func (b *Bilibili) Handler() *event.Handler {
	return b.handlers
}

func (b *Bilibili) Connect() bool {
	b.l().Info("Trying Connect Danmu Server")
	if b.client.InitRoom() && b.client.ConnectDanmuServer() {
		b.Handler().CallA(EventStatusChange, StatusChangeEvent{Connected: true, Client: b})
		b.l().Info("Connect Success")
		return true
	}
	b.l().Info("Connect Failed")
	return false
}

func (b *Bilibili) Disconnect() bool {
	b.l().Info("Disconnect from danmu server")
	b.client.Disconnect()
	b.Handler().CallA(EventStatusChange, StatusChangeEvent{Connected: false, Client: b})
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
		b.handlers.Call(&event.Event{
			Id:        EventMessageReceive,
			Cancelled: false,
			Data:      &dmsg,
		})
	}()
}
