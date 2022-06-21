package liveclient

import (
	"AynaLivePlayer/event"
	"AynaLivePlayer/logger"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestBilibili_Client(t *testing.T) {
	logger.Logger.SetLevel(logrus.DebugLevel)
	lc := NewBilibili(7777)
	//lc := NewBilibili(8524916587)
	lc.Handler().Register(&event.EventHandler{
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
