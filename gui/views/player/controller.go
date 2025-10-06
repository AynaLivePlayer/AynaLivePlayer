package player

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
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
		_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
			Position: 0,
			Absolute: true,
		})
	}
	PlayController.ButtonSwitch.OnTapped = func() {
		_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerToggleCmd, events.PlayerToggleCmdEvent{})
	}
	PlayController.ButtonNext.OnTapped = func() {
		_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
	}

	PlayController.ButtonLrc.OnTapped = func() {
		if !PlayController.LrcWindowOpen {
			PlayController.LrcWindowOpen = true
			createLyricWindowV2().Show()
		}
	}

	PlayController.ButtonPlayer.OnTapped = func() {
		showPlayerWindow()
	}

	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPropertyPauseUpdate, "gui.player.controller.paused", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		if event.Data.(events.PlayerPropertyPauseUpdateEvent).Paused {
			PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
		} else {
			PlayController.ButtonSwitch.Icon = theme.MediaPauseIcon()
		}
		PlayController.ButtonSwitch.Refresh()
	}))

	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPropertyPercentPosUpdate, "gui.player.controller.percent_pos", func(event *eventbus.Event) {
		if PlayController.Progress.Dragging {
			return
		}
		PlayController.Progress.Value = event.Data.(events.PlayerPropertyPercentPosUpdateEvent).PercentPos * 10
		gutil.RunInFyneThread(PlayController.Progress.Refresh)
	})

	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPropertyStateUpdate, "gui.player.controller.idle_active", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
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
		_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
			Position: f / 10,
			Absolute: false,
		})
	}

	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPropertyTimePosUpdate, "gui.player.controller.time_pos", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		PlayController.CurrentTime.SetText(util.FormatTime(int(event.Data.(events.PlayerPropertyTimePosUpdateEvent).TimePos)))
	}))

	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPropertyDurationUpdate, "gui.player.controller.duration", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		PlayController.TotalTime.SetText(util.FormatTime(int(event.Data.(events.PlayerPropertyDurationUpdateEvent).Duration)))
	}))

	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPropertyVolumeUpdate, "gui.player.controller.volume", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
		PlayController.Volume.Value = event.Data.(events.PlayerPropertyVolumeUpdateEvent).Volume
		PlayController.Volume.Refresh()
	}))

	PlayController.Volume.OnChanged = func(f float64) {
		_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerVolumeChangeCmd, events.PlayerVolumeChangeCmdEvent{
			Volume: f,
		})
	}

	// todo: double check cover loading for new thread model
	global.EventBus.Subscribe(gctx.EventChannel, events.PlayerPlayingUpdate, "gui.player.updateinfo", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
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
						gctx.Logger.Errorf("fail to load cover: %v", err)
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
	PlayController.CurrentTime = widget.NewLabel("00:00")
	PlayController.TotalTime = widget.NewLabel("00:00")
	progressItem := container.NewBorder(nil, nil,
		nil,
		PlayController.TotalTime,
		component.NewFixedHSplitContainer(PlayController.CurrentTime, PlayController.Progress, 0.1))

	PlayController.Title = widget.NewLabel("Title")
	PlayController.Title.Truncation = fyne.TextTruncateClip
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
