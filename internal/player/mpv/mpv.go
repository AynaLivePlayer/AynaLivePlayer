package mpv

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/logger"
	"AynaLivePlayer/pkg/util"
	"fmt"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/aynakeya/go-mpv"
	"github.com/tidwall/gjson"
)

var running bool = false
var libmpv *mpv.Mpv = nil
var log logger.ILogger = nil

func SetupPlayer() {
	running = true
	config.LoadConfig(cfg)
	libmpv = mpv.Create()
	log = global.Logger.WithPrefix("MPV Player")
	err := libmpv.Initialize()
	if err != nil {
		log.Error("initialize libmpv failed")
		return
	}
	_ = libmpv.SetOptionString("vo", "null")
	log.Info("initialize libmpv success")
	registerHandler()
	registerCmdHandler()
	restoreConfig()
	updateAudioDeviceList()
	log.Info("starting mpv player")
	go func() {
		for running {
			e := libmpv.WaitEvent(1)
			if e == nil {
				log.Warn("[MPV Player] event loop got nil event")
			}
			if e.EventId == mpv.EVENT_PROPERTY_CHANGE {
				eventProperty := e.Property()
				handler, ok := mpvPropertyHandler[eventProperty.Name]
				if !ok {
					continue
				}
				var value interface{} = nil
				if eventProperty.Data != nil {
					value = eventProperty.Data.(mpv.Node).Value
				}
				//log.Debugf("[MPV Player] property update %s %v", eventProperty.Name, value)
				handler(value)
			}
			if e.EventId == mpv.EVENT_SHUTDOWN {
				log.Info("[MPV Player] libmpv shutdown")
				StopPlayer()
			}
		}
	}()
}

func StopPlayer() {
	cfg.AudioDevice = libmpv.GetPropertyString("audio-device")
	log.Info("stopping mpv player")
	running = false
	// stop player async, should be closed in the end anyway
	go libmpv.TerminateDestroy()
	log.Info("mpv player stopped")
}

var mpvPropertyHandler = map[string]func(value any){
	"pause": func(value any) {
		var data events.PlayerPropertyPauseUpdateEvent
		log.Debugf("pause property update %v %T", value, value)
		data.Paused = value.(bool)
		global.EventManager.CallA(events.PlayerPropertyPauseUpdate, data)
	},
	"percent-pos": func(value any) {
		var data events.PlayerPropertyPercentPosUpdateEvent
		if value == nil {
			data.PercentPos = 0
		} else {
			data.PercentPos = value.(float64)
		}
		// ignore bug value
		if data.PercentPos < 0.1 {
			return
		}
		global.EventManager.CallA(events.PlayerPropertyPercentPosUpdate, data)

	},
	"idle-active": func(value any) {
		var data events.PlayerPropertyIdleActiveUpdateEvent
		if value == nil {
			data.IsIdle = true
		} else {
			data.IsIdle = value.(bool)
		}
		// if is idle, remove playing media
		if data.IsIdle {
			global.EventManager.CallA(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
				Media:   model.Media{},
				Removed: true,
			})
		}
		global.EventManager.CallA(events.PlayerPropertyIdleActiveUpdate, data)

	},
	"time-pos": func(value any) {
		var data events.PlayerPropertyTimePosUpdateEvent
		if value == nil {
			data.TimePos = 0
		} else {
			data.TimePos = value.(float64)
		}
		// ignore bug value
		if data.TimePos < 0.1 {
			return
		}
		global.EventManager.CallA(events.PlayerPropertyTimePosUpdate, data)
	},
	"duration": func(value any) {
		var data events.PlayerPropertyDurationUpdateEvent
		if value == nil {
			data.Duration = 0
		} else {
			data.Duration = value.(float64)
		}
		global.EventManager.CallA(events.PlayerPropertyDurationUpdate, data)
	},
	"volume": func(value any) {
		var data events.PlayerPropertyVolumeUpdateEvent
		if value == nil {
			data.Volume = 0
		} else {
			data.Volume = value.(float64)
		}
		global.EventManager.CallA(events.PlayerPropertyVolumeUpdate, data)
	},
}

func registerHandler() {
	var err error
	for property, _ := range mpvPropertyHandler {
		log.Infof("register handler for mpv property %s", property)
		err = libmpv.ObserveProperty(util.Hash64(property), property, mpv.FORMAT_NODE)
		if err != nil {
			log.Errorf("register handler for mpv property %s failed: %s", property, err)
		}
	}
}

func registerCmdHandler() {
	global.EventManager.RegisterA(events.PlayerPlayCmd, "player.play", func(evnt *event.Event) {
		mediaInfo := evnt.Data.(events.PlayerPlayCmdEvent).Media.Info
		log.Infof("[MPV Player] Play media %s", mediaInfo.Title)
		mediaUrls, err := miaosic.GetMediaUrl(mediaInfo.Meta, miaosic.QualityAny)
		if err != nil || len(mediaUrls) == 0 {
			log.Warn("[MPV PlayControl] get media url failed", err)
			global.EventManager.CallA(
				events.PlayerPlayErrorUpdate,
				events.PlayerPlayErrorUpdateEvent{
					Error: err,
				})
			return
		}
		mediaUrl := mediaUrls[0]
		if val, ok := mediaUrl.Header["User-Agent"]; ok {
			log.Debug("[MPV PlayControl] set user-agent for mpv player")
			err := libmpv.SetPropertyString("user-agent", val)
			if err != nil {
				log.Warn("[MPV PlayControl] set player user-agent failed", err)
				return
			}
		}

		if val, ok := mediaUrl.Header["Referer"]; ok {
			log.Debug("[MPV PlayControl] set referrer for mpv player")
			err := libmpv.SetPropertyString("referrer", val)
			if err != nil {
				log.Warn("[MPV PlayControl] set player referrer failed", err)
				return
			}
		}
		media := evnt.Data.(events.PlayerPlayCmdEvent).Media
		if m, err := miaosic.GetMediaInfo(media.Info.Meta); err == nil {
			media.Info = m
		}
		global.EventManager.CallA(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
			Media:   media,
			Removed: false,
		})
		log.Debugf("mpv command load file %s %s", mediaInfo.Title, mediaUrl.Url)
		if err := libmpv.Command([]string{"loadfile", mediaUrl.Url}); err != nil {
			log.Warn("[MPV PlayControl] mpv load media failed", mediaInfo)
			global.EventManager.CallA(
				events.PlayerPlayErrorUpdate,
				events.PlayerPlayErrorUpdateEvent{
					Error: err,
				})
			return
		}
	})
	global.EventManager.RegisterA(events.PlayerToggleCmd, "player.toggle", func(evnt *event.Event) {
		property, err := libmpv.GetProperty("pause", mpv.FORMAT_FLAG)
		if err != nil {
			log.Warn("[MPV PlayControl] get property pause failed", err)
			return
		}
		err = libmpv.SetProperty("pause", mpv.FORMAT_FLAG, !property.(bool))
		if err != nil {
			log.Warn("[MPV PlayControl] toggle pause failed", err)
		}
	})
	global.EventManager.RegisterA(events.PlayerSeekCmd, "player.seek", func(evnt *event.Event) {
		data := evnt.Data.(events.PlayerSeekCmdEvent)
		log.Debugf("seek to %f (absolute=%t)", data.Position, data.Absolute)
		var err error
		if data.Absolute {
			err = libmpv.SetProperty("time-pos", mpv.FORMAT_DOUBLE, data.Position)
		} else {
			err = libmpv.SetProperty("percent-pos", mpv.FORMAT_DOUBLE, data.Position)
		}
		if err != nil {
			log.Warn("seek failed", err)
		}
	})
	global.EventManager.RegisterA(events.PlayerVolumeChangeCmd, "player.volume", func(evnt *event.Event) {
		data := evnt.Data.(events.PlayerVolumeChangeCmdEvent)
		err := libmpv.SetProperty("volume", mpv.FORMAT_DOUBLE, data.Volume)
		if err != nil {
			log.Warn("set volume failed", err)
		}
	})

	global.EventManager.RegisterA(events.PlayerVideoPlayerSetWindowHandleCmd, "player.next", func(evnt *event.Event) {
		handle := evnt.Data.(events.PlayerVideoPlayerSetWindowHandleCmdEvent).Handle
		err := SetWindowHandle(handle)
		if err != nil {
			log.Warn("set window handle failed", err)
		}
	})

	global.EventManager.RegisterA(events.PlayerSetAudioDeviceCmd, "player.set_audio_device", func(evnt *event.Event) {
		device := evnt.Data.(events.PlayerSetAudioDeviceCmdEvent).Device
		err := libmpv.SetPropertyString("audio-device", device)
		if err != nil {
			global.EventManager.CallA(
				events.ErrorUpdate,
				events.ErrorUpdateEvent{
					Error: err,
				})
			log.Warn("set audio device failed", err)
		}
		log.Infof("set audio device to %s", device)
		return
	})
}

func SetWindowHandle(handle uintptr) error {
	log.Infof("set window handle %d", handle)
	_ = libmpv.SetOptionString("wid", fmt.Sprintf("%d", handle))
	return libmpv.SetOptionString("vo", "gpu")
}

// // updateAudioDeviceList get output device for mpv
// // return format is []AudioDevice
func updateAudioDeviceList() {
	property, err := libmpv.GetProperty("audio-device-list", mpv.FORMAT_STRING)
	if err != nil {
		return
	}
	ad := libmpv.GetPropertyString("audio-device")
	dl := make([]model.AudioDevice, 0)
	gjson.Parse(property.(string)).ForEach(func(key, value gjson.Result) bool {
		dl = append(dl, model.AudioDevice{
			Name:        value.Get("name").String(),
			Description: value.Get("description").String(),
		})
		return true
	})
	log.Infof("update audio device list %v, current %s", dl, ad)
	global.EventManager.CallA(events.PlayerAudioDeviceUpdate, events.PlayerAudioDeviceUpdateEvent{
		Current: ad,
		Devices: dl,
	})
	return
}
