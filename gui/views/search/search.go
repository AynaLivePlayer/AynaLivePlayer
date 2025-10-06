package search

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateView() fyne.CanvasObject {
	return container.NewBorder(createSearchBar(), nil, nil, nil, createSearchList())
}
