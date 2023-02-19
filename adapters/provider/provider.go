package provider

import (
	"AynaLivePlayer/core/adapter"
	model "AynaLivePlayer/core/model"
)

var RegisteredProviders = map[string]adapter.MediaProviderCtor{
	"netease":        NewNetease,
	"local":          NewLocalCtor,
	"kuwo":           NewKuwo,
	"bilibili":       NewBilibili,
	"bilibili-video": NewBilibiliVideo,
}

var Providers map[string]adapter.MediaProvider = make(map[string]adapter.MediaProvider)

func GetPlaylist(meta *model.Meta) ([]*model.Media, error) {
	if v, ok := Providers[meta.Name]; ok {
		return v.GetPlaylist(meta)
	}
	return nil, ErrorNoSuchProvider
}

func FormatPlaylistUrl(pname, uri string) (string, error) {
	if v, ok := Providers[pname]; ok {
		return v.FormatPlaylistUrl(uri), nil
	}
	return "", ErrorNoSuchProvider
}

func MatchMedia(provider string, keyword string) *model.Media {
	if v, ok := Providers[provider]; ok {
		return v.MatchMedia(keyword)
	}
	return nil
}

func Search(provider string, keyword string) ([]*model.Media, error) {
	if v, ok := Providers[provider]; ok {
		return v.Search(keyword)
	}
	return nil, ErrorNoSuchProvider
}

func UpdateMedia(media *model.Media) error {
	if v, ok := Providers[media.Meta.(model.Meta).Name]; ok {
		return v.UpdateMedia(media)
	}
	return ErrorNoSuchProvider
}

func UpdateMediaUrl(media *model.Media) error {
	if v, ok := Providers[media.Meta.(model.Meta).Name]; ok {
		return v.UpdateMediaUrl(media)
	}
	return ErrorNoSuchProvider
}

func UpdateMediaLyric(media *model.Media) error {
	if v, ok := Providers[media.Meta.(model.Meta).Name]; ok {
		return v.UpdateMediaLyric(media)
	}
	return ErrorNoSuchProvider
}
