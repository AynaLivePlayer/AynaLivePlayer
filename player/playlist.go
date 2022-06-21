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
	lock     sync.RWMutex
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
	p.lock.Lock()
	media := p.Playlist[0]
	p.Playlist = p.Playlist[1:]
	p.lock.Unlock()
	defer p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	return media
}

func (p *Playlist) Replace(medias []*Media) {
	p.lock.Lock()
	p.Playlist = medias
	p.Index = 0
	p.lock.Unlock()
	p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	return
}

func (p *Playlist) Push(media *Media) {
	p.Insert(-1, media)
	defer p.Handler.CallA(EventPlaylistUpdate, PlaylistUpdateEvent{Playlist: p})
	return
}

// Insert runtime in O(n) but i don't care
func (p *Playlist) Insert(index int, media *Media) {
	p.l().Infof("insert new meida to index %d", index)
	p.l().Debug("media=", *media)
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
	p.lock.Lock()
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
	p.lock.Unlock()
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
