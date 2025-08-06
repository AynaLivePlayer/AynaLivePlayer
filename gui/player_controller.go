package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/util"
	"AynaLivePlayer/resource"
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlayControllerContainer struct {
	Title         *widget.Label
	Artist        *widget.Label
	Username      *widget.Label
	Cover         *canvas.Image
	coverLoader   context.CancelFunc
	ButtonPrev    *widget.Button
	ButtonSwitch  *widget.Button
	ButtonNext    *widget.Button
	Progress      *component.SliderPlus
	Volume        *widget.Slider
	ButtonLrc     *widget.Button
	ButtonPlayer  *widget.Button
	LrcWindowOpen bool
	CurrentTime   *widget.Label
	TotalTime     *widget.Label
}

func (p *PlayControllerContainer) SetDefaultCover() {
	p.Cover.Resource = resource.ImageEmpty
	p.Cover.Refresh()
}

var PlayController = &PlayControllerContainer{}

func registerPlayControllerHandler() {
	PlayController.ButtonPrev.OnTapped = func() {
		global.EventManager.CallA(events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
			Position: 0,
			Absolute: true,
		})
	}
	PlayController.ButtonSwitch.OnTapped = func() {
		global.EventManager.CallA(events.PlayerToggleCmd, events.PlayerToggleCmdEvent{})
	}
	PlayController.ButtonNext.OnTapped = func() {
		global.EventManager.CallA(events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
	}

	PlayController.ButtonLrc.OnTapped = func() {
		if !PlayController.LrcWindowOpen {
			PlayController.LrcWindowOpen = true
			createLyricWindow().Show()
		}
	}

	PlayController.ButtonPlayer.OnTapped = func() {
		showPlayerWindow()
	}

	global.EventManager.RegisterA(events.PlayerPropertyPauseUpdate, "gui.player.controller.paused", gutil.ThreadSafeHandler(func(event *event.Event) {
		if event.Data.(events.PlayerPropertyPauseUpdateEvent).Paused {
			PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
		} else {
			PlayController.ButtonSwitch.Icon = theme.MediaPauseIcon()
		}
		PlayController.ButtonSwitch.Refresh()
	}))

	global.EventManager.RegisterA(events.PlayerPropertyPercentPosUpdate, "gui.player.controller.percent_pos", gutil.ThreadSafeHandler(func(event *event.Event) {
		if PlayController.Progress.Dragging {
			return
		}
		PlayController.Progress.Value = event.Data.(events.PlayerPropertyPercentPosUpdateEvent).PercentPos * 10
		PlayController.Progress.Refresh()
	}))

	global.EventManager.RegisterA(events.PlayerPropertyStateUpdate, "gui.player.controller.idle_active", gutil.ThreadSafeHandler(func(event *event.Event) {
		state := event.Data.(events.PlayerPropertyStateUpdateEvent).State
		if state == model.PlayerStateIdle || state == model.PlayerStateLoading {
			PlayController.Progress.Value = 0
			PlayController.Progress.Max = 0
			//PlayController.Title.SetText("Title")
			//PlayController.Artist.SetText("Artist")
			//PlayController.Username.SetText("Username")
			//PlayController.SetDefaultCover()
		} else {
			PlayController.Progress.Max = 1000
		}
	}))

	PlayController.Progress.Max = 0
	PlayController.Progress.OnDragEnd = func(f float64) {
		global.EventManager.CallA(events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
			Position: f / 10,
			Absolute: false,
		})
	}

	global.EventManager.RegisterA(events.PlayerPropertyTimePosUpdate, "gui.player.controller.time_pos", gutil.ThreadSafeHandler(func(event *event.Event) {
		PlayController.CurrentTime.SetText(util.FormatTime(int(event.Data.(events.PlayerPropertyTimePosUpdateEvent).TimePos)))
	}))

	global.EventManager.RegisterA(events.PlayerPropertyDurationUpdate, "gui.player.controller.duration", gutil.ThreadSafeHandler(func(event *event.Event) {
		PlayController.TotalTime.SetText(util.FormatTime(int(event.Data.(events.PlayerPropertyDurationUpdateEvent).Duration)))
	}))

	global.EventManager.RegisterA(events.PlayerPropertyVolumeUpdate, "gui.player.controller.volume", gutil.ThreadSafeHandler(func(event *event.Event) {
		PlayController.Volume.Value = event.Data.(events.PlayerPropertyVolumeUpdateEvent).Volume
		PlayController.Volume.Refresh()
	}))

	PlayController.Volume.OnChanged = func(f float64) {
		global.EventManager.CallA(events.PlayerVolumeChangeCmd, events.PlayerVolumeChangeCmdEvent{
			Volume: f,
		})
	}

	// todo: double check cover loading for new thread model
	global.EventManager.RegisterA(events.PlayerPlayingUpdate, "gui.player.updateinfo", gutil.ThreadSafeHandler(func(event *event.Event) {
		if event.Data.(events.PlayerPlayingUpdateEvent).Removed {
			PlayController.Progress.Value = 0
			PlayController.Progress.Max = 0
			PlayController.TotalTime.SetText("0:00")
			PlayController.CurrentTime.SetText("0:00")
			PlayController.Title.SetText("Title")
			PlayController.Artist.SetText("Artist")
			PlayController.Username.SetText("Username")
			PlayController.SetDefaultCover()
			return
		}
		media := event.Data.(events.PlayerPlayingUpdateEvent).Media
		//PlayController.Title.SetText(
		//	util.StringNormalize(media.Title, 16, 16))
		//PlayController.Artist.SetText(
		//	util.StringNormalize(media.Artist, 16, 16))
		PlayController.Title.SetText(
			media.Info.Title)
		PlayController.Artist.SetText(
			media.Info.Artist)
		PlayController.Username.SetText(media.ToUser().Name)
		if !media.Info.Cover.Exists() {
			PlayController.SetDefaultCover()
		} else {
			if PlayController.coverLoader != nil {
				PlayController.coverLoader()
			}
			var ctx context.Context
			ctx, PlayController.coverLoader = context.WithCancel(context.Background())
			go func() {
				ch := make(chan *canvas.Image)
				go func() {
					picture, err := gutil.NewImageFromPlayerPicture(media.Info.Cover)
					if err != nil {
						ch <- nil
						logger.Errorf("fail to load cover: %v", err)
						return
					}
					ch <- picture
				}()
				select {
				case <-ctx.Done():
					return
				case pic := <-ch:
					if pic == nil {
						PlayController.SetDefaultCover()
						return
					}
					PlayController.Cover.Resource = pic.Resource
					gutil.RunInFyneThread(PlayController.Cover.Refresh)
				}
			}()
		}
	}))
}

func createPlayControllerV2() fyne.CanvasObject {
	PlayController.Cover = canvas.NewImageFromResource(resource.ImageEmpty)
	PlayController.Cover.SetMinSize(fyne.NewSize(128, 128))
	PlayController.Cover.FillMode = canvas.ImageFillContain

	PlayController.ButtonPrev = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {})
	PlayController.ButtonSwitch = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	PlayController.ButtonNext = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {})

	PlayController.Volume = widget.NewSlider(0, 100)
	PlayController.ButtonLrc = widget.NewButton(i18n.T("gui.player.button.lrc"), func() {})
	PlayController.ButtonPlayer = widget.NewButton(i18n.T("gui.player.button.player"), func() {})

	buttonsBox := container.NewHBox(
		PlayController.ButtonPrev, PlayController.ButtonSwitch, PlayController.ButtonNext,
	)

	buttonBox2 := container.NewHBox(
		PlayController.ButtonLrc, PlayController.ButtonPlayer)

	volumeControl := component.NewFixedHSplitContainer(
		widget.NewLabel(""),
		container.NewBorder(nil, nil,
			widget.NewIcon(theme.VolumeUpIcon()),
			widget.NewLabel("       "),
			PlayController.Volume), 0.05)
	volumeControl.SeparatorThickness = 0

	controls := component.NewFixedHSplitContainer(
		container.NewBorder(nil, nil, nil, buttonBox2, buttonsBox),
		volumeControl,
		0.4)
	controls.SeparatorThickness = 0

	PlayController.Progress = component.NewSliderPlus(0, 1000)
	PlayController.CurrentTime = widget.NewLabel("0:00")
	PlayController.TotalTime = widget.NewLabel("0:00")
	progressItem := container.NewBorder(nil, nil,
		PlayController.CurrentTime,
		PlayController.TotalTime,
		PlayController.Progress)

	PlayController.Title = widget.NewLabel("Title")
	PlayController.Title.Wrapping = fyne.TextTruncate
	PlayController.Artist = widget.NewLabel("Artist")
	PlayController.Username = widget.NewLabel("Username")

	titleUser := component.NewFixedHSplitContainer(
		PlayController.Title, PlayController.Artist, 0.32)
	titleUser.SetSepThickness(0)

	playInfo := container.NewBorder(nil, nil, nil, PlayController.Username,
		titleUser)

	registerPlayControllerHandler()

	return container.NewBorder(nil, nil, container.NewHBox(PlayController.Cover, widget.NewSeparator()), nil,
		container.NewVBox(playInfo, progressItem, controls))
}
