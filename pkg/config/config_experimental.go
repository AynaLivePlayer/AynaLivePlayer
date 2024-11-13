package config

type _ExperimentalConfig struct {
	BaseConfig
	Headless bool
}

func (c *_ExperimentalConfig) Name() string {
	return "Experimental"
}

var Experimental = &_ExperimentalConfig{
	Headless: false,
}
