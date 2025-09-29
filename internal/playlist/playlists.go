package playlist

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"errors"
	"github.com/AynaLivePlayer/miaosic"
)

// todo: implement the playlist controller

var allPlaylists = make(map[string]*miaosic.Playlist)
var currentSelected string = ""

func createPlaylistManager() {
	allPlaylists = make(map[string]*miaosic.Playlist)
	for _, pl := range cfg.playlists {
		value := pl.Copy()
		allPlaylists[pl.Meta.ID()] = &value
	}
	currentSelected = ""
	if len(cfg.playlists) > 0 {
		currentSelected = cfg.playlists[0].Meta.ID()
	}

	_ = global.EventBus.Publish(
		events.PlaylistManagerCurrentUpdate,
		events.PlaylistManagerCurrentUpdateEvent{
			Medias: make([]model.Media, 0),
		})

	_ = global.EventBus.Publish(
		events.PlaylistManagerSetSystemCmd,
		events.PlaylistManagerSetSystemCmdEvent{
			PlaylistID: cfg.SystemPlaylistID,
		})

	global.EventBus.Subscribe("", events.PlaylistManagerSetSystemCmd,
		"internal.playlist.system_playlist.set",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerSetSystemCmdEvent)
			// default case
			if data.PlaylistID == "" {
				return
			}
			log.Infof("try to set system playlist %s", data.PlaylistID)
			pl, ok := allPlaylists[data.PlaylistID]
			if !ok {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("playlist not found"),
					})
				return
			}
			cfg.SystemPlaylistID = pl.Meta.ID()
			_ = global.EventBus.Publish(
				events.PlaylistManagerSystemUpdate,
				events.PlaylistManagerSystemUpdateEvent{
					Info: model.PlaylistInfo{
						Meta:  pl.Meta,
						Title: pl.DisplayName(),
					},
				})
			log.Infof("replace system playlist with %d medias", len(pl.Medias))
			medias := make([]model.Media, len(pl.Medias))
			for i, v := range pl.Medias {
				medias[i] = model.Media{
					Info: v,
					User: model.SystemUser,
				}
			}
			SystemPlaylist.Replace(medias)

		})
	global.EventBus.Subscribe("",
		events.PlaylistManagerRefreshCurrentCmd,
		"internal.playlist.current_playlist.refresh",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerRefreshCurrentCmdEvent)
			log.Infof("try to refresh playlist %s", data.PlaylistID)
			currentSelected = data.PlaylistID
			// default case
			if currentSelected == "" {
				return
			}
			pl, ok := allPlaylists[data.PlaylistID]
			if !ok {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("playlist not found"),
					})
				return
			}
			getPlaylist, err := miaosic.GetPlaylist(pl.Meta)
			if err != nil {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: err,
					})
				return
			}
			allPlaylists[pl.Meta.ID()] = getPlaylist
			updateCurrenMedias(getPlaylist)
			updatePlaylistManagerInfos()
		})

	global.EventBus.Subscribe("",
		events.PlaylistManagerGetCurrentCmd,
		"internal.playlist.current_playlist.get",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerGetCurrentCmdEvent)
			log.Infof("try to get playlist %s", data.PlaylistID)
			currentSelected = data.PlaylistID
			// default case
			if currentSelected == "" {
				return
			}
			pl, ok := allPlaylists[data.PlaylistID]
			if !ok {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("playlist not found"),
					})
				return
			}
			updateCurrenMedias(pl)
		})

	global.EventBus.Subscribe("",
		events.PlaylistManagerAddPlaylistCmd,
		"internal.playlist.add_playlist",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerAddPlaylistCmdEvent)
			log.Info("try to add playlist", data)
			meta, ok := miaosic.MatchPlaylistByProvider(data.Provider, data.URL)
			if !ok {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("not proper url"),
					})
				return
			}
			_, ok = allPlaylists[meta.ID()]
			if ok {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("playlist already exists"),
					})
				return
			}
			pl, err := miaosic.GetPlaylist(meta)
			if err != nil {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: err,
					})
				return
			}
			allPlaylists[meta.ID()] = pl
			updatePlaylistManagerInfos()
		})

	global.EventBus.Subscribe("",
		events.PlaylistManagerRemovePlaylistCmd,
		"internal.playlist.remove_playlist",
		func(event *eventbus.Event) {
			data := event.Data.(events.PlaylistManagerRemovePlaylistCmdEvent)
			if data.PlaylistID == cfg.SystemPlaylistID {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("cannot remove system playlist"),
					})
				return
			}
			_, ok := allPlaylists[data.PlaylistID]
			if !ok {
				_ = global.EventBus.Publish(
					events.ErrorUpdate,
					events.ErrorUpdateEvent{
						Error: errors.New("playlist not found"),
					})
				return
			}
			delete(allPlaylists, data.PlaylistID)
			updatePlaylistManagerInfos()
		})
	updatePlaylistManagerInfos()
}

func updateCurrenMedias(pl *miaosic.Playlist) {
	medias := make([]model.Media, len(pl.Medias))
	for i, v := range pl.Medias {
		medias[i] = model.Media{
			Info: v,
			User: model.SystemUser,
		}
	}
	_ = global.EventBus.Publish(
		events.PlaylistManagerCurrentUpdate,
		events.PlaylistManagerCurrentUpdateEvent{
			Medias: medias,
		})
}

func updatePlaylistManagerInfos() {
	playlists := make([]model.PlaylistInfo, 0)
	keys := make([]string, 0)
	for k := range allPlaylists {
		keys = append(keys, k)
	}
	for _, k := range keys {
		playlists = append(playlists, model.PlaylistInfo{
			Meta:  allPlaylists[k].Meta,
			Title: allPlaylists[k].DisplayName(),
		})
	}
	log.InfoW("update playlist manager infos")
	_ = global.EventBus.Publish(
		events.PlaylistManagerInfoUpdate,
		events.PlaylistManagerInfoUpdateEvent{
			Playlists: playlists,
		})
}
