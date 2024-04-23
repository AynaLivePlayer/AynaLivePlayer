package config

type _GeneralConfig struct {
	BaseConfig
	Width           float32
	Height          float32
	Language        string
	AutoCheckUpdate bool
	ShowSystemTray  bool
}

func (c *_GeneralConfig) Name() string {
	return "General"
}

var General = &_GeneralConfig{
	Language:        "zh-CN",
	ShowSystemTray:  false,
	AutoCheckUpdate: true,
	Width:           960,
	Height:          480,
}
