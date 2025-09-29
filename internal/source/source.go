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
	"github.com/AynaLivePlayer/miaosic/providers/qq"
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

func Initialize() {
	config.LoadConfig(sourceCfg)
	kugou.UseInstrumental()
	miaosic.RegisterProvider(local.NewLocal(sourceCfg.LocalSourcePath))
	if sourceCfg.QQChannel == "wechat" {
		qq.UseWechatLogin()
	} else {
		qq.UseQQLogin()
	}

	_ = global.EventBus.Publish(
		events.MediaProviderUpdate, events.MediaProviderUpdateEvent{
			Providers: miaosic.ListAvailableProviders(),
		})
}
