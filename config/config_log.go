package config

import "github.com/sirupsen/logrus"

type _LogConfig struct {
	Path           string
	Level          logrus.Level
	RedirectStderr bool
}

func (c *_LogConfig) OnLoad() {
}

func (c *_LogConfig) OnSave() {
}

func (c *_LogConfig) Name() string {
	return "Log"
}

var Log = &_LogConfig{
	Path:           "./log.txt",
	Level:          logrus.InfoLevel,
	RedirectStderr: false, // this should be true if it is in production mode.
}
