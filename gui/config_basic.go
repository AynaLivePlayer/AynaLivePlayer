package gui

import (
	"fyne.io/fyne/v2"
)

type bascicConfig struct{}

func (b bascicConfig) Title() string {
	return "Basic"
}

func (b bascicConfig) Description() string {
	return "Basic configuration"
}

func (b bascicConfig) Create() fyne.CanvasObject {
	//TODO implement me
	panic("implement me")
}
