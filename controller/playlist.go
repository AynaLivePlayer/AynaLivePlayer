package controller

import "AynaLivePlayer/player"

func AddToHistory(media *player.Media) {
	l().Tracef("add media %s (%s) to history", media.Title, media.Artist)
	media = media.Copy()
	if History.Size() >= 1024 {
		History.Replace([]*player.Media{})
	}
	History.Push(media)
	return
}

func ToHistoryMedia(media *player.Media) *player.Media {
	media = media.Copy()
	media.User = HistoryUser
	return media
}

func ToSystemMedia(media *player.Media) *player.Media {
	media = media.Copy()
	media.User = player.SystemUser
	return media
}
