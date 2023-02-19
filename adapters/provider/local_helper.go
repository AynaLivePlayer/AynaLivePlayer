package provider

import (
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/core/model"
	"github.com/dhowden/tag"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getPlaylistNames(localdir string) []string {
	names := make([]string, 0)
	items, _ := ioutil.ReadDir(localdir)
	for _, item := range items {
		if item.IsDir() {
			names = append(names, item.Name())
		}
	}
	return names
}

// readLocalPlaylist read files under a directory
// and return a _LocalPlaylist object.
// This function assume this directory exists
func readLocalPlaylist(localdir string, playlist *_LocalPlaylist) error {
	p1th := playlist.Name
	playlist.Medias = make([]*model.Media, 0)
	fullPath := filepath.Join(localdir, p1th)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return err
	}
	items, _ := ioutil.ReadDir(fullPath)
	for _, item := range items {
		// if item is a file, read file
		if !item.IsDir() {
			fn := item.Name()
			media := model.Media{
				Meta: model.Meta{
					Name: LocalAPI.GetName(),
					Id:   filepath.Join(fullPath, fn),
				},
			}
			if readMediaFile(&media) != nil {
				continue
			}
			playlist.Medias = append(playlist.Medias, &media)
		}
	}
	return nil
}

func readMediaFile(media *model.Media) error {
	p := media.Meta.(model.Meta).Id
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	meta, err := tag.ReadFrom(f)
	if err != nil {
		return err
	}
	media.Title = util.GetOrDefault(meta.Title(), filepath.Base(p))
	media.Artist = util.GetOrDefault(meta.Artist(), "Unknown")
	media.Album = util.GetOrDefault(meta.Album(), "Unknown")
	media.Lyric = meta.Lyrics()
	if meta.Picture() != nil {
		media.Cover.Data = meta.Picture().Data
	}
	return nil
}
