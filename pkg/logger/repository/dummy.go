package repository

import (
	"AynaLivePlayer/pkg/logger"
	"fmt"
)

type DummyLogger struct {
}

func (l *DummyLogger) DebugW(message string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) DebugS(message string, fields logger.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) InfoW(message string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) InfoS(message string, fields logger.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) WarnW(message string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) WarnS(message string, fields logger.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) ErrorW(message string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) ErrorS(message string, fields logger.LogField) {
	//TODO implement me
	panic("implement me")
}

func (l *DummyLogger) Debug(args ...interface{}) {
	fmt.Println(args...)
}

func (l *DummyLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *DummyLogger) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (l *DummyLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *DummyLogger) Warn(args ...interface{}) {
	fmt.Println(args...)
}

func (l *DummyLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *DummyLogger) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (l *DummyLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (l *DummyLogger) WithPrefix(prefix string) logger.ILogger {
	return l
}

func (l *DummyLogger) SetLogLevel(level logger.LogLevel) {

}
