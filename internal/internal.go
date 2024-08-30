package internal

import (
	"AynaLivePlayer/internal/controller"
	"AynaLivePlayer/internal/liveroom"
	"AynaLivePlayer/internal/player"
	"AynaLivePlayer/internal/playlist"
	"AynaLivePlayer/internal/plugins"
	"AynaLivePlayer/internal/source"
	"AynaLivePlayer/internal/updater"
	"AynaLivePlayer/plugin/diange"
	"AynaLivePlayer/plugin/durationmgmt"
	"AynaLivePlayer/plugin/qiege"
	"AynaLivePlayer/plugin/sourcelogin"
	"AynaLivePlayer/plugin/textinfo"
	"AynaLivePlayer/plugin/wshub"
)

func Initialize() {
	player.SetupMpvPlayer()
	source.Initialize()
	playlist.Initialize()
	controller.Initialize()
	liveroom.Initialize()
	plugins.Initialize()
	plugins.LoadPlugins(
		diange.NewDiange(), qiege.NewQiege(), sourcelogin.NewSourceLogin(),
		textinfo.NewTextInfo(),
		durationmgmt.NewMaxDuration(),
		wshub.NewWsHub(),
	)
	updater.Initialize()
	//sysmediacontrol.InitSystemMediaControl()
}

func Stop() {
	//sysmediacontrol.Destroy()
	liveroom.StopAndSave()
	playlist.Close()
	plugins.ClosePlugins()
	player.StopMpvPlayer()
}
