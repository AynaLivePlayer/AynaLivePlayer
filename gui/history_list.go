package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

var History = &struct {
	Medias []model.Media
	List   *widget.List
	mux    sync.RWMutex
}{}

func createHistoryList() fyne.CanvasObject {
	History.List = widget.NewList(
		func() int {
			return len(History.Medias)
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
					newLabelWithWrapping("user", fyne.TextTruncate)))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			m := History.Medias[len(History.Medias)-id-1]
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Info.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Info.Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				m.ToUser().Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			m.User = model.SystemUser
			btns[0].(*widget.Button).OnTapped = func() {
				global.EventManager.CallA(events.PlayerPlayCmd, events.PlayerPlayCmdEvent{
					Media: m,
				})
			}
			btns[1].(*widget.Button).OnTapped = func() {
				global.EventManager.CallA(events.PlaylistInsertCmd(model.PlaylistIDPlayer), events.PlaylistInsertCmdEvent{
					Media:    m,
					Position: -1,
				})
			}
		})
	registerHistoryHandler()
	return container.NewBorder(
		container.NewBorder(nil, nil,
			widget.NewLabel("#"), widget.NewLabel(i18n.T("gui.history.operation")),
			container.NewGridWithColumns(3,
				widget.NewLabel(i18n.T("gui.history.title")),
				widget.NewLabel(i18n.T("gui.history.artist")),
				widget.NewLabel(i18n.T("gui.history.user")))),
		nil, nil, nil,
		History.List,
	)
}

func registerHistoryHandler() {
	global.EventManager.RegisterA(
		events.PlaylistDetailUpdate(model.PlaylistIDHistory),
		"gui.history.update", func(event *event.Event) {
			History.mux.Lock()
			History.Medias = event.Data.(events.PlaylistDetailUpdateEvent).Medias
			History.List.Refresh()
			History.mux.Unlock()
		})
}
