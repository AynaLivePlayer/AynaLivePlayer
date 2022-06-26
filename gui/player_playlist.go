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

type playlistOperationButton struct {
	widget.Button
	Index int
	menu  *fyne.Menu
}

func (b *playlistOperationButton) Tapped(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func newPlaylistOperationButton() *playlistOperationButton {
	b := &playlistOperationButton{Index: 0}
	deleteItem := fyne.NewMenuItem(i18n.T("gui.player.playlist.op.delete"), func() {
		fmt.Println("delete", b.Index)
	})
	topItem := fyne.NewMenuItem(i18n.T("gui.player.playlist.op.top"), func() {
		controller.UserPlaylist.Move(b.Index, 0)
	})
	m := fyne.NewMenu("", deleteItem, topItem)
	b.menu = m
	b.Text = ""
	b.Icon = theme.MoreHorizontalIcon()
	b.ExtendBaseWidget(b)
	return b
}

type PlaylistContainer struct {
	Playlist *player.Playlist
	List     *widget.List
}

var UserPlaylist = &PlaylistContainer{}

func createPlaylist() fyne.CanvasObject {
	UserPlaylist.Playlist = controller.UserPlaylist
	UserPlaylist.List = widget.NewList(
		func() int {
			//debug.PrintStack()
			//todo: @4
			return UserPlaylist.Playlist.Size()
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewLabel("index"), newPlaylistOperationButton(),
				container.NewGridWithColumns(3,
					newLabelWithWrapping("title", fyne.TextTruncate),
					newLabelWithWrapping("artist", fyne.TextTruncate),
					newLabelWithWrapping("user", fyne.TextTruncate)))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				UserPlaylist.Playlist.Playlist[id].Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				UserPlaylist.Playlist.Playlist[id].Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				UserPlaylist.Playlist.Playlist[id].ToUser().Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			object.(*fyne.Container).Objects[2].(*playlistOperationButton).Index = id
		})
	registerPlaylistHandler()
	return container.NewBorder(
		container.NewBorder(nil, nil,
			widget.NewLabel("#"), widget.NewLabel(i18n.T("gui.player.playlist.ops")),
			container.NewGridWithColumns(3,
				widget.NewLabel(i18n.T("gui.player.playlist.title")),
				widget.NewLabel(i18n.T("gui.player.playlist.artist")),
				widget.NewLabel(i18n.T("gui.player.playlist.user")))),
		widget.NewSeparator(),
		nil, nil,
		UserPlaylist.List,
	)
}

func registerPlaylistHandler() {
	UserPlaylist.Playlist.Handler.RegisterA(player.EventPlaylistUpdate, "gui.playlist.update", func(event *event.Event) {
		// @6 Read lock Playlist when updating free after updating.
		UserPlaylist.Playlist.Lock.RLock()
		UserPlaylist.List.Refresh()
		UserPlaylist.Playlist.Lock.RUnlock()
	})
}
