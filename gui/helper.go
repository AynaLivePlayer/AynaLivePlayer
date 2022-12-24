package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

func showDialogIfError(err error) {
	if err != nil {
		dialog.ShowError(err, MainWindow)
	}
}

func newCheckInit(name string, changed func(bool), checked bool) *widget.Check {
	check := widget.NewCheck(name, changed)
	check.SetChecked(checked)
	return check
}
