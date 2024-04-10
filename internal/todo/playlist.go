package todo

import (
	"AynaLivePlayer/adapters/provider"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"errors"
	"fmt"
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
