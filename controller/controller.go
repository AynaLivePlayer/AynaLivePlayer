package controller

var Instance IController = nil

type IController interface {
	LiveRooms() ILiveRoomController
	PlayControl() IPlayController
	Playlists() IPlaylistController
	Provider() IProviderController
	Plugin() IPluginController
	LoadPlugins(plugins ...Plugin)
	CloseAndSave()
}
