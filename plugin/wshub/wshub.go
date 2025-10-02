package wshub

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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"net/url"
)

type WsHub struct {
	config.BaseConfig
	Enabled            bool
	Port               int
	LocalHostOnly      bool
	EnableWsHubControl bool
	panel              fyne.CanvasObject
	server             *wsServer
	log                logger.ILogger
}

func NewWsHub() *WsHub {
	return &WsHub{
		Enabled:            false,
		Port:               29629,
		LocalHostOnly:      true,
		EnableWsHubControl: false,
		log:                global.Logger.WithPrefix("plugin.wshub"),
	}
}

var globalEnableWsHubControl = false

func (w *WsHub) Enable() error {
	config.LoadConfig(w)
	// todo: should pass EnableWsHubControl to client instead of using global variable
	globalEnableWsHubControl = w.EnableWsHubControl
	w.server = newWsServer(&w.Port, &w.LocalHostOnly)
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
	localHostOnly := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.wshub.local_host_only")),
		component.NewCheckOneWayBinding("", &w.LocalHostOnly, w.LocalHostOnly))
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
	uri, _ := url.Parse("http://obsinfo.biliaudiobot.com/")
	infos := container.NewHBox(widget.NewLabel(i18n.T("plugin.wshub.webinfo_text")), widget.NewHyperlink("http://obsinfo.biliaudiobot.com", uri))
	w.panel = container.NewVBox(serverStatus, autoStart, localHostOnly, serverPreview, serverPort, ctrlBtns, infos)
	return nil
}

func (w *WsHub) registerEvents() {
	i := 0
	for eid, _ := range events.EventsMapping {
		eventCache = append(eventCache, &EventData{})
		currentIdx := i
		global.EventBus.Subscribe(eventChannel, eid,
			"plugin.wshub.event."+string(eid),
			func(e *eventbus.Event) {
				ed := EventData{
					EventID: e.Id,
					Data:    e.Data,
				}
				val, err := toCapitalizedJSON(ed)
				if err != nil {
					w.log.Errorf("failed to marshal event data %v", err)
					return
				}
				eventCache[currentIdx] = &ed
				w.server.broadcast(val)
			})
		i++
	}
}
