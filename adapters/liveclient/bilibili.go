package liveclient

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"encoding/json"
	"errors"
	"github.com/aynakeya/blivedm"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
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
	if !b.client.GetRoomInfo() {
		b.log.Info("[Bilibili LiveChatSDK] Connect Failed")
		return false
	}
	resp, err := resty.New().R().
		SetQueryParam("room_id", strconv.Itoa(b.client.RoomId)).
		Get("https://scene.aynakeya.com:3000/bilisrv/dminfo")
	if err != nil {
		b.log.Info("[Bilibili LiveChatSDK] Connect Failed")
		return false
	}
	gjresult := gjson.Parse(resp.String())
	if gjresult.Get("code").Int() != 0 {
		b.log.Info("[Bilibili LiveChatSDK] Connect Failed")
		return false
	}
	b.client.DanmuInfo = blivedm.DanmuInfoData{
		Token: gjresult.Get("data.token").String(),
	}
	b.client.Account.UID = int(gjresult.Get("data.uid").Int())
	if err := json.Unmarshal([]byte(gjresult.Get("data.host_list").String()), &b.client.DanmuInfo.HostList); err != nil {
		return false
	}
	if b.client.ConnectDanmuServer() {
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
