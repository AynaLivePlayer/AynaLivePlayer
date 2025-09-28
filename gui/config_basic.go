package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
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
			_ = global.EventBus.Publish(events.PlaylistModeChangeCmd(model.PlaylistIDPlayer),
				events.PlaylistModeChangeCmdEvent{
					Mode: mode,
				})
		})
	global.EventBus.Subscribe("", events.PlaylistModeChangeUpdate(model.PlaylistIDPlayer),
		"gui.config.basic.random_playlist.player",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistModeChangeUpdateEvent)
			playerRandomCheck.SetChecked(data.Mode == model.PlaylistModeRandom)
		}))

	systemRandomCheck := widget.NewCheck(i18n.T("gui.config.basic.random_playlist.system"),
		func(b bool) {
			mode := model.PlaylistModeNormal
			if b {
				mode = model.PlaylistModeRandom
			}
			_ = global.EventBus.Publish(events.PlaylistModeChangeCmd(model.PlaylistIDSystem),
				events.PlaylistModeChangeCmdEvent{
					Mode: mode,
				})
		})

	global.EventBus.Subscribe("", events.PlaylistModeChangeUpdate(model.PlaylistIDSystem),
		"gui.config.basic.random_playlist.system",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistModeChangeUpdateEvent)
			systemRandomCheck.SetChecked(data.Mode == model.PlaylistModeRandom)
		}))

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
		_ = global.EventBus.Publish(events.PlayerSetAudioDeviceCmd, events.PlayerSetAudioDeviceCmdEvent{
			Device: name,
		})
	})
	global.EventBus.Subscribe("",
		events.PlayerAudioDeviceUpdate,
		"gui.config.basic.audio_device.update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
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
		}))

	outputDevice := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("gui.config.basic.audio_device")), nil,
		deviceSel)
	skipWhenErr := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.skip_when_error")),
		component.NewCheckOneWayBinding(
			i18n.T("gui.config.basic.skip_when_error.prompt"),
			&config.General.PlayNextOnFail,
			config.General.PlayNextOnFail),
	)
	checkUpdateBox := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.auto_check_update")),
		component.NewCheckOneWayBinding(
			i18n.T("gui.config.basic.auto_check_update.prompt"),
			&config.General.AutoCheckUpdate,
			config.General.AutoCheckUpdate),
	)
	checkUpdateBtn := widget.NewButton(i18n.T("gui.config.basic.check_update"), func() {
		_ = global.EventBus.Publish(events.CheckUpdateCmd, events.CheckUpdateCmdEvent{})
	})
	useSysPlaylistBtn := container.NewHBox(
		widget.NewLabel(i18n.T("gui.config.basic.use_system_playlist")),
		component.NewCheckOneWayBinding(
			i18n.T("gui.config.basic.use_system_playlist.prompt"),
			&config.General.UseSystemPlaylist,
			config.General.UseSystemPlaylist),
	)
	b.panel = container.NewVBox(randomPlaylist, useSysPlaylistBtn, skipWhenErr, outputDevice, checkUpdateBox, checkUpdateBtn)
	return b.panel
}
