package webinfo

import "AynaLivePlayer/player"

type MediaInfo struct {
	Index    int
	Title    string
	Artist   string
	Album    string
	Username string
	Cover    player.Picture
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
	OutInfoPL = "Playlist"
)

type WebsocketData struct {
	Update string
	Data   OutInfo
}
