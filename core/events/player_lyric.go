package events

import "github.com/AynaLivePlayer/miaosic"

const CmdGetCurrentLyric = "cmd.player.lyric.request"

type CmdGetCurrentLyricData struct {
}

const UpdateCurrentLyric = "update.player.lyric.reload"

type UpdateCurrentLyricData struct {
	Lyrics miaosic.Lyrics
}

const PlayerLyricPosUpdate = "update.player.lyric.pos"

type PlayerLyricPosUpdateEvent struct {
	Time         float64
	CurrentIndex int // -1 means no lyric
	CurrentLine  miaosic.LyricLine
	Total        int // total lyric count
}
