package diange

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
	"strings"
	"time"
)

const MODULE_CMD_DIANGE = "CMD.DianGe"

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_CMD_DIANGE)
}

type Diange struct {
	UserPermission      bool
	PrivilegePermission bool
	AdminPermission     bool
	QueueMax            int
	UserCoolDown        int
	CustomCMD           string
	cooldowns           map[string]int
	panel               fyne.CanvasObject
}

func NewDiange() *Diange {
	return &Diange{
		UserPermission:      true,
		PrivilegePermission: true,
		AdminPermission:     true,
		QueueMax:            128,
		UserCoolDown:        -1,
		CustomCMD:           "add",
		cooldowns:           make(map[string]int),
	}
}

func (d *Diange) Name() string {
	return "Diange"
}

func (d *Diange) Enable() error {
	config.LoadConfig(d)
	controller.AddCommand(d)
	gui.AddConfigLayout(d)
	return nil
}

func (d *Diange) Match(command string) bool {
	for _, c := range []string{"点歌", d.CustomCMD} {
		if command == c {
			return true
		}
	}
	return false
}

func (d *Diange) Execute(command string, args []string, danmu *liveclient.DanmuMessage) {
	// if queue is full, return
	if controller.UserPlaylist.Size() >= d.QueueMax {
		l().Info("Queue is full, ignore diange")
		return
	}
	// if in user cool down, return
	ct := int(time.Now().Unix())
	if (ct - d.cooldowns[danmu.User.Uid]) <= d.UserCoolDown {
		l().Infof("User %s(%s) still in cool down period, diange failed", danmu.User.Username, danmu.User.Uid)
		return
	}
	keyword := strings.Join(args, " ")
	perm := d.UserPermission
	l().Trace("user permission check: ", perm)
	perm = perm || (d.PrivilegePermission && (danmu.User.Privilege > 0))
	l().Trace("privilege permission check: ", perm)
	perm = perm || (d.AdminPermission && (danmu.User.Admin))
	l().Trace("admin permission check: ", perm)
	if perm {
		// reset cool down
		d.cooldowns[danmu.User.Uid] = ct
		controller.Add(keyword, &danmu.User)
	}
}

func (d *Diange) Title() string {
	return "Diange"
}

func (d *Diange) Description() string {
	return "Basic diange configuration"
}

func (d *Diange) CreatePanel() fyne.CanvasObject {
	if d.panel != nil {
		return d.panel
	}
	dgPerm := container.NewHBox(
		widget.NewLabel("点歌权限"),
		widget.NewCheckWithData("用户", binding.BindBool(&d.UserPermission)),
		widget.NewCheckWithData("舰长", binding.BindBool(&d.PrivilegePermission)),
		widget.NewCheckWithData("管理员", binding.BindBool(&d.AdminPermission)),
	)
	dgQueue := container.NewBorder(nil, nil,
		widget.NewLabel("最大点歌数"), nil,
		widget.NewEntryWithData(binding.IntToString(binding.BindInt(&d.QueueMax))),
	)
	dgCoolDown := container.NewBorder(nil, nil,
		widget.NewLabel("点歌冷却时间"), nil,
		widget.NewEntryWithData(binding.IntToString(binding.BindInt(&d.UserCoolDown))),
	)
	dgShortCut := container.NewBorder(nil, nil,
		widget.NewLabel("自定义命令 (默认的依然可用)"), nil,
		widget.NewEntryWithData(binding.BindString(&d.CustomCMD)),
	)
	d.panel = container.NewVBox(dgPerm, dgQueue, dgCoolDown, dgShortCut)
	return d.panel
}
