package logger

import (
	"AynaLivePlayer/core/adapter"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/virtuald/go-paniclog"
	"io"
	"os"
)

type LogrusLogger struct {
	*logrus.Entry
	module string
}

func (l *LogrusLogger) SetLogLevel(level adapter.LogLevel) {
	switch level {
	case adapter.LogLevelDebug:
		l.Logger.SetLevel(logrus.DebugLevel)
	case adapter.LogLevelInfo:
		l.Logger.SetLevel(logrus.InfoLevel)
	case adapter.LogLevelWarn:
		l.Logger.SetLevel(logrus.WarnLevel)
	case adapter.LogLevelError:
		l.Logger.SetLevel(logrus.ErrorLevel)
	default:
		l.Logger.SetLevel(logrus.InfoLevel)
	}
}

func NewLogrusLogger(fileName string, redirectStderr bool, maxSize int64) *LogrusLogger {
	l := logrus.New()
	l.SetFormatter(
		&nested.Formatter{
			FieldsOrder: []string{"Module"},
			HideKeys:    true,
			NoColors:    true,
		})
	var file *os.File
	var err error
	if fileName != "" {
		fi, err := os.Stat(fileName)
		if err != nil {
			file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
		} else if fi.Size() > maxSize*1024*1024 {
			file, err = os.OpenFile(fileName, os.O_TRUNC|os.O_WRONLY, 0666)
		} else {
			file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
		}
		if err == nil {
			l.Out = io.MultiWriter(file, os.Stdout)
		} else {
			l.Info("Failed to log to file, using default stdout")
		}
	}
	if redirectStderr && file != nil {
		l.Info("panic/stderr redirect to log file")
		if _, err = paniclog.RedirectStderr(file); err != nil {
			l.Infof("Failed to redirect stderr to to file: %s", err)
		}
	}
	return &LogrusLogger{
		Entry: logrus.NewEntry(l),
	}
}

func (l *LogrusLogger) WithModule(prefix string) adapter.ILogger {
	return &LogrusLogger{
		Entry:  l.Entry.WithField("Module", prefix),
		module: prefix,
	}
}
