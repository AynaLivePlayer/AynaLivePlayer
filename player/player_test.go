package player

import (
	"AynaLivePlayer/model"
	"fmt"
	"github.com/aynakeya/go-mpv"
	"testing"
	"time"
)

func TestPlayer(t *testing.T) {
	player := NewPlayer()
	player.Start()
	defer player.Stop()

	player.ObserveProperty("time-pos", func(property *mpv.EventProperty) {
		fmt.Println(1, property.Data)
	})
	player.ObserveProperty("percent-pos", func(property *mpv.EventProperty) {
		fmt.Println(2, property.Data)
	})
	player.Play(&model.Media{
		Url: "https://ia600809.us.archive.org/19/items/VillagePeopleYMCAOFFICIALMusicVideo1978/Village%20People%20-%20YMCA%20OFFICIAL%20Music%20Video%201978.mp4",
	})
	time.Sleep(time.Second * 15)
}
