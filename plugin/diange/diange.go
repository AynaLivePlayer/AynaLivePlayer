package diange

import (
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"strings"
	"time"
)

const MODULE_CMD_DIANGE = "CMD.DianGe"

type Diange struct {
	config.BaseConfig
	UserPermission      bool
	PrivilegePermission bool
	AdminPermission     bool
	MedalName           string
	MedalPermission     int
	QueueMax            int
	UserCoolDown        int
	CustomCMD           string
	SourceCMD           []string
	cooldowns           map[string]int
	panel               fyne.CanvasObject
	contro              adapter.IControlBridge
	log                 adapter.ILogger
}

func NewDiange(contr adapter.IControlBridge) *Diange {
	return &Diange{
		UserPermission:      true,
		PrivilegePermission: true,
		AdminPermission:     true,
		QueueMax:            128,
		UserCoolDown:        -1,
		CustomCMD:           "add",
		SourceCMD:           make([]string, 0),
		cooldowns:           make(map[string]int),
		contro:              contr,
		log:                 contr.Logger().WithModule(MODULE_CMD_DIANGE),
	}
}

func (d *Diange) Name() string {
	return "Diange"
}

func (d *Diange) Enable() error {
	config.LoadConfig(d)
	d.initCMD()
	d.contro.LiveRooms().AddDanmuCommand(d)
	gui.AddConfigLayout(d)
	return nil
}

func (d *Diange) Disable() error {
	return nil
}

func (d *Diange) initCMD() {
	if len(d.SourceCMD) == len(d.contro.Provider().GetPriority()) {
		return
	}
	if len(d.SourceCMD) > len(d.contro.Provider().GetPriority()) {
		d.SourceCMD = d.SourceCMD[:len(d.contro.Provider().GetPriority())]
		return
	}
	for i := len(d.SourceCMD); i < len(d.contro.Provider().GetPriority()); i++ {
		d.SourceCMD = append(d.SourceCMD, "点歌"+d.contro.Provider().GetPriority()[i])
	}
}

// isCMD return int if the commmand name matches our command
// -1 = not match, 0 = normal command, 1+ = source command
func (d *Diange) isCMD(cmd string) int {
	if cmd == "点歌" || cmd == d.CustomCMD {
		return 0
	}
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

func (d *Diange) Execute(command string, args []string, danmu *model.DanmuMessage) {
	d.log.Infof("%s(%s) Execute command: %s %s", danmu.User.Username, danmu.User.Uid, command, args)
	// if queue is full, return
	if d.contro.Playlists().GetCurrent().Size() >= d.QueueMax {
		d.log.Info("Queue is full, ignore diange")
		return
	}
	// if in user cool down, return
	ct := int(time.Now().Unix())
	if (ct - d.cooldowns[danmu.User.Uid]) <= d.UserCoolDown {
		d.log.Infof("User %s(%s) still in cool down period, diange failed", danmu.User.Username, danmu.User.Uid)
		return
	}
	cmdType := d.isCMD(command)
	keyword := strings.Join(args, " ")
	perm := d.UserPermission
	d.log.Debugf("user permission check: ", perm)
	perm = perm || (d.PrivilegePermission && (danmu.User.Privilege > 0))
	d.log.Debugf("privilege permission check: ", perm)
	perm = perm || (d.AdminPermission && (danmu.User.Admin))
	d.log.Debugf("admin permission check: ", perm)
	// if use medal check
	if d.MedalName != "" && d.MedalPermission >= 0 {
		perm = perm || ((danmu.User.Medal.Name == d.MedalName) && danmu.User.Medal.Level >= d.MedalPermission)
	}
	if !perm {
		return
	}
	// reset cool down
	d.cooldowns[danmu.User.Uid] = ct
	if cmdType == 0 {
		d.contro.PlayControl().Add(keyword, &danmu.User)
	} else {
		d.contro.PlayControl().AddWithProvider(keyword, d.contro.Provider().GetPriority()[cmdType-1], &danmu.User)
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
		component.NewCheckOneWayBinding(
			i18n.T("plugin.diange.user"), &d.UserPermission, d.UserPermission),
		component.NewCheckOneWayBinding(
			i18n.T("plugin.diange.privilege"), &d.PrivilegePermission, d.PrivilegePermission),
		component.NewCheckOneWayBinding(
			i18n.T("plugin.diange.admin"), &d.AdminPermission, d.AdminPermission),
	)
	dgMdPerm := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.diange.medal.perm")), nil,
		container.NewGridWithColumns(2,
			container.NewBorder(nil, nil,
				widget.NewLabel(i18n.T("plugin.diange.medal.name")), nil,
				widget.NewEntryWithData(binding.BindString(&d.MedalName))),
			container.NewBorder(nil, nil,
				widget.NewLabel(i18n.T("plugin.diange.medal.level")), nil,
				widget.NewEntryWithData(binding.IntToString(binding.BindInt(&d.MedalPermission)))),
		),
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
				nil, nil, widget.NewLabel(d.contro.Provider().GetPriority()[i]), nil,
				widget.NewEntryWithData(binding.BindString(&d.SourceCMD[i]))))
	}
	dgSourceCMD := container.NewBorder(
		nil, nil, widget.NewLabel(i18n.T("plugin.diange.source_cmd")), nil,
		container.NewVBox(sourceCmds...))
	d.panel = container.NewVBox(dgPerm, dgMdPerm, dgQueue, dgCoolDown, dgShortCut, dgSourceCMD)
	return d.panel
}
