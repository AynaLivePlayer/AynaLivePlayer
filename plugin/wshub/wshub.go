package wshub

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type WsHub struct {
	config.BaseConfig
	Enabled bool
	Port    int
	panel   fyne.CanvasObject
	server  *wsServer
	log     logger.ILogger
}

func NewWsHub() *WsHub {
	return &WsHub{
		Enabled: false,
		Port:    29629,
		log:     global.Logger.WithPrefix("plugin.wshub"),
	}
}

func (w *WsHub) Enable() error {
	config.LoadConfig(w)
	w.server = newWsServer(&w.Port)
	gui.AddConfigLayout(w)
	w.registerEvents()
	w.log.Info("webinfo loaded")
	if w.Enabled {
		w.log.Info("starting web backend server")
		w.server.Start()
	}
	return nil
}

func (w *WsHub) Disable() error {
	if w.server.Running {
		err := w.server.Stop()
		if err != nil {
			w.log.Warnf("stop server have error: %s", err)
		}
	}
	return nil
}

func (w *WsHub) Name() string {
	return "WsHub"
}

func (w *WsHub) Title() string {
	return i18n.T("plugin.wshub.title")
}

func (w *WsHub) Description() string {
	return i18n.T("plugin.wshub.description")
}

func (w *WsHub) CreatePanel() fyne.CanvasObject {
	if w.panel != nil {
		return w.panel
	}
	statusText := widget.NewLabel("")
	freshStatusText := func() {
		if w.server.Running {
			statusText.SetText(i18n.T("plugin.wshub.server_status.running"))
			return
		} else {
			statusText.SetText(i18n.T("plugin.wshub.server_status.stopped"))
		}
	}
	serverStatus := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.wshub.server_status")),
		statusText,
	)
	autoStart := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.wshub.autostart")),
		component.NewCheckOneWayBinding("", &w.Enabled, w.Enabled))
	freshStatusText()
	serverPort := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.wshub.port")), nil,
		widget.NewEntryWithData(binding.IntToString(binding.BindInt(&w.Port))),
	)
	serverUrl := widget.NewEntry()
	serverUrl.SetText(w.server.getWsUrl())
	serverUrl.Disable()
	serverPreview := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.wshub.server_link")), nil,
		serverUrl,
	)
	refreshServerUrl := func() {
		serverUrl.SetText(w.server.getWsUrl())
	}
	stopBtn := component.NewAsyncButtonWithIcon(
		i18n.T("plugin.wshub.server_control.stop"),
		theme.MediaStopIcon(),
		func() {
			if !w.server.Running {
				return
			}
			w.log.Info("User try stop webinfo server")
			err := w.server.Stop()
			if err != nil {
				w.log.Warnf("stop server have error: %s", err)
			}
			freshStatusText()
		},
	)
	startBtn := component.NewAsyncButtonWithIcon(
		i18n.T("plugin.wshub.server_control.start"),
		theme.MediaPlayIcon(),
		func() {
			if w.server.Running {
				return
			}
			w.log.Infof("User try start webinfo server with port %d", w.Port)
			w.server.Start()
			freshStatusText()
			refreshServerUrl()
		},
	)
	restartBtn := component.NewAsyncButtonWithIcon(
		i18n.T("plugin.wshub.server_control.restart"),
		theme.MediaReplayIcon(),
		func() {
			w.log.Infof("User try restart webinfo server with port %d", w.Port)
			if w.server.Running {
				if err := w.server.Stop(); err != nil {
					w.log.Warnf("stop server have error: %s", err)
					return
				}
			}
			w.server.Start()
			freshStatusText()
			refreshServerUrl()
		},
	)
	ctrlBtns := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.wshub.server_control")),
		startBtn, stopBtn, restartBtn,
	)
	w.panel = container.NewVBox(serverStatus, autoStart, serverPreview, serverPort, ctrlBtns)
	return nil
}

func (w *WsHub) registerEvents() {
	for eid, _ := range events.EventsMapping {
		global.EventManager.RegisterA(eid,
			"plugin.wshub.event."+string(eid),
			func(e *event.Event) {
				val, err := json.Marshal(EventData{
					EventID: e.Id,
					Data:    e.Data,
				})
				if err != nil {
					w.log.Errorf("failed to marshal event data %v", err)
					return
				}
				w.server.broadcast(val)
			})
	}
}
