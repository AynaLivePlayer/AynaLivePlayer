package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createLyricObj(lyric *model.Lyric) []fyne.CanvasObject {
	lrcs := make([]fyne.CanvasObject, len(lyric.Lyrics))
	for i := 0; i < len(lrcs); i++ {
		lr := widget.NewLabelWithStyle(
			lyric.Lyrics[i].Lyric,
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
	w := App.NewWindow("Lyric")
	currentLrc := newLabelWithWrapping("", fyne.TextWrapBreak)
	currentLrc.Alignment = fyne.TextAlignCenter
	fullLrc := container.NewVBox(createLyricObj(API.PlayControl().GetLyric().Get())...)
	lrcWindow := container.NewVScroll(fullLrc)
	prevIndex := 0
	w.SetContent(container.NewBorder(nil,
		container.NewVBox(widget.NewSeparator(), currentLrc),
		nil, nil,
		lrcWindow))
	w.Resize(fyne.NewSize(360, 540))
	w.CenterOnScreen()

	// register handlers
	API.PlayControl().GetLyric().EventManager().RegisterA(
		events.EventLyricUpdate, "player.lyric.current_lyric", func(event *event.Event) {
			e := event.Data.(events.LyricUpdateEvent)
			if prevIndex >= len(fullLrc.Objects) || e.Lyric.Index >= len(fullLrc.Objects) {
				// fix race condition
				return
			}
			if e.Lyric == nil {
				currentLrc.SetText("")
				return
			}
			fullLrc.Objects[prevIndex].(*widget.Label).TextStyle.Bold = false
			fullLrc.Objects[prevIndex].Refresh()
			fullLrc.Objects[e.Lyric.Index].(*widget.Label).TextStyle.Bold = true
			fullLrc.Objects[e.Lyric.Index].Refresh()
			prevIndex = e.Lyric.Index
			currentLrc.SetText(e.Lyric.Now.Lyric)
			lrcWindow.Scrolled(&fyne.ScrollEvent{
				Scrolled: fyne.Delta{
					DX: 0,
					DY: lrcWindow.Offset.Y - float32(e.Lyric.Index-2)/float32(e.Lyric.Total)*lrcWindow.Content.Size().Height,
				},
			})
			fullLrc.Refresh()
		})
	API.PlayControl().GetLyric().EventManager().RegisterA(
		events.EventLyricReload, "player.lyric.new_media", func(event *event.Event) {
			e := event.Data.(events.LyricReloadEvent)
			lrcs := make([]string, len(e.Lyrics.Lyrics))
			for i := 0; i < len(lrcs); i++ {
				lrcs[i] = e.Lyrics.Lyrics[i].Lyric
			}
			fullLrc.Objects = createLyricObj(e.Lyrics)
			//fullLrc.SetText(strings.Join(lrcs, "\n"))
			//fullLrc.Segments[0] = &widget.TextSegment{
			//	Style: widget.RichTextStyleInline,
			//	Text:  strings.Join(lrcs, "\n\n"),
			//}
			lrcWindow.Refresh()
		})

	w.SetOnClosed(func() {
		API.PlayControl().GetLyric().EventManager().Unregister("player.lyric.current_lyric")
		API.PlayControl().GetLyric().EventManager().Unregister("player.lyric.new_media")
		PlayController.LrcWindowOpen = false
	})
	return w
}
