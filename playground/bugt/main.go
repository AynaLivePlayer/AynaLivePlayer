package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var w fyne.Window

func main() {
	a := app.NewWithID("io.fyne.demo")
	a.SetIcon(theme.FyneLogo())
	w = a.NewWindow("Fyne Demo")
	Regen(w)
	w.Resize(fyne.NewSize(1080, 720))
	w.ShowAndRun()
}

func Regen(w fyne.Window) {
	tabs := container.NewDocTabs()
	for _, datum := range generateData(100) {
		tabs.Append(newItemTab(&datum))
	}
	w.SetContent(tabs)
}

func generateData(n int) (result []int) {
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return
}

func newItemTab(i *int) *container.TabItem {
	c := container.NewVBox(
		BindIntWithEntry(i),
		widget.NewButton("Regen", func() {
			Regen(w)
		}),
	)
	return container.NewTabItemWithIcon(strconv.Itoa(*i), theme.MenuIcon(), c)
}

func BindIntWithLabel(k *int) *widget.Label {
	b := binding.BindInt(k)
	return widget.NewLabelWithData(binding.IntToString(b))
}

func BindIntWithEntry(k *int) *widget.Entry {
	b := binding.BindInt(k)
	return widget.NewEntryWithData(binding.IntToString(b))
}
