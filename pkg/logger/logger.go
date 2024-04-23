package logger

type LogLevel uint32

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

type LogField map[string]interface{}

func (f LogField) Flatten() []interface{} {
	var res []interface{}
	for k, v := range f {
		res = append(res, k, v)
	}
	return res
}

type ILogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	DebugW(message string, keysAndValues ...interface{})
	DebugS(message string, fields LogField)
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	InfoW(message string, keysAndValues ...interface{})
	InfoS(message string, fields LogField)
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	WarnW(message string, keysAndValues ...interface{})
	WarnS(message string, fields LogField)
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorW(message string, keysAndValues ...interface{})
	ErrorS(message string, fields LogField)
	WithPrefix(prefix string) ILogger
	SetLogLevel(level LogLevel)
}

type LogMessage struct {
	Timestamp int64
	Level     LogLevel
	Prefix    string
	Message   string
	Data      map[string]interface{}
}
