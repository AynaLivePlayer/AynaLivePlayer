package qiege

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/logger"
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
	UserPermission      bool
	PrivilegePermission bool
	AdminPermission     bool
	CustomCMD           string
	panel               fyne.CanvasObject
}

func NewQiege() *Qiege {
	return &Qiege{
		UserPermission:      true,
		PrivilegePermission: true,
		AdminPermission:     true,
		CustomCMD:           "skip",
	}
}

func (d *Qiege) Name() string {
	return "Qiege"
}

func (d *Qiege) Enable() error {
	config.LoadConfig(d)
	controller.AddCommand(d)
	gui.AddConfigLayout(d)
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
	if d.UserPermission && (controller.CurrentMedia != nil) {
		if controller.CurrentMedia.DanmuUser() != nil && controller.CurrentMedia.DanmuUser().Uid == danmu.User.Uid {
			controller.PlayNext()
			return
		}
	}
	if d.PrivilegePermission && danmu.User.Privilege > 0 {
		controller.PlayNext()
		return
	}
	if d.AdminPermission && danmu.User.Admin {
		controller.PlayNext()
		return
	}
}

func (d *Qiege) Title() string {
	return "Qiege"
}

func (d *Qiege) Description() string {
	return "Basic Qiege configuration"
}

func (d *Qiege) CreatePanel() fyne.CanvasObject {
	if d.panel != nil {
		return d.panel
	}

	dgPerm := container.NewHBox(
		widget.NewLabel("切歌权限"),
		widget.NewCheckWithData("切自己", binding.BindBool(&d.UserPermission)),
		widget.NewCheckWithData("舰长", binding.BindBool(&d.PrivilegePermission)),
		widget.NewCheckWithData("管理员", binding.BindBool(&d.AdminPermission)),
	)
	qgShortCut := container.NewBorder(nil, nil,
		widget.NewLabel("自定义命令 (默认的依然可用)"), nil,
		widget.NewEntryWithData(binding.BindString(&d.CustomCMD)),
	)
	d.panel = container.NewVBox(dgPerm, qgShortCut)
	return d.panel
}
