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
		_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistDeleteCmd(model.PlaylistIDPlayer), events.PlaylistDeleteCmdEvent{
			Index: b.Index,
		})
	})
	topItem := fyne.NewMenuItem(i18n.T("gui.player.playlist.op.top"), func() {
		_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistMoveCmd(model.PlaylistIDPlayer), events.PlaylistMoveCmdEvent{
			From: b.Index,
			To:   0,
		})
	})
	m := fyne.NewMenu("", deleteItem, topItem)
	b.menu = m
	b.Text = ""
	b.Icon = theme.MoreHorizontalIcon()
	b.ExtendBaseWidget(b)
	return b
}

var UserPlaylist = &struct {
	Medias []model.Media
	List   *widget.List
	mux    sync.RWMutex
}{}

func createPlaylist() fyne.CanvasObject {
	UserPlaylist.List = widget.NewList(
		func() int {
			//todo: @4
			return len(UserPlaylist.Medias)
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
				UserPlaylist.Medias[id].Info.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				UserPlaylist.Medias[id].Info.Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				UserPlaylist.Medias[id].ToUser().Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			object.(*fyne.Container).Objects[2].(*playlistOperationButton).Index = id
		})
	global.EventBus.Subscribe(eventChannel,  events.PlaylistDetailUpdate(model.PlaylistIDPlayer), "gui.player.playlist.update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		UserPlaylist.mux.Lock()
		UserPlaylist.Medias = event.Data.(events.PlaylistDetailUpdateEvent).Medias
		UserPlaylist.List.Refresh()
		UserPlaylist.mux.Unlock()
	}))
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
