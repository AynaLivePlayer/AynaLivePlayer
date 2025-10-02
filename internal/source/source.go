//go:build !nosource

package source

import (
	"github.com/AynaLivePlayer/miaosic"
	_ "github.com/AynaLivePlayer/miaosic/providers/bilivideo"
	"github.com/AynaLivePlayer/miaosic/providers/kugou"
	_ "github.com/AynaLivePlayer/miaosic/providers/kuwo"
	"github.com/AynaLivePlayer/miaosic/providers/local"
	_ "github.com/AynaLivePlayer/miaosic/providers/netease"
	"github.com/AynaLivePlayer/miaosic/providers/qq"
)

func loadMediaProvider() {
	kugou.UseInstrumental()
	miaosic.RegisterProvider(local.NewLocal(sourceCfg.LocalSourcePath))
	if sourceCfg.QQChannel == "wechat" {
		qq.UseWechatLogin()
	} else {
		qq.UseQQLogin()
	}
}
