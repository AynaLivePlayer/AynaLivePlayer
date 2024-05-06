package config

type _GeneralConfig struct {
	BaseConfig
	Width           float32
	Height          float32
	Language        string
	InfoApiServer   string
	AutoCheckUpdate bool
	ShowSystemTray  bool
	PlayNextOnFail  bool
}

func (c *_GeneralConfig) Name() string {
	return "General"
}

var General = &_GeneralConfig{
	Language:        "zh-CN",
	ShowSystemTray:  false,
	InfoApiServer:   "http://localhost:9090",
	AutoCheckUpdate: true,
	Width:           960,
	Height:          480,
	PlayNextOnFail:  false,
}
