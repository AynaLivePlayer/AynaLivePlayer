package plugins

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/logger"
)

var plugins []model.Plugin = make([]model.Plugin, 0)
var log logger.ILogger

func Initialize() {
	plugins = make([]model.Plugin, 0)
	log = global.Logger.WithPrefix("Plugin")
}

func LoadPlugin(plugin model.Plugin) {
	log.Info("[Plugin] Loading plugin: " + plugin.Name())
	if err := plugin.Enable(); err != nil {
		log.Warnf("[Plugin] Failed to load plugin: %s, %s", plugin.Name(), err)
		return
	}
	plugins = append(plugins, plugin)
}

func LoadPlugins(plugins ...model.Plugin) {
	for _, plugin := range plugins {
		LoadPlugin(plugin)
	}
}

func ClosePlugins() {
	for _, plugin := range plugins {
		if err := plugin.Disable(); err != nil {
			log.Warnf("[Plugin] Failed to close plugin: %s, %s", plugin.Name(), err)
			continue
		}
	}
}
