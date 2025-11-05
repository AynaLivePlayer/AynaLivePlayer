package mpv

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/logger"
	"AynaLivePlayer/pkg/util"
	"fmt"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/aynakeya/go-mpv"
	"github.com/tidwall/gjson"
	"math"
	"time"
)

var running bool = false
var libmpv *mpv.Mpv = nil
var log logger.ILogger = nil
var mpvClientVersion uint32 = 0

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
	mpvClientVersion = mpv.ClientApiVersion()
	log.Infof("libmpv version %d", mpv.ClientApiVersion())
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
				// should not call, otherwise StopPlayer gonna be call twice and cause panic
				// StopPlayer()
			}
		}
	}()
}

func StopPlayer() {
	cfg.AudioDevice = libmpv.GetPropertyString("audio-device")
	log.Debugf("successfully get audio-device and set config %s", cfg.AudioDevice)
	log.Info("stopping mpv player")
	running = false
	done := make(chan struct{})

	// Stop player async but wait for at most 1 second
	go func() {
		// todo: when call TerminateDestroy after wid has been set, a c code panic will arise.
		// maybe because the window mpv attach to has been closed. so handle was closed twice
		// for now. just don't destroy it. because it also might fix configuration
		// not properly saved issue
		libmpv.TerminateDestroy()
		close(done)
	}()

	select {
	case <-done:
		log.Info("mpv player stopped")
	case <-time.After(2 * time.Second):
		log.Error("mpv player stop timed out (2s) ")
	}
}

var prevPercentPos float64 = 0
var prevTimePos float64 = 0
var currentState = model.PlayerStateIdle

var mpvPropertyHandler = map[string]func(value any){
	"pause": func(value any) {
		var data events.PlayerPropertyPauseUpdateEvent
		log.Debugf("pause property update %v %T", value, value)
		data.Paused = value.(bool)
		_ = global.EventBus.Publish(events.PlayerPropertyPauseUpdate, data)
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
		// ignore small change
		if math.Abs(data.PercentPos-prevPercentPos) < 0.5 {
			return
		}
		prevPercentPos = data.PercentPos
		_ = global.EventBus.Publish(events.PlayerPropertyPercentPosUpdate, data)

	},
	"idle-active": func(value any) {
		var data events.PlayerPropertyStateUpdateEvent
		if value == nil {
			data.State = model.PlayerStateIdle
		} else {
			if value.(bool) {
				data.State = model.PlayerStateIdle
			} else {
				data.State = model.PlayerStatePlaying
			}
		}
		log.Debugf("mpv state update %v + %v = %v", currentState, data.State, currentState.NextState(data.State))
		currentState = currentState.NextState(data.State)
		data.State = currentState
		_ = global.EventBus.Publish(events.PlayerPropertyStateUpdate, data)

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
		// ignore small change
		if math.Abs(data.TimePos-prevTimePos) < 0.5 {
			return
		}
		prevTimePos = data.TimePos
		_ = global.EventBus.Publish(events.PlayerPropertyTimePosUpdate, data)
	},
	"duration": func(value any) {
		var data events.PlayerPropertyDurationUpdateEvent
		if value == nil {
			data.Duration = 0
		} else {
			data.Duration = value.(float64)
		}
		_ = global.EventBus.Publish(events.PlayerPropertyDurationUpdate, data)
	},
	"volume": func(value any) {
		var data events.PlayerPropertyVolumeUpdateEvent
		if value == nil {
			data.Volume = 0
		} else {
			data.Volume = value.(float64)
		}
		_ = global.EventBus.Publish(events.PlayerPropertyVolumeUpdate, data)
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
	global.EventBus.Subscribe("", events.PlayerPlayCmd, "player.play", func(evnt *eventbus.Event) {
		currentState = currentState.NextState(model.PlayerStateLoading)
		_ = global.EventBus.Publish(
			events.PlayerPropertyStateUpdate,
			events.PlayerPropertyStateUpdateEvent{
				State: currentState,
			})
		mediaInfo := evnt.Data.(events.PlayerPlayCmdEvent).Media.Info
		media := evnt.Data.(events.PlayerPlayCmdEvent).Media
		resp, err := global.EventBus.Call(events.CmdMiaosicGetMediaInfo, events.ReplyMiaosicGetMediaInfo,
			events.CmdMiaosicGetMediaInfoData{Meta: media.Info.Meta})
		if err == nil && resp.Data.(events.ReplyMiaosicGetMediaInfoData).Error == nil {
			media.Info = resp.Data.(events.ReplyMiaosicGetMediaInfoData).Info
		}
		_ = global.EventBus.Publish(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
			Media:   media,
			Removed: false,
		})
		log.Infof("[MPV Player] Play media %s", mediaInfo.Title)
		resp, err = global.EventBus.Call(events.CmdMiaosicGetMediaUrl, events.ReplyMiaosicGetMediaUrl,
			events.CmdMiaosicGetMediaUrlData{Meta: media.Info.Meta, Quality: miaosic.QualityAny})
		mediaUrls := resp.Data.(events.ReplyMiaosicGetMediaUrlData)
		if err != nil || mediaUrls.Error != nil || len(mediaUrls.Urls) == 0 {
			log.Warn("[MPV PlayControl] get media url failed ", mediaInfo.Meta.ID(), err)
			if err := libmpv.Command([]string{"stop"}); err != nil {
				log.Error("[MPV PlayControl] failed to stop", err)
			}
			_ = global.EventBus.Publish(
				events.PlayerPlayErrorUpdate,
				events.PlayerPlayErrorUpdateEvent{
					Error: err,
				})
			return
		}
		mediaUrl := mediaUrls.Urls[0]
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
		log.Debugf("mpv command loadfile %s %s", mediaInfo.Title, mediaUrl.Url)
		cmd := []string{"loadfile", mediaUrl.Url}
		if cfg.DisplayMusicCover && media.Info.Cover.Url != "" {
			// add media cover to video channel.
			// https://mpv.io/manual/master/#command-interface-[<options>]]]
			// api changes after client version 2.3 (0.38.0
			if mpvClientVersion >= ((2 << 16) | 3) {
				cmd = append(cmd, "replace", "0", "external-files-append=\""+media.Info.Cover.Url+"\",vid=1")
			} else {
				cmd = append(cmd, "replace", "external-files-append=\""+media.Info.Cover.Url+"\",vid=1")
			}
		}
		log.Debug("[MPV PlayControl] mpv command", cmd)
		if err := libmpv.Command(cmd); err != nil {
			log.Error("[MPV PlayControl] mpv load media failed", cmd, mediaInfo, err)
			_ = global.EventBus.Publish(
				events.PlayerPlayErrorUpdate,
				events.PlayerPlayErrorUpdateEvent{
					Error: err,
				})
			return
		}
		currentState = currentState.NextState(model.PlayerStatePlaying)
		_ = global.EventBus.Publish(
			events.PlayerPropertyStateUpdate,
			events.PlayerPropertyStateUpdateEvent{
				State: currentState,
			})
		_ = global.EventBus.Publish(events.PlayerPropertyTimePosUpdate, events.PlayerPropertyTimePosUpdateEvent{
			TimePos: 0,
		})
		_ = global.EventBus.Publish(events.PlayerPropertyPercentPosUpdate, events.PlayerPropertyPercentPosUpdateEvent{
			PercentPos: 0,
		})
	})
	global.EventBus.Subscribe("", events.PlayerToggleCmd, "player.toggle", func(evnt *eventbus.Event) {
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
	global.EventBus.Subscribe("", events.PlayerSetPauseCmd, "player.set_paused", func(evnt *eventbus.Event) {
		data := evnt.Data.(events.PlayerSetPauseCmdEvent)
		err := libmpv.SetProperty("pause", mpv.FORMAT_FLAG, data.Pause)
		if err != nil {
			log.Warn("[MPV PlayControl] set pause failed", err)
		}
	})
	global.EventBus.Subscribe("", events.PlayerSeekCmd, "player.seek", func(evnt *eventbus.Event) {
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
	global.EventBus.Subscribe("", events.PlayerVolumeChangeCmd, "player.volume", func(evnt *eventbus.Event) {
		data := evnt.Data.(events.PlayerVolumeChangeCmdEvent)
		err := libmpv.SetProperty("volume", mpv.FORMAT_DOUBLE, data.Volume)
		if err != nil {
			log.Warn("set volume failed", err)
		}
	})

	global.EventBus.Subscribe("", events.PlayerVideoPlayerSetWindowHandleCmd, "player.set_window_handle", func(evnt *eventbus.Event) {
		handle := evnt.Data.(events.PlayerVideoPlayerSetWindowHandleCmdEvent).Handle
		err := SetWindowHandle(handle)
		if err != nil {
			log.Warn("set window handle failed", err)
		}
	})

	global.EventBus.Subscribe("", events.PlayerSetAudioDeviceCmd, "player.set_audio_device", func(evnt *eventbus.Event) {
		device := evnt.Data.(events.PlayerSetAudioDeviceCmdEvent).Device
		err := libmpv.SetPropertyString("audio-device", device)
		if err != nil {
			_ = global.EventBus.Publish(
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
	_ = global.EventBus.Publish(events.PlayerAudioDeviceUpdate, events.PlayerAudioDeviceUpdateEvent{
		Current: ad,
		Devices: dl,
	})
	return
}
