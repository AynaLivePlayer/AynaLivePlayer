package internal

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
)

type LyricLoader struct {
	Lyric   *model.Lyric
	Handler *event.Manager
	prev    float64
}

func NewLyricLoader() *LyricLoader {
	return &LyricLoader{
		Lyric:   model.LoadLyric(""),
		Handler: event.MainManager.NewChildManager(),
		prev:    -1,
	}
}

func (l *LyricLoader) EventManager() *event.Manager {
	return l.Handler
}

func (l *LyricLoader) Get() *model.Lyric {
	return l.Lyric
}

func (l *LyricLoader) Reload(lyric string) {
	l.Lyric = model.LoadLyric(lyric)
	l.Handler.CallA(
		events.EventLyricReload,
		events.LyricReloadEvent{
			Lyrics: l.Lyric,
		})
}

func (l *LyricLoader) Update(time float64) {
	lrc := l.Lyric.FindContext(time, 1, 3)
	if lrc == nil {
		return
	}
	if l.prev == lrc.Now.Time {
		return
	}
	l.prev = lrc.Now.Time
	l.Handler.CallA(
		events.EventLyricUpdate,
		events.LyricUpdateEvent{
			Lyrics: l.Lyric,
			Time:   time,
			Lyric:  lrc,
		})
	return
}
