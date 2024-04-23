package repository

import "testing"

func TestNewZapLogger(t *testing.T) {
	l := NewZapColoredLogger()
	l.Infof("asdfasdf %s", "aaa")
	l.WithPrefix("prefix").Infof("111 %s", "aaa")
}
