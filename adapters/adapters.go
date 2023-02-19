package adapters

import (
	"AynaLivePlayer/adapters/liveclient"
	"AynaLivePlayer/adapters/logger"
	"AynaLivePlayer/adapters/player"
	"AynaLivePlayer/common/event"
	"AynaLivePlayer/core/adapter"
)

var Logger = &logger.LoggerFactory{}

var LiveClient = &liveclient.LiveClientFactory{
	LiveClients: map[string]adapter.LiveClientCtor{
		"bilibili": liveclient.BilibiliCtor,
	},
	EventManager: event.MainManager,
	Logger:       &logger.EmptyLogger{},
}

var Player = &player.PlayerFactory{
	EventManager: event.MainManager,
	Logger:       &logger.EmptyLogger{},
}
