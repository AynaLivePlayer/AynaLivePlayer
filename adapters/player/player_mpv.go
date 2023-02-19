package player

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"fmt"
	"github.com/aynakeya/go-mpv"
	"github.com/tidwall/gjson"
)

var mpvPropertyMap = map[model.PlayerProperty]string{
	model.PlayerPropDuration:   "duration",
	model.PlayerPropTimePos:    "time-pos",
	model.PlayerPropIdleActive: "idle-active",
	model.PlayerPropPercentPos: "percent-pos",
	model.PlayerPropPause:      "pause",
	model.PlayerPropVolume:     "volume",
}

var mpvPropertyMapInv = map[string]model.PlayerProperty{
	"duration":    model.PlayerPropDuration,
	"time-pos":    model.PlayerPropTimePos,
	"idle-active": model.PlayerPropIdleActive,
	"percent-pos": model.PlayerPropPercentPos,
	"pause":       model.PlayerPropPause,
	"volume":      model.PlayerPropVolume,
}

type MpvPlayer struct {
	running             bool
	libmpv              *mpv.Mpv
	Playing             *model.Media
	propertyWatchedFlag map[model.PlayerProperty]int
	eventManager        *event.Manager
	log                 adapter.ILogger
}

func NewMpvPlayer(em *event.Manager, log adapter.ILogger) adapter.IPlayer {
	player := &MpvPlayer{
		running:             true,
		libmpv:              mpv.Create(),
		propertyWatchedFlag: make(map[model.PlayerProperty]int),
		eventManager:        em,
		log:                 log,
	}
	err := player.libmpv.Initialize()
	if err != nil {
		player.log.Error("[MPV Player] initialize libmpv failed")
		return nil
	}
	_ = player.libmpv.SetOptionString("vo", "null")
	player.log.Info("[MPV Player] initialize libmpv success")
	_ = player.ObserveProperty(model.PlayerPropIdleActive, "player.setplaying", func(evnt *event.Event) {
		isIdle := evnt.Data.(events.PlayerPropertyUpdateEvent).Value.(bool)
		if isIdle {
			player.Playing = nil
		}
	})
	player.Start()
	return player
}

func (p *MpvPlayer) Start() {
	p.log.Info("[MPV Player] starting mpv player")
	go func() {
		for p.running {
			e := p.libmpv.WaitEvent(1)
			if e == nil {
				p.log.Warn("[MPV Player] event loop got nil event")
			}
			//p.log.Trace("[MPV Player] new event", e)
			if e.EventId == mpv.EVENT_PROPERTY_CHANGE {
				eventProperty := e.Property()
				property, ok := mpvPropertyMapInv[eventProperty.Name]
				if !ok {
					continue
				}
				var value interface{} = nil
				if eventProperty.Data != nil {
					value = eventProperty.Data.(mpv.Node).Value
				}
				p.eventManager.CallA(
					events.EventPlayerPropertyUpdate(property),
					events.PlayerPropertyUpdateEvent{
						Property: property,
						Value:    value,
					})

			}
			if e.EventId == mpv.EVENT_SHUTDOWN {
				p.log.Info("[MPV Player] libmpv shutdown")
				p.Stop()
			}
		}
	}()
	return
}

func (p *MpvPlayer) Stop() {
	p.log.Info("[MPV Player] stopping mpv player")
	p.running = false
	p.libmpv.TerminateDestroy()
}

func (p *MpvPlayer) GetPlaying() *model.Media {
	return p.Playing
}

func (p *MpvPlayer) SetWindowHandle(handle uintptr) error {
	p.log.Infof("[MPV Player] set window handle %d", handle)
	_ = p.libmpv.SetOptionString("wid", fmt.Sprintf("%d", handle))
	return p.libmpv.SetOptionString("vo", "gpu")
}

func (p *MpvPlayer) Play(media *model.Media) error {
	p.log.Infof("[MPV Player] Play media %s", media.Url)
	if val, ok := media.Header["User-Agent"]; ok {
		p.log.Debug("[MPV PlayControl] set user-agent for mpv player")
		err := p.libmpv.SetPropertyString("user-agent", val)
		if err != nil {
			p.log.Warn("[MPV PlayControl] set player user-agent failed", err)
			return err
		}
	}

	if val, ok := media.Header["Referer"]; ok {
		p.log.Debug("[MPV PlayControl] set referrer for mpv player")
		err := p.libmpv.SetPropertyString("referrer", val)
		if err != nil {
			p.log.Warn("[MPV PlayControl] set player referrer failed", err)
			return err
		}
	}
	p.log.Debugf("mpv command load file %s %s", media.Title, media.Url)
	if err := p.libmpv.Command([]string{"loadfile", media.Url}); err != nil {
		p.log.Warn("[MPV PlayControl] mpv load media failed", media)
		return err
	}
	p.Playing = media
	return nil
}

func (p *MpvPlayer) IsPaused() bool {
	property, err := p.libmpv.GetProperty("pause", mpv.FORMAT_FLAG)
	if err != nil {
		p.log.Warn("[MPV PlayControl] get property pause failed", err)
		return false
	}
	return property.(bool)
}

func (p *MpvPlayer) Pause() error {
	p.log.Debugf("[MPV Player] pause")
	return p.libmpv.SetProperty("pause", mpv.FORMAT_FLAG, true)
}

func (p *MpvPlayer) Unpause() error {
	p.log.Debugf("[MPV Player] unpause")
	return p.libmpv.SetProperty("pause", mpv.FORMAT_FLAG, false)
}

// SetVolume set mpv volume, from 0.0 - 100.0
func (p *MpvPlayer) SetVolume(volume float64) error {
	p.log.Debugf("[MPV Player] set volume to %f", volume)
	return p.libmpv.SetProperty("volume", mpv.FORMAT_DOUBLE, volume)
}

func (p *MpvPlayer) IsIdle() bool {
	property, err := p.libmpv.GetProperty("idle-active", mpv.FORMAT_FLAG)
	if err != nil {
		p.log.Warn("[MPV Player] get property idle-active failed", err)
		return false
	}
	return property.(bool)
}

// Seek change position for current file
// absolute = true : position is the time in second
// absolute = false: position is in percentage eg 0.1 0.2
func (p *MpvPlayer) Seek(position float64, absolute bool) error {
	p.log.Debugf("[MPV Player] seek to %f (absolute=%t)", position, absolute)
	if absolute {
		return p.libmpv.SetProperty("time-pos", mpv.FORMAT_DOUBLE, position)
	} else {
		return p.libmpv.SetProperty("percent-pos", mpv.FORMAT_DOUBLE, position)
	}
}

func (p *MpvPlayer) ObserveProperty(property model.PlayerProperty, name string, handler event.HandlerFunc) error {
	p.log.Debugf("[MPV Player] add property observer for mpv")
	p.eventManager.RegisterA(
		events.EventPlayerPropertyUpdate(property),
		name, handler)
	if _, ok := p.propertyWatchedFlag[property]; !ok {
		p.propertyWatchedFlag[property] = 1
		return p.libmpv.ObserveProperty(util.Hash64(mpvPropertyMap[property]), mpvPropertyMap[property], mpv.FORMAT_NODE)
	}
	return nil
}

// GetAudioDeviceList get output device for mpv
// return format is []AudioDevice
func (p *MpvPlayer) GetAudioDeviceList() ([]model.AudioDevice, error) {
	p.log.Debugf("[MPV Player] getting audio device list for mpv")
	property, err := p.libmpv.GetProperty("audio-device-list", mpv.FORMAT_STRING)
	if err != nil {
		return nil, err
	}
	dl := make([]model.AudioDevice, 0)
	gjson.Parse(property.(string)).ForEach(func(key, value gjson.Result) bool {
		dl = append(dl, model.AudioDevice{
			Name:        value.Get("name").String(),
			Description: value.Get("description").String(),
		})
		return true
	})
	return dl, nil
}

func (p *MpvPlayer) SetAudioDevice(device string) error {
	p.log.Debugf("[MPV Player] set audio device %s for mpv", device)
	return p.libmpv.SetPropertyString("audio-device", device)
}
