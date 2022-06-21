package logger

import (
	"AynaLivePlayer/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetLevel(config.Log.Level)
	file, err := os.OpenFile(config.Log.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Logger.Out = io.MultiWriter(file, os.Stdout)
	} else {
		Logger.Info("Failed to log to file, using default stdout")
	}
	Logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"Module"},
		HideKeys:    true,
		NoColors:    true,
	})
}
