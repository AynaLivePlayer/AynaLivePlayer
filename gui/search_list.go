package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

var SearchResult = &struct {
	List  *widget.List
	Items []model.Media
	mux   sync.Mutex
}{
	Items: []model.Media{},
	mux:   sync.Mutex{},
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
				SearchResult.Items[id].Info.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				SearchResult.Items[id].Info.Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				SearchResult.Items[id].Info.Meta.Provider)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			btns[0].(*widget.Button).OnTapped = func() {
				_ = global.EventBus.PublishToChannel(eventChannel, events.PlayerPlayCmd, events.PlayerPlayCmdEvent{
					Media: SearchResult.Items[id],
				})
			}
			btns[1].(*widget.Button).OnTapped = func() {
				_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistInsertCmd(model.PlaylistIDPlayer), events.PlaylistInsertCmdEvent{
					Media:    SearchResult.Items[id],
					Position: -1,
				})
			}
		})
	global.EventBus.Subscribe(eventChannel,  events.SearchResultUpdate, "gui.search.update_result", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		items := event.Data.(events.SearchResultUpdateEvent).Medias
		SearchResult.Items = items
		SearchResult.mux.Lock()
		SearchResult.List.Refresh()
		SearchResult.mux.Unlock()
	}))
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
