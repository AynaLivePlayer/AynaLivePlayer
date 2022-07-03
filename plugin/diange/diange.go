package diange

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/liveclient"
	"AynaLivePlayer/logger"
	"fmt"
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
	SourceCMD           []string
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
		SourceCMD:           make([]string, 0),
		cooldowns:           make(map[string]int),
	}
}

func (d *Diange) Name() string {
	return "Diange"
}

func (d *Diange) Enable() error {
	config.LoadConfig(d)
	d.initCMD()
	controller.AddCommand(d)
	gui.AddConfigLayout(d)
	return nil
}

func (d *Diange) initCMD() {
	if len(d.SourceCMD) == len(config.Provider.Priority) {
		return
	}
	if len(d.SourceCMD) > len(config.Provider.Priority) {
		d.SourceCMD = d.SourceCMD[:len(config.Provider.Priority)]
		return
	}
	for i := len(d.SourceCMD); i < len(config.Provider.Priority); i++ {
		d.SourceCMD = append(d.SourceCMD, "点歌"+config.Provider.Priority[i])
	}
}

// isCMD return int if the commmand name matches our command
// -1 = not match, 0 = normal command, 1+ = source command
func (d *Diange) isCMD(cmd string) int {
	if cmd == "点歌" || cmd == d.CustomCMD {
		return 0
	}
	fmt.Println(d.SourceCMD)
	for index, c := range d.SourceCMD {
		if cmd == c {
			return index + 1
		}
	}
	return -1
}

func (d *Diange) Match(command string) bool {
	return d.isCMD(command) >= 0
}

func (d *Diange) Execute(command string, args []string, danmu *liveclient.DanmuMessage) {
	l().Infof("%s(%s) Execute command: %s %s", danmu.User.Username, danmu.User.Uid, command, args)
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
	cmdType := d.isCMD(command)
	keyword := strings.Join(args, " ")
	perm := d.UserPermission
	l().Trace("user permission check: ", perm)
	perm = perm || (d.PrivilegePermission && (danmu.User.Privilege > 0))
	l().Trace("privilege permission check: ", perm)
	perm = perm || (d.AdminPermission && (danmu.User.Admin))
	l().Trace("admin permission check: ", perm)
	if !perm {
		return
	}
	// reset cool down
	d.cooldowns[danmu.User.Uid] = ct
	if cmdType == 0 {
		controller.Add(keyword, &danmu.User)
	} else {
		controller.AddWithProvider(keyword, config.Provider.Priority[cmdType-1], &danmu.User)
	}
}

func (d *Diange) Title() string {
	return i18n.T("plugin.diange.title")
}

func (d *Diange) Description() string {
	return i18n.T("plugin.diange.description")
}

func (d *Diange) CreatePanel() fyne.CanvasObject {
	if d.panel != nil {
		return d.panel
	}
	dgPerm := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.diange.permission")),
		widget.NewCheckWithData(i18n.T("plugin.diange.user"), binding.BindBool(&d.UserPermission)),
		widget.NewCheckWithData(i18n.T("plugin.diange.privilege"), binding.BindBool(&d.PrivilegePermission)),
		widget.NewCheckWithData(i18n.T("plugin.diange.admin"), binding.BindBool(&d.AdminPermission)),
	)
	dgQueue := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.diange.queue_max")), nil,
		widget.NewEntryWithData(binding.IntToString(binding.BindInt(&d.QueueMax))),
	)
	dgCoolDown := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.diange.cooldown")), nil,
		widget.NewEntryWithData(binding.IntToString(binding.BindInt(&d.UserCoolDown))),
	)
	dgShortCut := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.diange.custom_cmd")), nil,
		widget.NewEntryWithData(binding.BindString(&d.CustomCMD)),
	)
	sourceCmds := []fyne.CanvasObject{}
	for i, _ := range d.SourceCMD {
		sourceCmds = append(
			sourceCmds,
			container.NewBorder(
				nil, nil, widget.NewLabel(config.Provider.Priority[i]), nil,
				widget.NewEntryWithData(binding.BindString(&d.SourceCMD[i]))))
	}
	dgSourceCMD := container.NewBorder(
		nil, nil, widget.NewLabel(i18n.T("plugin.diange.source_cmd")), nil,
		container.NewVBox(sourceCmds...))
	d.panel = container.NewVBox(dgPerm, dgQueue, dgCoolDown, dgShortCut, dgSourceCMD)
	return d.panel
}
