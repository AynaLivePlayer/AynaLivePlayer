package gui

import (
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui/component"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlaylistsTab struct {
	Playlists             *widget.List
	PlaylistMedia         *widget.List
	Index                 int
	AddBtn                *widget.Button
	RemoveBtn             *widget.Button
	SetAsSystemBtn        *component.AsyncButton
	RefreshBtn            *component.AsyncButton
	CurrentSystemPlaylist *widget.Label
}

func (p *PlaylistsTab) UpdateCurrentSystemPlaylist() {
	p.CurrentSystemPlaylist.SetText(i18n.T("gui.playlist.current") + controller.Instance.Playlists().GetDefault().Name())
}

var PlaylistManager = &PlaylistsTab{}

func createPlaylists() fyne.CanvasObject {
	PlaylistManager.Playlists = widget.NewList(
		func() int {
			return controller.Instance.Playlists().Size()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(
				controller.Instance.Playlists().Get(id).Name())
		})
	PlaylistManager.AddBtn = widget.NewButton(i18n.T("gui.playlist.button.add"), func() {
		providerEntry := widget.NewSelect(controller.Instance.Provider().GetPriority(), nil)
		idEntry := widget.NewEntry()
		dia := dialog.NewCustomConfirm(
			i18n.T("gui.playlist.add.title"),
			i18n.T("gui.playlist.add.confirm"),
			i18n.T("gui.playlist.add.cancel"),
			container.NewVBox(
				container.New(
					layout.NewFormLayout(),
					widget.NewLabel(i18n.T("gui.playlist.add.source")),
					providerEntry,
					widget.NewLabel(i18n.T("gui.playlist.add.id_url")),
					idEntry,
				),
				widget.NewLabel(i18n.T("gui.playlist.add.prompt")),
			),
			func(b bool) {
				if b && len(providerEntry.Selected) > 0 && len(idEntry.Text) > 0 {
					controller.Instance.Playlists().Add(providerEntry.Selected, idEntry.Text)
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
		controller.Instance.Playlists().Remove(PlaylistManager.Index)
		//PlaylistManager.Index = 0
		PlaylistManager.Playlists.Select(0)
		PlaylistManager.Playlists.Refresh()
		PlaylistManager.PlaylistMedia.Refresh()
	})
	PlaylistManager.Playlists.OnSelected = func(id widget.ListItemID) {
		PlaylistManager.Index = id
		PlaylistManager.PlaylistMedia.Refresh()
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
	PlaylistManager.RefreshBtn = component.NewAsyncButtonWithIcon(
		i18n.T("gui.playlist.button.refresh"), theme.ViewRefreshIcon(),
		func() {
			showDialogIfError(controller.Instance.Playlists().PreparePlaylistByIndex(PlaylistManager.Index))
			PlaylistManager.PlaylistMedia.Refresh()
		})
	PlaylistManager.SetAsSystemBtn = component.NewAsyncButton(
		i18n.T("gui.playlist.button.set_as_system"),
		func() {
			showDialogIfError(controller.Instance.Playlists().SetDefault(PlaylistManager.Index))
			PlaylistManager.PlaylistMedia.Refresh()
			PlaylistManager.UpdateCurrentSystemPlaylist()
		})

	PlaylistManager.CurrentSystemPlaylist = widget.NewLabel("Current: ")
	PlaylistManager.UpdateCurrentSystemPlaylist()
	PlaylistManager.PlaylistMedia = widget.NewList(
		func() int {
			if controller.Instance.Playlists().Size() == 0 {
				return 0
			}
			return controller.Instance.Playlists().Get(PlaylistManager.Index).Size()
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
			m := controller.Instance.Playlists().Get(PlaylistManager.Index).Get(id).Copy()
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Artist)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			m.User = controller.SystemUser
			btns[0].(*widget.Button).OnTapped = func() {
				controller.Instance.PlayControl().Play(m)
			}
			btns[1].(*widget.Button).OnTapped = func() {
				controller.Instance.Playlists().GetCurrent().Push(m)
			}
		})
	return container.NewBorder(
		container.NewHBox(PlaylistManager.RefreshBtn, PlaylistManager.SetAsSystemBtn, PlaylistManager.CurrentSystemPlaylist), nil,
		nil, nil,
		PlaylistManager.PlaylistMedia)
}
