package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// AsyncButton is a Button that handle OnTapped handler asynchronously.
type AsyncButton struct {
	widget.Button
	anim *fyne.Animation
}

func NewAsyncButton(label string, tapped func()) *AsyncButton {
	button := &AsyncButton{
		Button: widget.Button{
			Text:     label,
			OnTapped: tapped,
		},
	}
	button.ExtendBaseWidget(button)
	return button
}

func NewAsyncButtonWithIcon(label string, icon fyne.Resource, tapped func()) *AsyncButton {
	button := &AsyncButton{
		Button: widget.Button{
			Text:     label,
			Icon:     icon,
			OnTapped: tapped,
		},
	}
	button.ExtendBaseWidget(button)
	return button
}

func (b *AsyncButton) Tapped(e *fyne.PointEvent) {
	if b.Disabled() {
		return
	}

	// missing animation
	b.Refresh()

	if b.OnTapped != nil {
		b.Disable()
		go func() {
			b.OnTapped()
			b.Enable()
		}()
	}
}
