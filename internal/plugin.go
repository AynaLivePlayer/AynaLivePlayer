package internal

import (
	"AynaLivePlayer/core/adapter"
	"github.com/sirupsen/logrus"
)

type PluginController struct {
	plugins map[string]adapter.Plugin
	log     adapter.ILogger
}

func NewPluginController(log adapter.ILogger) adapter.IPluginController {
	return &PluginController{
		plugins: make(map[string]adapter.Plugin),
		log:     log,
	}
}

func (p *PluginController) LoadPlugin(plugin adapter.Plugin) {
	p.log.Info("[Plugin] Loading plugin: " + plugin.Name())
	if _, ok := p.plugins[plugin.Name()]; ok {
		logrus.Warnf("[Plugin] plugin with same name already exists, skip")
		return
	}
	if err := plugin.Enable(); err != nil {
		p.log.Warnf("[Plugin] Failed to load plugin: %s, %s", plugin.Name(), err)
		return
	}
	p.plugins[plugin.Name()] = plugin
}

func (p *PluginController) LoadPlugins(plugins ...adapter.Plugin) {
	for _, plugin := range plugins {
		p.LoadPlugin(plugin)
	}
}

func (p *PluginController) ClosePlugins() {
	for _, plugin := range p.plugins {
		if err := plugin.Disable(); err != nil {
			p.log.Warnf("[Plugin] Failed to close plugin: %s, %s", plugin.Name(), err)
			continue
		}
	}
}
