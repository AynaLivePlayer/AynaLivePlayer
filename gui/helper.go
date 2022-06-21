package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func newPaddedBoarder(top, bottom, left, right fyne.CanvasObject, objects ...fyne.CanvasObject) *fyne.Container {
	return container.NewPadded(container.NewBorder(top, bottom, left, right, objects...))
}

func newLabelWithWrapping(text string, wrapping fyne.TextWrap) *widget.Label {
	w := widget.NewLabel(text)
	w.Wrapping = wrapping

	return w
}

func createAsyncOnTapped(btn *widget.Button, f func()) func() {
	return func() {
		btn.Disable()
		go func() {
			f()
			btn.Enable()
		}()
	}
}

func createAsyncButton(btn *widget.Button, tapped func()) *widget.Button {
	btn.OnTapped = createAsyncOnTapped(btn, tapped)
	return btn
}
