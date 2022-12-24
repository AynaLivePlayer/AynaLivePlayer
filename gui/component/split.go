package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FixedSplit struct {
	widget.BaseWidget
	Offset     float64
	Horizontal bool
	Leading    fyne.CanvasObject
	Trailing   fyne.CanvasObject
}

func NewFixedHSplitContainer(leading, trailing fyne.CanvasObject, offset float64) *FixedSplit {
	return NewFixedSplitContainer(leading, trailing, true, offset)

}

func NewFixedVSplitContainer(top, bottom fyne.CanvasObject, offset float64) *FixedSplit {
	return NewFixedSplitContainer(top, bottom, false, offset)
}

func NewFixedSplitContainer(leading, trailing fyne.CanvasObject, horizontal bool, offset float64) *FixedSplit {
	s := &FixedSplit{
		Offset:     offset, // Sensible default, can be overridden with SetOffset
		Horizontal: horizontal,
		Leading:    leading,
		Trailing:   trailing,
	}
	s.BaseWidget.ExtendBaseWidget(s)
	return s
}

func (s *FixedSplit) CreateRenderer() fyne.WidgetRenderer {
	s.BaseWidget.ExtendBaseWidget(s)
	d := widget.NewSeparator()
	return &fixedSplitContainerRenderer{
		split:   s,
		divider: d,
		objects: []fyne.CanvasObject{s.Leading, d, s.Trailing},
	}
}

func (s *FixedSplit) SetOffset(offset float64) {
	if s.Offset == offset {
		return
	}
	s.Offset = offset
	s.Refresh()
}

type fixedSplitContainerRenderer struct {
	split   *FixedSplit
	divider *widget.Separator
	objects []fyne.CanvasObject
}

func (r *fixedSplitContainerRenderer) Destroy() {
}

func (r *fixedSplitContainerRenderer) Layout(size fyne.Size) {
	var dividerPos, leadingPos, trailingPos fyne.Position
	var dividerSize, leadingSize, trailingSize fyne.Size

	if r.split.Horizontal {
		lw, tw := r.computeSplitLengths(size.Width, r.split.Leading.MinSize().Width, r.split.Trailing.MinSize().Width)
		leadingPos.X = 0
		leadingSize.Width = lw
		leadingSize.Height = size.Height
		dividerPos.X = lw
		dividerSize.Width = theme.SeparatorThicknessSize()
		dividerSize.Height = size.Height
		trailingPos.X = lw + dividerSize.Width
		trailingSize.Width = tw
		trailingSize.Height = size.Height
	} else {
		lh, th := r.computeSplitLengths(size.Height, r.split.Leading.MinSize().Height, r.split.Trailing.MinSize().Height)
		leadingPos.Y = 0
		leadingSize.Width = size.Width
		leadingSize.Height = lh
		dividerPos.Y = lh
		dividerSize.Width = size.Width
		dividerSize.Height = theme.SeparatorThicknessSize()
		trailingPos.Y = lh + dividerSize.Height
		trailingSize.Width = size.Width
		trailingSize.Height = th
	}

	r.divider.Move(dividerPos)
	r.divider.Resize(dividerSize)
	r.split.Leading.Move(leadingPos)
	r.split.Leading.Resize(leadingSize)
	r.split.Trailing.Move(trailingPos)
	r.split.Trailing.Resize(trailingSize)
	canvas.Refresh(r.divider)
	canvas.Refresh(r.split.Leading)
	canvas.Refresh(r.split.Trailing)
}

func (r *fixedSplitContainerRenderer) MinSize() fyne.Size {
	s := fyne.NewSize(0, 0)
	for _, o := range r.objects {
		min := o.MinSize()
		if r.split.Horizontal {
			s.Width += min.Width
			s.Height = fyne.Max(s.Height, min.Height)
		} else {
			s.Width = fyne.Max(s.Width, min.Width)
			s.Height += min.Height
		}
	}
	return s
}

func (r *fixedSplitContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *fixedSplitContainerRenderer) Refresh() {
	r.objects[0] = r.split.Leading
	// [1] is divider which doesn't change
	r.objects[2] = r.split.Trailing
	r.Layout(r.split.Size())
	canvas.Refresh(r.split)
}

func (r *fixedSplitContainerRenderer) computeSplitLengths(total, lMin, tMin float32) (float32, float32) {
	available := float64(total - theme.SeparatorThicknessSize())
	if available <= 0 {
		return 0, 0
	}
	ld := float64(lMin)
	tr := float64(tMin)
	offset := r.split.Offset

	min := ld / available
	max := 1 - tr/available
	if min <= max {
		if offset < min {
			offset = min
		}
		if offset > max {
			offset = max
		}
	} else {
		offset = ld / (ld + tr)
	}

	ld = offset * available
	tr = available - ld
	return float32(ld), float32(tr)
}
