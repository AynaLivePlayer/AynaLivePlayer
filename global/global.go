package global

import (
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/logger"
)

var Logger logger.ILogger = nil

var EventManager *event.Manager = nil
