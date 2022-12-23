package player

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/model"
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
}

func NewMpvPlayer() IPlayer {
	player := &MpvPlayer{
		running:             true,
		libmpv:              mpv.Create(),
		propertyWatchedFlag: make(map[model.PlayerProperty]int),
		eventManager:        event.MainManager.NewChildManager(),
	}
	err := player.libmpv.Initialize()
	if err != nil {
		lg.Error("[MPV PlayControl] initialize libmpv failed")
		return nil
	}
	_ = player.libmpv.SetOptionString("vo", "null")
	lg.Info("[MPV PlayControl] initialize libmpv success")
	player.Start()
	return player
}

func (p *MpvPlayer) Start() {
	lg.Info("[MPV PlayControl] starting mpv player")
	go func() {
		for p.running {
			e := p.libmpv.WaitEvent(1)
			if e == nil {
				lg.Warn("[MPV PlayControl] event loop got nil event")
			}
			lg.Trace("[MPV PlayControl] new event", e)
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
					model.EventPlayerPropertyUpdate(property),
					model.PlayerPropertyUpdateEvent{
						Property: property,
						Value:    value,
					})

			}
			if e.EventId == mpv.EVENT_SHUTDOWN {
				lg.Info("[MPV PlayControl] libmpv shutdown")
				p.Stop()
			}
		}
	}()
}

func (p *MpvPlayer) Stop() {
	lg.Info("[MPV PlayControl] stopping mpv player")
	p.running = false
	p.libmpv.TerminateDestroy()
}

func (p *MpvPlayer) Play(media *model.Media) error {
	lg.Infof("[MPV PlayControl] Play media %s", media.Url)
	if val, ok := media.Header["User-Agent"]; ok {
		lg.Debug("[MPV PlayControl] set user-agent for mpv player")
		err := p.libmpv.SetPropertyString("user-agent", val)
		if err != nil {
			lg.Warn("[MPV PlayControl] set player user-agent failed", err)
			return err
		}
	}

	if val, ok := media.Header["Referer"]; ok {
		lg.Debug("[MPV PlayControl] set referrer for mpv player")
		err := p.libmpv.SetPropertyString("referrer", val)
		if err != nil {
			lg.Warn("[MPV PlayControl] set player referrer failed", err)
			return err
		}
	}
	lg.Debugf("mpv command load file %s %s", media.Title, media.Url)
	if err := p.libmpv.Command([]string{"loadfile", media.Url}); err != nil {
		lg.Warn("[MPV PlayControl] mpv load media failed", media)
		return err
	}
	p.Playing = media
	return nil
}

func (p *MpvPlayer) IsPaused() bool {
	property, err := p.libmpv.GetProperty("pause", mpv.FORMAT_FLAG)
	if err != nil {
		lg.Warn("[MPV PlayControl] get property pause failed", err)
		return false
	}
	return property.(bool)
}

func (p *MpvPlayer) Pause() error {
	lg.Tracef("[MPV PlayControl] pause")
	return p.libmpv.SetProperty("pause", mpv.FORMAT_FLAG, true)
}

func (p *MpvPlayer) Unpause() error {
	lg.Tracef("[MPV PlayControl] unpause")
	return p.libmpv.SetProperty("pause", mpv.FORMAT_FLAG, false)
}

// SetVolume set mpv volume, from 0.0 - 100.0
func (p *MpvPlayer) SetVolume(volume float64) error {
	lg.Tracef("[MPV PlayControl] set volume to %f", volume)
	return p.libmpv.SetProperty("volume", mpv.FORMAT_DOUBLE, volume)
}

func (p *MpvPlayer) IsIdle() bool {
	property, err := p.libmpv.GetProperty("idle-active", mpv.FORMAT_FLAG)
	if err != nil {
		lg.Warn("[MPV PlayControl] get property idle-active failed", err)
		return false
	}
	return property.(bool)
}

// Seek change position for current file
// absolute = true : position is the time in second
// absolute = false: position is in percentage eg 0.1 0.2
func (p *MpvPlayer) Seek(position float64, absolute bool) error {
	lg.Tracef("[MPV PlayControl] seek to %f (absolute=%t)", position, absolute)
	if absolute {
		return p.libmpv.SetProperty("time-pos", mpv.FORMAT_DOUBLE, position)
	} else {
		return p.libmpv.SetProperty("percent-pos", mpv.FORMAT_DOUBLE, position)
	}
}

func (p *MpvPlayer) ObserveProperty(property model.PlayerProperty, name string, handler event.HandlerFunc) error {
	lg.Trace("[MPV PlayControl] add property observer for mpv")
	p.eventManager.RegisterA(
		model.EventPlayerPropertyUpdate(property),
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
	lg.Trace("[MPV PlayControl] getting audio device list for mpv")
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
	lg.Tracef("[MPV PlayControl] set audio device %s for mpv", device)
	return p.libmpv.SetPropertyString("audio-device", device)
}
