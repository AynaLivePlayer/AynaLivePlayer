package player

import (
	"AynaLivePlayer/event"
	"AynaLivePlayer/logger"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

const MODULE_PLAYLIST = "Player.Playlist"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type PlaylistConfig struct {
	RandomNext bool
}

type Playlist struct {
	Index    int
	Name     string
	Config   PlaylistConfig
	Playlist []*Media
	Handler  *event.Handler
	Meta     interface{}
	Lock     sync.RWMutex
}

func NewPlaylist(name string, config PlaylistConfig) *Playlist {
	return &Playlist{
		Index:    0,
		Name:     name,
		Config:   config,
		Playlist: make([]*Media, 0),
		Handler:  event.NewHandler(),
	}
}

func (p *Playlist) l() *logrus.Entry {
	return logger.Logger.WithFields(logrus.Fields{
		"Module": MODULE_PLAYLIST,
		"Name":   p.Name,
	})
}

func (p *Playlist) Size() int {
	p.l().Tracef("getting size=%d", len(p.Playlist))
	return len(p.Playlist)
}

func (p *Playlist) Pop() *Media {
	p.l().Infof("pop first media")
	if p.Size() == 0 {
		p.l().Warn("pop first media failed, no media left in the playlist")
		return nil
	}
	p.Lock.Lock()
	index := 0
	if p.Config.RandomNext {
		index = rand.Intn(p.Size())
	}
	media := p.Playlist[index]
	for i := index; i > 0; i-- {
		p.Playlist[i] = p.Playlist[i-1]
	}
	p.Playlist = p.Playlist[1:]
	p.Lock.Unlock()
	defer p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	return media
}

func (p *Playlist) Replace(medias []*Media) {
	p.Lock.Lock()
	p.Playlist = medias
	p.Index = 0
	p.Lock.Unlock()
	p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	return
}

func (p *Playlist) Push(media *Media) {
	p.Insert(-1, media)
	return
}

// Insert runtime in O(n) but i don't care
func (p *Playlist) Insert(index int, media *Media) {
	p.l().Infof("insert new meida to index %d", index)
	p.l().Debugf("media= %s", media.Title)
	e := event.Event{
		Id:        EventPlaylistPreInsert,
		Cancelled: false,
		Data: PlaylistInsertEvent{
			Playlist: p,
			Index:    index,
			Media:    media,
		},
	}
	p.Handler.Call(&e)
	if e.Cancelled {
		p.l().Info("insert new media has been cancelled by handler")
		return
	}
	p.Lock.Lock()
	if index > p.Size() {
		index = p.Size()
	}
	if index < 0 {
		index = p.Size() + index + 1
	}
	p.Playlist = append(p.Playlist, nil)
	for i := p.Size() - 1; i > index; i-- {
		p.Playlist[i] = p.Playlist[i-1]
	}
	p.Playlist[index] = media
	p.Lock.Unlock()
	defer func() {
		p.Handler.Call(&event.Event{
			Id:        EventPlaylistInsert,
			Cancelled: false,
			Data: PlaylistInsertEvent{
				Playlist: p,
				Index:    index,
				Media:    media,
			},
		})
		p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	}()
}

func (p *Playlist) Delete(index int) {
	p.l().Infof("from media at index %d", index)
	p.Lock.Lock()
	if index >= p.Size() || index < 0 {
		p.l().Warnf("media at index %d does not exist", index)
		p.Lock.Unlock()
		return
	}
	// todo: @5 delete optimization
	p.Playlist = append(p.Playlist[:index], p.Playlist[index+1:]...)
	p.Lock.Unlock()
	defer p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
}

func (p *Playlist) Move(src int, dest int) {
	p.l().Infof("from media from index %d to %d", src, dest)
	p.Lock.Lock()
	if src >= p.Size() || src < 0 {
		p.l().Warnf("media at index %d does not exist", src)
		p.Lock.Unlock()
		return
	}
	if dest >= p.Size() {
		dest = p.Size() - 1
	}
	if dest < 0 {
		dest = 0
	}
	if dest == src {
		p.l().Warn("src and dest are same, operation not perform")
		p.Lock.Unlock()
		return
	}
	step := 1
	if dest < src {
		step = -1
	}
	tmp := p.Playlist[src]
	for i := src; i != dest; i += step {
		p.Playlist[i] = p.Playlist[i+step]
	}
	p.Playlist[dest] = tmp
	p.Lock.Unlock()
	defer p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
}

func (p *Playlist) Next() *Media {
	p.l().Infof("get next media with random=%t", p.Config.RandomNext)
	if p.Size() == 0 {
		p.l().Info("get next media failed, no media left in the playlist")
		return nil
	}
	var index int
	index = p.Index
	if p.Config.RandomNext {
		p.Index = rand.Intn(p.Size())
	} else {
		p.Index = (p.Index + 1) % p.Size()
	}
	p.l().Tracef("return index %d, new index %d", index, p.Index)
	defer p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	return p.Playlist[index]
}
