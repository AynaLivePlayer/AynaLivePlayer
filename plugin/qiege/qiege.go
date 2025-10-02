package qiege

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"

	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type Qiege struct {
	config.BaseConfig
	UserPermission      bool
	PrivilegePermission bool
	AdminPermission     bool
	CustomCMD           string
	currentUid          string
	panel               fyne.CanvasObject
	log                 logger.ILogger
}

func NewQiege() *Qiege {
	return &Qiege{
		UserPermission:      true,
		PrivilegePermission: true,
		AdminPermission:     true,
		CustomCMD:           "切歌",
		log:                 global.Logger.WithPrefix("plugin.qiege"),
	}
}

func (d *Qiege) Name() string {
	return "Qiege"
}

func (d *Qiege) Enable() error {
	config.LoadConfig(d)
	gui.AddConfigLayout(d)
	global.EventBus.Subscribe("",
		events.LiveRoomMessageReceive,
		"plugin.qiege.message",
		d.handleMessage)
	global.EventBus.Subscribe("",
		events.PlayerPlayingUpdate,
		"plugin.qiege.playing",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlayerPlayingUpdateEvent)
			if data.Removed {
				d.currentUid = ""
			}
			if !data.Media.IsLiveRoomUser() {
				d.currentUid = ""
			}
			d.currentUid = data.Media.DanmuUser().Uid
		})
	return nil
}

func (d *Qiege) Disable() error {
	return nil
}

func (d *Qiege) handleMessage(event *eventbus.Event) {
	message := event.Data.(events.LiveRoomMessageReceiveEvent).Message
	msgs := strings.Split(message.Message, " ")
	if len(msgs) < 1 || msgs[0] != d.CustomCMD {
		return
	}
	d.log.Infof("recieve qiege command")
	if d.UserPermission {
		if d.currentUid == message.User.Uid {
			_ = global.EventBus.Publish(
				events.PlayerPlayNextCmd,
				events.PlayerPlayNextCmdEvent{})
			return
		}
	}
	if d.PrivilegePermission && message.User.Privilege > 0 {
		_ = global.EventBus.Publish(
			events.PlayerPlayNextCmd,
			events.PlayerPlayNextCmdEvent{})
		return
	}
	if d.AdminPermission && message.User.Admin {
		_ = global.EventBus.Publish(
			events.PlayerPlayNextCmd,
			events.PlayerPlayNextCmdEvent{})
		return
	}
}

func (d *Qiege) Title() string {
	return i18n.T("plugin.qiege.title")
}

func (d *Qiege) Description() string {
	return i18n.T("plugin.qiege.description")
}

func (d *Qiege) CreatePanel() fyne.CanvasObject {
	if d.panel != nil {
		return d.panel
	}
	dgPerm := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.qiege.permission")),
		component.NewCheckOneWayBinding(
			i18n.T("plugin.qiege.user"), &d.UserPermission, d.UserPermission),
		component.NewCheckOneWayBinding(
			i18n.T("plugin.qiege.privilege"), &d.PrivilegePermission, d.PrivilegePermission),
		component.NewCheckOneWayBinding(
			i18n.T("plugin.qiege.admin"), &d.AdminPermission, d.AdminPermission),
	)
	qgShortCut := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.qiege.custom_cmd")), nil,
		widget.NewEntryWithData(binding.BindString(&d.CustomCMD)),
	)
	d.panel = container.NewVBox(dgPerm, qgShortCut)
	return d.panel
}
