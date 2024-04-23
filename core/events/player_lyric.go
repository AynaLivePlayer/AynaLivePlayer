package events

import "github.com/AynaLivePlayer/miaosic"

const PlayerLyricRequestCmd = "cmd.player.lyric.request"

type PlayerLyricRequestCmdEvent struct {
}

const PlayerLyricReload = "update.player.lyric.reload"

type PlayerLyricReloadEvent struct {
	Lyrics miaosic.Lyrics
}

const PlayerLyricPosUpdate = "update.player.lyric.pos"

type PlayerLyricPosUpdateEvent struct {
	Time         float64
	CurrentIndex int // -1 means no lyric
	CurrentLine  miaosic.LyricLine
	Total        int // total lyric count
}
