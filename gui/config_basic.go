package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
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

	playerRandomCheck := widget.NewCheck(i18n.T("gui.config.basic.random_playlist.user"),
		func(b bool) {
			mode := model.PlaylistModeNormal
			if b {
				mode = model.PlaylistModeRandom
			}
			logger.Infof("Set player playlist mode to %d", mode)
			global.EventManager.CallA(events.PlaylistModeChangeCmd(model.PlaylistIDPlayer),
				events.PlaylistModeChangeCmdEvent{
					Mode: mode,
				})
		})
	global.EventManager.RegisterA(events.PlaylistModeChangeUpdate(model.PlaylistIDPlayer),
		"gui.config.basic.random_playlist.player",
		func(event *event.Event) {
			data := event.Data.(events.PlaylistModeChangeUpdateEvent)
			playerRandomCheck.SetChecked(data.Mode == model.PlaylistModeRandom)
		})

	systemRandomCheck := widget.NewCheck(i18n.T("gui.config.basic.random_playlist.system"),
		func(b bool) {
			mode := model.PlaylistModeNormal
			if b {
				mode = model.PlaylistModeRandom
			}
			global.EventManager.CallA(events.PlaylistModeChangeCmd(model.PlaylistIDSystem),
				events.PlaylistModeChangeCmdEvent{
					Mode: mode,
				})
		})

	global.EventManager.RegisterA(events.PlaylistModeChangeUpdate(model.PlaylistIDSystem),
		"gui.config.basic.random_playlist.system",
		func(event *event.Event) {
			data := event.Data.(events.PlaylistModeChangeUpdateEvent)
			systemRandomCheck.SetChecked(data.Mode == model.PlaylistModeRandom)
		})

	randomPlaylist := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.random_playlist")),
		playerRandomCheck,
		systemRandomCheck,
	)
	deviceDesc2Name := make(map[string]string)
	deviceSel := widget.NewSelect(make([]string, 0), func(s string) {
		name, ok := deviceDesc2Name[s]
		if !ok {
			return
		}
		global.EventManager.CallA(events.PlayerSetAudioDeviceCmd, events.PlayerSetAudioDeviceCmdEvent{
			Device: name,
		})
	})
	global.EventManager.RegisterA(
		events.PlayerAudioDeviceUpdate,
		"gui.config.basic.audio_device.update",
		func(event *event.Event) {
			data := event.Data.(events.PlayerAudioDeviceUpdateEvent)
			devices := make([]string, len(data.Devices))
			deviceDesc2Name = make(map[string]string)
			currentDevice := ""
			for i, device := range data.Devices {
				devices[i] = device.Description
				deviceDesc2Name[device.Description] = device.Name
				if device.Name == data.Current {
					currentDevice = device.Description
				}
			}
			logger.Infof("update audio device. set current to %s (%s)", data.Current, deviceDesc2Name[data.Current])
			deviceSel.Options = devices
			deviceSel.Selected = currentDevice
			deviceSel.Refresh()
		})

	outputDevice := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("gui.config.basic.audio_device")), nil,
		deviceSel)
	//skipWhenErr := container.NewHBox(
	//	widget.NewLabel(i18n.T("gui.config.basic.skip_when_error")),
	//	component.NewCheckOneWayBinding(
	//		i18n.T("gui.config.basic.skip_when_error.prompt"),
	//		&API.PlayControl().Config().AutoNextWhenFail,
	//		API.PlayControl().Config().AutoNextWhenFail),
	//)
	//checkUpdateBox := container.NewHBox(
	//	widget.NewLabel(i18n.T("gui.config.basic.auto_check_update")),
	//	component.NewCheckOneWayBinding(
	//		i18n.T("gui.config.basic.auto_check_update.prompt"),
	//		&config.General.AutoCheckUpdate,
	//		config.General.AutoCheckUpdate),
	//)
	//checkUpdateBtn := widget.NewButton(i18n.T("gui.config.basic.check_update"), func() {
	//	err := API.App().CheckUpdate()
	//	if err != nil {
	//		showDialogIfError(err)
	//		return
	//	}
	//	if API.App().LatestVersion().Version > API.App().Version().Version {
	//		dialog.ShowCustom(
	//			i18n.T("gui.update.new_version"),
	//			"OK",
	//			widget.NewRichTextFromMarkdown(API.App().LatestVersion().Info),
	//			MainWindow)
	//	}
	//})
	//b.panel = container.NewVBox(randomPlaylist, outputDevice, skipPlaylist, skipWhenErr, checkUpdateBox, checkUpdateBtn)
	b.panel = container.NewVBox(randomPlaylist, outputDevice)
	return b.panel
}
