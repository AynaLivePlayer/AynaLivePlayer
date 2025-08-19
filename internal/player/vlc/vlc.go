package vlc

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/logger"
	"fmt"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/adrg/libvlc-go/v3"
	"math"
	"runtime"
	"strings"
	"sync"
)

var running bool = false
var log logger.ILogger = nil
var player *vlc.Player
var eventManager *vlc.EventManager
var lock sync.Mutex

// 状态变量
var prevPercentPos float64 = 0
var prevTimePos float64 = 0
var duration float64 = 0
var currentMedia model.Media
var currentWindowHandle uintptr

var audioDevices []model.AudioDevice
var currentAudioDevice string

var videoOptions = map[string][]string{
	"windows": {"--video-filter=adjust", "--directx-hwnd"},
	"darwin":  {"--vout=macosx"},
	"linux":   {"--vout=x11", "--x11-display=:0"},
}

func setWindowHandle(handle uintptr) error {
	return nil
	os := runtime.GOOS
	switch os {
	case "windows":
		// Windows 平台使用 DirectX
		player.SetHWND(uintptr(handle))
	case "darwin":
		// macOS 平台使用 NSView
		player.SetNSObject(handle)
	case "linux":
		// Linux 平台使用 XWindow
		player.SetXWindow(uint32(handle))
	default:
		return fmt.Errorf("unsupported platform: %s", os)
	}

	currentWindowHandle = handle

	// 如果当前有媒体在播放，需要重新加载视频输出
	if player.IsPlaying() {
		player.Stop()
		player.Play()
	}

	return nil
}

func SetupPlayer() {
	running = true
	config.LoadConfig(cfg)
	log = global.Logger.WithPrefix("VLC Player")

	opts := []string{"--no-video", "--quiet"}
	//os := runtime.GOOS
	//if platformOpts, ok := videoOptions[os]; ok {
	//	opts = append(opts, platformOpts...)
	//}

	// 初始化libvlc
	if err := vlc.Init(opts...); err != nil {
		log.Error("initialize libvlc failed: ", err)
		return
	}

	// 创建播放器
	var err error
	player, err = vlc.NewPlayer()
	if err != nil {
		log.Error("create player failed: ", err)
		return
	}

	// 获取事件管理器
	eventManager, err = player.EventManager()
	if err != nil {
		log.Error("get event manager failed: ", err)
		return
	}

	// 注册事件
	registerEvents()
	registerCmdHandler()
	updateAudioDeviceList()
	restoreConfig()
	log.Info("VLC player initialized")
}

func StopPlayer() {
	log.Info("stopping VLC player")
	if currentAudioDevice != "" {
		cfg.AudioDevice = currentAudioDevice
		log.Infof("save audio device config: %s", cfg.AudioDevice)
	}
	running = false
	if player != nil {
		err := player.Stop()
		if err != nil {
			log.Error("stop player failed: ", err)
		}
		err = player.Release()
		if err != nil {
			log.Error("release player failed: ", err)
		}
	}
	err := vlc.Release()
	if err != nil {
		log.Error("release player failed: ", err)
	}
	log.Info("VLC player stopped")
}

func registerEvents() {
	// 播放结束事件
	_, err := eventManager.Attach(vlc.MediaPlayerEndReached, func(e vlc.Event, userData interface{}) {
		global.EventManager.CallA(events.PlayerPropertyStateUpdate, events.PlayerPropertyStateUpdateEvent{
			State: model.PlayerStateIdle,
		})
		global.EventManager.CallA(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
			Media:   model.Media{},
			Removed: true,
		})
	}, nil)
	if err != nil {
		log.Error("register MediaPlayerEndReached event failed: ", err)
	}

	// 播放位置改变事件
	_, err = eventManager.Attach(vlc.MediaPlayerPositionChanged, func(e vlc.Event, userData interface{}) {
		pos32, _ := player.MediaPosition()
		pos := float64(pos32)
		if duration > 0 {
			timePos := pos * duration
			percentPos := pos * 100
			// 忽略小变化
			if math.Abs(timePos-prevTimePos) < 0.5 && math.Abs(percentPos-prevPercentPos) < 0.5 {
				return
			}
			prevTimePos = timePos
			prevPercentPos = percentPos
			global.EventManager.CallA(events.PlayerPropertyTimePosUpdate, events.PlayerPropertyTimePosUpdateEvent{
				TimePos: timePos,
			})
			global.EventManager.CallA(events.PlayerPropertyPercentPosUpdate, events.PlayerPropertyPercentPosUpdateEvent{
				PercentPos: percentPos,
			})
		}
	}, nil)
	if err != nil {
		log.Error("register MediaPlayerPositionChanged event failed: ", err)
	}

	// 时间改变事件（获取时长）
	_, err = eventManager.Attach(vlc.MediaPlayerTimeChanged, func(e vlc.Event, userData interface{}) {
		dur, _ := player.MediaLength()
		duration = float64(dur) / 1000.0 // 转换为秒
		global.EventManager.CallA(events.PlayerPropertyDurationUpdate, events.PlayerPropertyDurationUpdateEvent{
			Duration: duration,
		})
	}, nil)
	if err != nil {
		log.Error("register MediaPlayerTimeChanged event failed: ", err)
	}

	// 暂停状态改变
	_, err = eventManager.Attach(vlc.MediaPlayerPaused, func(e vlc.Event, userData interface{}) {
		log.Debug("VLC player paused")
		global.EventManager.CallA(events.PlayerPropertyPauseUpdate, events.PlayerPropertyPauseUpdateEvent{
			Paused: true,
		})
	}, nil)
	if err != nil {
		log.Error("register MediaPlayerPaused event failed: ", err)
	}

	_, err = eventManager.Attach(vlc.MediaPlayerPlaying, func(e vlc.Event, userData interface{}) {
		log.Debug("VLC player playing")
		global.EventManager.CallA(events.PlayerPropertyPauseUpdate, events.PlayerPropertyPauseUpdateEvent{
			Paused: false,
		})
	}, nil)
	if err != nil {
		log.Error("register MediaPlayerPlaying event failed: ", err)
	}

	_, err = eventManager.Attach(vlc.MediaPlayerAudioVolume, func(e vlc.Event, userData interface{}) {
		volume, _ := player.Volume()
		log.Debug("VLC player audio volume: ", volume)
		global.EventManager.CallA(events.PlayerPropertyVolumeUpdate, events.PlayerPropertyVolumeUpdateEvent{
			Volume: float64(volume),
		})
	}, nil)
}

func registerCmdHandler() {
	global.EventManager.RegisterA(events.PlayerPlayCmd, "player.play", func(evnt *event.Event) {
		lock.Lock()
		defer lock.Unlock()

		mediaInfo := evnt.Data.(events.PlayerPlayCmdEvent).Media.Info
		log.Infof("[VLC Player] Play media %s", mediaInfo.Title)

		mediaUrls, err := miaosic.GetMediaUrl(mediaInfo.Meta, miaosic.QualityAny)
		if err != nil || len(mediaUrls) == 0 {
			log.Warn("[VLC PlayControl] get media url failed ", mediaInfo.Meta.ID(), err)
			global.EventManager.CallA(
				events.PlayerPlayErrorUpdate,
				events.PlayerPlayErrorUpdateEvent{
					Error: err,
				})
			return
		}

		// 创建媒体对象
		var media *vlc.Media
		log.Debugf("[VLC PlayControl] get player media %s", mediaUrls[0].Url)
		if strings.HasPrefix(mediaUrls[0].Url, "http") {
			media, err = vlc.NewMediaFromURL(mediaUrls[0].Url)
		} else {
			media, err = vlc.NewMediaFromPath(mediaUrls[0].Url)
		}
		if err != nil {
			log.Error("create media failed: ", err)
			return
		}

		// 设置HTTP头
		if val, ok := mediaUrls[0].Header["User-Agent"]; ok {
			err = media.AddOptions(":http-user-agent=" + val)
			if err != nil {
				log.Warn("add http-user-agent options failed: ", err)
			}
		}
		if val, ok := mediaUrls[0].Header["Referer"]; ok {
			err = media.AddOptions(":http-referrer=" + val)
			if err != nil {
				log.Warn("add http-referrer options failed: ", err)
			}
		}

		// 更新媒体信息
		mediaData := evnt.Data.(events.PlayerPlayCmdEvent).Media
		if m, err := miaosic.GetMediaInfo(mediaData.Info.Meta); err == nil {
			mediaData.Info = m
		}
		currentMedia = mediaData

		global.EventManager.CallA(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
			Media:   mediaData,
			Removed: false,
		})

		// 播放
		if err := player.SetMedia(media); err != nil {
			log.Error("set media failed: ", err)
			return
		}

		if currentWindowHandle != 0 {
			if err := setWindowHandle(currentWindowHandle); err != nil {
				log.Error("apply window handle failed: ", err)
			}
		}

		if err := player.Play(); err != nil {
			log.Error("play failed: ", err)
			return
		}

		// 重置位置信息
		prevPercentPos = 0
		prevTimePos = 0
		global.EventManager.CallA(events.PlayerPropertyTimePosUpdate, events.PlayerPropertyTimePosUpdateEvent{
			TimePos: 0,
		})
		global.EventManager.CallA(events.PlayerPropertyPercentPosUpdate, events.PlayerPropertyPercentPosUpdateEvent{
			PercentPos: 0,
		})
		global.EventManager.CallA(events.PlayerPropertyStateUpdate, events.PlayerPropertyStateUpdateEvent{
			State: model.PlayerStatePlaying,
		})
	})

	global.EventManager.RegisterA(events.PlayerToggleCmd, "player.toggle", func(evnt *event.Event) {
		lock.Lock()
		defer lock.Unlock()
		err := player.TogglePause()
		if err != nil {
			log.Errorf("[VLC Player] Toggle pause failed: %v", err)
			return
		}
	})

	global.EventManager.RegisterA(events.PlayerSetPauseCmd, "player.set_paused", func(evnt *event.Event) {
		lock.Lock()
		defer lock.Unlock()
		data := evnt.Data.(events.PlayerSetPauseCmdEvent)
		err := player.SetPause(data.Pause)
		if err != nil {
			log.Errorf("[VLC Player] SetPause failed: %v", err)
			return
		}
	})

	global.EventManager.RegisterA(events.PlayerSeekCmd, "player.seek", func(evnt *event.Event) {
		lock.Lock()
		defer lock.Unlock()
		data := evnt.Data.(events.PlayerSeekCmdEvent)
		if data.Absolute {
			player.SetMediaTime(int(data.Position * 1000)) // 转换为毫秒
		} else {
			player.SetMediaPosition(float32(data.Position / 100))
		}
	})

	global.EventManager.RegisterA(events.PlayerVolumeChangeCmd, "player.volume", func(evnt *event.Event) {
		lock.Lock()
		defer lock.Unlock()
		data := evnt.Data.(events.PlayerVolumeChangeCmdEvent)
		err := player.SetVolume(int(data.Volume))
		if err != nil {
			log.Errorf("[VLC Player] SetVolume failed: %v", err)
		}
	})

	global.EventManager.RegisterA(events.PlayerVideoPlayerSetWindowHandleCmd, "player.set_window_handle", func(evnt *event.Event) {
		handle := evnt.Data.(events.PlayerVideoPlayerSetWindowHandleCmdEvent).Handle
		setWindowHandle(handle)
	})

	global.EventManager.RegisterA(events.PlayerSetAudioDeviceCmd, "player.set_audio_device", func(evnt *event.Event) {
		device := evnt.Data.(events.PlayerSetAudioDeviceCmdEvent).Device
		if err := setAudioDevice(device); err != nil {
			log.Warn("set audio device failed", err)
			global.EventManager.CallA(
				events.ErrorUpdate,
				events.ErrorUpdateEvent{
					Error: err,
				})
		}
	})
}

// setAudioDevice 设置音频输出设备
func setAudioDevice(deviceID string) error {
	lock.Lock()
	defer lock.Unlock()

	log.Infof("set audio device to: %s", deviceID)

	// 验证设备是否在列表中
	found := false
	for _, dev := range audioDevices {
		if dev.Name == deviceID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("audio device not found: %s", deviceID)
	}

	// 设置音频设备
	if err := player.SetAudioOutputDevice(deviceID, ""); err != nil {
		log.Error("set audio device failed: ", err)
		return err
	}

	currentAudioDevice = deviceID

	// 更新配置
	cfg.AudioDevice = deviceID

	// 发送更新事件
	global.EventManager.CallA(events.PlayerAudioDeviceUpdate, events.PlayerAudioDeviceUpdateEvent{
		Current: currentAudioDevice,
		Devices: audioDevices,
	})

	return nil
}

// updateAudioDeviceList 获取并更新音频设备列表
func updateAudioDeviceList() {
	lock.Lock()
	defer lock.Unlock()

	// 获取所有音频设备
	devices, err := player.AudioOutputDevices()
	if err != nil {
		log.Error("get audio device list failed: ", err)
		return
	}

	// 获取当前音频设备
	currentDevice, err := player.AudioOutputDevice()
	if err != nil {
		log.Warn("get current audio device failed: ", err)
		currentDevice = ""
	}
	log.Debugf("current audio device list: %s", devices)
	log.Debugf("current audio device: %s", currentDevice)

	// 转换设备格式
	audioDevices = make([]model.AudioDevice, 0, len(devices))
	for _, device := range devices {
		audioDevices = append(audioDevices, model.AudioDevice{
			Name:        device.Name,
			Description: device.Description,
		})
	}

	currentAudioDevice = currentDevice

	log.Infof("update audio device list: %d devices, current: %s",
		len(audioDevices), currentAudioDevice)

	// 发送事件通知
	global.EventManager.CallA(events.PlayerAudioDeviceUpdate, events.PlayerAudioDeviceUpdateEvent{
		Current: currentAudioDevice,
		Devices: audioDevices,
	})
}
