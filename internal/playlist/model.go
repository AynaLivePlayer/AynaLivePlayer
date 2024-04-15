package playlist

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
	"math/rand"
	"sync"
)

type playlist struct {
	Index      int
	playlistId model.PlaylistID
	mode       model.PlaylistMode
	Medias     []model.Media
	Lock       sync.RWMutex
}

func newPlaylist(id model.PlaylistID) *playlist {
	pl := &playlist{
		playlistId: id,
		Medias:     make([]model.Media, 0),
		Lock:       sync.RWMutex{},
		Index:      0,
	}
	global.EventManager.RegisterA(events.PlaylistMoveCmd(id), "internal.playlist.move", func(event *event.Event) {
		e := event.Data.(events.PlaylistMoveCmdEvent)
		pl.Move(e.From, e.To)
	})
	global.EventManager.RegisterA(events.PlaylistInsertCmd(id), "internal.playlist.insert", func(event *event.Event) {
		e := event.Data.(events.PlaylistInsertCmdEvent)
		pl.Insert(e.Position, e.Media)
	})
	global.EventManager.RegisterA(events.PlaylistDeleteCmd(id), "internal.playlist.delete", func(event *event.Event) {
		e := event.Data.(events.PlaylistDeleteCmdEvent)
		pl.Delete(e.Index)
	})
	global.EventManager.RegisterA(events.PlaylistNextCmd(id), "internal.playlist.next", func(event *event.Event) {
		pl.Next(event.Data.(events.PlaylistNextCmdEvent).Remove)
	})
	global.EventManager.RegisterA(events.PlaylistModeChangeCmd(id), "internal.playlist.mode", func(event *event.Event) {
		pl.mode = event.Data.(events.PlaylistModeChangeCmdEvent).Mode
		log.Infof("Playlist %s mode changed to %d", id, pl.mode)
		global.EventManager.CallA(events.PlaylistModeChangeUpdate(id), events.PlaylistModeChangeUpdateEvent{
			Mode: pl.mode,
		})
	})
	return pl
}

func (p *playlist) CopyMedia() []model.Media {
	medias := make([]model.Media, len(p.Medias))
	copy(medias, p.Medias)
	return medias
}

func (p *playlist) Size() int {
	return len(p.Medias)
}

func (p *playlist) Replace(medias []model.Media) {
	p.Lock.Lock()
	p.Medias = medias
	p.Index = 0
	p.Lock.Unlock()
	global.EventManager.CallA(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
		Medias: p.CopyMedia(),
	})
}

func (p *playlist) Insert(index int, media model.Media) {
	p.Lock.Lock()
	if index > p.Size() {
		index = p.Size()
	}
	if index < 0 {
		index = p.Size() + index + 1
	}
	p.Medias = append(p.Medias, model.Media{})
	for i := p.Size() - 1; i > index; i-- {
		p.Medias[i] = p.Medias[i-1]
	}
	p.Medias[index] = media
	p.Lock.Unlock()
	global.EventManager.CallA(events.PlaylistInsertUpdate(p.playlistId), events.PlaylistInsertUpdateEvent{
		Position: index,
		Media:    media,
	})
	global.EventManager.CallA(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
		Medias: p.CopyMedia(),
	})
}

func (p *playlist) Delete(index int) {
	p.Lock.Lock()
	if index >= p.Size() || index < 0 {
		p.Lock.Unlock()
		return
	}
	// todo: @5 delete optimization
	p.Medias = append(p.Medias[:index], p.Medias[index+1:]...)
	p.Lock.Unlock()
	global.EventManager.CallA(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
		Medias: p.CopyMedia(),
	})
}

func (p *playlist) Move(src int, dst int) {
	if src >= p.Size() || src < 0 {
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
	global.EventManager.CallA(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
		Medias: p.CopyMedia(),
	})
}

func (p *playlist) Next(delete bool) {
	if p.Size() == 0 {
		// no media in the playlist
		// do not issue any event
		return
	}
	var index int
	index = p.Index
	if p.mode == model.PlaylistModeRandom {
		p.Index = rand.Intn(p.Size())
	} else if p.mode == model.PlaylistModeNormal {
		p.Index = (p.Index + 1) % p.Size()
	} else {
		p.Index = index
	}
	m := p.Medias[index]
	global.EventManager.CallA(events.PlaylistNextUpdate(p.playlistId), events.PlaylistNextUpdateEvent{
		Media: m,
	})
	if delete {
		p.Delete(index)
		if p.mode == model.PlaylistModeRandom {
			p.Index = rand.Intn(p.Size())
		} else if p.mode == model.PlaylistModeNormal {
			p.Index = index
		} else {
			p.Index = index
		}
	}
}
