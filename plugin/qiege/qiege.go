package qiege

import (
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/liveclient"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

const MODULE_CMD_QieGE = "CMD.QieGe"

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_CMD_QieGE)
}

type Qiege struct {
	config.BaseConfig
	UserPermission      bool
	PrivilegePermission bool
	AdminPermission     bool
	CustomCMD           string
	panel               fyne.CanvasObject
	ctr                 controller.IController
}

func NewQiege(ctr controller.IController) *Qiege {
	return &Qiege{
		UserPermission:      true,
		PrivilegePermission: true,
		AdminPermission:     true,
		CustomCMD:           "skip",
		ctr:                 ctr,
	}
}

func (d *Qiege) Name() string {
	return "Qiege"
}

func (d *Qiege) Enable() error {
	config.LoadConfig(d)
	d.ctr.LiveRooms().AddDanmuCommand(d)
	gui.AddConfigLayout(d)
	return nil
}

func (d *Qiege) Disable() error {
	return nil
}

func (d *Qiege) Match(command string) bool {
	for _, c := range []string{"切歌", d.CustomCMD} {
		if command == c {
			return true
		}
	}
	return false
}

func (d *Qiege) Execute(command string, args []string, danmu *liveclient.DanmuMessage) {
	if d.UserPermission && (d.ctr.PlayControl().GetPlaying() != nil) {
		if d.ctr.PlayControl().GetPlaying().DanmuUser() != nil && d.ctr.PlayControl().GetPlaying().DanmuUser().Uid == danmu.User.Uid {
			d.ctr.PlayControl().PlayNext()
			return
		}
	}
	if d.PrivilegePermission && danmu.User.Privilege > 0 {
		d.ctr.PlayControl().PlayNext()
		return
	}
	if d.AdminPermission && danmu.User.Admin {
		d.ctr.PlayControl().PlayNext()
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
		widget.NewCheckWithData(i18n.T("plugin.qiege.user"), binding.BindBool(&d.UserPermission)),
		widget.NewCheckWithData(i18n.T("plugin.qiege.privilege"), binding.BindBool(&d.PrivilegePermission)),
		widget.NewCheckWithData(i18n.T("plugin.qiege.admin"), binding.BindBool(&d.AdminPermission)),
	)
	qgShortCut := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.qiege.custom_cmd")), nil,
		widget.NewEntryWithData(binding.BindString(&d.CustomCMD)),
	)
	d.panel = container.NewVBox(dgPerm, qgShortCut)
	return d.panel
}
