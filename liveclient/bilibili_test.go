package liveclient

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/logger"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestBilibili_Client(t *testing.T) {
	logger.Logger.SetLevel(logrus.DebugLevel)
	lc := NewBilibili(7777, event.NewManger())
	//lc := NewBilibili(8524916587)
	lc.Handler().Register(&event.Handler{
		EventId: EventMessageReceive,
		Name:    "test.receivemsg",
		Handler: func(event *event.Event) {
			fmt.Println(event.Data.(*DanmuMessage).Message)
		},
	})
	lc.Connect()
	time.Sleep(time.Second * 60)
	lc.Disconnect()
}
