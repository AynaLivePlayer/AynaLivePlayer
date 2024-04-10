package main

import (
	"AynaLivePlayer/gui/xfyne"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"time"
)

var w fyne.Window

func main() {
	a := app.NewWithID("io.fyne.demo")
	w = a.NewWindow("Fyne Demo")
	go func() {
		time.Sleep(5 * time.Second)
		println("Window handle:", xfyne.GetWindowHandle(w))
	}()
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}
