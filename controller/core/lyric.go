package core

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/model"
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
		model.EventLyricReload,
		model.LyricReloadEvent{
			Lyrics: l.Lyric,
		})
}

func (l *LyricLoader) Update(time float64) {
	lrc := l.Lyric.Find(time)
	if lrc == nil {
		return
	}
	if l.prev == lrc.Time {
		return
	}
	l.prev = lrc.Time
	l.Handler.CallA(
		model.EventLyricUpdate,
		model.LyricUpdateEvent{
			Lyrics: l.Lyric,
			Time:   time,
			Lyric:  lrc,
		})
	return
}
