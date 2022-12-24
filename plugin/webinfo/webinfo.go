package webinfo

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/model"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const MODULE_PLGUIN_WEBINFO = "plugin.webinfo"

var lg = logger.Logger.WithField("Module", MODULE_PLGUIN_WEBINFO)

type WebInfo struct {
	config.BaseConfig
	Enabled bool
	Port    int
	server  *WebInfoServer
	panel   fyne.CanvasObject
	ctr     controller.IController
}

func NewWebInfo(ctr controller.IController) *WebInfo {
	return &WebInfo{
		Enabled: true,
		Port:    4000,
		ctr:     ctr,
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
	w.server = NewWebInfoServer(w.Port)
	w.registerHandlers()
	gui.AddConfigLayout(w)
	lg.Info("webinfo loaded")
	if w.Enabled {
		lg.Info("starting web backend server")
		w.server.Start()
	}
	return nil
}

func (w *WebInfo) Disable() error {
	lg.Info("closing webinfo backend server")
	if err := w.server.Stop(); err != nil {
		lg.Warnf("stop webinfo server encouter an error: %s", err)
	}
	return nil
}

func (t *WebInfo) registerHandlers() {
	t.ctr.PlayControl().EventManager().RegisterA(model.EventPlay, "plugin.webinfo.current", func(event *event.Event) {
		t.server.Info.Current = MediaInfo{
			Index:    0,
			Title:    event.Data.(model.PlayEvent).Media.Title,
			Artist:   event.Data.(model.PlayEvent).Media.Artist,
			Album:    event.Data.(model.PlayEvent).Media.Album,
			Cover:    event.Data.(model.PlayEvent).Media.Cover,
			Username: event.Data.(model.PlayEvent).Media.ToUser().Name,
		}
		t.server.SendInfo(
			OutInfoC,
			OutInfo{Current: t.server.Info.Current},
		)
	})
	if t.ctr.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropTimePos, "plugin.webinfo.timepos", func(event *event.Event) {
			data := event.Data.(model.PlayerPropertyUpdateEvent).Value
			if data == nil {
				t.server.Info.CurrentTime = 0
				return
			}
			ct := int(data.(float64))
			if ct == t.server.Info.CurrentTime {
				return
			}
			t.server.Info.CurrentTime = ct
			t.server.SendInfo(
				OutInfoCT,
				OutInfo{CurrentTime: t.server.Info.CurrentTime},
			)
		}) != nil {
		lg.Error("register time-pos handler failed")
	}
	if t.ctr.PlayControl().GetPlayer().ObserveProperty(
		model.PlayerPropDuration, "plugin.webinfo.duration", func(event *event.Event) {
			data := event.Data.(model.PlayerPropertyUpdateEvent).Value
			if data == nil {
				t.server.Info.TotalTime = 0
				return
			}
			t.server.Info.TotalTime = int(data.(float64))
			t.server.SendInfo(
				OutInfoTT,
				OutInfo{TotalTime: t.server.Info.TotalTime},
			)
		}) != nil {
		lg.Error("fail to register handler for total time with property duration")
	}
	t.ctr.Playlists().GetCurrent().EventManager().RegisterA(
		model.EventPlaylistUpdate, "plugin.webinfo.playlist", func(event *event.Event) {
			pl := make([]MediaInfo, 0)
			e := event.Data.(model.PlaylistUpdateEvent)
			for index, m := range e.Playlist.Medias {
				pl = append(pl, MediaInfo{
					Index:    index,
					Title:    m.Title,
					Artist:   m.Artist,
					Album:    m.Album,
					Username: m.ToUser().Name,
				})
			}
			t.server.Info.Playlist = pl
			t.server.SendInfo(
				OutInfoPL,
				OutInfo{Playlist: t.server.Info.Playlist},
			)
		})
	t.ctr.PlayControl().GetLyric().EventManager().RegisterA(
		model.EventLyricUpdate, "plugin.webinfo.lyric", func(event *event.Event) {
			lrcLine := event.Data.(model.LyricUpdateEvent).Lyric
			t.server.Info.Lyric = lrcLine.Now.Lyric
			t.server.SendInfo(
				OutInfoL,
				OutInfo{Lyric: t.server.Info.Lyric},
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
		widget.NewCheckWithData("", binding.BindBool(&w.Enabled)))
	statusText.SetText(w.getServerStatusText())
	serverPort := container.NewBorder(nil, nil,
		widget.NewLabel(i18n.T("plugin.webinfo.port")), nil,
		widget.NewEntryWithData(binding.IntToString(binding.BindInt(&w.Port))),
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
			lg.Info("User try stop webinfo server")
			err := w.server.Stop()
			if err != nil {
				lg.Warnf("stop server have error: %s", err)
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
			lg.Infof("User try start webinfo server with port %d", w.Port)
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
			lg.Infof("User try restart webinfo server with port %d", w.Port)
			if w.server.Running {
				if err := w.server.Stop(); err != nil {
					lg.Warnf("stop server have error: %s", err)
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
