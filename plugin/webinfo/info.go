package webinfo

import (
	"AynaLivePlayer/model"
)

type MediaInfo struct {
	Index    int
	Title    string
	Artist   string
	Album    string
	Username string
	Cover    model.Picture
}

type OutInfo struct {
	Current     MediaInfo
	CurrentTime int
	TotalTime   int
	Lyric       string
	Playlist    []MediaInfo
}

const (
	OutInfoC  = "Current"
	OutInfoCT = "CurrentTime"
	OutInfoTT = "TotalTime"
	OutInfoL  = "Lyric"
	OutInfoPL = "Playlists"
)

type WebsocketData struct {
	Update string
	Data   OutInfo
}
