package textinfo

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
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ajstarks/svgo"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const MODULE_PLUGIN_TEXTINFO = "plugin.textinfo"

const Template_Path = "./template/"
const Out_Path = "./txtinfo/"

type Template struct {
	Name string
	Text string
	Tmpl *template.Template
}

type TextInfo struct {
	config.BaseConfig
	Rendering  bool
	info       OutInfo
	templates  []*Template
	emptyCover []byte
	panel      fyne.CanvasObject
	log        logger.ILogger
}

func NewTextInfo() *TextInfo {
	buf := bytes.NewBuffer([]byte{})
	canvas := svg.New(buf)
	canvas.Start(256, 256)
	canvas.Image(0, 0, 256, 256, "cover.jpg")
	canvas.End()
	return &TextInfo{Rendering: true,
		emptyCover: buf.Bytes(),
		log:        global.Logger.WithPrefix(MODULE_PLUGIN_TEXTINFO),
	}
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
		component.NewCheckOneWayBinding(
			i18n.T("plugin.textinfo.checkbox"),
			&t.Rendering, t.Rendering),
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

func (d *TextInfo) Disable() error {
	return nil
}

func (t *TextInfo) reloadTemplates() {
	var err error
	t.templates = make([]*Template, 0)
	files, err := ioutil.ReadDir(Template_Path)
	if err != nil {
		t.log.Warn("read template directory failed: ", err)
		return
	}
	for _, f := range files {
		t.log.Info("loading template: ", f.Name())
		content, err := ioutil.ReadFile(filepath.Join(Template_Path, f.Name()))
		if err != nil {
			t.log.Warnf("read template file %s failed: %s", f.Name(), err)
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
			t.log.Warnf("parse template %s failed: %s", f.Name, err)
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
		t.log.Debug("rendering template: ", tmpl.Name)
		out, err := os.Create(filepath.Join(Out_Path, tmpl.Name))
		defer out.Close()
		if err != nil {
			t.log.Warnf("create output file %s failed: %s", tmpl.Name, err)
			continue
		}
		if err = tmpl.Tmpl.Execute(out, t.info); err != nil {
			t.log.Warnf("rendering template %s failed: %s", tmpl.Name, err)
			return
		}
	}
}

func (t *TextInfo) OutputCover() {
	if !t.Rendering {
		return
	}
	if !t.info.Current.Cover.Exists() {
		err := os.WriteFile(filepath.Join(Out_Path, "cover.jpg"), t.emptyCover, 0666)
		if err != nil {
			t.log.Warnf("write cover file failed: %s", err)
		}
		return
	}
	if t.info.Current.Cover.Data != nil {
		err := os.WriteFile(filepath.Join(Out_Path, "cover.jpg"), t.info.Current.Cover.Data, 0666)
		if err != nil {
			t.log.Warnf("write cover file failed: %s", err)
		}
		return
	}
	go func() {
		resp, err := resty.New().R().
			Get(t.info.Current.Cover.Url)
		if err != nil {
			t.log.Warnf("get cover %s content failed: %s", t.info.Current.Cover.Url, err)
			return
		}
		err = os.WriteFile(filepath.Join(Out_Path, "cover.jpg"), resp.Body(), 0666)
		if err != nil {
			t.log.Warnf("write cover file failed: %s", err)
		}
	}()
}

func (t *TextInfo) registerHandlers() {
	global.EventManager.RegisterA(
		events.PlayerPlayingUpdate, "plugin.textinfo.playing", func(event *event.Event) {
			data := event.Data.(events.PlayerPlayingUpdateEvent)
			if data.Removed {
				t.info.Current = MediaInfo{}
			} else {
				t.info.Current = NewMediaInfo(0, data.Media)
			}
			t.RenderTemplates()
			t.OutputCover()
		})
	global.EventManager.RegisterA(
		events.PlayerPropertyTimePosUpdate, "plugin.txtinfo.timepos", func(event *event.Event) {
			data := event.Data.(events.PlayerPropertyTimePosUpdateEvent).TimePos
			ct := int(data)
			if ct == t.info.CurrentTime.TotalSeconds {
				return
			}
			t.info.CurrentTime = NewTimeFromSec(ct)
			t.RenderTemplates()
		})
	global.EventManager.RegisterA(
		events.PlayerPropertyDurationUpdate, "plugin.txtinfo.duration", func(event *event.Event) {
			data := event.Data.(events.PlayerPropertyDurationUpdateEvent).Duration
			ct := int(data)
			if ct == t.info.TotalTime.TotalSeconds {
				return
			}
			t.info.TotalTime = NewTimeFromSec(ct)
			t.RenderTemplates()
		})
	global.EventManager.RegisterA(
		events.PlaylistDetailUpdate(model.PlaylistIDPlayer), "plugin.textinfo.playlist", func(event *event.Event) {
			pl := make([]MediaInfo, 0)
			data := event.Data.(events.PlaylistDetailUpdateEvent)
			for index, m := range data.Medias {
				pl = append(pl, NewMediaInfo(index, m))
			}
			t.info.Playlist = pl
			t.info.PlaylistLength = len(pl)
			t.RenderTemplates()
		})
	global.EventManager.RegisterA(
		events.PlayerLyricPosUpdate, "plugin.textinfo.lyricpos", func(event *event.Event) {
			data := event.Data.(events.PlayerLyricPosUpdateEvent)
			t.info.Lyric = data.CurrentLine.Lyric
			t.RenderTemplates()
		})

}
