package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/player"
	"AynaLivePlayer/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/aynakeya/go-mpv"
)

type PlayControllerContainer struct {
	Title         *widget.Label
	Artist        *widget.Label
	Username      *widget.Label
	Cover         *canvas.Image
	ButtonPrev    *widget.Button
	ButtonSwitch  *widget.Button
	ButtonNext    *widget.Button
	Progress      *widget.Slider
	Volume        *widget.Slider
	ButtonLrc     *widget.Button
	LrcWindowOpen bool
	CurrentTime   *widget.Label
	TotalTime     *widget.Label
}

func (p *PlayControllerContainer) SetDefaultCover() {
	p.Cover.Resource = nil
	p.Cover.File = config.GetAssetPath("empty.png")
	p.Cover.Refresh()
}

var PlayController = &PlayControllerContainer{}

func createPlayController() fyne.CanvasObject {
	PlayController.Cover = canvas.NewImageFromFile(config.GetAssetPath("empty.png"))
	PlayController.Cover.SetMinSize(fyne.NewSize(128, 128))
	PlayController.Cover.FillMode = canvas.ImageFillContain

	PlayController.ButtonPrev = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {})
	PlayController.ButtonSwitch = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	PlayController.ButtonNext = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {})

	buttonsBox := container.NewCenter(
		container.NewHBox(PlayController.ButtonPrev, PlayController.ButtonSwitch, PlayController.ButtonNext))

	PlayController.Progress = widget.NewSlider(0, 1000)
	PlayController.CurrentTime = widget.NewLabel("0:00")
	PlayController.TotalTime = widget.NewLabel("0:00")
	progressItem := container.NewBorder(nil, nil, PlayController.CurrentTime, PlayController.TotalTime, PlayController.Progress)

	PlayController.Title = widget.NewLabel("Title")
	PlayController.Artist = widget.NewLabel("Artist")
	PlayController.Username = widget.NewLabel("Username")

	playInfo := container.NewVBox(PlayController.Title, PlayController.Artist, PlayController.Username)

	PlayController.Volume = widget.NewSlider(0, 100)
	volumeIcon := widget.NewIcon(theme.VolumeMuteIcon())
	PlayController.ButtonLrc = widget.NewButton(i18n.T("gui.player.button.lrc"), func() {})

	volumeControl := container.NewBorder(nil, nil, container.NewHBox(widget.NewLabel(" "), volumeIcon), nil,
		container.NewGridWithColumns(3, container.NewMax(PlayController.Volume), PlayController.ButtonLrc))

	registerPlayControllerHandler()

	return container.NewBorder(nil, nil, container.NewHBox(PlayController.Cover, playInfo, widget.NewSeparator()), nil,
		container.NewVBox(buttonsBox, progressItem, volumeControl))
}

func registerPlayControllerHandler() {
	PlayController.ButtonPrev.OnTapped = func() {
		controller.Seek(0, true)
	}
	PlayController.ButtonSwitch.OnTapped = func() {
		if controller.Toggle() {
			PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
		} else {
			PlayController.ButtonSwitch.Icon = theme.MediaStopIcon()
		}
	}
	PlayController.ButtonNext.OnTapped = func() {
		controller.PlayNext()
	}

	PlayController.ButtonLrc.OnTapped = func() {
		if !PlayController.LrcWindowOpen {
			PlayController.LrcWindowOpen = true
			createLyricWindow().Show()
		}
	}

	if controller.MainPlayer.ObserveProperty("pause", func(property *mpv.EventProperty) {
		if property.Data == nil {
			PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
			return
		}
		if property.Data.(mpv.Node).Value.(bool) {
			PlayController.ButtonSwitch.Icon = theme.MediaPlayIcon()
		} else {
			PlayController.ButtonSwitch.Icon = theme.MediaStopIcon()
		}
	}) != nil {
		l().Error("fail to register handler for switch button with property pause")
	}

	if controller.MainPlayer.ObserveProperty("percent-pos", func(property *mpv.EventProperty) {
		if property.Data == nil {
			PlayController.Progress.Value = 0
		} else {
			PlayController.Progress.Value = property.Data.(mpv.Node).Value.(float64) * 10
		}
		PlayController.Progress.Refresh()
	}) != nil {
		l().Error("fail to register handler for progress bar with property percent-pos")
	}

	if controller.MainPlayer.ObserveProperty("idle-active", func(property *mpv.EventProperty) {
		isIdle := property.Data.(mpv.Node).Value.(bool)
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
	PlayController.Progress.OnChanged = func(f float64) {
		controller.Seek(f/10, false)
	}

	if controller.MainPlayer.ObserveProperty("time-pos", func(property *mpv.EventProperty) {
		if property.Data == nil {
			PlayController.CurrentTime.SetText("0:00")
			return
		}
		PlayController.CurrentTime.SetText(util.FormatTime(int(property.Data.(mpv.Node).Value.(float64))))
	}) != nil {
		l().Error("fail to register handler for current time with property time-pos")
	}

	if controller.MainPlayer.ObserveProperty("duration", func(property *mpv.EventProperty) {
		if property.Data == nil {
			PlayController.TotalTime.SetText("0:00")
			return
		}
		PlayController.TotalTime.SetText(util.FormatTime(int(property.Data.(mpv.Node).Value.(float64))))
	}) != nil {
		l().Error("fail to register handler for total time with property duration")
	}

	if controller.MainPlayer.ObserveProperty("volume", func(property *mpv.EventProperty) {
		l().Trace("receive volume change event", *property)
		if property.Data == nil {
			PlayController.Volume.Value = 0
		} else {
			PlayController.Volume.Value = property.Data.(mpv.Node).Value.(float64)
		}

		PlayController.Volume.Refresh()
	}) != nil {
		l().Error("fail to register handler for progress bar with property percent-pos")
	}

	PlayController.Volume.OnChanged = func(f float64) {
		controller.SetVolume(f)
	}

	controller.MainPlayer.EventHandler.RegisterA(player.EventPlay, "gui.player.updateinfo", func(event *event.Event) {
		l().Debug("receive EventPlay update player info")
		media := event.Data.(player.PlayEvent).Media
		PlayController.Title.SetText(
			util.StringNormalize(media.Title, 16, 16))
		PlayController.Artist.SetText(
			util.StringNormalize(media.Artist, 16, 16))
		PlayController.Username.SetText(media.ToUser().Name)
		if media.Cover == "" {
			PlayController.SetDefaultCover()
		} else {
			uri, err := storage.ParseURI(media.Cover)
			if err != nil {
				l().Warn("fail to load parse cover url", media.Cover)
			}
			// async update
			go func() {
				img := canvas.NewImageFromURI(uri)
				if img == nil {
					l().Warn("fail to load parse cover url", media.Cover)
					PlayController.SetDefaultCover()
					return
				}
				PlayController.Cover.Resource = img.Resource
				PlayController.Cover.Refresh()
			}()
		}
	})
	return
}

func createPlayControllerV2() fyne.CanvasObject {
	PlayController.Cover = canvas.NewImageFromFile(config.GetAssetPath("empty.png"))
	PlayController.Cover.SetMinSize(fyne.NewSize(128, 128))
	PlayController.Cover.FillMode = canvas.ImageFillContain

	PlayController.ButtonPrev = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {})
	PlayController.ButtonSwitch = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	PlayController.ButtonNext = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {})

	buttonsBox := container.NewCenter(
		container.NewHBox(PlayController.ButtonPrev, PlayController.ButtonSwitch, PlayController.ButtonNext))

	PlayController.Progress = widget.NewSlider(0, 1000)
	PlayController.CurrentTime = widget.NewLabel("0:00")
	PlayController.TotalTime = widget.NewLabel("0:00")
	progressItem := container.NewBorder(nil, nil, PlayController.CurrentTime, PlayController.TotalTime, PlayController.Progress)

	PlayController.Title = widget.NewLabel("Title")
	PlayController.Title.TextStyle.Bold = true
	//a := canvas.NewText("asdf", color.Black)
	//a.TextSize = 12
	PlayController.Artist = widget.NewLabel("Artist")
	PlayController.Username = widget.NewLabel("Username")

	playInfo := container.NewBorder(PlayController.Username, nil, nil, PlayController.Artist, PlayController.Title)

	PlayController.Volume = widget.NewSlider(0, 100)
	volumeIcon := widget.NewIcon(theme.VolumeMuteIcon())
	PlayController.ButtonLrc = widget.NewButton(i18n.T("gui.player.button.lrc"), func() {})

	volumeControl := container.NewBorder(nil, nil, container.NewHBox(widget.NewLabel(" "), volumeIcon), nil,
		container.NewGridWithColumns(3, container.NewMax(PlayController.Volume), PlayController.ButtonLrc))

	registerPlayControllerHandler()

	return container.NewBorder(nil, nil, container.NewHBox(PlayController.Cover, widget.NewSeparator()), nil,
		container.NewVBox(playInfo, buttonsBox, progressItem, volumeControl))
}
