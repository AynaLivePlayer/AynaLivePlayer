package webinfo

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/player"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/aynakeya/go-mpv"
)

const MODULE_PLGUIN_WEBINFO = "plugin.webinfo"

var lg = logger.Logger.WithField("Module", MODULE_PLGUIN_WEBINFO)

type WebInfo struct {
	Port   int
	server *WebInfoServer
}

func NewWebInfo() *WebInfo {
	return &WebInfo{
		Port: 4000,
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
	lg.Info("starting web backend server")
	w.server.Start()
	w.registerHandlers()
	gui.AddConfigLayout(w)
	lg.Info("webinfo loaded")
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
	controller.MainPlayer.EventHandler.RegisterA(player.EventPlay, "plugin.webinfo.current", func(event *event.Event) {
		t.server.Info.Current = MediaInfo{
			Index:    0,
			Title:    event.Data.(player.PlayEvent).Media.Title,
			Artist:   event.Data.(player.PlayEvent).Media.Artist,
			Album:    event.Data.(player.PlayEvent).Media.Album,
			Cover:    event.Data.(player.PlayEvent).Media.Cover,
			Username: event.Data.(player.PlayEvent).Media.ToUser().Name,
		}
		t.server.SendInfo(
			OutInfoC,
			OutInfo{Current: t.server.Info.Current},
		)
	})
	if controller.MainPlayer.ObserveProperty("time-pos", func(property *mpv.EventProperty) {
		if property.Data == nil {
			t.server.Info.CurrentTime = 0
			return
		}
		ct := int(property.Data.(mpv.Node).Value.(float64))
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
	if controller.MainPlayer.ObserveProperty("duration", func(property *mpv.EventProperty) {
		if property.Data == nil {
			t.server.Info.TotalTime = 0
			return
		}
		t.server.Info.TotalTime = int(property.Data.(mpv.Node).Value.(float64))
		t.server.SendInfo(
			OutInfoTT,
			OutInfo{TotalTime: t.server.Info.TotalTime},
		)
	}) != nil {
		lg.Error("fail to register handler for total time with property duration")
	}
	controller.UserPlaylist.Handler.RegisterA(player.EventPlaylistUpdate, "plugin.webinfo.playlist", func(event *event.Event) {
		pl := make([]MediaInfo, 0)
		e := event.Data.(player.PlaylistUpdateEvent)
		e.Playlist.Lock.RLock()
		for index, m := range e.Playlist.Playlist {
			pl = append(pl, MediaInfo{
				Index:    index,
				Title:    m.Title,
				Artist:   m.Artist,
				Album:    m.Album,
				Username: m.ToUser().Name,
			})
		}
		e.Playlist.Lock.RUnlock()
		t.server.Info.Playlist = pl
		t.server.SendInfo(
			OutInfoPL,
			OutInfo{Playlist: t.server.Info.Playlist},
		)
	})
	controller.CurrentLyric.Handler.RegisterA(player.EventLyricUpdate, "plugin.webinfo.lyric", func(event *event.Event) {
		lrcLine := event.Data.(player.LyricUpdateEvent).Lyric
		t.server.Info.Lyric = lrcLine.Lyric
		t.server.SendInfo(
			OutInfoL,
			OutInfo{Lyric: t.server.Info.Lyric},
		)
	})

}

func (w *WebInfo) CreatePanel() fyne.CanvasObject {
	return widget.NewLabel("No Setting")
}
