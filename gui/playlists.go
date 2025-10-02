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
	SetAsSystemBtn        *widget.Button
	RefreshBtn            *widget.Button
	CurrentSystemPlaylist *widget.Label
	currentMedias         []model.Media
	currentPlaylists      []model.PlaylistInfo
	providers             []string
}

var PlaylistManager = &PlaylistsTab{}

func createPlaylists() fyne.CanvasObject {
	PlaylistManager.Playlists = widget.NewList(
		func() int {
			return len(PlaylistManager.currentPlaylists)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(PlaylistManager.currentPlaylists[id].DisplayName())
		})
	PlaylistManager.AddBtn = widget.NewButton(i18n.T("gui.playlist.button.add"), func() {
		providerEntry := widget.NewSelect(PlaylistManager.providers, nil)
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
					logger.Infof("add playlists %s %s", providerEntry.Selected, idEntry.Text)
					_ = global.EventBus.PublishToChannel(eventChannel,
						events.PlaylistManagerAddPlaylistCmd,
						events.PlaylistManagerAddPlaylistCmdEvent{
							Provider: providerEntry.Selected,
							URL:      idEntry.Text,
						})
				}
			},
			MainWindow,
		)
		dia.Resize(fyne.NewSize(512, 256))
		dia.Show()
	})
	PlaylistManager.RemoveBtn = widget.NewButton(i18n.T("gui.playlist.button.remove"), func() {
		if PlaylistManager.Index >= len(PlaylistManager.currentPlaylists) {
			return
		}
		logger.Infof("remove playlists %s", PlaylistManager.currentPlaylists[PlaylistManager.Index].Meta.ID())
		_ = global.EventBus.PublishToChannel(eventChannel,
			events.PlaylistManagerRemovePlaylistCmd,
			events.PlaylistManagerRemovePlaylistCmdEvent{
				PlaylistID: PlaylistManager.currentPlaylists[PlaylistManager.Index].Meta.ID(),
			})
	})
	PlaylistManager.Playlists.OnSelected = func(id widget.ListItemID) {
		if id >= len(PlaylistManager.currentPlaylists) {
			return
		}
		PlaylistManager.Index = id
		_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistManagerGetCurrentCmd, events.PlaylistManagerGetCurrentCmdEvent{
			PlaylistID: PlaylistManager.currentPlaylists[id].Meta.ID(),
		})
	}
	global.EventBus.Subscribe(eventChannel,  events.MediaProviderUpdate,
		"gui.playlists.provider.update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			providers := event.Data.(events.MediaProviderUpdateEvent)
			s := make([]string, len(providers.Providers))
			copy(s, providers.Providers)
			PlaylistManager.providers = s
		}))
	global.EventBus.Subscribe(eventChannel,  events.PlaylistManagerInfoUpdate,
		"gui.playlists.info.update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerInfoUpdateEvent)
			prevLen := len(PlaylistManager.currentPlaylists)
			PlaylistManager.currentPlaylists = data.Playlists
			logger.Infof("receive playlist info update, try to refresh playlists. prevLen=%d, newLen=%d", prevLen, len(PlaylistManager.currentPlaylists))
			PlaylistManager.Playlists.Refresh()
			if prevLen != len(PlaylistManager.currentPlaylists) {
				PlaylistManager.Playlists.Select(0)
			}
		}))
	global.EventBus.Subscribe(eventChannel,  events.PlaylistManagerSystemUpdate,
		"gui.playlists.system.update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerSystemUpdateEvent)
			PlaylistManager.CurrentSystemPlaylist.SetText(i18n.T("gui.playlist.current") + data.Info.DisplayName())
		}))
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
	PlaylistManager.RefreshBtn = widget.NewButtonWithIcon(
		i18n.T("gui.playlist.button.refresh"), theme.ViewRefreshIcon(),
		func() {
			if PlaylistManager.Index >= len(PlaylistManager.currentPlaylists) {
				return
			}
			_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistManagerRefreshCurrentCmd, events.PlaylistManagerRefreshCurrentCmdEvent{
				PlaylistID: PlaylistManager.currentPlaylists[PlaylistManager.Index].Meta.ID(),
			})
		})
	PlaylistManager.SetAsSystemBtn = widget.NewButton(
		i18n.T("gui.playlist.button.set_as_system"),
		func() {
			if PlaylistManager.Index >= len(PlaylistManager.currentPlaylists) {
				return
			}
			logger.Infof("set playlist %s as system", PlaylistManager.currentPlaylists[PlaylistManager.Index].Meta.ID())
			_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistManagerSetSystemCmd, events.PlaylistManagerSetSystemCmdEvent{
				PlaylistID: PlaylistManager.currentPlaylists[PlaylistManager.Index].Meta.ID(),
			})
		})

	PlaylistManager.CurrentSystemPlaylist = widget.NewLabel("Current: ")
	PlaylistManager.PlaylistMedia = widget.NewList(
		func() int {
			return len(PlaylistManager.currentMedias)
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
			m := PlaylistManager.currentMedias[id]
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Info.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Info.Artist)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			m.User = model.SystemUser
			btns[0].(*widget.Button).OnTapped = func() {
				_ = global.EventBus.PublishToChannel(eventChannel, events.PlayerPlayCmd, events.PlayerPlayCmdEvent{
					Media: m,
				})
			}
			btns[1].(*widget.Button).OnTapped = func() {
				_ = global.EventBus.PublishToChannel(eventChannel, events.PlaylistInsertCmd(model.PlaylistIDPlayer), events.PlaylistInsertCmdEvent{
					Media:    m,
					Position: -1,
				})
			}
		})
	global.EventBus.Subscribe(eventChannel,  events.PlaylistManagerCurrentUpdate,
		"gui.playlists.current.update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			logger.Infof("receive current playlist update, try to refresh playlist medias")
			data := event.Data.(events.PlaylistManagerCurrentUpdateEvent)
			PlaylistManager.currentMedias = data.Medias
			PlaylistManager.PlaylistMedia.Refresh()
		}))
	return container.NewBorder(
		container.NewHBox(PlaylistManager.RefreshBtn, PlaylistManager.SetAsSystemBtn, PlaylistManager.CurrentSystemPlaylist), nil,
		nil, nil,
		PlaylistManager.PlaylistMedia)
}
