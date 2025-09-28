//go:build !nosource

package source

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"github.com/AynaLivePlayer/miaosic"
	_ "github.com/AynaLivePlayer/miaosic/providers/bilivideo"
	"github.com/AynaLivePlayer/miaosic/providers/kugou"
	_ "github.com/AynaLivePlayer/miaosic/providers/kuwo"
	"github.com/AynaLivePlayer/miaosic/providers/local"
	_ "github.com/AynaLivePlayer/miaosic/providers/netease"
	_ "github.com/AynaLivePlayer/miaosic/providers/qq"
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
	kugou.UseInstrumental()
	miaosic.RegisterProvider(local.NewLocal(sourceCfg.LocalSourcePath))

	_ = global.EventBus.Publish(
		events.MediaProviderUpdate, events.MediaProviderUpdateEvent{
			Providers: miaosic.ListAvailableProviders(),
		})
}
