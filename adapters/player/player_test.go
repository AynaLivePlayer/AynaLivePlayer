package player

import (
	"AynaLivePlayer/adapters/logger"
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/model"
	"fmt"
	"testing"
	"time"
)

func TestPlayer(t *testing.T) {
	player := NewMpvPlayer(event.MainManager, &logger.EmptyLogger{})
	player.Start()
	defer player.Stop()

	player.ObserveProperty("time-pos", "testplayer.timepos", func(evnt *event.Event) {
		fmt.Println(1, evnt.Data)
	})
	player.ObserveProperty("percent-pos", "testplayer.percentpos", func(evnt *event.Event) {
		fmt.Println(2, evnt.Data)
	})
	player.Play(&model.Media{
		Url: "https://ia600809.us.archive.org/19/items/VillagePeopleYMCAOFFICIALMusicVideo1978/Village%20People%20-%20YMCA%20OFFICIAL%20Music%20Video%201978.mp4",
	})
	time.Sleep(time.Second * 15)
}
