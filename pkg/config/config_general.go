package config

type _GeneralConfig struct {
	BaseConfig
	Language        string
	AutoCheckUpdate bool
}

func (c *_GeneralConfig) Name() string {
	return "General"
}

var General = &_GeneralConfig{
	Language:        "zh-CN",
	AutoCheckUpdate: true,
}
