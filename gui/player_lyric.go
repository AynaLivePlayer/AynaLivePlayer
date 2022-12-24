package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createLyricObj(lyric *model.Lyric) []fyne.CanvasObject {
	lrcs := make([]fyne.CanvasObject, len(lyric.Lyrics))
	for i := 0; i < len(lrcs); i++ {
		l := widget.NewLabelWithStyle(
			lyric.Lyrics[i].Lyric,
			fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
		l.Wrapping = fyne.TextWrapWord
		lrcs[i] = l
	}
	return lrcs
}

func createLyricWindow() fyne.Window {

	// create widgets
	w := App.NewWindow("Lyric")
	currentLrc := newLabelWithWrapping("", fyne.TextWrapBreak)
	currentLrc.Alignment = fyne.TextAlignCenter
	fullLrc := container.NewVBox(createLyricObj(controller.Instance.PlayControl().GetLyric().Get())...)
	lrcWindow := container.NewVScroll(fullLrc)
	prevIndex := 0
	w.SetContent(container.NewBorder(nil,
		container.NewVBox(widget.NewSeparator(), currentLrc),
		nil, nil,
		lrcWindow))
	w.Resize(fyne.NewSize(360, 540))
	w.CenterOnScreen()

	// register handlers
	controller.Instance.PlayControl().GetLyric().EventManager().RegisterA(
		model.EventLyricUpdate, "player.lyric.current_lyric", func(event *event.Event) {
			e := event.Data.(model.LyricUpdateEvent)
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
	controller.Instance.PlayControl().GetLyric().EventManager().RegisterA(
		model.EventLyricReload, "player.lyric.new_media", func(event *event.Event) {
			e := event.Data.(model.LyricReloadEvent)
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
		controller.Instance.PlayControl().GetLyric().EventManager().Unregister("player.lyric.current_lyric")
		controller.Instance.PlayControl().GetLyric().EventManager().Unregister("player.lyric.new_media")
		PlayController.LrcWindowOpen = false
	})
	return w
}
