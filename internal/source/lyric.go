package source

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"github.com/AynaLivePlayer/miaosic"
)

type lyricLoader struct {
	Lyric     miaosic.Lyrics
	prev      float64
	prevIndex int
}

var lyricManager = &lyricLoader{}

func createLyricLoader() {
	var err error
	err = global.EventBus.Subscribe("", events.PlayerPlayingUpdate, "internal.lyric.update", func(event *eventbus.Event) {
		data := event.Data.(events.PlayerPlayingUpdateEvent)
		if data.Removed {
			log.Debugf("current media removed, clear lyric")
			lyricManager.Lyric = miaosic.ParseLyrics("", "")
			return
		}
		log.Infof("update lyric for %s", data.Media.Info.Title)
		lyric, err := miaosic.GetMediaLyric(data.Media.Info.Meta)
		if err == nil && len(lyric) > 0 {
			lyricManager.Lyric = lyric[0]
		} else {
			lyricManager.Lyric = miaosic.ParseLyrics("", "")
			log.Errorf("failed to get lyric for %s (%s): %s", data.Media.Info.Title, data.Media.Info.Meta.ID(), err)
		}
		_ = global.EventBus.Publish(events.UpdateCurrentLyric, events.UpdateCurrentLyricData{
			Lyrics: lyricManager.Lyric,
		})
	})
	if err != nil {
		log.ErrorW("Subscribe player playing update event failed", "error", err)
	}
	err = global.EventBus.Subscribe("", events.PlayerPropertyTimePosUpdate, "internal.lyric.update_current", func(event *eventbus.Event) {
		time := event.Data.(events.PlayerPropertyTimePosUpdateEvent).TimePos
		idx := lyricManager.Lyric.FindIndex(time)
		if idx == lyricManager.prevIndex {
			return
		}
		lyricManager.prevIndex = idx
		_ = global.EventBus.Publish(
			events.PlayerLyricPosUpdate,
			events.PlayerLyricPosUpdateEvent{
				CurrentIndex: idx,
				Time:         time,
				CurrentLine:  lyricManager.Lyric.Find(time),
				Total:        len(lyricManager.Lyric.Content),
			})
		return
	})
	if err != nil {
		log.ErrorW("Subscribe player time position update event failed", "error", err)
	}
	err = global.EventBus.Subscribe("", events.CmdGetCurrentLyric, "internal.lyric.request", func(event *eventbus.Event) {
		_ = global.EventBus.Reply(event, events.UpdateCurrentLyric, events.UpdateCurrentLyricData{
			Lyrics: lyricManager.Lyric,
		})
	})
	if err != nil {
		log.ErrorW("Subscribe player lyric request command event failed", "error", err)
	}
}
