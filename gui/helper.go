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

type ContextMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

func (b *ContextMenuButton) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newContextMenuButton(label string, menu *fyne.Menu) *ContextMenuButton {
	b := &ContextMenuButton{menu: menu}
	b.Text = label

	b.ExtendBaseWidget(b)
	return b
}

type FixedSplitContainer struct {
	*container.Split
}

func (f *FixedSplitContainer) Dragged(event *fyne.DragEvent) {
	// do nothing
}

func (f *FixedSplitContainer) DragEnd() {
	// do nothing
}

func newFixedSplitContainer(horizontal bool, leading, trailing fyne.CanvasObject) *FixedSplitContainer {
	s := &container.Split{
		Offset:     0.5, // Sensible default, can be overridden with SetOffset
		Horizontal: horizontal,
		Leading:    leading,
		Trailing:   trailing,
	}
	fs := &FixedSplitContainer{
		s,
	}
	fs.Split.BaseWidget.ExtendBaseWidget(s)
	return fs
}
