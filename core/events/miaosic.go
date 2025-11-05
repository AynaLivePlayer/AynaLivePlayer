package events

import "github.com/AynaLivePlayer/miaosic"

const CmdMiaosicGetMediaInfo = "cmd.miaosic.getMediaInfo"

type CmdMiaosicGetMediaInfoData struct {
	Meta miaosic.MetaData `json:"meta"`
}

const ReplyMiaosicGetMediaInfo = "reply.miaosic.getMediaInfo"

type ReplyMiaosicGetMediaInfoData struct {
	Info  miaosic.MediaInfo `json:"info"`
	Error error             `json:"error"`
}

const CmdMiaosicGetMediaUrl = "cmd.miaosic.getMediaUrl"

type CmdMiaosicGetMediaUrlData struct {
	Meta    miaosic.MetaData `json:"meta"`
	Quality miaosic.Quality  `json:"quality"`
}

const ReplyMiaosicGetMediaUrl = "reply.miaosic.getMediaUrl"

type ReplyMiaosicGetMediaUrlData struct {
	Urls  []miaosic.MediaUrl `json:"urls"`
	Error error              `json:"error"`
}

const CmdMiaosicQrLogin = "cmd.miaosic.qrLogin"

type CmdMiaosicQrLoginData struct {
	Provider string `json:"provider"`
}

const ReplyMiaosicQrLogin = "reply.miaosic.qrLogin"

type ReplyMiaosicQrLoginData struct {
	Session miaosic.QrLoginSession `json:"session"`
	Error   error                  `json:"error"`
}

const CmdMiaosicQrLoginVerify = "cmd.miaosic.qrLoginVerify"

type CmdMiaosicQrLoginVerifyData struct {
	Session miaosic.QrLoginSession `json:"session"`
}

const ReplyMiaosicQrLoginVerify = "reply.miaosic.qrLoginVerify"

type ReplyMiaosicQrLoginVerifyData struct {
	Result miaosic.QrLoginResult `json:"result"`
	Error  error                 `json:"error"`
}
