package core

import (
	"AynaLivePlayer/common/logger"
	"AynaLivePlayer/controller"
)

var lg = logger.Logger.WithField("Module", "CoreController")

type Controller struct {
	liveroom controller.ILiveRoomController `ini:"-"`
	player   controller.IPlayController     `ini:"-"`
	lyric    controller.ILyricLoader        `ini:"-"`
	playlist controller.IPlaylistController `ini:"-"`
	provider controller.IProviderController `ini:"-"`
	plugin   controller.IPluginController   `ini:"-"`
}

func NewController(
	liveroom controller.ILiveRoomController, player controller.IPlayController,
	playlist controller.IPlaylistController,
	provider controller.IProviderController, plugin controller.IPluginController) controller.IController {
	cc := &Controller{
		liveroom: liveroom,
		player:   player,
		playlist: playlist,
		provider: provider,
		plugin:   plugin,
	}
	return cc
}

func (c *Controller) LiveRooms() controller.ILiveRoomController {
	return c.liveroom
}

func (c *Controller) PlayControl() controller.IPlayController {
	return c.player
}

func (c *Controller) Playlists() controller.IPlaylistController {
	return c.playlist
}

func (c *Controller) Provider() controller.IProviderController {
	return c.provider
}

func (c *Controller) Plugin() controller.IPluginController {
	return c.plugin
}

func (c *Controller) LoadPlugins(plugins ...controller.Plugin) {
	c.plugin.LoadPlugins(plugins...)
}

func (c *Controller) CloseAndSave() {
	c.plugin.ClosePlugins()
}
