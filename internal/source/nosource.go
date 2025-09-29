//go:build nosource

package source

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"github.com/AynaLivePlayer/miaosic"
)

type _sourceConfig struct {
	LocalSourcePath string
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
}

func Initialize() {
	config.LoadConfig(sourceCfg)
	miaosic.RegisterProvider(&dummySource{})
	_ = global.EventBus.Publish(
		events.MediaProviderUpdate, events.MediaProviderUpdateEvent{
			Providers: miaosic.ListAvailableProviders(),
		})
}
