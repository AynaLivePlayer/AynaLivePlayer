package gui

import (
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/player"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
			// todo: @4
			return UserPlaylist.Playlist.Size()
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewLabel("index"), widget.NewLabel("user"),
				container.NewGridWithColumns(2,
					newLabelWithWrapping("title", fyne.TextTruncate),
					newLabelWithWrapping("artist", fyne.TextTruncate)))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				UserPlaylist.Playlist.Playlist[id].Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				UserPlaylist.Playlist.Playlist[id].Artist)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			object.(*fyne.Container).Objects[2].(*widget.Label).SetText(UserPlaylist.Playlist.Playlist[id].ToUser().Name)
		})
	registerPlaylistHandler()
	return container.NewBorder(
		container.NewBorder(nil, nil,
			widget.NewLabel("#"), widget.NewLabel("User"),
			container.NewGridWithColumns(2, widget.NewLabel("Title"), widget.NewLabel("Artist"))),
		widget.NewSeparator(),
		nil, nil,
		UserPlaylist.List,
	)
}

func registerPlaylistHandler() {
	UserPlaylist.Playlist.Handler.RegisterA(player.EventPlaylistUpdate, "gui.playlist.update", func(event *event.Event) {
		UserPlaylist.List.Refresh()
	})
}
