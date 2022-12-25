package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

var History = &struct {
	Playlist *model.Playlist
	List     *widget.List
	mux      sync.RWMutex
}{}

func createHistoryList() fyne.CanvasObject {
	History.Playlist = controller.Instance.Playlists().GetHistory().Model().Copy()
	History.List = widget.NewList(
		func() int {
			return History.Playlist.Size()
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
			m := History.Playlist.Medias[History.Playlist.Size()-id-1].Copy()
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				m.ToUser().Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			m.User = controller.HistoryUser
			btns[0].(*widget.Button).OnTapped = func() {
				showDialogIfError(controller.Instance.PlayControl().Play(m))
			}
			btns[1].(*widget.Button).OnTapped = func() {
				controller.Instance.Playlists().GetCurrent().Push(m)
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
	controller.Instance.Playlists().GetHistory().EventManager().RegisterA(model.EventPlaylistUpdate, "gui.history.update", func(event *event.Event) {
		History.mux.RLock()
		History.Playlist = event.Data.(model.PlaylistUpdateEvent).Playlist
		History.List.Refresh()
		History.mux.RUnlock()
	})
}
