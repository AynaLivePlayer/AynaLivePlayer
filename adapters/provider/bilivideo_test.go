package provider

import (
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/model"
	"fmt"
	"regexp"
	"testing"
)

func TestBV_GetMusicMeta(t *testing.T) {
	var api adapter.MediaProvider = BilibiliVideoAPI

	media := model.Media{
		Meta: model.Meta{
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
	var api adapter.MediaProvider = BilibiliVideoAPI
	media := model.Media{
		Meta: model.Meta{
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

func TestBV_Regex(t *testing.T) {
	fmt.Println(regexp.MustCompile("^BV[0-9A-Za-z]+(\\?p=[0-9]+)?").FindString("BV1gA411P7ir?p=3"))
}

func TestBV_GetMusicMeta2(t *testing.T) {
	var api adapter.MediaProvider = BilibiliVideoAPI

	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "BV1gA411P7ir?p=3",
		},
	}
	err := api.UpdateMedia(&media)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(media)
}

func TestBV_GetMusic2(t *testing.T) {
	var api adapter.MediaProvider = BilibiliVideoAPI
	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "BV1gA411P7ir?p=3",
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

func TestBV_Search(t *testing.T) {
	var api adapter.MediaProvider = BilibiliVideoAPI
	result, err := api.Search("家有女友")
	if err != nil {
		fmt.Println(1, err)
		return
	}
	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r.Artist)
	}
}
