package core

import (
	"AynaLivePlayer/controller"
	"github.com/sirupsen/logrus"
)

type PluginController struct {
	plugins map[string]controller.Plugin
}

func NewPluginController() controller.IPluginController {
	return &PluginController{
		plugins: make(map[string]controller.Plugin),
	}
}

func (p *PluginController) LoadPlugin(plugin controller.Plugin) {
	lg.Info("[Plugin] Loading plugin: " + plugin.Name())
	if _, ok := p.plugins[plugin.Name()]; ok {
		logrus.Warnf("[Plugin] plugin with same name already exists, skip")
		return
	}
	if err := plugin.Enable(); err != nil {
		lg.Warnf("Failed to load plugin: %s, %s", plugin.Name(), err)
		return
	}
	p.plugins[plugin.Name()] = plugin
}

func (p *PluginController) LoadPlugins(plugins ...controller.Plugin) {
	for _, plugin := range plugins {
		p.LoadPlugin(plugin)
	}
}

func (p *PluginController) ClosePlugins() {
	for _, plugin := range p.plugins {
		if err := plugin.Disable(); err != nil {
			lg.Warnf("Failed to close plugin: %s, %s", plugin.Name(), err)
			continue
		}
	}
}
