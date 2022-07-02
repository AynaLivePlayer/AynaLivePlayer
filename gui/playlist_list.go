package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/i18n"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlaylistManagerContainer struct {
	Playlists             *widget.List
	PlaylistMedia         *widget.List
	Index                 int
	AddBtn                *widget.Button
	RemoveBtn             *widget.Button
	SetAsSystemBtn        *widget.Button
	RefreshBtn            *widget.Button
	CurrentSystemPlaylist *widget.Label
}

func (p *PlaylistManagerContainer) UpdateCurrentSystemPlaylist() {
	if config.Player.PlaylistIndex >= len(controller.PlaylistManager) {
		p.CurrentSystemPlaylist.SetText(i18n.T("gui.playlist.current.none"))
	}
	p.CurrentSystemPlaylist.SetText(i18n.T("gui.playlist.current") + controller.PlaylistManager[config.Player.PlaylistIndex].Name)
}

var PlaylistManager = &PlaylistManagerContainer{}

func createPlaylists() fyne.CanvasObject {
	PlaylistManager.Playlists = widget.NewList(
		func() int {
			return len(controller.PlaylistManager)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(
				controller.PlaylistManager[id].Name)
		})
	PlaylistManager.AddBtn = widget.NewButton(i18n.T("gui.playlist.button.add"), func() {
		providerEntry := widget.NewSelect(config.Provider.Priority, nil)
		idEntry := widget.NewEntry()
		dia := dialog.NewCustomConfirm(
			i18n.T("gui.playlist.add.title"),
			i18n.T("gui.playlist.add.confirm"),
			i18n.T("gui.playlist.add.cancel"),
			container.NewVBox(
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel(i18n.T("gui.playlist.add.confirm")),
					providerEntry,
					widget.NewLabel(i18n.T("gui.playlist.add.id_url")),
					idEntry,
				),
				widget.NewLabel(i18n.T("gui.playlist.add.prompt")),
			),
			func(b bool) {
				if b && len(providerEntry.Selected) > 0 && len(idEntry.Text) > 0 {
					controller.AddPlaylist(providerEntry.Selected, idEntry.Text)
					PlaylistManager.Playlists.Refresh()
					PlaylistManager.PlaylistMedia.Refresh()
				}
			},
			MainWindow,
		)
		dia.Resize(fyne.NewSize(512, 256))
		dia.Show()
	})
	PlaylistManager.RemoveBtn = widget.NewButton(i18n.T("gui.playlist.button.remove"), func() {
		controller.RemovePlaylist(PlaylistManager.Index)
		//PlaylistManager.Index = 0
		PlaylistManager.Playlists.Select(0)
		PlaylistManager.Playlists.Refresh()
		PlaylistManager.PlaylistMedia.Refresh()
	})
	PlaylistManager.Playlists.OnSelected = func(id widget.ListItemID) {
		PlaylistManager.Index = id
	}
	return container.NewHBox(
		container.NewBorder(
			nil, container.NewCenter(container.NewHBox(PlaylistManager.AddBtn, PlaylistManager.RemoveBtn)),
			nil, nil,
			PlaylistManager.Playlists,
		),
		widget.NewSeparator(),
	)
}

func createPlaylistMedias() fyne.CanvasObject {
	PlaylistManager.RefreshBtn = createAsyncButton(
		widget.NewButtonWithIcon(i18n.T("gui.playlist.button.refresh"), theme.ViewRefreshIcon(), nil),
		func() {
			controller.PreparePlaylistByIndex(PlaylistManager.Index)
			PlaylistManager.PlaylistMedia.Refresh()
		})
	PlaylistManager.SetAsSystemBtn = createAsyncButton(
		widget.NewButton(i18n.T("gui.playlist.button.set_as_system"), nil),
		func() {
			controller.SetSystemPlaylist(PlaylistManager.Index)
			PlaylistManager.PlaylistMedia.Refresh()
			PlaylistManager.UpdateCurrentSystemPlaylist()
		})
	PlaylistManager.CurrentSystemPlaylist = widget.NewLabel("Current: ")
	PlaylistManager.UpdateCurrentSystemPlaylist()
	PlaylistManager.PlaylistMedia = widget.NewList(
		func() int {
			if len(controller.PlaylistManager) == 0 {
				return 0
			}
			return controller.PlaylistManager[PlaylistManager.Index].Size()
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabel("index"),
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
					widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
				),
				container.NewGridWithColumns(2,
					newLabelWithWrapping("title", fyne.TextTruncate),
					newLabelWithWrapping("artist", fyne.TextTruncate)))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			m := controller.PlaylistManager[PlaylistManager.Index].Playlist[id]
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Artist)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			btns[0].(*widget.Button).OnTapped = func() {
				controller.Play(controller.ToSystemMedia(m))
			}
			btns[1].(*widget.Button).OnTapped = func() {
				controller.UserPlaylist.Push(controller.ToSystemMedia(m))
			}
		})
	return container.NewBorder(
		container.NewHBox(PlaylistManager.RefreshBtn, PlaylistManager.SetAsSystemBtn, PlaylistManager.CurrentSystemPlaylist), nil,
		nil, nil,
		PlaylistManager.PlaylistMedia)
}
