package controller

type Plugin interface {
	Name() string
	Enable() error
}

func LoadPlugin(plugin Plugin) {
	l().Info("Loading plugin: " + plugin.Name())
	if err := plugin.Enable(); err != nil {
		l().Warnf("Failed to load plugin: %s, %s", plugin.Name(), err)
	}
}

func LoadPlugins(plugins ...Plugin) {
	for _, plugin := range plugins {
		LoadPlugin(plugin)
	}
}
