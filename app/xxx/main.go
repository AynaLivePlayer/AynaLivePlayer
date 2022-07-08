package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

func main() {
	var app = app.New()

	var (
		labelText       = ""
		bindedLabelText = binding.BindString(&labelText)

		label = widget.NewLabelWithData(bindedLabelText)
	)

	var window = app.NewWindow("Canvas")

	var verticalBox = container.NewVBox(label)
	window.SetContent(verticalBox)

	go func() {
		for i := 0; ; i++ {
			var newLabelText = strconv.Itoa(i)
			if err := bindedLabelText.Set(newLabelText); err != nil {
				panic(err)
			}

			time.Sleep(time.Microsecond)

			// NOTE: the only thing, that helps prevent UI updates from freezes, except for the direct manipulation with window size, e.g. update window size from 499x499 -> 500x500 and vice-versa for each iteration
			canvas.Refresh(label)
		}
	}()

	window.ShowAndRun()
}
