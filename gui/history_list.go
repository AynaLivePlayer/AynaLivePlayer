package gui

import (
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/player"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HistoryContainer struct {
	Playlist *player.Playlist
	List     *widget.List
}

var History = &PlaylistContainer{}

func createHistoryList() fyne.CanvasObject {
	History.Playlist = controller.History
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
			m := History.Playlist.Playlist[History.Playlist.Size()-id-1]
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				m.ToUser().Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			btns[0].(*widget.Button).OnTapped = func() {
				controller.Play(controller.ToHistoryMedia(m))
			}
			btns[1].(*widget.Button).OnTapped = func() {
				controller.UserPlaylist.Push(controller.ToHistoryMedia(m))
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
	History.Playlist.Handler.RegisterA(player.EventPlaylistUpdate, "gui.history.update", func(event *event.Event) {
		History.Playlist.Lock.RLock()
		History.List.Refresh()
		History.Playlist.Lock.RUnlock()
	})
}
