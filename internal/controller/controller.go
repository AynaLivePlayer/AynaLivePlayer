package controller

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/internal/playlist"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
)

func Initialize() {
	handleSearch()
	createLyricLoader()
	handlePlayNext()
}

func Stop() {

}

func handlePlayNext() {
	log := global.Logger.WithPrefix("Controller")
	playerState := model.PlayerStatePlaying
	global.EventManager.RegisterA(
		events.PlayerPropertyStateUpdate,
		"internal.controller.playcontrol.idleplaynext",
		func(event *event.Event) {
			data := event.Data.(events.PlayerPropertyStateUpdateEvent)
			log.Debug("[MPV PlayControl] update player to state", playerState, "->", data.State)
			playerState = data.State
			if playerState == model.PlayerStateIdle {
				log.Info("mpv went idle, try play next")
				global.EventManager.CallA(events.PlayerPlayNextCmd,
					events.PlayerPlayNextCmdEvent{})
			}
		})

	global.EventManager.RegisterA(
		events.PlayerPropertyStateUpdate,
		"internal.controller.playcontrol.clear_when_idle", func(event *event.Event) {
			data := event.Data.(events.PlayerPropertyStateUpdateEvent)
			// if is idle, remove playing media
			if data.State == model.PlayerStateIdle {
				global.EventManager.CallA(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
					Media:   model.Media{},
					Removed: true,
				})
			}
		})

	global.EventManager.RegisterA(
		events.PlaylistInsertUpdate(model.PlaylistIDPlayer),
		"internal.controller.playcontrol.playnext_when_insert.player",
		func(event *event.Event) {
			if playerState == model.PlayerStateIdle {
				global.EventManager.CallA(events.PlayerPlayNextCmd,
					events.PlayerPlayNextCmdEvent{})
			}
		})

	global.EventManager.RegisterA(
		events.PlaylistInsertUpdate(model.PlaylistIDSystem),
		"internal.controller.playcontrol.playnext_when_insert.system",
		func(event *event.Event) {
			if playerState == model.PlayerStateIdle {
				global.EventManager.CallA(events.PlayerPlayNextCmd,
					events.PlayerPlayNextCmdEvent{})
			}
		})

	global.EventManager.RegisterA(
		events.PlayerPlayNextCmd,
		"internal.controller.playcontrol.playnext",
		func(event *event.Event) {
			if playlist.PlayerPlaylist.Size() > 0 {
				log.Infof("Try to play next media in player playlist")
				global.EventManager.CallA(events.PlaylistNextCmd(model.PlaylistIDPlayer),
					events.PlaylistNextCmdEvent{
						Remove: true,
					})
				return
			}
			if !config.General.UseSystemPlaylist {
				// do not play system playlist
				return
			}
			log.Infof("Try to play next media in system playlist")
			global.EventManager.CallA(events.PlaylistNextCmd(model.PlaylistIDSystem),
				events.PlaylistNextCmdEvent{
					Remove: false,
				})
		})

	global.EventManager.RegisterA(
		events.PlayerPlayErrorUpdate,
		"internal.controller.playcontrol.playnext_on_error",
		func(event *event.Event) {
			if config.General.PlayNextOnFail {
				global.EventManager.CallA(events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
				return
			}
		})

	global.EventManager.RegisterA(events.PlaylistNextUpdate(model.PlaylistIDPlayer),
		"internal.controller.playcontrol.play_when_next", func(event *event.Event) {
			data := event.Data.(events.PlaylistNextUpdateEvent)
			global.EventManager.CallA(
				events.PlayerPlayCmd,
				events.PlayerPlayCmdEvent{
					Media: data.Media,
				})
		})

	global.EventManager.RegisterA(events.PlaylistNextUpdate(model.PlaylistIDSystem),
		"internal.controller.playcontrol.play_when_next.system_playlist", func(event *event.Event) {
			data := event.Data.(events.PlaylistNextUpdateEvent)
			global.EventManager.CallA(
				events.PlayerPlayCmd,
				events.PlayerPlayCmdEvent{
					Media: data.Media,
				})
		})
}
