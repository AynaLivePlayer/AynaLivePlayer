package config

import "github.com/sirupsen/logrus"

type _LogConfig struct {
	Path  string
	Level logrus.Level
}

func (c *_LogConfig) Name() string {
	return "Log"
}

var Log = &_LogConfig{
	Path:  "./log.txt",
	Level: logrus.InfoLevel,
}
