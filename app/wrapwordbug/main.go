package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	texts := make([]fyne.CanvasObject, 1)
	for i := 0; i < len(texts); i++ {
		l := widget.NewLabelWithStyle(
			" AAAA",
			fyne.TextAlignCenter, fyne.TextStyle{})
		l.Wrapping = fyne.TextWrapWord
		texts[i] = l
	}
	vbox := container.NewVBox(texts...)
	scroll := container.NewScroll(vbox)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(360, 540))
	w.ShowAndRun()
}
