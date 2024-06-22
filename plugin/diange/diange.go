package diange

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/AynaLivePlayer/miaosic"
	"golang.org/x/exp/slices"
	"sort"
	"strings"
	"time"
)

type sourceConfig struct {
	Enable   bool   `json:"enable"`
	Command  string `json:"command"`
	Priority int    `json:"priority"`
}

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
	SourceConfigPath    string
	BlackListItemPath   string
	SkipSystemPlaylist  bool

	currentQueueLength int
	isCurrentSystem    bool
	sourceConfigs      map[string]*sourceConfig
	blacklist          []blacklistItem
	cooldowns          map[string]int
	panel              fyne.CanvasObject
	log                logger.ILogger
}

var diange *Diange

func NewDiange() *Diange {
	diange = &Diange{
		UserPermission:      true,
		PrivilegePermission: true,
		AdminPermission:     true,
		QueueMax:            128,
		UserCoolDown:        -1,
		CustomCMD:           "点歌",
		SourceConfigPath:    "./config/diange.json",
		BlackListItemPath:   "./config/diange_blacklist.json",

		currentQueueLength: 0,
		sourceConfigs: map[string]*sourceConfig{
			"netease": {
				Enable:   true,
				Command:  "点w歌",
				Priority: 1,
			},
			"kuwo": {
				Enable:   true,
				Command:  "点k歌",
				Priority: 2,
			},
			"bilibili-video": {
				Enable:   true,
				Command:  "点b歌",
				Priority: 3,
			},
			"local": {
				Enable:   true,
				Command:  "点l歌",
				Priority: 4,
			},
			"kugou": {
				Enable:   true,
				Command:  "点kg歌",
				Priority: 5,
			},
		},
		cooldowns: make(map[string]int),
		log:       global.Logger.WithPrefix("Plugin.Diange"),
	}
	return diange
}

func (d *Diange) Name() string {
	return "Diange"
}

func (c *Diange) OnLoad() {
	_ = config.LoadJson(c.SourceConfigPath, &c.sourceConfigs)
	_ = config.LoadJson(c.BlackListItemPath, &c.blacklist)
}

func (c *Diange) OnSave() {
	_ = config.SaveJson(c.SourceConfigPath, c.sourceConfigs)
	_ = config.SaveJson(c.BlackListItemPath, c.blacklist)
}

func (d *Diange) Enable() error {
	config.LoadConfig(d)
	gui.AddConfigLayout(d)
	gui.AddConfigLayout(&blacklist{})
	global.EventManager.RegisterA(
		events.LiveRoomMessageReceive,
		"plugin.diange.message",
		d.handleMessage)
	global.EventManager.RegisterA(
		events.PlaylistDetailUpdate(model.PlaylistIDPlayer),
		"plugin.diange.queue.update",
		func(event *event.Event) {
			d.currentQueueLength = len(event.Data.(events.PlaylistDetailUpdateEvent).Medias)
		})
	global.EventManager.RegisterA(
		events.PlayerPlayingUpdate,
		"plugin.diange.check_playing",
		func(event *event.Event) {
			data := event.Data.(events.PlayerPlayingUpdateEvent)
			if data.Removed {
				d.isCurrentSystem = true
				return
			}
			d.isCurrentSystem = (!data.Media.IsLiveRoomUser()) && (data.Media.ToUser().Name == model.SystemUser.Name)
		})
	return nil
}

func (d *Diange) Disable() error {
	return nil
}

func (d *Diange) getSources() []string {
	sources := []string{}
	for source, c := range d.sourceConfigs {
		if c.Enable {
			sources = append(sources, source)
		}
	}
	sort.Slice(sources, func(i, j int) bool {
		return d.sourceConfigs[sources[i]].Priority < d.sourceConfigs[sources[j]].Priority
	})
	return sources
}

func (d *Diange) getSource(cmd string) []string {
	customCmds := strings.Split(d.CustomCMD, "|")
	if slices.Contains(customCmds, cmd) {
		return d.getSources()
	}
	sources := []string{}
	for source, c := range d.sourceConfigs {
		if c.Command == cmd {
			sources = append(sources, source)
		}
	}
	return sources
}

func (d *Diange) handleMessage(event *event.Event) {
	message := event.Data.(events.LiveRoomMessageReceiveEvent).Message
	msgs := strings.Split(message.Message, " ")
	if len(msgs) < 2 || len(msgs[0]) == 0 || len(msgs[1]) == 0 {
		return
	}
	sources := d.getSource(msgs[0])
	if len(sources) == 0 {
		return
	}
	// if queue is full, return
	if d.currentQueueLength >= d.QueueMax {
		d.log.Info("Queue is full, ignore diange")
		return
	}

	// if in user cool down, return
	ct := int(time.Now().Unix())
	if (ct - d.cooldowns[message.User.Uid]) <= d.UserCoolDown {
		d.log.Infof("User %s(%s) still in cool down period, diange failed", message.User.Username, message.User.Uid)
		return
	}
	perm := d.UserPermission
	d.log.Debug("user permission check: ", perm)
	perm = perm || (d.PrivilegePermission && (message.User.Privilege > 0))
	d.log.Debug("privilege permission check: ", perm)
	perm = perm || (d.AdminPermission && (message.User.Admin))
	d.log.Debug("admin permission check: ", perm)
	// if use medal check
	if d.MedalName != "" && d.MedalPermission >= 0 {
		perm = perm || ((message.User.Medal.Name == d.MedalName) && message.User.Medal.Level >= d.MedalPermission)
	}
	if !perm {
		return
	}
	keywords := strings.Join(msgs[1:], " ")
	// blacklist check
	for _, item := range d.blacklist {
		if item.Exact && item.Value == keywords {
			d.log.Warnf("User %s(%s) diange %s is in blacklist %s, ignore", message.User.Username, message.User.Uid, keywords, item.Value)
			return
		}
		if !item.Exact && strings.Contains(keywords, item.Value) {
			d.log.Warnf("User %s(%s) diange %s is in blacklist %s, ignore", message.User.Username, message.User.Uid, keywords, item.Value)
			return
		}
	}

	d.cooldowns[message.User.Uid] = ct

	// match media first

	var mediaMeta miaosic.MetaData
	found := false
	for _, source := range sources {
		mediaMeta, found = miaosic.MatchMediaByProvider(source, keywords)
		if found {
			break
		}
	}

	var media miaosic.MediaInfo

	if !found {
		for _, source := range sources {
			medias, err := miaosic.SearchByProvider(source, keywords, 1, 10)
			if len(medias) == 0 || err != nil {
				continue
			}
			media = medias[0]
			found = true
			break
		}
	} else {
		d.log.Info("Match media: ", mediaMeta)
		m, err := miaosic.GetMediaInfo(mediaMeta)
		if err != nil {
			d.log.Error("Get media info failed: ", err)
			found = false
		}
		media = m
	}

	if found {
		// double check blacklist
		for _, item := range d.blacklist {
			if item.Exact && item.Value == media.Title {
				d.log.Warnf("User %s(%s) diange %s is in blacklist %s, ignore", message.User.Username, message.User.Uid, keywords, item.Value)
				return
			}
			if !item.Exact && strings.Contains(media.Title, item.Value) {
				d.log.Warnf("User %s(%s) diange %s is in blacklist %s, ignore", message.User.Username, message.User.Uid, keywords, item.Value)
				return
			}
		}
		if d.SkipSystemPlaylist && d.isCurrentSystem {
			global.EventManager.CallA(
				events.PlayerPlayCmd,
				events.PlayerPlayCmdEvent{
					Media: model.Media{
						Info: media,
						User: message.User,
					},
				})
			return
		}
		global.EventManager.CallA(
			events.PlaylistInsertCmd(model.PlaylistIDPlayer),
			events.PlaylistInsertCmdEvent{
				Position: -1,
				Media: model.Media{
					Info: media,
					User: message.User,
				},
			})
		return
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
	skipPlaylistCheck := widget.NewCheckWithData(i18n.T("plugin.diange.skip_playlist.prompt"), binding.BindBool(&d.SkipSystemPlaylist))
	skipPlaylist := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.diange.skip_playlist")),
		skipPlaylistCheck,
	)
	sourceCfgs := []fyne.CanvasObject{}
	for source, cfg := range d.sourceConfigs {
		sourceCfgs = append(
			sourceCfgs, container.NewGridWithColumns(2,
				widget.NewLabel(source),
				widget.NewCheckWithData(i18n.T("plugin.diange.source.enable"), binding.BindBool(&cfg.Enable)),
				widget.NewLabel(i18n.T("plugin.diange.source.priority")),
				widget.NewEntryWithData(binding.IntToString(binding.BindInt(&cfg.Priority))),
				widget.NewLabel(i18n.T("plugin.diange.source.command")),
				widget.NewEntryWithData(binding.BindString(&cfg.Command)),
			),
		)
	}
	dgSourceCMD := container.NewBorder(
		nil, nil, widget.NewLabel(i18n.T("plugin.diange.source_cmd")), nil,
		container.NewVBox(sourceCfgs...))
	d.panel = container.NewVBox(dgPerm, dgMdPerm, dgQueue, dgCoolDown, dgShortCut, skipPlaylist, dgSourceCMD)
	return d.panel
}
