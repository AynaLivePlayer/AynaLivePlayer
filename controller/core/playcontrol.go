package core

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
)

type PlayController struct {
	eventManager *event.Manager                 `ini:"-"`
	player       player.IPlayer                 `ini:"-"`
	playlist     controller.IPlaylistController `ini:"-"`
	provider     controller.IProviderController `ini:"-"`
	lyric        controller.ILyricLoader        `ini:"-"`
	playing      *model.Media                   `ini:"-"`
	AudioDevice  string
	Volume       float64
	SkipPlaylist bool
}

func (pc *PlayController) GetSkipPlaylist() bool {
	return pc.SkipPlaylist
}

func (pc *PlayController) SetSkipPlaylist(b bool) {
	pc.SkipPlaylist = b
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
	player player.IPlayer,
	playlist controller.IPlaylistController,
	lyric controller.ILyricLoader,
	provider controller.IProviderController) controller.IPlayController {
	pc := &PlayController{
		eventManager: event.MainManager.NewChildManager(),
		player:       player,
		playlist:     playlist,
		lyric:        lyric,
		provider:     provider,
		playing:      &model.Media{},
		AudioDevice:  "auto",
		Volume:       100,
		SkipPlaylist: false,
	}
	config.LoadConfig(pc)
	pc.SetVolume(pc.Volume)
	pc.SetAudioDevice(pc.AudioDevice)
	pc.player.ObserveProperty(model.PlayerPropIdleActive, "controller.playcontrol.idleplaynext", pc.handleMpvIdlePlayNext)
	pc.playlist.GetCurrent().EventManager().RegisterA(model.EventPlaylistInsert, "controller.playcontrol.playlistadd", pc.handlePlaylistAdd)
	pc.player.ObserveProperty(model.PlayerPropTimePos, "controller.playcontrol.updatelyric", pc.handleLyricUpdate)
	return pc
}

func (pc *PlayController) handleMpvIdlePlayNext(event *event.Event) {
	isIdle := event.Data.(model.PlayerPropertyUpdateEvent).Value.(bool)
	if isIdle {
		lg.Info("[Controller] mpv went idle, try play next")
		pc.PlayNext()
	}
}

func (pc *PlayController) handlePlaylistAdd(event *event.Event) {
	if pc.player.IsIdle() {
		pc.PlayNext()
		return
	}
	lg.Debugf("[PlayController] playlist add event, SkipPlaylist=%t", pc.SkipPlaylist)
	if pc.SkipPlaylist && pc.playing != nil && pc.playing.User == controller.PlaylistUser {
		pc.PlayNext()
		return
	}
}

func (pc *PlayController) handleLyricUpdate(event *event.Event) {
	data := event.Data.(model.PlayerPropertyUpdateEvent).Value
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

func (pc *PlayController) GetPlayer() player.IPlayer {
	return pc.player
}

func (pc *PlayController) GetLyric() controller.ILyricLoader {
	return pc.lyric
}

func (pc *PlayController) PlayNext() {
	lg.Infof("[PlayController] try to play next possible media")
	if pc.playlist.GetCurrent().Size() == 0 && pc.playlist.GetDefault().Size() == 0 {
		return
	}
	var media *model.Media
	if pc.playlist.GetCurrent().Size() != 0 {
		media = pc.playlist.GetCurrent().Pop().Copy()
	} else if pc.playlist.GetDefault().Size() != 0 {
		media = pc.playlist.GetDefault().Next().Copy()
		media.User = controller.PlaylistUser
	}
	pc.Play(media)
}

func (pc *PlayController) Play(media *model.Media) {
	lg.Infof("[PlayController] prepare media %s", media.Title)
	err := pc.provider.PrepareMedia(media)
	if err != nil {
		lg.Warn("[PlayController] prepare media failed. try play next")
		pc.PlayNext()
		return
	}
	pc.playing = media
	pc.playlist.AddToHistory(media)
	if err := pc.player.Play(media); err != nil {
		lg.Warn("[PlayController] play failed", err)
		return
	}
	pc.eventManager.CallA(model.EventPlay, model.PlayEvent{
		Media: media,
	})
	pc.lyric.Reload(media.Lyric)
	// reset
	media.Url = ""
}

func (pc *PlayController) Add(keyword string, user interface{}) {
	media := pc.provider.MediaMatch(keyword)
	if media == nil {
		medias, err := pc.provider.Search(keyword)
		if err != nil {
			lg.Warnf("[PlayController] search for %s, got error %s", keyword, err)
			return
		}
		if len(medias) == 0 {
			lg.Info("[PlayController] search for %s, got no result", keyword)
			return
		}
		media = medias[0]
	}
	media.User = user
	lg.Infof("[PlayController] add media %s (%s)", media.Title, media.Artist)
	pc.playlist.GetCurrent().Insert(-1, media)
}

func (pc *PlayController) AddWithProvider(keyword string, pname string, user interface{}) {
	media := provider.MatchMedia(pname, keyword)
	if media == nil {
		medias, err := provider.Search(pname, keyword)
		if err != nil {
			lg.Warnf("[PlayController] search for %s, got error %s", keyword, err)
			return
		}
		if len(medias) == 0 {
			lg.Infof("[PlayController] search for %s, got no result", keyword)
			return
		}
		media = medias[0]
	}
	media.User = user
	lg.Infof("[PlayController] add media %s (%s)", media.Title, media.Artist)
	pc.playlist.GetCurrent().Insert(-1, media)
}

func (pc *PlayController) Seek(position float64, absolute bool) {
	if err := pc.player.Seek(position, absolute); err != nil {
		lg.Warnf("[PlayController] seek to position %f (%t) failed, %s", position, absolute, err)
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
		lg.Warn("[PlayController] toggle failed", err)
	}
	return
}

func (pc *PlayController) SetVolume(volume float64) {
	if pc.player.SetVolume(volume) != nil {
		lg.Warnf("[PlayController] set mpv volume to %f failed", volume)
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
	lg.Infof("[PlayController] set audio device to %s", device)
	if err := pc.player.SetAudioDevice(device); err != nil {
		lg.Warnf("[PlayController] set mpv audio device to %s failed, %s", device, err)
		_ = pc.player.SetAudioDevice("auto")
		pc.AudioDevice = "auto"
		return
	}
	pc.AudioDevice = device
}
