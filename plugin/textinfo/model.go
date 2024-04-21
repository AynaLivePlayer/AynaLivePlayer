package textinfo

import (
	"AynaLivePlayer/core/model"
	"github.com/AynaLivePlayer/miaosic"
)

type Time struct {
	Seconds      int
	Minutes      int
	TotalSeconds int
}

func NewTimeFromSec(sec int) Time {
	return Time{
		Seconds:      sec % 60,
		Minutes:      sec / 60,
		TotalSeconds: sec,
	}
}

type MediaInfo struct {
	Index    int
	Title    string
	Artist   string
	Album    string
	Username string
	Cover    miaosic.Picture
}

func NewMediaInfo(idx int, media model.Media) MediaInfo {
	return MediaInfo{
		Index:    idx,
		Title:    media.Info.Title,
		Artist:   media.Info.Artist,
		Album:    media.Info.Album,
		Username: media.ToUser().Name,
		Cover:    media.Info.Cover,
	}
}

type OutInfo struct {
	// ============== Current ==============
	Current     MediaInfo
	CurrentTime Time
	TotalTime   Time

	Lyric      string
	NextLyrics []string

	// ============== Playlist ==============
	Playlist       []MediaInfo
	PlaylistLength int
}
