package repository

import (
	"AynaLivePlayer/pkg/logger"
	"errors"
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

func (l *LogrusLogger) SetLogLevel(level logger.LogLevel) {
	switch level {
	case logger.LogLevelDebug:
		l.Logger.SetLevel(logrus.DebugLevel)
	case logger.LogLevelInfo:
		l.Logger.SetLevel(logrus.InfoLevel)
	case logger.LogLevelWarn:
		l.Logger.SetLevel(logrus.WarnLevel)
	case logger.LogLevelError:
		l.Logger.SetLevel(logrus.ErrorLevel)
	default:
		l.Logger.SetLevel(logrus.InfoLevel)
	}
}

func getLogOut(filename string, maxSize int64) (*os.File, error) {
	if filename == "" {
		return nil, errors.New("failed to log to file, using default stdout")
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.New("failed to log to file, using default stdout")
	}
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, errors.New("failed to get log file stat, using default stdout")
	}
	if fileinfo.Size() > maxSize {
		_ = file.Truncate(0)
	}
	return file, nil
}

func NewLogrusLogger(fileName string, maxSize int64, redirectStderr bool) *LogrusLogger {
	l := logrus.New()
	l.SetFormatter(
		&nested.Formatter{
			FieldsOrder: []string{"Module"},
			HideKeys:    true,
			NoColors:    false,
		})
	var file *os.File
	var err error
	if fileName != "" {
		file, err = getLogOut(fileName, maxSize)
		if err == nil {
			l.Out = io.MultiWriter(file, os.Stdout)
		} else {
			l.Warnf(err.Error())
		}
	}
	if redirectStderr && file != nil {
		l.Info("panic/stderr redirect to log file")
		if _, err = paniclog.RedirectStderr(file); err != nil {
			l.Warnf("Failed to redirect stderr to to file: %s", err)
		}
	}
	return &LogrusLogger{
		Entry: logrus.NewEntry(l),
	}
}
