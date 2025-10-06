package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type LabelFixedSize struct {
	*widget.Label
	fixedSize fyne.Size
}

func (t *LabelFixedSize) MinSize() fyne.Size {
	return t.fixedSize
}

func NewLabelFixedSize(label *widget.Label) *LabelFixedSize {
	return &LabelFixedSize{
		Label:     label,
		fixedSize: label.MinSize(),
	}
}
