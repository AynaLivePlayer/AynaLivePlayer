package config

type _GeneralConfig struct {
	BaseConfig
	Width             float32
	Height            float32
	Language          string
	InfoApiServer     string
	AutoCheckUpdate   bool
	ShowSystemTray    bool
	PlayNextOnFail    bool
	UseSystemPlaylist bool
	FixedSize         bool
	EnableSMC         bool // enable system media control
	UseSystemFonts    bool // using system fonts

}

func (c *_GeneralConfig) Name() string {
	return "General"
}

var General = &_GeneralConfig{
	Language:          "zh-CN",
	ShowSystemTray:    false,
	InfoApiServer:     "http://localhost:9090",
	AutoCheckUpdate:   true,
	Width:             960,
	Height:            480,
	PlayNextOnFail:    false,
	UseSystemPlaylist: true,
	FixedSize:         true,
	EnableSMC:         true,
	UseSystemFonts:    true,
}
