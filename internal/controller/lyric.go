package controller

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
	"github.com/AynaLivePlayer/miaosic"
)

type lyricLoader struct {
	Lyric     miaosic.Lyrics
	prev      float64
	prevIndex int
}

var lyricManager = &lyricLoader{}

func createLyricLoader() {
	log := global.Logger.WithPrefix("LyricLoader")
	global.EventManager.RegisterA(events.PlayerPlayingUpdate, "internal.lyric.update", func(event *event.Event) {
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
			log.Errorf("failed to get lyric for %s (%s): %s", data.Media.Info.Title, data.Media.Info.Meta.ID(), err)
		}
		global.EventManager.CallA(events.PlayerLyricReload, events.PlayerLyricReloadEvent{
			Lyrics: lyricManager.Lyric,
		})
	})
	global.EventManager.RegisterA(events.PlayerPropertyTimePosUpdate, "internal.lyric.update_current", func(event *event.Event) {
		time := event.Data.(events.PlayerPropertyTimePosUpdateEvent).TimePos
		idx := lyricManager.Lyric.FindIndex(time)
		if idx == lyricManager.prevIndex {
			return
		}
		lyricManager.prevIndex = idx
		global.EventManager.CallA(
			events.PlayerLyricPosUpdate,
			events.PlayerLyricPosUpdateEvent{
				CurrentIndex: idx,
				Time:         time,
				CurrentLine:  lyricManager.Lyric.Find(time),
				Total:        len(lyricManager.Lyric.Content),
			})
		return
	})
	global.EventManager.RegisterA(events.PlayerLyricRequestCmd, "internal.lyric.request", func(event *event.Event) {
		global.EventManager.CallA(events.PlayerLyricReload, events.PlayerLyricReloadEvent{
			Lyrics: lyricManager.Lyric,
		})
	})
}
