package internal

import (
	"AynaLivePlayer/adapters/provider"
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"errors"
)

type PlayController struct {
	playing      *model.Media `ini:"-"`
	AudioDevice  string
	Volume       float64
	Cfg          adapter.IPlayControlConfig
	eventManager *event.Manager              `ini:"-"`
	player       adapter.IPlayer             `ini:"-"`
	playlist     adapter.IPlaylistController `ini:"-"`
	provider     adapter.IProviderController `ini:"-"`
	lyric        adapter.ILyricLoader        `ini:"-"`
	log          adapter.ILogger             `ini:"-"`
}

func (pc *PlayController) Config() *adapter.IPlayControlConfig {
	return &pc.Cfg
}

func (pc *PlayController) Name() string {
	return "PlayController"
}

func (pc *PlayController) OnLoad() {
	return
}

func (pc *PlayController) OnSave() {
	return
}

func NewPlayerController(
	player adapter.IPlayer,
	playlist adapter.IPlaylistController,
	lyric adapter.ILyricLoader,
	provider adapter.IProviderController,
	log adapter.ILogger) adapter.IPlayController {
	pc := &PlayController{
		eventManager: event.MainManager.NewChildManager(),
		player:       player,
		playlist:     playlist,
		lyric:        lyric,
		provider:     provider,
		playing:      &model.Media{},
		AudioDevice:  "auto",
		Volume:       100,
		Cfg: adapter.IPlayControlConfig{
			SkipPlaylist:     false,
			AutoNextWhenFail: true,
		},
		log: log,
	}
	config.LoadConfig(pc)
	pc.SetVolume(pc.Volume)
	pc.SetAudioDevice(pc.AudioDevice)
	pc.player.ObserveProperty(model.PlayerPropIdleActive, "controller.playcontrol.idleplaynext", pc.handleMpvIdlePlayNext)
	pc.playlist.GetCurrent().EventManager().RegisterA(events.EventPlaylistInsert, "controller.playcontrol.playlistadd", pc.handlePlaylistAdd)
	pc.player.ObserveProperty(model.PlayerPropTimePos, "controller.playcontrol.updatelyric", pc.handleLyricUpdate)
	return pc
}

func (pc *PlayController) handleMpvIdlePlayNext(event *event.Event) {
	isIdle := event.Data.(events.PlayerPropertyUpdateEvent).Value.(bool)
	if isIdle {
		pc.log.Info("[Controller] mpv went idle, try play next")
		pc.PlayNext()
	}
}

func (pc *PlayController) handlePlaylistAdd(event *event.Event) {
	if pc.player.IsIdle() {
		pc.PlayNext()
		return
	}
	pc.log.Debugf("[PlayController] playlist add event, SkipPlaylist=%t", pc.Cfg.SkipPlaylist)
	if pc.Cfg.SkipPlaylist && pc.playing != nil && pc.playing.User == PlaylistUser {
		pc.PlayNext()
		return
	}
}

func (pc *PlayController) handleLyricUpdate(event *event.Event) {
	data := event.Data.(events.PlayerPropertyUpdateEvent).Value
	if data == nil {
		return
	}
	pc.lyric.Update(data.(float64))
}

func (pc *PlayController) EventManager() *event.Manager {
	return pc.eventManager
}

func (pc *PlayController) GetPlaying() *model.Media {
	return pc.playing
}

func (pc *PlayController) GetPlayer() adapter.IPlayer {
	return pc.player
}

func (pc *PlayController) GetLyric() adapter.ILyricLoader {
	return pc.lyric
}

func (pc *PlayController) PlayNext() {
	pc.log.Infof("[PlayController] try to play next possible media")
	if pc.playlist.GetCurrent().Size() == 0 && pc.playlist.GetDefault().Size() == 0 {
		return
	}
	var media *model.Media
	if pc.playlist.GetCurrent().Size() != 0 {
		media = pc.playlist.GetCurrent().Pop().Copy()
	} else if pc.playlist.GetDefault().Size() != 0 {
		media = pc.playlist.GetDefault().Next().Copy()
		media.User = PlaylistUser
	}
	_ = pc.Play(media)
}

func (pc *PlayController) Play(media *model.Media) error {
	pc.log.Infof("[PlayController] prepare media %s", media.Title)
	err := pc.provider.PrepareMedia(media)
	if err != nil {
		pc.log.Warn("[PlayController] prepare media failed, try play next")
		if pc.Cfg.AutoNextWhenFail {
			go pc.PlayNext()
		}
		//pc.PlayNext()
		return errors.New("prepare media failed")
	}
	pc.eventManager.CallA(events.EventPlay, events.PlayEvent{
		Media: media,
	})
	pc.playing = media
	pc.playlist.AddToHistory(media)
	if err := pc.player.Play(media); err != nil {
		pc.log.Warn("[PlayController] play failed", err)
		return errors.New("player play failed")
	}
	pc.eventManager.CallA(events.EventPlayed, events.PlayEvent{
		Media: media,
	})
	pc.lyric.Reload(media.Lyric)
	// reset
	media.Url = ""
	return nil
}

func (pc *PlayController) Add(keyword string, user interface{}) {
	media := pc.provider.MediaMatch(keyword)
	if media == nil {
		medias, err := pc.provider.Search(keyword)
		if err != nil {
			pc.log.Warnf("[PlayController] search for %s, got error %s", keyword, err)
			return
		}
		if len(medias) == 0 {
			pc.log.Infof("[PlayController] search for %s, got no result", keyword)
			return
		}
		media = medias[0]
	}
	media.User = user
	pc.log.Infof("[PlayController] add media %s (%s)", media.Title, media.Artist)
	pc.playlist.GetCurrent().Insert(-1, media)
}

func (pc *PlayController) AddWithProvider(keyword string, pname string, user interface{}) {
	media := provider.MatchMedia(pname, keyword)
	if media == nil {
		medias, err := provider.Search(pname, keyword)
		if err != nil {
			pc.log.Warnf("[PlayController] search for %s, got error %s", keyword, err)
			return
		}
		if len(medias) == 0 {
			pc.log.Infof("[PlayController] search for %s, got no result", keyword)
			return
		}

		media = medias[0]
	}
	media.User = user
	pc.log.Infof("[PlayController] add media %s (%s)", media.Title, media.Artist)
	pc.playlist.GetCurrent().Insert(-1, media)
}

func (pc *PlayController) Seek(position float64, absolute bool) {
	if err := pc.player.Seek(position, absolute); err != nil {
		pc.log.Warnf("[PlayController] seek to position %f (%t) failed, %s", position, absolute, err)
	}
}

func (pc *PlayController) Toggle() (b bool) {
	var err error
	if pc.player.IsPaused() {
		err = pc.player.Unpause()
		b = false
	} else {
		err = pc.player.Pause()
		b = true
	}
	if err != nil {
		pc.log.Warn("[PlayController] toggle failed", err)
	}
	return
}

func (pc *PlayController) SetVolume(volume float64) {
	if pc.player.SetVolume(volume) != nil {
		pc.log.Warnf("[PlayController] set mpv volume to %f failed", volume)
		return
	}
	pc.Volume = volume
}

func (pc *PlayController) Destroy() {
	pc.player.Stop()
}

func (pc *PlayController) GetCurrentAudioDevice() string {
	return pc.AudioDevice
}

func (pc *PlayController) GetAudioDevices() []model.AudioDevice {
	dl, err := pc.player.GetAudioDeviceList()
	if err != nil {
		return make([]model.AudioDevice, 0)
	}
	return dl
}

func (pc *PlayController) SetAudioDevice(device string) {
	pc.log.Infof("[PlayController] set audio device to %s", device)
	if err := pc.player.SetAudioDevice(device); err != nil {
		pc.log.Warnf("[PlayController] set mpv audio device to %s failed, %s", device, err)
		_ = pc.player.SetAudioDevice("auto")
		pc.AudioDevice = "auto"
		return
	}
	pc.AudioDevice = device
}
