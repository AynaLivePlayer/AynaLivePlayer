package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type LabelOpt func(*widget.Label)

func LabelWrapping(wrapping fyne.TextWrap) LabelOpt {
	return func(l *widget.Label) {
		l.Wrapping = wrapping
	}
}

func LabelAlignment(align fyne.TextAlign) LabelOpt {
	return func(l *widget.Label) {
		l.Alignment = align
	}
}

func LabelTextStyle(style fyne.TextStyle) LabelOpt {
	return func(l *widget.Label) {
		l.TextStyle = style
	}
}

func LabelTruncation(truncation fyne.TextTruncation) LabelOpt {
	return func(l *widget.Label) {
		l.Truncation = truncation
	}
}

func NewLabelWithOpts(text string, opts ...LabelOpt) *widget.Label {
	l := widget.NewLabel(text)
	for _, opt := range opts {
		opt(l)
	}
	return l
}
