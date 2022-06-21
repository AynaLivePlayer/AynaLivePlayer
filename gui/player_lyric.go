package gui

import (
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/player"
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
	lrcs := make([]string, len(controller.CurrentLyric.Lyrics))
	for i := 0; i < len(lrcs); i++ {
		lrcs[i] = controller.CurrentLyric.Lyrics[i].Lyric
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
	controller.CurrentLyric.Handler.RegisterA(player.EventLyricUpdate, "player.lyric.current_lyric", func(event *event.Event) {
		e := event.Data.(player.LyricUpdateEvent)
		if e.Lyric == nil {
			currentLrc.SetText("")
			return
		}
		currentLrc.SetText(e.Lyric.Lyric)
	})
	controller.CurrentLyric.Handler.RegisterA(player.EventLyricReload, "player.lyric.new_media", func(event *event.Event) {
		e := event.Data.(player.LyricReloadEvent)
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
		controller.CurrentLyric.Handler.Unregister("player.lyric.current_lyric")
		controller.CurrentLyric.Handler.Unregister("player.lyric.new_media")
		PlayController.LrcWindowOpen = false
	})
	return w
}
