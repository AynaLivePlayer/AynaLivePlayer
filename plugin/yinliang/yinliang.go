package yinliang

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	config2 "AynaLivePlayer/gui/views/config"

	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Yinliang struct {
	config.BaseConfig
	Enabled         bool
	AdminPermission bool
	VolumeUpCMD     string
	VolumeDownCMD   string
	VolumeStep      float64
	MaxVolume       float64
	currentVolume   float64
	panel           fyne.CanvasObject
	log             logger.ILogger
}

func NewYinliang() *Yinliang {
	return &Yinliang{
		Enabled:         false,
		AdminPermission: true,
		VolumeUpCMD:     "音量调大",
		VolumeDownCMD:   "音量调小",
		VolumeStep:      5.0,
		MaxVolume:       75.0,
		currentVolume:   50.0,
		log:             global.Logger.WithPrefix("Plugin.Yiliang"),
	}
}

func (y *Yinliang) Name() string {
	return "Yinliang"
}

func (y *Yinliang) Enable() error {
	config.LoadConfig(y)

	// 配置校验
	if y.VolumeStep > 25 {
		y.VolumeStep = 25
	} else if y.VolumeStep < 0 {
		y.VolumeStep = 5
	}
	if y.MaxVolume > 100 {
		y.MaxVolume = 100
	} else if y.MaxVolume < 0 {
		y.MaxVolume = 0
	}

	config2.AddConfigLayout(y)

	_ = global.EventBus.Subscribe("",
		events.LiveRoomMessageReceive,
		"plugin.yinliang.message",
		y.handleMessage)

	_ = global.EventBus.Subscribe("",
		events.PlayerVolumeChangeCmd,
		"plugin.yinliang.volume_tracker",
		func(e *eventbus.Event) {
			data := e.Data.(events.PlayerVolumeChangeCmdEvent)
			y.currentVolume = data.Volume
		})
	return nil
}

func (y *Yinliang) Disable() error {
	return nil
}

func (y *Yinliang) handleMessage(event *eventbus.Event) {
	if !y.Enabled {
		return
	}
	message := event.Data.(events.LiveRoomMessageReceiveEvent).Message
	cmd := strings.TrimSpace(message.Message)

	if cmd != y.VolumeUpCMD && cmd != y.VolumeDownCMD {
		return
	}

	if !y.AdminPermission || !message.User.Admin {
		y.log.Infof("User <%s> modify volume failed: no permission", message.User.Username)
		return
	}

	delta := y.VolumeStep
	if cmd == y.VolumeDownCMD {
		delta = -y.VolumeStep
	}

	newVolume := y.currentVolume + delta
	if newVolume > y.MaxVolume {
		newVolume = y.MaxVolume
	} else if newVolume < 0 {
		newVolume = 0
	}
	y.log.Infof("User <%s> modify volume from %.2f to %.2f", message.User.Username, y.currentVolume, newVolume)

	_ = global.EventBus.Publish(
		events.PlayerVolumeChangeCmd,
		events.PlayerVolumeChangeCmdEvent{
			Volume: newVolume,
		})
}

func (y *Yinliang) Title() string {
	return i18n.T("plugin.yinliang.title")
}

func (y *Yinliang) Description() string {
	return i18n.T("plugin.yinliang.description")
}

// 在CreatePanel方法中修改控件创建方式
func (y *Yinliang) CreatePanel() fyne.CanvasObject {
	if y.panel != nil {
		return y.panel
	}

	enabledCheck := component.NewCheckOneWayBinding(i18n.T("plugin.yinliang.enabled"),
		&y.Enabled,
		y.Enabled)

	permCheck := component.NewCheckOneWayBinding(
		i18n.T("plugin.yinliang.admin_permission"),
		&y.AdminPermission,
		y.AdminPermission)

	cmdConfig := container.NewGridWithColumns(2,
		widget.NewLabel(i18n.T("plugin.yinliang.volume_up_cmd")),
		widget.NewEntryWithData(binding.BindString(&y.VolumeUpCMD)),
		widget.NewLabel(i18n.T("plugin.yinliang.volume_down_cmd")),
		widget.NewEntryWithData(binding.BindString(&y.VolumeDownCMD)),
	)

	stepEntry := widget.NewEntryWithData(binding.FloatToStringWithFormat(binding.BindFloat(&y.VolumeStep), "%.1f"))
	stepEntry.Validator = createFloatValidator(0, 25)
	stepEntry.OnChanged = func(s string) {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			if v > 25 {
				y.VolumeStep = 25
				stepEntry.SetText("25")
			} else if v < 0 {
				y.VolumeStep = 5
				stepEntry.SetText("5")
			}
		}
	}

	maxVolEntry := widget.NewEntryWithData(binding.FloatToStringWithFormat(binding.BindFloat(&y.MaxVolume), "%.1f"))
	maxVolEntry.Validator = createFloatValidator(0, 100)
	maxVolEntry.OnChanged = func(s string) {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			if v > 100 {
				y.MaxVolume = 100
				maxVolEntry.SetText("100")
			} else if v < 0 {
				y.MaxVolume = 0
				maxVolEntry.SetText("0")
			}
		}
	}

	volumeControlConfig := container.NewGridWithColumns(2,
		widget.NewLabel(i18n.T("plugin.yinliang.volume_step")),
		stepEntry,
		widget.NewLabel(i18n.T("plugin.yinliang.max_volume")),
		maxVolEntry,
	)

	y.panel = container.NewVBox(
		enabledCheck,
		permCheck,
		cmdConfig,
		volumeControlConfig,
	)
	return y.panel
}

func createFloatValidator(min, max float64) func(string) error {
	return func(s string) error {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf(i18n.T("validation.number_required"))
		}
		if v < min || v > max {
			return fmt.Errorf(i18n.T("validation.range_error"))
		}
		return nil
	}
}
