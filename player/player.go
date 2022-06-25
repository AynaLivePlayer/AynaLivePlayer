package player

import (
	"AynaLivePlayer/event"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/util"
	"github.com/aynakeya/go-mpv"
	"github.com/sirupsen/logrus"
)

const MODULE_PLAYER = "Player.Player"

type PropertyHandlerFunc func(property *mpv.EventProperty)

type Player struct {
	running         bool
	libmpv          *mpv.Mpv
	Playing         *Media
	PropertyHandler map[string][]PropertyHandlerFunc
	EventHandler    *event.Handler
}

func NewPlayer() *Player {
	player := &Player{
		running:         true,
		libmpv:          mpv.Create(),
		PropertyHandler: make(map[string][]PropertyHandlerFunc),
		EventHandler:    event.NewHandler(),
	}
	err := player.libmpv.Initialize()
	if err != nil {
		player.l().Error("initialize libmpv failed")
		return nil
	}
	player.libmpv.SetOptionString("vo", "null")
	player.l().Info("initialize libmpv success")
	return player
}

func (p *Player) Start() {
	p.l().Info("starting mpv player")
	go func() {
		for p.running {
			e := p.libmpv.WaitEvent(1)
			if e == nil {
				p.l().Warn("event loop got nil event")
			}
			p.l().Trace("new event", e)
			if e.EventId == mpv.EVENT_PROPERTY_CHANGE {
				property := e.Property()
				p.l().Trace("receive property change event", property)
				for _, handler := range p.PropertyHandler[property.Name] {
					// todo: @3
					go handler(&property)
				}
			}
			if e.EventId == mpv.EVENT_SHUTDOWN {
				p.l().Info("libmpv shutdown")
				p.Stop()
			}
		}
	}()
}

func (p *Player) Stop() {
	p.l().Info("stopping mpv player")
	p.running = false
	p.libmpv.TerminateDestroy()
}

func (p *Player) l() *logrus.Entry {
	return logger.Logger.WithField("Module", MODULE_PLAYER)
}

func (p *Player) Play(media *Media) error {
	p.l().Infof("Play media %s", media.Url)
	p.l().Trace("set user-agent for mpv player")
	if val, ok := media.Header["user-agent"]; ok {
		err := p.libmpv.SetPropertyString("user-agent", val)
		if err != nil {
			p.l().Warn("set player user-agent failed", err)
			return err
		}
	}
	p.l().Trace("set referrer for mpv player")
	if val, ok := media.Header["referrer"]; ok {
		err := p.libmpv.SetPropertyString("referrer", val)
		if err != nil {
			p.l().Warn("set player referrer failed", err)
			return err
		}
	}
	p.l().Debugf("mpv command load file %s %s", media.Title, media.Url)
	if err := p.libmpv.Command([]string{"loadfile", media.Url}); err != nil {
		p.l().Warn("mpv load media failed", media)
		return err
	}
	p.Playing = media
	p.EventHandler.CallA(EventPlay, PlayEvent{Media: media})
	return nil
}

func (p *Player) IsPaused() bool {
	property, err := p.libmpv.GetProperty("pause", mpv.FORMAT_FLAG)
	if err != nil {
		p.l().Warn("get property pause failed", err)
		return false
	}
	return property.(bool)
}

func (p *Player) Pause() error {
	p.l().Tracef("pause")
	return p.libmpv.SetProperty("pause", mpv.FORMAT_FLAG, true)
}

func (p *Player) Unpause() error {
	p.l().Tracef("unpause")
	return p.libmpv.SetProperty("pause", mpv.FORMAT_FLAG, false)
}

// SetVolume set mpv volume, from 0.0 - 100.0
func (p *Player) SetVolume(volume float64) error {
	p.l().Tracef("set volume to %f", volume)
	return p.libmpv.SetProperty("volume", mpv.FORMAT_DOUBLE, volume)
}

func (p *Player) IsIdle() bool {
	property, err := p.libmpv.GetProperty("idle-active", mpv.FORMAT_FLAG)
	if err != nil {
		p.l().Warn("get property idle-active failed", err)
		return false
	}
	return property.(bool)
}

// Seek change position for current file
// absolute = true : position is the time in second
// absolute = false: position is in percentage eg 0.1 0.2
func (p *Player) Seek(position float64, absolute bool) error {
	p.l().Tracef("seek to %f (absolute=%t)", position, absolute)
	if absolute {
		return p.libmpv.SetProperty("time-pos", mpv.FORMAT_DOUBLE, position)
	} else {
		return p.libmpv.SetProperty("percent-pos", mpv.FORMAT_DOUBLE, position)
	}
}

func (p *Player) ObserveProperty(property string, handler ...PropertyHandlerFunc) error {
	p.l().Trace("add property observer for mpv")
	p.PropertyHandler[property] = append(p.PropertyHandler[property], handler...)
	if len(p.PropertyHandler[property]) == 1 {
		return p.libmpv.ObserveProperty(util.Hash64(property), property, mpv.FORMAT_NODE)
	}
	return nil
}
