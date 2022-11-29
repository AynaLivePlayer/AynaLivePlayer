package config

type _ProviderConfig struct {
	BaseConfig
	Priority []string
	LocalDir string
}

func (c *_ProviderConfig) Name() string {
	return "Provider"
}

var Provider = &_ProviderConfig{
	Priority: []string{"netease", "kuwo", "bilibili", "local", "bilibili-video"},
	LocalDir: "./music",
}
