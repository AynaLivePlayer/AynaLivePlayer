package internal

import (
	"AynaLivePlayer/internal/controller"
	"AynaLivePlayer/internal/liveroom"
	"AynaLivePlayer/internal/player"
	"AynaLivePlayer/internal/playlist"
	"AynaLivePlayer/internal/plugins"
	"AynaLivePlayer/internal/source"
	"AynaLivePlayer/internal/sysmediacontrol"
	"AynaLivePlayer/internal/updater"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/durationmgmt"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/sourcelogin"
	"AynaLivePlayer/plugin/textinfo"
	"AynaLivePlayer/plugin/wshub"
	"AynaLivePlayer/plugin/yinliang"
)

func Initialize() {
	player.SetupMpvPlayer()
	source.Initialize()
	playlist.Initialize()
	controller.Initialize()
	liveroom.Initialize()
	plugins.Initialize()
	plugins.LoadPlugins(
		diange.NewDiange(), qiege.NewQiege(), yinliang.NewYinliang(), sourcelogin.NewSourceLogin(),
		textinfo.NewTextInfo(),
		durationmgmt.NewMaxDuration(),
		wshub.NewWsHub(),
	)
	updater.Initialize()
	if config.General.EnableSMC {
		sysmediacontrol.InitSystemMediaControl()
	}
}

func Stop() {
	if config.General.EnableSMC {
		sysmediacontrol.Destroy()
	}
	liveroom.StopAndSave()
	playlist.Close()
	plugins.ClosePlugins()
	player.StopMpvPlayer()
}
