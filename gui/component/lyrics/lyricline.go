package lyrics

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type lyricLine struct {
	widget.BaseWidget

	Text             string
	SizeName         fyne.ThemeSizeName
	ColorName        fyne.ThemeColorName
	HoveredColorName fyne.ThemeColorName
	Alignment        fyne.TextAlign
	Tappable         bool

	onTapped func()
	hovered  bool
	richtext *widget.RichText
}

func newLyricLine(text string, onTapped func()) *lyricLine {
	l := &lyricLine{
		Text:      text,
		SizeName:  theme.SizeNameSubHeadingText,
		ColorName: theme.ColorNameForeground,
		Alignment: fyne.TextAlignLeading,
		onTapped:  onTapped,
	}
	l.ExtendBaseWidget(l)
	return l
}

var _ desktop.Hoverable = (*lyricLine)(nil)

func (l *lyricLine) MouseIn(*desktop.MouseEvent) {
	if l.Tappable {
		l.hovered = true
		l.Refresh()
	}
}

func (l *lyricLine) MouseMoved(*desktop.MouseEvent) {
}

func (l *lyricLine) MouseOut() {
	l.hovered = false
	l.Refresh()
}

var _ desktop.Cursorable = (*lyricLine)(nil)

func (l *lyricLine) Cursor() desktop.Cursor {
	if l.Tappable {
		return desktop.PointerCursor
	}
	return desktop.DefaultCursor
}

var _ fyne.Tappable = (*lyricLine)(nil)

func (l *lyricLine) Tapped(*fyne.PointEvent) {
	if l.Tappable {
		l.onTapped()
	}
}

func (l *lyricLine) updateRichText() {
	if l.richtext == nil {
		l.richtext = widget.NewRichText(&widget.TextSegment{
			Style: widget.RichTextStyleSubHeading,
		})
		l.richtext.Wrapping = fyne.TextWrapWord
	}
	seg := l.richtext.Segments[0].(*widget.TextSegment)
	seg.Text = l.Text
	seg.Style.Alignment = l.Alignment
	if l.hovered {
		seg.Style.ColorName = l.HoveredColorName
	} else {
		seg.Style.ColorName = l.ColorName
	}
	seg.Style.SizeName = l.SizeName
}

func (l *lyricLine) Refresh() {
	l.updateRichText()
	l.richtext.Refresh()
}

func (l *lyricLine) CreateRenderer() fyne.WidgetRenderer {
	l.updateRichText()
	return widget.NewSimpleRenderer(l.richtext)
}
