package core

import (
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"AynaLivePlayer/repo/provider"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type PlaylistController struct {
	PlaylistPath          string
	provider              controller.IProviderController
	History               controller.IPlaylist   `ini:"-"`
	Current               controller.IPlaylist   `ini:"-"`
	Default               controller.IPlaylist   `ini:"-"`
	Playlists             []controller.IPlaylist `ini:"-"`
	DefaultIndex          int
	CurrentPlaylistRandom bool
	DefaultPlaylistRandom bool
}

func NewPlaylistController(
	provider controller.IProviderController) controller.IPlaylistController {
	pc := &PlaylistController{
		PlaylistPath:          "playlists.json",
		provider:              provider,
		History:               NewPlaylist("history"),
		Default:               NewPlaylist("default"),
		Current:               NewPlaylist("current"),
		Playlists:             make([]controller.IPlaylist, 0),
		DefaultIndex:          0,
		CurrentPlaylistRandom: false,
		DefaultPlaylistRandom: true,
	}
	config.LoadConfig(pc)
	if pc.DefaultIndex < 0 || pc.DefaultIndex >= len(pc.Playlists) {
		pc.DefaultIndex = 0
		lg.Warn("playlist index did not find")
	}
	go func() {
		_ = pc.SetDefault(pc.DefaultIndex)
	}()
	return pc
}

func (pc *PlaylistController) Name() string {
	return "Playlists"
}

func (pc *PlaylistController) OnLoad() {
	var metas = []model.Meta{
		{
			"netease",
			"2382819181",
		},
		{"netease",
			"4987059624",
		},
		{"local",
			"list1",
		},
	}
	_ = config.LoadJson(pc.PlaylistPath, &metas)
	for _, m := range metas {
		p := NewPlaylist(fmt.Sprintf("%s-%s", m.Name, m.Id))
		p.Model().Meta = m
		pc.Playlists = append(pc.Playlists, p)
	}
	if pc.CurrentPlaylistRandom {
		pc.Current.Model().Mode = model.PlaylistModeRandom
	}
	if pc.DefaultPlaylistRandom {
		pc.Default.Model().Mode = model.PlaylistModeRandom
	}
}

func (pc *PlaylistController) OnSave() {
	var metas = make([]model.Meta, 0)
	for _, pl := range pc.Playlists {
		metas = append(metas, pl.Model().Meta)
	}
	_ = config.SaveJson(pc.PlaylistPath, &metas)

	if pc.Current.Model().Mode == model.PlaylistModeRandom {
		pc.CurrentPlaylistRandom = true
	} else {
		pc.CurrentPlaylistRandom = false
	}
	if pc.Default.Model().Mode == model.PlaylistModeRandom {
		pc.DefaultPlaylistRandom = true
	} else {
		pc.DefaultPlaylistRandom = false
	}
}

func (pc *PlaylistController) Size() int {
	return len(pc.Playlists)
}

func (pc *PlaylistController) GetHistory() controller.IPlaylist {
	return pc.History
}

func (pc *PlaylistController) GetDefault() controller.IPlaylist {
	return pc.Default
}

func (pc *PlaylistController) GetCurrent() controller.IPlaylist {
	return pc.Current
}

func (pc *PlaylistController) AddToHistory(media *model.Media) {
	lg.Tracef("add media %s (%s) to history", media.Title, media.Artist)
	media = media.Copy()
	// reset url for future use
	media.Url = ""
	if pc.History.Size() >= 1024 {
		pc.History.Replace([]*model.Media{})
	}
	media.User = controller.HistoryUser
	pc.History.Push(media)
	return
}

func (pc *PlaylistController) Get(index int) controller.IPlaylist {
	if index < 0 || index >= len(pc.Playlists) {
		lg.Warnf("playlist.index=%d not found", index)
		return nil
	}
	return pc.Playlists[index]
}

func (pc *PlaylistController) Add(pname string, uri string) controller.IPlaylist {
	lg.Infof("try add playlist %s with provider %s", uri, pname)
	id, err := provider.FormatPlaylistUrl(pname, uri)
	if err != nil || id == "" {
		lg.Warnf("fail to format %s playlist id for %s", uri, pname)
		return nil
	}
	p := NewPlaylist(fmt.Sprintf("%s-%s", pname, id))
	p.Model().Meta = model.Meta{
		Name: pname,
		Id:   id,
	}
	pc.Playlists = append(pc.Playlists, p)
	return p
}

func (pc *PlaylistController) Remove(index int) controller.IPlaylist {
	lg.Infof("Try to remove playlist.index=%d", index)
	if index < 0 || index >= len(pc.Playlists) {
		lg.Warnf("playlist.index=%d not found", index)
		return nil
	}
	if index == pc.DefaultIndex {
		lg.Info("Delete current system playlist, reset system playlist to index = 0")
		_ = pc.SetDefault(0)
	}
	if index < pc.DefaultIndex {
		lg.Debugf("Delete playlist before system playlist (index=%d), reduce system playlist index by 1", pc.DefaultIndex)
		pc.DefaultIndex = pc.DefaultIndex - 1
	}
	pl := pc.Playlists[index]
	pc.Playlists = append(pc.Playlists[:index], pc.Playlists[index+1:]...)
	return pl
}

func (pc *PlaylistController) SetDefault(index int) error {
	lg.Infof("try set system playlist to playlist.id=%d", index)
	if index < 0 || index >= len(pc.Playlists) {
		lg.Warn("playlist.index=%d not found", index)
		return errors.New("playlist.index not found")
	}
	err := pc.provider.PreparePlaylist(pc.Playlists[index])
	if err != nil {
		return err
	}
	pl := pc.Playlists[index].Model().Copy()
	pc.DefaultIndex = index
	controller.ApplyUser(pl.Medias, controller.PlaylistUser)
	pc.Default.Replace(pl.Medias)
	pc.Default.Model().Name = pl.Name
	return nil
}

func (pc *PlaylistController) PreparePlaylistByIndex(index int) error {
	lg.Infof("try prepare playlist.id=%d", index)
	if index < 0 || index >= len(pc.Playlists) {
		lg.Warn("playlist.id=%d not found", index)
		return nil
	}
	return pc.provider.PreparePlaylist(pc.Playlists[index])
}

type corePlaylist struct {
	model.Playlist
	Index        int
	Lock         sync.RWMutex
	eventManager *event.Manager
}

func NewPlaylist(name string) controller.IPlaylist {
	return &corePlaylist{
		Index: 0,
		Playlist: model.Playlist{
			Name:   name,
			Medias: make([]*model.Media, 0),
			Mode:   model.PlaylistModeNormal,
			Meta:   model.Meta{},
		},
		eventManager: event.MainManager.NewChildManager(),
	}
}

func (p *corePlaylist) Model() *model.Playlist {
	return &p.Playlist
}

func (p *corePlaylist) EventManager() *event.Manager {
	return p.eventManager
}

func (p *corePlaylist) Name() string {
	return p.Playlist.Name
}

func (p *corePlaylist) Size() int {
	return p.Playlist.Size()
}

func (p *corePlaylist) Get(index int) *model.Media {
	if index < 0 || index >= p.Playlist.Size() {
		return nil
	}
	return p.Playlist.Medias[index]
}

func (p *corePlaylist) Pop() *model.Media {
	lg.Info("[Playlists] %s pop first media", p.Playlist)
	if p.Size() == 0 {
		return nil
	}
	p.Lock.Lock()
	index := 0
	if p.Mode == model.PlaylistModeRandom {
		index = rand.Intn(p.Size())
	}
	m := p.Medias[index]
	for i := index; i > 0; i-- {
		p.Medias[i] = p.Medias[i-1]
	}
	p.Medias = p.Medias[1:]
	p.Lock.Unlock()
	if m == nil {
		lg.Warn("[Playlists] pop first media failed, no media left in the playlist")
		return nil
	}
	p.eventManager.CallA(
		model.EventPlaylistUpdate,
		model.PlaylistUpdateEvent{Playlist: p.Playlist.Copy()},
	)
	return m
}

func (p *corePlaylist) Replace(medias []*model.Media) {
	lg.Infof("[Playlists] %s replace all media", &p.Playlist)
	p.Lock.Lock()
	p.Playlist.Medias = medias
	p.Index = 0
	p.Lock.Unlock()
	p.eventManager.CallA(
		model.EventPlaylistUpdate,
		model.PlaylistUpdateEvent{Playlist: p.Playlist.Copy()},
	)
}

func (p *corePlaylist) Push(media *model.Media) {
	p.Insert(-1, media)
}

func (p *corePlaylist) Insert(index int, media *model.Media) {
	lg.Infof("[Playlists]insert media into new index %d at %s ", index, p.Playlist)
	lg.Debugf("media=%s %v", media.Title, media.Meta)
	e := event.Event{
		Id:        model.EventPlaylistPreInsert,
		Cancelled: false,
		Data: model.PlaylistInsertEvent{
			Playlist: p.Playlist.Copy(),
			Index:    index,
			Media:    media,
		},
	}
	p.eventManager.Call(&e)
	if e.Cancelled {
		lg.Info("[Playlists] media insertion has been cancelled by handler")
		return
	}
	p.Lock.Lock()
	if index > p.Size() {
		index = p.Size()
	}
	if index < 0 {
		index = p.Size() + index + 1
	}
	p.Medias = append(p.Medias, nil)
	for i := p.Size() - 1; i > index; i-- {
		p.Medias[i] = p.Medias[i-1]
	}
	p.Medias[index] = media
	p.Lock.Unlock()
	p.eventManager.CallA(
		model.EventPlaylistUpdate,
		model.PlaylistUpdateEvent{Playlist: p.Playlist.Copy()},
	)
	p.eventManager.CallA(
		model.EventPlaylistInsert,
		model.PlaylistInsertEvent{
			Playlist: p.Playlist.Copy(),
			Index:    index,
			Media:    media,
		},
	)
}

func (p *corePlaylist) Delete(index int) *model.Media {
	lg.Infof("from media at index %d from %s", index, p.Playlist)
	if index >= p.Size() || index < 0 {
		p.Lock.Unlock()
		return nil
	}
	m := p.Medias[index]
	p.Lock.Lock()
	// todo: @5 delete optimization
	p.Medias = append(p.Medias[:index], p.Medias[index+1:]...)
	p.Lock.Unlock()
	if m == nil {
		lg.Warnf("media at index %d does not exist", index)
	}
	p.eventManager.CallA(
		model.EventPlaylistUpdate,
		model.PlaylistUpdateEvent{Playlist: p.Playlist.Copy()})
	return m
}

func (p *corePlaylist) Move(src int, dst int) {
	lg.Infof("from media from index %d to %d", src, dst)
	if src >= p.Size() || src < 0 {
		lg.Warnf("media at index %d does not exist", src)
		return
	}
	p.Lock.Lock()
	if dst >= p.Size() {
		dst = p.Size() - 1
	}
	if dst < 0 {
		dst = 0
	}
	if dst == src {
		p.Lock.Unlock()
		return
	}
	step := 1
	if dst < src {
		step = -1
	}
	tmp := p.Medias[src]
	for i := src; i != dst; i += step {
		p.Medias[i] = p.Medias[i+step]
	}
	p.Medias[dst] = tmp
	p.Lock.Unlock()
	p.eventManager.CallA(
		model.EventPlaylistUpdate,
		model.PlaylistUpdateEvent{Playlist: p.Playlist.Copy()})
}

func (p *corePlaylist) Next() *model.Media {
	lg.Infof("[Playlists] %s get next media with random=%t", p, p.Mode == model.PlaylistModeRandom)
	if p.Size() == 0 {
		lg.Warn("[Playlists] get next media failed, no media left in the playlist")
		return nil
	}
	var index int
	index = p.Index
	if p.Mode == model.PlaylistModeRandom {
		p.Index = rand.Intn(p.Size())
	} else {
		p.Index = (p.Index + 1) % p.Size()
	}
	m := p.Medias[index]
	return m
}
