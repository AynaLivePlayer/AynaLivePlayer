package events

import (
	"AynaLivePlayer/core/model"
)

const SearchCmd = "cmd.search"

type SearchCmdEvent struct {
	Keyword  string
	Provider string
}

const SearchResultUpdate = "update.search_result"

type SearchResultUpdateEvent struct {
	Medias []model.Media
}
