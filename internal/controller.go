package internal

import (
	"AynaLivePlayer/core/adapter"
)

type Controller struct {
	app      adapter.IApplication        `ini:"-"`
	liveroom adapter.ILiveRoomController `ini:"-"`
	player   adapter.IPlayController     `ini:"-"`
	lyric    adapter.ILyricLoader        `ini:"-"`
	playlist adapter.IPlaylistController `ini:"-"`
	provider adapter.IProviderController `ini:"-"`
	plugin   adapter.IPluginController   `ini:"-"`
	log      adapter.ILogger             `ini:"-"`
}

func (c *Controller) Logger() adapter.ILogger {
	return c.log
}

func NewController(
	liveroom adapter.ILiveRoomController, player adapter.IPlayController,
	playlist adapter.IPlaylistController,
	provider adapter.IProviderController, plugin adapter.IPluginController,
	log adapter.ILogger) adapter.IControlBridge {
	cc := &Controller{
		app:      &AppBilibiliChannel{},
		liveroom: liveroom,
		player:   player,
		playlist: playlist,
		provider: provider,
		plugin:   plugin,
		log:      log,
	}
	return cc
}

func (c *Controller) App() adapter.IApplication {
	return c.app
}

func (c *Controller) LiveRooms() adapter.ILiveRoomController {
	return c.liveroom
}

func (c *Controller) PlayControl() adapter.IPlayController {
	return c.player
}

func (c *Controller) Playlists() adapter.IPlaylistController {
	return c.playlist
}

func (c *Controller) Provider() adapter.IProviderController {
	return c.provider
}

func (c *Controller) Plugin() adapter.IPluginController {
	return c.plugin
}

func (c *Controller) LoadPlugins(plugins ...adapter.Plugin) {
	c.plugin.LoadPlugins(plugins...)
}

func (c *Controller) CloseAndSave() {
	c.plugin.ClosePlugins()
}
