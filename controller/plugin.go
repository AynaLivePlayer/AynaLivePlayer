package controller

type Plugin interface {
	Name() string
	Enable() error
	Disable() error
}

type IPluginController interface {
	LoadPlugin(plugin Plugin)
	LoadPlugins(plugins ...Plugin)
	ClosePlugins()
}
