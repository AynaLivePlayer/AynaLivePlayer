package source

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/logger"
	"github.com/AynaLivePlayer/miaosic"
)

type _sourceConfig struct {
	LocalSourcePath string
	QQChannel       string
}

func (_ _sourceConfig) Name() string {
	return "Source"
}

func (_ _sourceConfig) OnLoad() {
}

func (_ _sourceConfig) OnSave() {
}

var sourceCfg = &_sourceConfig{
	LocalSourcePath: "./music",
	QQChannel:       "qq",
}

var log logger.ILogger = nil

func Initialize() {
	config.LoadConfig(sourceCfg)

	log = global.Logger.WithPrefix("MediaProvider")

	loadMediaProvider()
	handleSearch()
	handleInfo()
	createLyricLoader()

	_ = global.EventBus.Publish(
		events.MediaProviderUpdate, events.MediaProviderUpdateEvent{
			Providers: miaosic.ListAvailableProviders(),
		})
}
