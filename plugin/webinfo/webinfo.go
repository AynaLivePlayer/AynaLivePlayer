package webinfo

import (
	"AynaLivePlayer/adapters/logger"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/xfyne"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/util"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const MODULE_PLGUIN_WEBINFO = "plugin.webinfo"

var lg adapter.ILogger = &logger.EmptyLogger{}

type WebInfo struct {
	config.BaseConfig
	Enabled bool
	Port    int
	server  *WebInfoServer
	panel   fyne.CanvasObject
	ctr     adapter.IControlBridge
	log     adapter.ILogger
}

func NewWebInfo(ctr adapter.IControlBridge) *WebInfo {
	lg = ctr.Logger().WithModule(MODULE_PLGUIN_WEBINFO)
	return &WebInfo{
		Enabled: true,
		Port:    4000,
		ctr:     ctr,
		log:     ctr.Logger().WithModule(MODULE_PLGUIN_WEBINFO),
	}
}

func (w *WebInfo) Name() string {
	return "WebInfo"
}

func (w *WebInfo) Title() string {
	return i18n.T("plugin.webinfo.title")
}

func (w *WebInfo) Description() string {
	return i18n.T("plugin.webinfo.description")
}

func (w *WebInfo) Enable() error {
	config.LoadConfig(w)
	w.server = NewWebInfoServer(w.Port, w.log)
	w.registerHandlers()
	gui.AddConfigLayout(w)
	w.log.Info("webinfo loaded")
	if w.Enabled {
		w.log.Info("starting web backend server")
		w.server.Start()
	}
	return nil
}

func (w *WebInfo) Disable() error {
	w.log.Info("closing webinfo backend server")
	if err := w.server.Stop(); err != nil {
		w.log.Warnf("stop webinfo server encouter an error: %s", err)
	}
	return nil
}

func (w *WebInfo) registerHandlers() {
	w.ctr.PlayControl().EventManager().RegisterA(events.EventPlay, "plugin.webinfo.current", func(event *event.Event) {
		w.server.Info.Current = MediaInfo{
			Index:    0,
			Title:    event.Data.(events.PlayEvent).Media.Title,
			Artist:   event.Data.(events.PlayEvent).Media.Artist,
			Album:    event.Data.(events.PlayEvent).Media.Album,
			Cover:    event.Data.(events.PlayEvent).Media.Cover,
			Username: event.Data.(events.PlayEvent).Media.ToUser().Name,
		}
		w.server.SendInfo(
			OutInfoC,
			OutInfo{Current: w.server.Info.Current},
		)
	})
	if w.ctr.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropTimePos, "plugin.webinfo.timepos", func(event *event.Event) {
			data := event.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				w.server.Info.CurrentTime = 0
				return
			}
			ct := int(data.(float64))
			if ct == w.server.Info.CurrentTime {
				return
			}
			w.server.Info.CurrentTime = ct
			w.server.SendInfo(
				OutInfoCT,
				OutInfo{CurrentTime: w.server.Info.CurrentTime},
			)
		}) != nil {
		w.log.Error("register time-pos handler failed")
	}
	if w.ctr.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropDuration, "plugin.webinfo.duration", func(event *event.Event) {
			data := event.Data.(events.PlayerPropertyUpdateEvent).Value
			if data == nil {
				w.server.Info.TotalTime = 0
				return
			}
			w.server.Info.TotalTime = int(data.(float64))
			w.server.SendInfo(
				OutInfoTT,
				OutInfo{TotalTime: w.server.Info.TotalTime},
			)
		}) != nil {
		w.log.Error("fail to register handler for total time with property duration")
	}
	w.ctr.Playlists().GetCurrent().EventManager().RegisterA(
		events.EventPlaylistUpdate, "plugin.webinfo.playlist", func(event *event.Event) {
			pl := make([]MediaInfo, 0)
			e := event.Data.(events.PlaylistUpdateEvent)
			for index, m := range e.Playlist.Medias {
				pl = append(pl, MediaInfo{
					Index:    index,
					Title:    m.Title,
					Artist:   m.Artist,
					Album:    m.Album,
					Username: m.ToUser().Name,
				})
			}
			w.server.Info.Playlist = pl
			w.server.SendInfo(
				OutInfoPL,
				OutInfo{Playlist: w.server.Info.Playlist},
			)
		})
	w.ctr.PlayControl().GetLyric().EventManager().RegisterA(
		events.EventLyricUpdate, "plugin.webinfo.lyric", func(event *event.Event) {
			lrcLine := event.Data.(events.LyricUpdateEvent).Lyric
			w.server.Info.Lyric = lrcLine.Now.Lyric
			w.server.SendInfo(
				OutInfoL,
				OutInfo{Lyric: w.server.Info.Lyric},
			)
		})
}

func (w *WebInfo) getServerStatusText() string {
	if w.server.Running {
		return i18n.T("plugin.webinfo.server_status.running")
	} else {
		return i18n.T("plugin.webinfo.server_status.stopped")
	}
}

func (w *WebInfo) getServerUrl() string {
	return fmt.Sprintf("http://localhost:%d/#/previewV2", w.Port)
}

func (w *WebInfo) CreatePanel() fyne.CanvasObject {
	if w.panel != nil {
		return w.panel
	}

	statusText := widget.NewLabel("")
	serverStatus := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.webinfo.server_status")),
		statusText,
	)
	autoStart := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.webinfo.autostart")),
		component.NewCheckOneWayBinding("", &w.Enabled, w.Enabled))
	statusText.SetText(w.getServerStatusText())
	serverPort := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.webinfo.port")), nil,
		xfyne.EntryDisableUndoRedo(widget.NewEntryWithData(binding.IntToString(binding.BindInt(&w.Port)))),
	)
	serverUrl := widget.NewHyperlink(w.getServerUrl(), util.UrlMustParse(w.getServerUrl()))
	serverPreview := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.webinfo.server_preview")),
		serverUrl,
	)
	stopBtn := component.NewAsyncButtonWithIcon(
		i18n.T("plugin.webinfo.server_control.stop"),
		theme.MediaStopIcon(),
		func() {
			if !w.server.Running {
				return
			}
			w.log.Info("User try stop webinfo server")
			err := w.server.Stop()
			if err != nil {
				w.log.Warnf("stop server have error: %s", err)
				return
			}
			statusText.SetText(w.getServerStatusText())
		},
	)
	startBtn := component.NewAsyncButtonWithIcon(
		i18n.T("plugin.webinfo.server_control.start"),
		theme.MediaPlayIcon(),
		func() {
			if w.server.Running {
				return
			}
			w.log.Infof("User try start webinfo server with port %d", w.Port)
			w.server.Port = w.Port
			w.server.Start()
			statusText.SetText(w.getServerStatusText())
			serverUrl.SetText(w.getServerUrl())
			_ = serverUrl.SetURLFromString(w.getServerUrl())
		},
	)
	restartBtn := component.NewAsyncButtonWithIcon(
		i18n.T("plugin.webinfo.server_control.restart"),
		theme.MediaReplayIcon(),
		func() {
			w.log.Infof("User try restart webinfo server with port %d", w.Port)
			if w.server.Running {
				if err := w.server.Stop(); err != nil {
					w.log.Warnf("stop server have error: %s", err)
					return
				}
			}
			w.server.Port = w.Port
			w.server.Start()
			statusText.SetText(w.getServerStatusText())
			serverUrl.SetText(w.getServerUrl())
			_ = serverUrl.SetURLFromString(w.getServerUrl())
		},
	)
	ctrlBtns := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.webinfo.server_control")),
		startBtn, stopBtn, restartBtn,
	)
	w.panel = container.NewVBox(serverStatus, autoStart, serverPreview, serverPort, ctrlBtns)
	return w.panel
}
