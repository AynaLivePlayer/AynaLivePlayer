package playlist

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"math/rand"
	"sync"
	"time"
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().Unix()))
}

type playlist struct {
	Index       int
	playlistId  model.PlaylistID
	mode        model.PlaylistMode
	Medias      []model.Media
	randomIndex []int
	Lock        sync.RWMutex
}

// resetRandomIndex reset the random index
// this function is not locked, should be called in a locked context
func (p *playlist) resetRandomIndex() {
	p.randomIndex = make([]int, p.Size())
	for i := 0; i < p.Size(); i++ {
		p.randomIndex[i] = i
	}
	rand.Shuffle(p.Size(), func(i, j int) {
		p.randomIndex[i], p.randomIndex[j] = p.randomIndex[j], p.randomIndex[i]
	})
}

func newPlaylist(id model.PlaylistID) *playlist {
	pl := &playlist{
		playlistId: id,
		Medias:     make([]model.Media, 0),
		Lock:       sync.RWMutex{},
		Index:      0,
	}
	global.EventBus.Subscribe("", events.PlaylistMoveCmd(id), "internal.playlist.move", func(event *eventbus.Event) {
		e := event.Data.(events.PlaylistMoveCmdEvent)
		pl.Move(e.From, e.To)
	})
	global.EventBus.Subscribe("", events.PlaylistInsertCmd(id), "internal.playlist.insert", func(event *eventbus.Event) {
		e := event.Data.(events.PlaylistInsertCmdEvent)
		pl.Insert(e.Position, e.Media)
	})
	global.EventBus.Subscribe("", events.PlaylistDeleteCmd(id), "internal.playlist.delete", func(event *eventbus.Event) {
		e := event.Data.(events.PlaylistDeleteCmdEvent)
		pl.Delete(e.Index)
	})
	global.EventBus.Subscribe("", events.PlaylistNextCmd(id), "internal.playlist.next", func(event *eventbus.Event) {
		log.Infof("Playlist %s recieve next", id)
		pl.Next(event.Data.(events.PlaylistNextCmdEvent).Remove)
	})
	global.EventBus.Subscribe("", events.PlaylistModeChangeCmd(id), "internal.playlist.mode", func(event *eventbus.Event) {
		pl.Lock.Lock()
		pl.mode = event.Data.(events.PlaylistModeChangeCmdEvent).Mode
		pl.Index = 0
		pl.resetRandomIndex()
		pl.Lock.Unlock()
		log.Infof("Playlist %s mode changed to %d", id, pl.mode)
		_ = global.EventBus.Publish(events.PlaylistModeChangeUpdate(id), events.PlaylistModeChangeUpdateEvent{
			Mode: pl.mode,
		})
	})
	global.EventBus.Subscribe("", events.PlaylistSetIndexCmd(id), "internal.playlist.setindex", func(event *eventbus.Event) {
		index := event.Data.(events.PlaylistSetIndexCmdEvent).Index
		if index >= pl.Size() || index < 0 {
			index = 0
		}
		pl.Index = index
	})
	pl.resetRandomIndex()
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
	p.resetRandomIndex()
	p.Lock.Unlock()
	_ = global.EventBus.Publish(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
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
	p.resetRandomIndex()
	p.Lock.Unlock()
	_ = global.EventBus.Publish(events.PlaylistInsertUpdate(p.playlistId), events.PlaylistInsertUpdateEvent{
		Position: index,
		Media:    media,
	})
	_ = global.EventBus.Publish(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
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
	if p.Index >= p.Size() {
		p.Index = 0
	}
	p.resetRandomIndex()
	p.Lock.Unlock()
	_ = global.EventBus.Publish(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
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
	_ = global.EventBus.Publish(events.PlaylistDetailUpdate(p.playlistId), events.PlaylistDetailUpdateEvent{
		Medias: p.CopyMedia(),
	})
}

func (p *playlist) Next(delete bool) {
	p.Lock.Lock()
	if p.Size() == 0 {
		// no media in the playlist
		// do not dispatch any event
		p.Lock.Unlock()
		return
	}
	var index int
	index = p.Index
	// add guard, i don't know why this is needed
	// but it seems to fix a bug
	if index >= p.Size() {
		index = 0
	}
	if (index == 0) && p.mode == model.PlaylistModeRandom {
		p.resetRandomIndex()
	}
	var m model.Media
	if p.mode == model.PlaylistModeRandom {
		m = p.Medias[p.randomIndex[index]]
	} else {
		m = p.Medias[index]
	}
	//// fix race condition
	//currentSize := p.Size() - 1
	//if delete {
	//	if p.mode == model.PlaylistModeRandom {
	//		if currentSize == 0 {
	//			p.Index = 0
	//		} else {
	//			p.Index = rand.Intn(currentSize)
	//		}
	//	} else if p.mode == model.PlaylistModeNormal {
	//		p.Index = index
	//	} else {
	//		p.Index = index
	//	}
	//}
	p.Lock.Unlock()
	_ = global.EventBus.Publish(events.PlaylistNextUpdate(p.playlistId), events.PlaylistNextUpdateEvent{
		Media: m,
	})
	if delete {
		if p.mode == model.PlaylistModeRandom {
			p.Delete(p.randomIndex[index])
		} else {
			p.Delete(index)
		}
	} else {
		p.Index = (p.Index + 1) % p.Size()
	}
}
