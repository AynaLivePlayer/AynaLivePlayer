package gui

import "fyne.io/fyne/v2"

var ConfigList = []ConfigLayout{&bascicConfig{}}

type ConfigLayout interface {
	Title() string
	Description() string
	CreatePanel() fyne.CanvasObject
}

func AddConfigLayout(cfgs ...ConfigLayout) {
	ConfigList = append(ConfigList, cfgs...)
}
