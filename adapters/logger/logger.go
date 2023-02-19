package logger

import (
	"AynaLivePlayer/core/adapter"
)

type LoggerFactory struct {
	LiveClients map[string]adapter.LiveClientCtor
}

func (f *LoggerFactory) NewLogrus(filename string, redirectStderr bool) adapter.ILogger {
	return NewLogrusLogger(filename, redirectStderr)
}
