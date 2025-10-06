package events

import "github.com/AynaLivePlayer/miaosic"

const CmdMiaosicGetMediaInfo = "cmd.miaosic.getMediaInfo"

type CmdMiaosicGetMediaInfoData struct {
	Meta miaosic.MetaData `json:"meta"`
}

const ReplyMiaosicGetMediaInfo = "reply.miaosic.getMediaInfo"

type ReplyMiaosicGetMediaInfoData struct {
	Info  miaosic.MediaInfo `json:"info"`
	Error error
}

const CmdMiaosicGetMediaUrl = "cmd.miaosic.getMediaUrl"

type CmdMiaosicGetMediaUrlData struct {
	Meta    miaosic.MetaData `json:"meta"`
	Quality miaosic.Quality  `json:"quality"`
}

const ReplyMiaosicGetMediaUrl = "reply.miaosic.getMediaUrl"

type ReplyMiaosicGetMediaUrlData struct {
	Urls  []miaosic.MediaUrl `json:"urls"`
	Error error
}
