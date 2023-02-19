package internal

import (
	"AynaLivePlayer/adapters/provider"
	"AynaLivePlayer/common/config"
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type PlaylistController struct {
	PlaylistPath          string
	DefaultIndex          int
	CurrentPlaylistRandom bool
	DefaultPlaylistRandom bool
	History               adapter.IPlaylist   `ini:"-"`
	Current               adapter.IPlaylist   `ini:"-"`
	Default               adapter.IPlaylist   `ini:"-"`
	Playlists             []adapter.IPlaylist `ini:"-"`
	eventManager          *event.Manager
	log                   adapter.ILogger
	provider              adapter.IProviderController
}

func NewPlaylistController(
	em *event.Manager, log adapter.ILogger,
	provider adapter.IProviderController) adapter.IPlaylistController {
	pc := &PlaylistController{
		PlaylistPath:          "playlists.json",
		History:               newPlaylistImpl("history", em.NewChildManager(), log),
		Default:               newPlaylistImpl("default", em.NewChildManager(), log),
		Current:               newPlaylistImpl("current", em.NewChildManager(), log),
		Playlists:             make([]adapter.IPlaylist, 0),
		DefaultIndex:          0,
		CurrentPlaylistRandom: false,
		DefaultPlaylistRandom: true,
		eventManager:          em,
		log:                   log,
		provider:              provider,
	}
	config.LoadConfig(pc)
	if pc.DefaultIndex < 0 || pc.DefaultIndex >= len(pc.Playlists) {
		pc.DefaultIndex = 0
		log.Warn("playlist index did not find")
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
	var playlists = []*model.Playlist{
		{
			Meta: model.Meta{
				Name: "netease",
				Id:   "2382819181",
			},
		},
		{
			Meta: model.Meta{
				Name: "netease",
				Id:   "4987059624",
			},
		},
		{
			Meta: model.Meta{
				Name: "local",
				Id:   "list1",
			},
		},
	}
	_ = config.LoadJson(pc.PlaylistPath, &playlists)
	for _, pl := range playlists {
		pc.Playlists = append(pc.Playlists, &playlistImpl{
			Index:        0,
			Playlist:     *pl,
			eventManager: pc.eventManager.NewChildManager(),
			log:          pc.log,
		})
	}
	if pc.CurrentPlaylistRandom {
		pc.Current.Model().Mode = model.PlaylistModeRandom
	}
	if pc.DefaultPlaylistRandom {
		pc.Default.Model().Mode = model.PlaylistModeRandom
	}
}

func (pc *PlaylistController) OnSave() {
	var playlists = make([]*model.Playlist, 0)
	for _, pl := range pc.Playlists {
		playlists = append(playlists, pl.Model())
	}
	_ = config.SaveJson(pc.PlaylistPath, &playlists)

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

func (pc *PlaylistController) GetHistory() adapter.IPlaylist {
	return pc.History
}

func (pc *PlaylistController) GetDefault() adapter.IPlaylist {
	return pc.Default
}

func (pc *PlaylistController) GetCurrent() adapter.IPlaylist {
	return pc.Current
}

func (pc *PlaylistController) AddToHistory(media *model.Media) {
	pc.log.Debugf("add media %s (%s) to history", media.Title, media.Artist)
	media = media.Copy()
	// reset url for future use
	media.Url = ""
	if pc.History.Size() >= 1024 {
		pc.History.Replace([]*model.Media{})
	}
	media.User = HistoryUser
	pc.History.Push(media)
	return
}

func (pc *PlaylistController) Get(index int) adapter.IPlaylist {
	if index < 0 || index >= len(pc.Playlists) {
		pc.log.Warnf("playlist.index=%d not found", index)
		return nil
	}
	return pc.Playlists[index]
}

func (pc *PlaylistController) Add(pname string, uri string) (adapter.IPlaylist, error) {
	pc.log.Infof("try add playlist %s with provider %s", uri, pname)
	id, err := provider.FormatPlaylistUrl(pname, uri)
	if err != nil || id == "" {
		pc.log.Warnf("[PlaylistController] fail to format %s playlist id for %s: %s", uri, pname, err)
		return nil, errors.New(fmt.Sprintf("fail to format playlist id: %s", err))
	}
	p := newPlaylistImpl("", pc.eventManager.NewChildManager(), pc.log)
	p.Model().Meta = model.Meta{
		Name: pname,
		Id:   id,
	}
	pc.Playlists = append(pc.Playlists, p)
	return p, nil
}

func (pc *PlaylistController) Remove(index int) (adapter.IPlaylist, error) {
	pc.log.Infof("Try to remove playlist.index=%d", index)
	if index < 0 || index >= len(pc.Playlists) {
		pc.log.Warnf("playlist.index=%d not found", index)
		return nil, fmt.Errorf("playlist.index=%d not found", index)
	}
	if index == pc.DefaultIndex {
		pc.log.Info("Delete current system playlist, reset system playlist to index = 0")
		_ = pc.SetDefault(0)
	}
	if index < pc.DefaultIndex {
		pc.log.Debugf("Delete playlist before system playlist (index=%d), reduce system playlist index by 1", pc.DefaultIndex)
		pc.DefaultIndex = pc.DefaultIndex - 1
	}
	pl := pc.Playlists[index]
	pc.Playlists = append(pc.Playlists[:index], pc.Playlists[index+1:]...)
	return pl, nil
}

func (pc *PlaylistController) SetDefault(index int) error {
	pc.log.Infof("try set system playlist to playlist.id=%d", index)
	if index < 0 || index >= len(pc.Playlists) {
		pc.log.Warn("playlist.index=%d not found", index)
		return errors.New("playlist.index not found")
	}
	err := pc.provider.PreparePlaylist(pc.Playlists[index])
	if err != nil {
		return err
	}
	pl := pc.Playlists[index].Model().Copy()
	pc.DefaultIndex = index
	model.ApplyUser(pl.Medias, PlaylistUser)
	pc.Default.Replace(pl.Medias)
	pc.Default.Model().Title = pl.Title
	pc.Default.Model().Meta = pl.Meta
	return nil
}

func (pc *PlaylistController) PreparePlaylistByIndex(index int) error {
	pc.log.Infof("try prepare playlist.id=%d", index)
	if index < 0 || index >= len(pc.Playlists) {
		pc.log.Warn("playlist.id=%d not found", index)
		return nil
	}
	return pc.provider.PreparePlaylist(pc.Playlists[index])
}

type playlistImpl struct {
	model.Playlist
	Index        int
	Lock         sync.RWMutex
	eventManager *event.Manager
	log          adapter.ILogger
}

func newPlaylistImpl(title string, em *event.Manager, log adapter.ILogger) adapter.IPlaylist {
	return &playlistImpl{
		Index: 0,
		Playlist: model.Playlist{
			Title:  title,
			Medias: make([]*model.Media, 0),
			Mode:   model.PlaylistModeNormal,
			Meta:   model.Meta{},
		},
		eventManager: em,
		log:          log,
	}
}

func (p *playlistImpl) Model() *model.Playlist {
	return &p.Playlist
}

func (p *playlistImpl) EventManager() *event.Manager {
	return p.eventManager
}

func (p *playlistImpl) Identifier() string {
	return p.Playlist.Meta.Identifier()
}

func (p *playlistImpl) DisplayName() string {
	if p.Playlist.Title != "" {
		return p.Playlist.Title
	}
	return p.Playlist.Meta.Identifier()
}

func (p *playlistImpl) Size() int {
	return p.Playlist.Size()
}

func (p *playlistImpl) Get(index int) *model.Media {
	if index < 0 || index >= p.Playlist.Size() {
		return nil
	}
	return p.Playlist.Medias[index]
}

func (p *playlistImpl) Pop() *model.Media {
	p.log.Debugf("[Playlists] %s pop first media", p.Playlist)
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
		p.log.Warn("[Playlists] pop first media failed, no media left in the playlist")
		return nil
	}
	p.eventManager.CallA(
		events.EventPlaylistUpdate,
		events.PlaylistUpdateEvent{
			Playlist: p.Playlist.Copy(),
		},
	)
	return m
}

func (p *playlistImpl) Replace(medias []*model.Media) {
	p.log.Infof("[Playlists] %s replace all media", &p.Playlist)
	p.Lock.Lock()
	p.Playlist.Medias = medias
	p.Index = 0
	p.Lock.Unlock()
	p.eventManager.CallA(
		events.EventPlaylistUpdate,
		events.PlaylistUpdateEvent{
			Playlist: p.Playlist.Copy(),
		},
	)
}

func (p *playlistImpl) Push(media *model.Media) {
	p.Insert(-1, media)
}

func (p *playlistImpl) Insert(index int, media *model.Media) {
	p.log.Infof("[Playlists] insert media into new index %d at %s ", index, p.Playlist)
	p.log.Debugf("media=%s %v", media.Title, media.Meta)
	e := event.Event{
		Id:        events.EventPlaylistPreInsert,
		Cancelled: false,
		Data: events.PlaylistInsertEvent{
			Playlist: p.Playlist.Copy(),
			Index:    index,
			Media:    media,
		},
	}
	p.eventManager.Call(&e)
	if e.Cancelled {
		p.log.Info("[Playlists] media insertion has been cancelled by handler")
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
		events.EventPlaylistUpdate,
		events.PlaylistUpdateEvent{
			Playlist: p.Playlist.Copy(),
		},
	)
	p.eventManager.CallA(
		events.EventPlaylistInsert,
		events.PlaylistInsertEvent{
			Playlist: p.Playlist.Copy(),
			Index:    index,
			Media:    media,
		},
	)
}

func (p *playlistImpl) Delete(index int) *model.Media {
	p.log.Infof("from media at index %d from %s", index, p.Playlist)
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
		p.log.Warnf("[Playlists] media at index %d does not exist", index)
	}
	p.eventManager.CallA(
		events.EventPlaylistUpdate,
		events.PlaylistUpdateEvent{
			Playlist: p.Playlist.Copy(),
		})
	return m
}

func (p *playlistImpl) Move(src int, dst int) {
	p.log.Infof("[Playlists] from media from index %d to %d", src, dst)
	if src >= p.Size() || src < 0 {
		p.log.Warnf("[Playlists] media at index %d does not exist", src)
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
		events.EventPlaylistUpdate,
		events.PlaylistUpdateEvent{
			Playlist: p.Playlist.Copy(),
		})
}

func (p *playlistImpl) Next() *model.Media {
	p.log.Infof("[Playlists] %s get next media with random=%t", p, p.Mode == model.PlaylistModeRandom)
	if p.Size() == 0 {
		p.log.Warn("[Playlists] get next media failed, no media left in the playlist")
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
