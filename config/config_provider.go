package config

type _ProviderConfig struct {
	Priority []string
	LocalDir string
}

func (c *_ProviderConfig) Name() string {
	return "Provider"
}

var Provider = &_ProviderConfig{
	Priority: []string{"netease", "kuwo", "bilibili"},
	LocalDir: "./music",
}
