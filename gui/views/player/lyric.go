package player

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component/lyrics"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/AynaLivePlayer/miaosic"
	"sync"
)

var lyricWindow fyne.Window = nil
var lyricViewer *lyrics.LyricsViewer = nil
var currLyrics []string
var currentLrcObj miaosic.Lyrics = miaosic.Lyrics{}
var lrcmux sync.RWMutex

func setupLyricViewer() {
	if lyricWindow != nil {
		return
	}
	lyricViewer = lyrics.NewLyricsViewer()
	lyricViewer.ActiveLyricPosition = lyrics.ActiveLyricPositionUpperMiddle
	lyricViewer.Alignment = fyne.TextAlignCenter
	lyricViewer.HoveredLyricColorName = theme.ColorNameDisabled
	lyricViewer.SetLyrics([]string{""}, true)
	lyricViewer.OnLyricTapped = func(lineNum int) {
		lineNum = lineNum - 1
		if lineNum < 0 {
			return
		}
		lrcmux.Lock()
		if lineNum >= len(currentLrcObj.Content) {
			lrcmux.Unlock()
			return
		}
		line := currentLrcObj.Content[lineNum]
		lrcmux.Unlock()
		_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
			Position: line.Time,
			Absolute: true,
		})
	}

	global.EventBus.Subscribe(gctx.EventChannel, events.UpdateCurrentLyric, "player.lyric.current_lyric", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		e := event.Data.(events.UpdateCurrentLyricData)
		tmpLyric := make([]string, 0)
		for _, l := range e.Lyrics.Content {
			tmpLyric = append(tmpLyric, l.Lyric)
		}
		// ensure at least one line
		if len(tmpLyric) == 0 {
			tmpLyric = append(tmpLyric, "")
		}
		lrcmux.Lock()
		currentLrcObj = event.Data.(events.UpdateCurrentLyricData).Lyrics
		currLyrics = tmpLyric
		lyricViewer.SetLyrics(currLyrics, true)
		lyricViewer.SetCurrentLine(0)
		lrcmux.Unlock()
	}))

	// register handlers
	global.EventBus.Subscribe(gctx.EventChannel,
		events.PlayerLyricPosUpdate, "player.lyric.lyric_pos_update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			e := event.Data.(events.PlayerLyricPosUpdateEvent)
			gctx.Logger.Debug("lyric update", e)
			lrcmux.Lock()
			if e.CurrentIndex >= len(currLyrics) {
				// fix race condition
				lrcmux.Unlock()
				return
			}
			index := 0
			if e.CurrentIndex != -1 {
				index = e.CurrentIndex
			}
			lyricViewer.SetCurrentLine(index + 1)
			lrcmux.Unlock()
		}))
}

func createLyricWindowV2() fyne.Window {
	// create widgets
	lyricWindow = gctx.Context.App.NewWindow(i18n.T("gui.lyric.title"))
	lyricWindow.SetContent(lyricViewer)
	lyricWindow.Resize(fyne.NewSize(360, 540))
	lyricWindow.CenterOnScreen()
	lyricWindow.SetOnClosed(func() {
		PlayController.LrcWindowOpen = false
	})
	return lyricWindow
}
