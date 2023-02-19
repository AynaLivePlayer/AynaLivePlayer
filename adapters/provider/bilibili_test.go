package provider

import (
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/model"
	"fmt"
	"testing"
)

func TestBilibili_Search(t *testing.T) {
	var api adapter.MediaProvider = BilibiliAPI
	result, err := api.Search("染 reol")
	if err != nil {
		fmt.Println(1, err)
		return
	}
	fmt.Println(result)
	media := result[0]
	fmt.Println(*media)
	err = api.UpdateMediaUrl(media)
	fmt.Println(err)
	fmt.Println(media.Url)
}

func TestBilibili_GetMusicMeta(t *testing.T) {
	var api adapter.MediaProvider = BilibiliAPI

	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "1560601",
		},
	}
	err := api.UpdateMedia(&media)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(media)
}

func TestBilibili_GetMusic(t *testing.T) {
	var api adapter.MediaProvider = BilibiliAPI
	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "1560601",
		},
	}
	err := api.UpdateMedia(&media)
	if err != nil {
		return
	}
	err = api.UpdateMediaUrl(&media)
	if err != nil {
		return
	}
	fmt.Println(media)
	fmt.Println(media.Url)
}
