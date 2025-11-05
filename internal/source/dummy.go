package source

import (
	"github.com/AynaLivePlayer/miaosic"
)

// dummySource is placeholder source for bypassing copyright requirement
type dummySource struct{}

func (d *dummySource) GetName() string {
	return "dummy"
}

func (d *dummySource) Qualities() []miaosic.Quality {
	return []miaosic.Quality{}
}

func (d *dummySource) Search(keyword string, page, size int) ([]miaosic.MediaInfo, error) {
	return []miaosic.MediaInfo{
		miaosic.MediaInfo{
			Title:  keyword,
			Artist: "Unknown",
			Album:  "Unknown",
			Meta: miaosic.MetaData{
				Provider:   "dummy",
				Identifier: keyword,
			},
		},
	}, nil
}

func (d *dummySource) MatchMedia(uri string) (miaosic.MetaData, bool) {
	return miaosic.MetaData{}, false
}

func (d *dummySource) GetMediaInfo(meta miaosic.MetaData) (miaosic.MediaInfo, error) {
	return miaosic.MediaInfo{
		Title:  meta.Identifier,
		Artist: "Unknown",
		Album:  "Unknown",
		Meta: miaosic.MetaData{
			Provider:   "dummy",
			Identifier: meta.Identifier,
		},
	}, nil
}

func (d *dummySource) GetMediaUrl(meta miaosic.MetaData, quality miaosic.Quality) ([]miaosic.MediaUrl, error) {
	return []miaosic.MediaUrl{}, miaosic.ErrNotImplemented
}

func (d *dummySource) GetMediaLyric(meta miaosic.MetaData) ([]miaosic.Lyrics, error) {
	return []miaosic.Lyrics{}, miaosic.ErrNotImplemented
}

func (d *dummySource) MatchPlaylist(uri string) (miaosic.MetaData, bool) {
	return miaosic.MetaData{}, false
}

func (d *dummySource) GetPlaylist(meta miaosic.MetaData) (*miaosic.Playlist, error) {
	return &miaosic.Playlist{}, miaosic.ErrNotImplemented
}
