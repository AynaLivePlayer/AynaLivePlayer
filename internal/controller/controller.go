package controller

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/internal/playlist"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/eventbus"
)

func Initialize() {
	handlePlayNext()
}

func Stop() {

}

func handlePlayNext() {
	log := global.Logger.WithPrefix("Controller")
	playerState := model.PlayerStatePlaying
	global.EventBus.Subscribe("",
		events.PlayerPropertyStateUpdate,
		"internal.controller.playcontrol.idleplaynext",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlayerPropertyStateUpdateEvent)
			log.Debug("[MPV PlayControl] update player to state", playerState, "->", data.State)
			playerState = data.State
			if playerState == model.PlayerStateIdle {
				log.Info("mpv went idle, try play next")
				_ = global.EventBus.Publish(events.PlayerPlayNextCmd,
					events.PlayerPlayNextCmdEvent{})
			}
		})

	global.EventBus.Subscribe("",
		events.PlayerPropertyStateUpdate,
		"internal.controller.playcontrol.clear_when_idle", func(event *eventbus.Event) {
			data := event.Data.(events.PlayerPropertyStateUpdateEvent)
			// if is idle, remove playing media
			if data.State == model.PlayerStateIdle {
				_ = global.EventBus.Publish(events.PlayerPlayingUpdate, events.PlayerPlayingUpdateEvent{
					Media:   model.Media{},
					Removed: true,
				})
			}
		})

	global.EventBus.Subscribe("",
		events.PlaylistInsertUpdate(model.PlaylistIDPlayer),
		"internal.controller.playcontrol.playnext_when_insert.player",
		func(event *eventbus.Event) {
			if playerState == model.PlayerStateIdle {
				_ = global.EventBus.Publish(events.PlayerPlayNextCmd,
					events.PlayerPlayNextCmdEvent{})
			}
		})

	global.EventBus.Subscribe("",
		events.PlaylistInsertUpdate(model.PlaylistIDSystem),
		"internal.controller.playcontrol.playnext_when_insert.system",
		func(event *eventbus.Event) {
			if playerState == model.PlayerStateIdle {
				_ = global.EventBus.Publish(events.PlayerPlayNextCmd,
					events.PlayerPlayNextCmdEvent{})
			}
		})

	global.EventBus.Subscribe("",
		events.PlayerPlayNextCmd,
		"internal.controller.playcontrol.playnext",
		func(event *eventbus.Event) {
			if playlist.PlayerPlaylist.Size() > 0 {
				log.Infof("Try to play next media in player playlist")
				_ = global.EventBus.Publish(events.PlaylistNextCmd(model.PlaylistIDPlayer),
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
			_ = global.EventBus.Publish(events.PlaylistNextCmd(model.PlaylistIDSystem),
				events.PlaylistNextCmdEvent{
					Remove: false,
				})
		})

	global.EventBus.Subscribe("",
		events.PlayerPlayErrorUpdate,
		"internal.controller.playcontrol.playnext_on_error",
		func(event *eventbus.Event) {
			if config.General.PlayNextOnFail {
				_ = global.EventBus.Publish(events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
				return
			}
		})

	global.EventBus.Subscribe("", events.PlaylistNextUpdate(model.PlaylistIDPlayer),
		"internal.controller.playcontrol.play_when_next", func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistNextUpdateEvent)
			_ = global.EventBus.Publish(
				events.PlayerPlayCmd,
				events.PlayerPlayCmdEvent{
					Media: data.Media,
				})
		})

	global.EventBus.Subscribe("", events.PlaylistNextUpdate(model.PlaylistIDSystem),
		"internal.controller.playcontrol.play_when_next.system_playlist", func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistNextUpdateEvent)
			_ = global.EventBus.Publish(
				events.PlayerPlayCmd,
				events.PlayerPlayCmdEvent{
					Media: data.Media,
				})
		})
}
