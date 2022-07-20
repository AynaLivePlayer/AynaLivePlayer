package logger

import (
	"AynaLivePlayer/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/virtuald/go-paniclog"
	"io"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetLevel(config.Log.Level)
	Logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"Module"},
		HideKeys:    true,
		NoColors:    true,
	})
	file, err := os.OpenFile(config.Log.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = io.MultiWriter(file, os.Stdout)
	} else {
		Logger.Info("Failed to log to file, using default stdout")
	}
	if config.Log.RedirectStderr {
		Logger.Info("panic/stderr redirect to log file")
		if _, err = paniclog.RedirectStderr(file); err != nil {
			Logger.Infof("Failed to redirect stderr to to file: %s", err)
		}
	}
}
