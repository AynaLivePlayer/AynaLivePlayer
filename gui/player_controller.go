package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/gutil"
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
		API.PlayControl().Seek(0, true)
	}
	PlayController.ButtonSwitch.OnTapped = func() {
		API.PlayControl().Toggle()
	}
	PlayController.ButtonNext.OnTapped = func() {
		API.PlayControl().PlayNext()
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

	if API.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropPause, "gui.play_controller.pause", func(ev *event.Event) {
			data := ev.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
				return
			}
			if data.(bool) {
				PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
			} else {
				PlayController.ButtonSwitch.Icon = theme.MediaPauseIcon()
			}
		}) != nil {
		l().Error("fail to register handler for switch button with property pause")
	}

	if API.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropPercentPos, "gui.play_controller.percent_pos", func(ev *event.Event) {
			if PlayController.Progress.Dragging {
				return
			}
			data := ev.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.Progress.Value = 0
			} else {
				PlayController.Progress.Value = data.(float64) * 10
			}
			PlayController.Progress.Refresh()
		}) != nil {
		l().Error("fail to register handler for progress bar with property percent-pos")
	}

	if API.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropIdleActive, "gui.play_controller.idle_active", func(ev *event.Event) {
			isIdle := ev.Data.(events.PlayerPropertyUpdateEvent).Value.(bool)
			l().Debug("receive idle active ", isIdle, " set/reset info")
			// todo: @3
			if isIdle {
				PlayController.Progress.Value = 0
				PlayController.Progress.Max = 0
				//PlayController.Title.SetText("Title")
				//PlayController.Artist.SetText("Artist")
				//PlayController.Username.SetText("Username")
				//PlayController.SetDefaultCover()
			} else {
				PlayController.Progress.Max = 1000
			}
		}) != nil {
		l().Error("fail to register handler for progress bar with property idle-active")
	}

	PlayController.Progress.Max = 0
	PlayController.Progress.OnDragEnd = func(f float64) {
		API.PlayControl().Seek(f/10, false)
	}

	if API.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropTimePos, "gui.play_controller.time_pos", func(ev *event.Event) {
			data := ev.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.CurrentTime.SetText("0:00")
				return
			}
			PlayController.CurrentTime.SetText(util.FormatTime(int(data.(float64))))
		}) != nil {
		l().Error("fail to register handler for current time with property time-pos")
	}

	if API.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropDuration, "gui.play_controller.duration", func(ev *event.Event) {
			data := ev.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.TotalTime.SetText("0:00")
				return
			}
			PlayController.TotalTime.SetText(util.FormatTime(int(data.(float64))))
		}) != nil {
		l().Error("fail to register handler for total time with property duration")
	}

	if API.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropVolume, "gui.play_controller.volume", func(ev *event.Event) {
			data := ev.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.Volume.Value = 0
			} else {
				PlayController.Volume.Value = data.(float64)
			}

			PlayController.Volume.Refresh()
		}) != nil {
		l().Error("fail to register handler for progress bar with property percent-pos")
	}

	PlayController.Volume.OnChanged = func(f float64) {
		API.PlayControl().SetVolume(f)
	}

	API.PlayControl().EventManager().RegisterA(events.EventPlay, "gui.player.updateinfo", func(event *event.Event) {
		l().Debug("receive EventPlay update player info")
		media := event.Data.(events.PlayEvent).Media
		//PlayController.Title.SetText(
		//	util.StringNormalize(media.Title, 16, 16))
		//PlayController.Artist.SetText(
		//	util.StringNormalize(media.Artist, 16, 16))
		PlayController.Title.SetText(
			media.Title)
		PlayController.Artist.SetText(
			media.Artist)
		PlayController.Username.SetText(media.ToUser().Name)
		if !media.Cover.Exists() {
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
					picture, err := gutil.NewImageFromPlayerPicture(media.Cover)
					if err != nil {
						ch <- nil
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
					PlayController.Cover.Refresh()
				}

			}()
		}
	})
	return
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
			widget.NewIcon(theme.VolumeMuteIcon()),
			widget.NewLabel("       "),
			PlayController.Volume), 0.05)
	volumeControl.SeparatorThickness = 0

	controls := component.NewFixedHSplitContainer(
		container.NewBorder(nil, nil, nil, buttonBox2, buttonsBox),
		volumeControl,
		0.4)
	controls.SeparatorThickness = 0

	//controls := container.NewPadded(container.NewBorder(nil, nil,
	//	buttonsBox, nil,
	//	container.NewGridWithColumns(
	//		2,
	//		container.NewMax(),
	//		container.NewBorder(nil, nil, widget.NewIcon(theme.VolumeMuteIcon()), PlayController.ButtonLrc,
	//			PlayController.Volume)),
	//))

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
