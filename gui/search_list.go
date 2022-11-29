package gui

import (
	"AynaLivePlayer/controller"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var SearchResult = &struct {
	List  *widget.List
	Items []*player.Media
}{
	Items: []*player.Media{},
}

func createSearchList() fyne.CanvasObject {
	SearchResult.List = widget.NewList(
		func() int {
			return len(SearchResult.Items)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabel("index"),
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
					widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
				),
				container.NewGridWithColumns(3,
					newLabelWithWrapping("title", fyne.TextTruncate),
					newLabelWithWrapping("artist", fyne.TextTruncate),
					newLabelWithWrapping("source", fyne.TextTruncate)))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				SearchResult.Items[id].Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				SearchResult.Items[id].Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				SearchResult.Items[id].Meta.(provider.Meta).Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			btns[0].(*widget.Button).OnTapped = func() {
				controller.Play(SearchResult.Items[id])
			}
			btns[1].(*widget.Button).OnTapped = func() {
				controller.UserPlaylist.Push(SearchResult.Items[id])
			}
		})
	return container.NewBorder(
		container.NewBorder(nil, nil,
			widget.NewLabel("#"), widget.NewLabel(i18n.T("gui.search.operation")),
			container.NewGridWithColumns(3,
				widget.NewLabel(i18n.T("gui.search.title")),
				widget.NewLabel(i18n.T("gui.search.artist")),
				widget.NewLabel(i18n.T("gui.search.source")))),
		nil, nil, nil,
		SearchResult.List,
	)
}
