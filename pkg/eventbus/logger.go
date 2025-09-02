package eventbus

import "log"

type Logger interface {
	Printf(string, ...interface{})
}

type loggerImpl struct{}

func (l loggerImpl) Printf(s string, i ...interface{}) {
	log.Printf(s, i...)
}

// Log replace with your own logger if needed
var Log Logger = &loggerImpl{}
