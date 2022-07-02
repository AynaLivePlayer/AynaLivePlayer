package provider

import (
	"AynaLivePlayer/player"
	"fmt"
	"testing"
)

func TestBV_GetMusicMeta(t *testing.T) {
	var api MediaProvider = BilibiliVideoAPI

	media := player.Media{
		Meta: Meta{
			Name: api.GetName(),
			Id:   "BV1434y1q71P",
		},
	}
	err := api.UpdateMedia(&media)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(media)
}

func TestBV_GetMusic(t *testing.T) {
	var api MediaProvider = BilibiliVideoAPI
	media := player.Media{
		Meta: Meta{
			Name: api.GetName(),
			Id:   "BV1434y1q71P",
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
	//fmt.Println(media)
	fmt.Println(media.Url)
}
