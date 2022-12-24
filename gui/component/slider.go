package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SliderPlus struct {
	widget.Slider
	OnDragEnd func(value float64)
	Dragging  bool // during dragging
}

func NewSliderPlus(min, max float64) *SliderPlus {
	slider := &SliderPlus{
		Slider: widget.Slider{
			Value:       0,
			Min:         min,
			Max:         max,
			Step:        1,
			Orientation: widget.Horizontal,
		},
	}
	slider.ExtendBaseWidget(slider)
	return slider
}

func (s *SliderPlus) DragEnd() {
	if s.OnDragEnd != nil {
		s.OnDragEnd(s.Value)
	}
	s.Dragging = false
}

func (s *SliderPlus) Dragged(e *fyne.DragEvent) {
	s.Dragging = true
	s.Slider.Dragged(e)
}
