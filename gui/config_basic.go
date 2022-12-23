package gui

import (
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
					controller.Instance.Playlists().GetCurrent().Model().Mode = model.PlaylistModeRandom
				} else {
					controller.Instance.Playlists().GetCurrent().Model().Mode = model.PlaylistModeNormal
				}
			},
			controller.Instance.Playlists().GetCurrent().Model().Mode == model.PlaylistModeRandom),
		newCheckInit(
			i18n.T("gui.config.basic.random_playlist.system"),
			func(b bool) {
				l().Infof("Set random playlist for system: %t", b)
				if b {
					controller.Instance.Playlists().GetDefault().Model().Mode = model.PlaylistModeRandom
				} else {
					controller.Instance.Playlists().GetDefault().Model().Mode = model.PlaylistModeNormal
				}
			},
			controller.Instance.Playlists().GetDefault().Model().Mode == model.PlaylistModeRandom),
	)
	devices := controller.Instance.PlayControl().GetAudioDevices()
	deviceDesc := make([]string, len(devices))
	deviceDesc2Name := make(map[string]string)
	for i, device := range devices {
		deviceDesc[i] = device.Description
		deviceDesc2Name[device.Description] = device.Name
	}
	deviceSel := widget.NewSelect(deviceDesc, func(s string) {
		controller.Instance.PlayControl().SetAudioDevice(deviceDesc2Name[s])
	})
	deviceSel.Selected = controller.Instance.PlayControl().GetCurrentAudioDevice()
	outputDevice := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("gui.config.basic.audio_device")), nil,
		deviceSel)
	skipPlaylist := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.skip_playlist")),
		newCheckInit(
			i18n.T("gui.config.basic.skip_playlist.prompt"),
			func(b bool) {
				controller.Instance.PlayControl().SetSkipPlaylist(b)
			},
			controller.Instance.PlayControl().GetSkipPlaylist(),
		),
	)
	b.panel = container.NewVBox(randomPlaylist, outputDevice, skipPlaylist)
	return b.panel
}
