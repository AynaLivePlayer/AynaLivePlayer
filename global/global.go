package global

import (
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/logger"
)

var Logger logger.ILogger = nil

var EventBus eventbus.Bus = nil
