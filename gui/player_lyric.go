package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/AynaLivePlayer/miaosic"
)

func createLyricObj(lyric *miaosic.Lyrics) []fyne.CanvasObject {
	lrcs := make([]fyne.CanvasObject, len(lyric.Content))
	for i := 0; i < len(lrcs); i++ {
		lr := widget.NewLabelWithStyle(
			lyric.Content[i].Lyric,
			fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
		//lr.Wrapping = fyne.TextWrapWord
		// todo fix fyne bug
		lr.Wrapping = fyne.TextWrapBreak
		lrcs[i] = lr
	}
	return lrcs
}

func createLyricWindow() fyne.Window {
	// create widgets
	w := App.NewWindow(i18n.T("gui.lyric.title"))
	currentLrc := newLabelWithWrapping("", fyne.TextWrapBreak)
	currentLrc.Alignment = fyne.TextAlignCenter
	fullLrc := container.NewVBox()
	lrcWindow := container.NewVScroll(fullLrc)
	prevIndex := 0
	w.SetContent(container.NewBorder(nil,
		container.NewVBox(widget.NewSeparator(), currentLrc),
		nil, nil,
		lrcWindow))
	w.Resize(fyne.NewSize(360, 540))
	w.CenterOnScreen()

	// register handlers
	global.EventBus.Subscribe("",
		events.PlayerLyricPosUpdate, "player.lyric.current_lyric", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			e := event.Data.(events.PlayerLyricPosUpdateEvent)
			logger.Debug("lyric update", e)
			if prevIndex >= len(fullLrc.Objects) || e.CurrentIndex >= len(fullLrc.Objects) {
				// fix race condition
				return
			}
			if e.CurrentIndex == -1 {
				currentLrc.SetText("")
				return
			}
			fullLrc.Objects[prevIndex].(*widget.Label).TextStyle.Bold = false
			fullLrc.Objects[prevIndex].Refresh()
			fullLrc.Objects[e.CurrentIndex].(*widget.Label).TextStyle.Bold = true
			fullLrc.Objects[e.CurrentIndex].Refresh()
			prevIndex = e.CurrentIndex
			currentLrc.SetText(e.CurrentLine.Lyric)
			lrcWindow.Scrolled(&fyne.ScrollEvent{
				Scrolled: fyne.Delta{
					DX: 0,
					DY: lrcWindow.Offset.Y - float32(e.CurrentIndex-2)/float32(e.Total)*lrcWindow.Content.Size().Height,
				},
			})
			fullLrc.Refresh()
		}))

	global.EventBus.Subscribe("", events.PlayerLyricReload, "player.lyric.current_lyric", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		e := event.Data.(events.PlayerLyricReloadEvent)
		fullLrc.Objects = createLyricObj(&e.Lyrics)
		lrcWindow.Refresh()
	}))

	_ = global.EventBus.Publish(events.PlayerLyricRequestCmd, events.PlayerLyricRequestCmdEvent{})

	w.SetOnClosed(func() {
		global.EventBus.Unsubscribe(events.PlayerLyricReload, "player.lyric.current_lyric")
		PlayController.LrcWindowOpen = false
	})
	return w
}
