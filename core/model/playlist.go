package model

import "fmt"

type PlaylistMode int

const (
	PlaylistModeNormal PlaylistMode = iota
	PlaylistModeRandom
	PlaylistModeRepeat
)

type Playlist struct {
	Title  string // can be same, display name
	Medias []*Media
	Mode   PlaylistMode
	Meta   Meta
}

func (p Playlist) String() string {
	return fmt.Sprintf("<Playlist %s len:%d>", p.Title, len(p.Medias))
}

func (p *Playlist) Size() int {
	return len(p.Medias)
}

func (p *Playlist) Copy() *Playlist {
	medias := make([]*Media, len(p.Medias))
	copy(medias, p.Medias)
	return &Playlist{
		Title:  p.Title,
		Medias: medias,
		Mode:   p.Mode,
		Meta:   p.Meta,
	}
}
