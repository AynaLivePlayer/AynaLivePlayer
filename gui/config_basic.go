package gui

import (
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/gui/component"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type bascicConfig struct {
	panel fyne.CanvasObject
}

func (b *bascicConfig) Title() string {
	return i18n.T("gui.config.basic.title")
}

func (b *bascicConfig) Description() string {
	return i18n.T("gui.config.basic.description")
}

func (b *bascicConfig) CreatePanel() fyne.CanvasObject {
	if b.panel != nil {
		return b.panel
	}

	randomPlaylist := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.random_playlist")),
		newCheckInit(
			i18n.T("gui.config.basic.random_playlist.user"),
			func(b bool) {
				l().Infof("Set random playlist for user: %t", b)
				if b {
					API.Playlists().GetCurrent().Model().Mode = model.PlaylistModeRandom
				} else {
					API.Playlists().GetCurrent().Model().Mode = model.PlaylistModeNormal
				}
			},
			API.Playlists().GetCurrent().Model().Mode == model.PlaylistModeRandom),
		newCheckInit(
			i18n.T("gui.config.basic.random_playlist.system"),
			func(b bool) {
				l().Infof("Set random playlist for system: %t", b)
				if b {
					API.Playlists().GetDefault().Model().Mode = model.PlaylistModeRandom
				} else {
					API.Playlists().GetDefault().Model().Mode = model.PlaylistModeNormal
				}
			},
			API.Playlists().GetDefault().Model().Mode == model.PlaylistModeRandom),
	)
	devices := API.PlayControl().GetAudioDevices()
	deviceDesc := make([]string, len(devices))
	deviceDesc2Name := make(map[string]string)
	for i, device := range devices {
		deviceDesc[i] = device.Description
		deviceDesc2Name[device.Description] = device.Name
	}
	deviceSel := widget.NewSelect(deviceDesc, func(s string) {
		API.PlayControl().SetAudioDevice(deviceDesc2Name[s])
	})
	deviceSel.Selected = API.PlayControl().GetCurrentAudioDevice()
	outputDevice := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("gui.config.basic.audio_device")), nil,
		deviceSel)
	skipPlaylist := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.skip_playlist")),
		component.NewCheckOneWayBinding(
			i18n.T("gui.config.basic.skip_playlist.prompt"),
			&API.PlayControl().Config().SkipPlaylist,
			API.PlayControl().Config().SkipPlaylist),
	)
	skipWhenErr := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.skip_when_error")),
		component.NewCheckOneWayBinding(
			i18n.T("gui.config.basic.skip_when_error.prompt"),
			&API.PlayControl().Config().AutoNextWhenFail,
			API.PlayControl().Config().AutoNextWhenFail),
	)
	checkUpdateBox := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.auto_check_update")),
		component.NewCheckOneWayBinding(
			i18n.T("gui.config.basic.auto_check_update.prompt"),
			&config.General.AutoCheckUpdate,
			config.General.AutoCheckUpdate),
	)
	checkUpdateBtn := widget.NewButton(i18n.T("gui.config.basic.check_update"), func() {
		err := API.App().CheckUpdate()
		if err != nil {
			showDialogIfError(err)
			return
		}
		dialog.ShowCustom(
			i18n.T("gui.update.new_version"),
			"OK",
			widget.NewRichTextFromMarkdown(API.App().LatestVersion().Info),
			MainWindow)
	})
	b.panel = container.NewVBox(randomPlaylist, outputDevice, skipPlaylist, skipWhenErr, checkUpdateBox, checkUpdateBtn)
	return b.panel
}
