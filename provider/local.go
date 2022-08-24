package provider

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/player"
	"os"
	"sort"
	"strings"
)

type _LocalPlaylist struct {
	Name   string
	Medias []*player.Media
}

type Local struct {
	Playlists []*_LocalPlaylist
}

var LocalAPI *Local

func init() {
	LocalAPI = _newLocal()
	Providers[LocalAPI.GetName()] = LocalAPI
}

func _newLocal() *Local {
	l := &Local{Playlists: make([]*_LocalPlaylist, 0)}
	if err := os.MkdirAll(config.Provider.LocalDir, 0755); err != nil {
		return l
	}
	for _, n := range getPlaylistNames() {
		l.Playlists = append(l.Playlists, &_LocalPlaylist{Name: n})
	}
	for i, _ := range l.Playlists {
		_ = readLocalPlaylist(l.Playlists[i])
	}
	return l
}

func (l *Local) GetName() string {
	return "local"
}

func (l *Local) MatchMedia(keyword string) *player.Media {
	return nil
}

func (l *Local) UpdateMediaLyric(media *player.Media) error {
	// already update in UpdateMedia, do nothing
	return nil
}

func (l *Local) FormatPlaylistUrl(uri string) string {
	return uri
}

func (l *Local) GetPlaylist(playlist Meta) ([]*player.Media, error) {
	var pl *_LocalPlaylist = nil
	for _, p := range l.Playlists {
		if p.Name == playlist.Id {
			pl = p
		}
	}
	if pl == nil {
		l.Playlists = append(l.Playlists, &_LocalPlaylist{Name: playlist.Id})
		pl = l.Playlists[len(l.Playlists)-1]
	}
	if readLocalPlaylist(pl) != nil {
		return nil, ErrorExternalApi
	}
	return pl.Medias, nil
}

func (l *Local) Search(keyword string) ([]*player.Media, error) {
	result := make([]struct {
		M *player.Media
		N int
	}, 0)
	keywords := strings.Split(keyword, " ")
	for _, p := range l.Playlists {
		for _, m := range p.Medias {
			title := strings.ToLower(m.Title)
			artist := strings.ToLower(m.Artist)
			n := 0
			for _, k := range keywords {
				kw := strings.ToLower(k)
				if strings.Contains(title, kw) || strings.Contains(artist, kw) {
					n++
				}
				if kw == title {
					n += 3
				}
			}
			if n > 0 {
				result = append(result, struct {
					M *player.Media
					N int
				}{M: m, N: n})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].N > result[j].N
	})
	medias := make([]*player.Media, len(result))
	for i, r := range result {
		medias[i] = r.M.Copy()
	}
	return medias, nil
}

func (l *Local) UpdateMedia(media *player.Media) error {
	mediaPath := media.Meta.(Meta).Id
	_, err := os.Stat(mediaPath)
	if err != nil {
		return err
	}
	return readMediaFile(media)
}

func (l *Local) UpdateMediaUrl(media *player.Media) error {
	mediaPath := media.Meta.(Meta).Id
	_, err := os.Stat(mediaPath)
	if err != nil {
		return err
	}
	media.Url = mediaPath
	return nil
}
