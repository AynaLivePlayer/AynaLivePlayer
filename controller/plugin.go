package controller

type Plugin interface {
	Name() string
	Enable() error
	Disable() error
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

func ClosePlugins(plugins ...Plugin) {
	for _, plugin := range plugins {
		err := plugin.Disable()
		if err != nil {
			l().Warnf("Failed to close plugin: %s, %s", plugin.Name(), err)
			return
		}
	}
}
