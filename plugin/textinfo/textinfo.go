package textinfo

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/event"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/player"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/aynakeya/go-mpv"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const MODULE_PLUGIN_TEXTINFO = "plugin.textinfo"

func l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_PLUGIN_TEXTINFO)
}

const Template_Path = "./template/"
const Out_Path = "./txtinfo/"

type Template struct {
	Name string
	Text string
	Tmpl *template.Template
}

type MediaInfo struct {
	Index    int
	Title    string
	Artist   string
	Album    string
	Username string
	Cover    player.Picture
}

type OutInfo struct {
	Current       MediaInfo
	CurrentTime   int
	TotalTime     int
	Lyric         string
	Playlist      []MediaInfo
	PlaylistCount int
}

type TextInfo struct {
	Rendering  bool
	info       OutInfo
	templates  []*Template
	emptyCover []byte
	panel      fyne.CanvasObject
}

func NewTextInfo() *TextInfo {
	b, _ := ioutil.ReadFile(config.GetAssetPath("empty.png"))
	return &TextInfo{Rendering: true, emptyCover: b}
}

func (t *TextInfo) Title() string {
	return i18n.T("plugin.textinfo.title")
}

func (t *TextInfo) Description() string {
	return i18n.T("plugin.textinfo.description")
}

func (t *TextInfo) CreatePanel() fyne.CanvasObject {
	if t.panel != nil {
		return t.panel
	}
	enableRendering := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.textinfo.prompt")),
		widget.NewCheckWithData(i18n.T("plugin.textinfo.checkbox"), binding.BindBool(&t.Rendering)),
	)
	t.panel = container.NewVBox(enableRendering)
	return t.panel
}

func (t *TextInfo) Name() string {
	return "TextInfo"
}

func (t *TextInfo) Enable() (err error) {
	// ensure the output/input directory exists
	if err = os.MkdirAll(Template_Path, 0755); err != nil {
		return
	}
	if err = os.MkdirAll(Out_Path, 0755); err != nil {
		return
	}
	config.LoadConfig(t)
	t.reloadTemplates()
	t.registerHandlers()
	gui.AddConfigLayout(t)
	return nil
}

func (t *TextInfo) reloadTemplates() {
	var err error
	t.templates = make([]*Template, 0)
	files, err := ioutil.ReadDir(Template_Path)
	if err != nil {
		l().Warn("read template directory failed: ", err)
		return
	}
	for _, f := range files {
		l().Info("loading template: ", f.Name())
		content, err := ioutil.ReadFile(filepath.Join(Template_Path, f.Name()))
		if err != nil {
			l().Warnf("read template file %s failed: %s", f.Name(), err)
			continue
		}
		parse, err := template.New("info").
			Funcs(template.FuncMap{
				"GetSeconds": func(t int) int {
					return t % 60
				},
				"GetMinutes": func(t int) int {
					return t / 60
				},
			}).
			Parse(string(content))
		if err != nil {
			l().Warnf("parse template %s failed: %s", f.Name, err)
			continue
		}
		t.templates = append(t.templates, &Template{
			Name: f.Name(),
			Text: string(content),
			Tmpl: parse,
		})
	}
}

// RenderTemplates render the template to the output file
func (t *TextInfo) RenderTemplates() {
	if !t.Rendering {
		return
	}
	for _, tmpl := range t.templates {
		l().Trace("rendering template: ", tmpl.Name)
		out, err := os.Create(filepath.Join(Out_Path, tmpl.Name))
		defer out.Close()
		if err != nil {
			l().Warnf("create output file %s failed: %s", tmpl.Name, err)
			continue
		}
		if err = tmpl.Tmpl.Execute(out, t.info); err != nil {
			l().Warnf("rendering template %s failed: %s", tmpl.Name, err)
			return
		}
	}
}

func (t *TextInfo) OutputCover() {
	if !t.Rendering {
		return
	}
	if !t.info.Current.Cover.Exists() {
		err := ioutil.WriteFile(filepath.Join(Out_Path, "cover.jpg"), t.emptyCover, 0666)
		if err != nil {
			l().Warnf("write cover file failed: %s", err)
		}
		return
	}
	if t.info.Current.Cover.Data != nil {
		err := ioutil.WriteFile(filepath.Join(Out_Path, "cover.jpg"), t.info.Current.Cover.Data, 0666)
		if err != nil {
			l().Warnf("write cover file failed: %s", err)
		}
		return
	}
	go func() {
		resp, err := resty.New().R().
			Get(t.info.Current.Cover.Url)
		if err != nil {
			l().Warnf("get cover %s content failed: %s", t.info.Current.Cover.Url, err)
			return
		}
		err = ioutil.WriteFile(filepath.Join(Out_Path, "cover.jpg"), resp.Body(), 0666)
		if err != nil {
			l().Warnf("write cover file failed: %s", err)
		}
	}()
}

func (t *TextInfo) registerHandlers() {
	controller.MainPlayer.EventHandler.RegisterA(player.EventPlay, "plugin.textinfo.current", func(event *event.Event) {
		t.info.Current = MediaInfo{
			Index:    0,
			Title:    event.Data.(player.PlayEvent).Media.Title,
			Artist:   event.Data.(player.PlayEvent).Media.Artist,
			Album:    event.Data.(player.PlayEvent).Media.Album,
			Cover:    event.Data.(player.PlayEvent).Media.Cover,
			Username: event.Data.(player.PlayEvent).Media.ToUser().Name,
		}
		t.RenderTemplates()
		t.OutputCover()
	})
	if controller.MainPlayer.ObserveProperty("time-pos", func(property *mpv.EventProperty) {
		if property.Data == nil {
			t.info.CurrentTime = 0
			return
		}
		ct := int(property.Data.(mpv.Node).Value.(float64))
		if ct == t.info.CurrentTime {
			return
		}
		t.info.CurrentTime = ct
		t.RenderTemplates()
	}) != nil {
		l().Error("register time-pos handler failed")
	}
	if controller.MainPlayer.ObserveProperty("duration", func(property *mpv.EventProperty) {
		if property.Data == nil {
			t.info.TotalTime = 0
			return
		}
		t.info.TotalTime = int(property.Data.(mpv.Node).Value.(float64))
	}) != nil {
		l().Error("fail to register handler for total time with property duration")
	}
	controller.UserPlaylist.Handler.RegisterA(player.EventPlaylistUpdate, "plugin.textinfo.playlist", func(event *event.Event) {
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
		t.info.Playlist = pl
		t.RenderTemplates()
	})
	controller.CurrentLyric.Handler.RegisterA(player.EventLyricUpdate, "plugin.textinfo.lyric", func(event *event.Event) {
		lrcLine := event.Data.(player.LyricUpdateEvent).Lyric
		t.info.Lyric = lrcLine.Lyric
		t.RenderTemplates()
	})

}
