package provider

import (
	"AynaLivePlayer/model"
	"fmt"
	"testing"
)

func TestNetease_Search(t *testing.T) {
	var api MediaProvider = NeteaseAPI
	result, err := api.Search("染 reol")
	if err != nil {
		return
	}
	fmt.Println(result)
	media := result[0]
	fmt.Println(media)
	err = api.UpdateMediaUrl(media)
	fmt.Println(err)
	fmt.Println(media.Url)
}

func TestNetease_Search2(t *testing.T) {
	var api MediaProvider = NeteaseAPI
	result, err := api.Search("出山")
	if err != nil {
		return
	}
	t.Log(result)
	media := result[0]
	t.Log(media)
	err = api.UpdateMediaUrl(media)
	t.Log(err)
	t.Log(media.Url)
}

func TestNetease_GetMusicMeta(t *testing.T) {
	var api MediaProvider = NeteaseAPI

	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "33516503",
		},
	}
	err := api.UpdateMedia(&media)
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println(media)
}

func TestNetease_GetMusic(t *testing.T) {
	var api MediaProvider = NeteaseAPI
	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "33516503",
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

func TestNetease_GetPlaylist(t *testing.T) {
	var api MediaProvider = NeteaseAPI
	playlist, err := api.GetPlaylist(&model.Meta{
		Name: api.GetName(),
		//Id:   "2520739691",
		Id: "2382819181",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(playlist))
	for _, media := range playlist {
		fmt.Println(media.Title, media.Artist, media.Album)
	}

}

func TestNetease_UpdateMediaLyric(t *testing.T) {
	var api MediaProvider = NeteaseAPI
	media := model.Media{
		Meta: model.Meta{
			Name: api.GetName(),
			Id:   "33516503",
		},
	}
	err := api.UpdateMediaLyric(&media)
	fmt.Println(err)
	fmt.Println(media.Lyric)
}
