package durationmgmt

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/xfyne"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MaxDuration struct {
	config.BaseConfig
	MaxDuration int
	SkipOnPlay  bool
	SkipOnReach bool
	skipped     bool
	panel       fyne.CanvasObject
	log         logger.ILogger
}

func NewMaxDuration() *MaxDuration {
	return &MaxDuration{
		MaxDuration: 60 * 10,
		SkipOnPlay:  false,
		SkipOnReach: false,
		skipped:     false,
		log:         global.Logger.WithPrefix("plugin.maxduration"),
	}
}

func (d *MaxDuration) Name() string {
	return "MaxDuration"
}

func (d *MaxDuration) Enable() error {
	config.LoadConfig(d)
	gui.AddConfigLayout(d)
	global.EventBus.Subscribe("",
		events.PlayerPropertyDurationUpdate,
		"plugin.maxduration.duration",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlayerPropertyDurationUpdateEvent)
			if int(data.Duration) > d.MaxDuration && d.SkipOnPlay {
				d.log.Infof("Skip on reach max duration %.2f/%d (on play)", data.Duration, d.MaxDuration)
				_ = global.EventBus.Publish(
					events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
			}
		})
	global.EventBus.Subscribe("",
		events.PlayerPropertyTimePosUpdate,
		"plugin.maxduration.timepos",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlayerPropertyTimePosUpdateEvent)
			if int(data.TimePos) > d.MaxDuration && d.SkipOnReach && !d.skipped {
				d.log.Infof("Skip on reach max duration %.2f/%d (on time pos reach)", data.TimePos, d.MaxDuration)
				d.skipped = true
				_ = global.EventBus.Publish(
					events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
			}
		})
	global.EventBus.Subscribe("",
		events.PlayerPlayingUpdate,
		"plugin.maxduration.play",
		func(event *eventbus.Event) {
			d.skipped = false
		})
	return nil
}

func (d *MaxDuration) Disable() error {
	return nil
}

func (d *MaxDuration) Title() string {
	return i18n.T("plugin.maxduration.title")
}

func (d *MaxDuration) Description() string {
	return i18n.T("plugin.maxduration.description")
}

func (d *MaxDuration) CreatePanel() fyne.CanvasObject {
	if d.panel != nil {
		return d.panel
	}
	maxDurationInput := xfyne.EntryDisableUndoRedo(widget.NewEntryWithData(binding.IntToString(binding.BindInt(&d.MaxDuration))))
	skipOnPlayCheckbox := widget.NewCheckWithData(i18n.T("plugin.maxduration.enable"), binding.BindBool(&d.SkipOnPlay))
	skipOnReachCheckbox := widget.NewCheckWithData(i18n.T("plugin.maxduration.enable"), binding.BindBool(&d.SkipOnReach))
	d.panel = container.New(
		layout.NewFormLayout(),
		widget.NewLabel(i18n.T("plugin.maxduration.maxduration")),
		maxDurationInput,
		widget.NewLabel(i18n.T("plugin.maxduration.skiponplay")),
		skipOnPlayCheckbox,
		widget.NewLabel(i18n.T("plugin.maxduration.skiponreach")),
		skipOnReachCheckbox,
	)
	return d.panel
}
