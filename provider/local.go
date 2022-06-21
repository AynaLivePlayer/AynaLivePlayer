package provider

import "AynaLivePlayer/player"

type Local struct {
}

var LocalAPI *Local

func init() {
	LocalAPI = _newLocal()
	//Providers[LocalAPI.GetName()] = LocalAPI
}

func _newLocal() *Local {
	return &Local{}
}

func (l *Local) GetName() string {
	return "local"
}

func (l *Local) FormatPlaylistUrl(uri string) string {
	return ""
}

func (l *Local) GetPlaylist(playlist string) ([]*player.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Local) Search(keyword string) ([]*player.Media, error) {
	//TODO implement me
	panic("implement me")
}

func (l *Local) UpdateMedia(media *player.Media) error {
	//TODO implement me
	panic("implement me")
}

func (l *Local) UpdateMediaUrl(media *player.Media) error {
	//TODO implement me
	panic("implement me")
}
