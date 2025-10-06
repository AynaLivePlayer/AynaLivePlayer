package events

import (
	"AynaLivePlayer/core/model"
)

const CmdMiaosicSearch = "cmd.search"

type CmdMiaosicSearchData struct {
	Keyword  string
	Provider string
}

const ReplyMiaosicSearch = "update.search_result"

type ReplyMiaosicSearchData struct {
	Medias []model.Media
}
