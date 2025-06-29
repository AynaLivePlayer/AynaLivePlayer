package config

type _ExperimentalConfig struct {
	BaseConfig
	Headless   bool
	PlayerCore string
}

func (c *_ExperimentalConfig) Name() string {
	return "Experimental"
}

var Experimental = &_ExperimentalConfig{
	Headless:   false,
	PlayerCore: "mpv",
}
