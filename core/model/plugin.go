package model

type Plugin interface {
	Name() string
	Enable() error
	Disable() error
}
