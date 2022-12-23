package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func createLyricWindow() fyne.Window {

	// create widgets
	w := App.NewWindow("Lyric")
	currentLrc := newLabelWithWrapping("", fyne.TextWrapBreak)
	currentLrc.Alignment = fyne.TextAlignCenter
	lrcs := make([]string, len(controller.Instance.PlayControl().GetLyric().Get().Lyrics))
	for i := 0; i < len(lrcs); i++ {
		lrcs[i] = controller.Instance.PlayControl().GetLyric().Get().Lyrics[i].Lyric
	}
	fullLrc := widget.NewRichTextWithText(strings.Join(lrcs, "\n\n"))
	fullLrc.Scroll = container.ScrollVerticalOnly
	fullLrc.Wrapping = fyne.TextWrapWord
	w.SetContent(container.NewBorder(nil,
		container.NewVBox(widget.NewSeparator(), currentLrc),
		nil, nil,
		fullLrc))
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
			currentLrc.SetText(e.Lyric.Lyric)
		})
	controller.Instance.PlayControl().GetLyric().EventManager().RegisterA(
		model.EventLyricReload, "player.lyric.new_media", func(event *event.Event) {
			e := event.Data.(model.LyricReloadEvent)
			lrcs := make([]string, len(e.Lyrics.Lyrics))
			for i := 0; i < len(lrcs); i++ {
				lrcs[i] = e.Lyrics.Lyrics[i].Lyric
			}
			fullLrc.Segments[0] = &widget.TextSegment{
				Style: widget.RichTextStyleInline,
				Text:  strings.Join(lrcs, "\n\n"),
			}
			fullLrc.Refresh()
		})

	w.SetOnClosed(func() {
		controller.Instance.PlayControl().GetLyric().EventManager().Unregister("player.lyric.current_lyric")
		controller.Instance.PlayControl().GetLyric().EventManager().Unregister("player.lyric.new_media")
		PlayController.LrcWindowOpen = false
	})
	return w
}
