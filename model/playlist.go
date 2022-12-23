package model

import "fmt"

type PlaylistMode int

const (
	PlaylistModeNormal PlaylistMode = iota
	PlaylistModeRandom
)

type Playlist struct {
	Name   string
	Medias []*Media
	Mode   PlaylistMode
	Meta   Meta
}

func (p Playlist) String() string {
	return fmt.Sprintf("<Playlist %s>", p.Name)
}

func (p *Playlist) Size() int {
	return len(p.Medias)
}

func (p *Playlist) Copy() *Playlist {
	medias := make([]*Media, len(p.Medias))
	copy(medias, p.Medias)
	return &Playlist{
		Name:   p.Name,
		Medias: medias,
		Mode:   p.Mode,
		Meta:   p.Meta,
	}
}
