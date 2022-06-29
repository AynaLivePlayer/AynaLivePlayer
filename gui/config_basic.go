package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
		widget.NewCheckWithData(
			i18n.T("gui.config.basic.random_playlist.user"),
			binding.BindBool(&controller.UserPlaylist.Config.RandomNext)),
		widget.NewCheckWithData(
			i18n.T("gui.config.basic.random_playlist.system"),
			binding.BindBool(&controller.SystemPlaylist.Config.RandomNext)),
	)
	devices := controller.GetAudioDevices()
	deviceDesc := make([]string, len(devices))
	deviceDesc2Name := make(map[string]string)
	for i, device := range devices {
		deviceDesc[i] = device.Description
		deviceDesc2Name[device.Description] = device.Name
	}
	deviceSel := widget.NewSelect(deviceDesc, func(s string) {
		controller.SetAudioDevice(deviceDesc2Name[s])
	})
	deviceSel.Selected = config.Player.AudioDevice
	outputDevice := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("gui.config.basic.audio_device")), nil,
		deviceSel)
	b.panel = container.NewVBox(randomPlaylist, outputDevice)
	return b.panel
}
