package logger

import (
	"AynaLivePlayer/core/adapter"
)

type EmptyLogger struct {
}

func (e EmptyLogger) Debug(args ...interface{}) {
	return
}

func (e EmptyLogger) Debugf(format string, args ...interface{}) {
	return
}

func (e EmptyLogger) Info(args ...interface{}) {
	return
}

func (e EmptyLogger) Infof(format string, args ...interface{}) {
	return
}

func (e EmptyLogger) Warn(args ...interface{}) {
	return
}

func (e EmptyLogger) Warnf(format string, args ...interface{}) {
	return
}

func (e EmptyLogger) Error(args ...interface{}) {
	return
}

func (e EmptyLogger) Errorf(format string, args ...interface{}) {
	return
}

func (e EmptyLogger) WithModule(prefix string) adapter.ILogger {
	return &EmptyLogger{}
}

func (e EmptyLogger) SetLogLevel(level adapter.LogLevel) {
	return
}
