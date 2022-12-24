package provider

import (
	"AynaLivePlayer/model"
	"os"
	"sort"
	"strings"
)

type _LocalPlaylist struct {
	Name   string
	Medias []*model.Media
}

type Local struct {
	localDir  string
	Playlists []*_LocalPlaylist
}

func NewLocalCtor(config MediaProviderConfig) MediaProvider {
	localDir, ok := config["local_dir"]
	if !ok {
		localDir = "./local"
	}
	l := &Local{Playlists: make([]*_LocalPlaylist, 0), localDir: localDir}
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return l
	}
	for _, n := range getPlaylistNames(localDir) {
		l.Playlists = append(l.Playlists, &_LocalPlaylist{Name: n})
	}
	for i, _ := range l.Playlists {
		_ = readLocalPlaylist(localDir, l.Playlists[i])
	}
	LocalAPI = l
	Providers[LocalAPI.GetName()] = LocalAPI
	return l
}

var LocalAPI *Local

func NewLocal(localdir string) *Local {
	l := &Local{Playlists: make([]*_LocalPlaylist, 0), localDir: localdir}
	if err := os.MkdirAll(localdir, 0755); err != nil {
		return l
	}
	for _, n := range getPlaylistNames(localdir) {
		l.Playlists = append(l.Playlists, &_LocalPlaylist{Name: n})
	}
	for i, _ := range l.Playlists {
		_ = readLocalPlaylist(localdir, l.Playlists[i])
	}
	LocalAPI = l
	Providers[LocalAPI.GetName()] = LocalAPI
	return l
}

func (l *Local) GetName() string {
	return "local"
}

func (l *Local) MatchMedia(keyword string) *model.Media {
	return nil
}

func (l *Local) UpdateMediaLyric(media *model.Media) error {
	// already update in UpdateMedia, do nothing
	return nil
}

func (l *Local) FormatPlaylistUrl(uri string) string {
	return uri
}

func (l *Local) GetPlaylist(playlist *model.Meta) ([]*model.Media, error) {
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
	if readLocalPlaylist(l.localDir, pl) != nil {
		return nil, ErrorExternalApi
	}
	return pl.Medias, nil
}

func (l *Local) Search(keyword string) ([]*model.Media, error) {
	result := make([]struct {
		M *model.Media
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
					M *model.Media
					N int
				}{M: m, N: n})
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].N > result[j].N
	})
	medias := make([]*model.Media, len(result))
	for i, r := range result {
		medias[i] = r.M.Copy()
	}
	return medias, nil
}

func (l *Local) UpdateMedia(media *model.Media) error {
	mediaPath := media.Meta.(model.Meta).Id
	_, err := os.Stat(mediaPath)
	if err != nil {
		return err
	}
	return readMediaFile(media)
}

func (l *Local) UpdateMediaUrl(media *model.Media) error {
	mediaPath := media.Meta.(model.Meta).Id
	_, err := os.Stat(mediaPath)
	if err != nil {
		return err
	}
	media.Url = mediaPath
	return nil
}
