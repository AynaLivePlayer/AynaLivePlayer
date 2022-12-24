package gui

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/model"
	"AynaLivePlayer/resource"
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
	ButtonPrev    *widget.Button
	ButtonSwitch  *widget.Button
	ButtonNext    *widget.Button
	Progress      *component.SliderPlus
	Volume        *widget.Slider
	ButtonLrc     *widget.Button
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
		controller.Instance.PlayControl().Seek(0, true)
	}
	PlayController.ButtonSwitch.OnTapped = func() {
		controller.Instance.PlayControl().Toggle()
	}
	PlayController.ButtonNext.OnTapped = func() {
		controller.Instance.PlayControl().PlayNext()
	}

	PlayController.ButtonLrc.OnTapped = func() {
		if !PlayController.LrcWindowOpen {
			PlayController.LrcWindowOpen = true
			createLyricWindow().Show()
		}
	}

	if controller.Instance.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropPause, "gui.play_controller.pause", func(ev *event.Event) {
			data := ev.Data.(model.PlayerPropertyUpdateEvent).Value
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

	if controller.Instance.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropPercentPos, "gui.play_controller.percent_pos", func(ev *event.Event) {
			if PlayController.Progress.Dragging {
				return
			}
			data := ev.Data.(model.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.Progress.Value = 0
			} else {
				PlayController.Progress.Value = data.(float64) * 10
			}
			PlayController.Progress.Refresh()
		}) != nil {
		l().Error("fail to register handler for progress bar with property percent-pos")
	}

	if controller.Instance.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropIdleActive, "gui.play_controller.idle_active", func(ev *event.Event) {
			isIdle := ev.Data.(model.PlayerPropertyUpdateEvent).Value.(bool)
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
		controller.Instance.PlayControl().Seek(f/10, false)
	}

	if controller.Instance.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropTimePos, "gui.play_controller.time_pos", func(ev *event.Event) {
			data := ev.Data.(model.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.CurrentTime.SetText("0:00")
				return
			}
			PlayController.CurrentTime.SetText(util.FormatTime(int(data.(float64))))
		}) != nil {
		l().Error("fail to register handler for current time with property time-pos")
	}

	if controller.Instance.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropDuration, "gui.play_controller.duration", func(ev *event.Event) {
			data := ev.Data.(model.PlayerPropertyUpdateEvent).Value
			if data == nil {
				PlayController.TotalTime.SetText("0:00")
				return
			}
			PlayController.TotalTime.SetText(util.FormatTime(int(data.(float64))))
		}) != nil {
		l().Error("fail to register handler for total time with property duration")
	}

	if controller.Instance.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropVolume, "gui.play_controller.volume", func(ev *event.Event) {
			data := ev.Data.(model.PlayerPropertyUpdateEvent).Value
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
		controller.Instance.PlayControl().SetVolume(f)
	}

	controller.Instance.PlayControl().EventManager().RegisterA(model.EventPlay, "gui.player.updateinfo", func(event *event.Event) {
		l().Debug("receive EventPlay update player info")
		media := event.Data.(model.PlayEvent).Media
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
			go func() {
				picture, err := gutil.NewImageFromPlayerPicture(media.Cover)
				if err != nil {
					l().Warn("fail to load parse cover url", media.Cover)
					PlayController.SetDefaultCover()
					return
				}
				PlayController.Cover.Resource = picture.Resource
				PlayController.Cover.Refresh()
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

	buttonsBox := container.NewHBox(PlayController.ButtonPrev, PlayController.ButtonSwitch, PlayController.ButtonNext)

	PlayController.Volume = widget.NewSlider(0, 100)
	PlayController.ButtonLrc = widget.NewButton(i18n.T("gui.player.button.lrc"), func() {})

	controls := container.NewPadded(container.NewBorder(nil, nil,
		buttonsBox, nil,
		container.NewGridWithColumns(
			2,
			container.NewMax(),
			container.NewBorder(nil, nil, widget.NewIcon(theme.VolumeMuteIcon()), PlayController.ButtonLrc,
				PlayController.Volume)),
	))

	PlayController.Progress = component.NewSliderPlus(0, 1000)
	PlayController.CurrentTime = widget.NewLabel("0:00")
	PlayController.TotalTime = widget.NewLabel("0:00")
	progressItem := container.NewBorder(nil, nil,
		PlayController.CurrentTime,
		PlayController.TotalTime,
		PlayController.Progress)

	PlayController.Title = widget.NewLabel("Title")
	PlayController.Artist = widget.NewLabel("Artist")
	PlayController.Username = widget.NewLabel("Username")

	playInfo := container.NewBorder(nil, nil, nil, PlayController.Username,
		container.NewHBox(PlayController.Title, PlayController.Artist))

	registerPlayControllerHandler()

	return container.NewBorder(nil, nil, container.NewHBox(PlayController.Cover, widget.NewSeparator()), nil,
		container.NewVBox(playInfo, progressItem, controls))
}
