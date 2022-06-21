package provider

import (
	"AynaLivePlayer/player"
	"AynaLivePlayer/util"
	"fmt"
	"github.com/tidwall/gjson"
	"html"
	"net/url"
	"regexp"
)

type Kuwo struct {
	InfoApi      string
	FileApi      string
	SearchCookie string
	SearchApi    string
}

func _newKuwo() *Kuwo {
	return &Kuwo{
		InfoApi: "http://www.kuwo.cn/api/www/music/musicInfo?mid=%s&httpsStatus=1",
		//FileApi:      "http://www.kuwo.cn/api/v1/www/music/playUrl?mid=%d&type=music&httpsStatus=1",
		FileApi:      "http://antiserver.kuwo.cn/anti.s?type=convert_url&format=mp3&response=url&rid=MUSIC_%s",
		SearchCookie: "http://kuwo.cn/search/list?key=%s",
		SearchApi:    "http://www.kuwo.cn/api/www/search/searchMusicBykeyWord?key=%s&pn=%d&rn=%d",
	}
}

var KuwoAPI *Kuwo

func init() {
	KuwoAPI = _newKuwo()
	Providers[KuwoAPI.GetName()] = KuwoAPI
}

func (k *Kuwo) GetName() string {
	return "kuwo"
}

func (k *Kuwo) FormatPlaylistUrl(uri string) string {
	return ""
}

func (k *Kuwo) _kuwoGet(url string) string {
	searchCookie, err := httpHead(fmt.Sprintf(k.SearchCookie, "any"), nil)
	if err != nil {
		return ""
	}
	kwToken, ok := util.SliceString(regexp.MustCompile("kw_token=([^;])*;").FindString(searchCookie.Header().Get("set-cookie")), 9, -1)
	if !ok {
		return ""
	}
	return httpGetString(url, map[string]string{
		"cookie":  "kw_token=" + kwToken,
		"csrf":    kwToken,
		"referer": "http://www.kuwo.cn/",
	})
}

func (k *Kuwo) Search(keyword string) ([]*player.Media, error) {
	resp := k._kuwoGet(fmt.Sprintf(k.SearchApi, url.QueryEscape(keyword), 1, 64))
	if resp == "" {
		return nil, ErrorExternalApi
	}
	result := make([]*player.Media, 0)
	gjson.Parse(resp).Get("data.list").ForEach(func(key, value gjson.Result) bool {
		result = append(result, &player.Media{
			Title:  html.UnescapeString(value.Get("name").String()),
			Cover:  value.Get("pic").String(),
			Artist: value.Get("artist").String(),
			Album:  value.Get("album").String(),
			Meta: Meta{
				Name: k.GetName(),
				Id:   value.Get("rid").String(),
			},
		})
		return true
	})
	return result, nil
}

func (k *Kuwo) UpdateMedia(media *player.Media) error {
	resp := k._kuwoGet(fmt.Sprintf(k.InfoApi, media.Meta.(Meta).Id))
	if resp == "" {
		return ErrorExternalApi
	}
	jresp := gjson.Parse(resp)
	if jresp.Get("data.musicrid").String() == "" {
		return ErrorExternalApi
	}
	media.Title = html.UnescapeString(jresp.Get("data.name").String())
	media.Cover = jresp.Get("data.pic").String()
	media.Artist = jresp.Get("data.artist").String()
	media.Album = jresp.Get("data.album").String()
	return nil
}

func (k *Kuwo) UpdateMediaUrl(media *player.Media) error {
	result := httpGetString(fmt.Sprintf(k.FileApi, media.Meta.(Meta).Id), nil)
	if result == "" {
		return ErrorExternalApi
	}
	media.Url = result
	return nil
}

func (k *Kuwo) UpdateMediaLyric(media *player.Media) error {
	fmt.Println(k._kuwoGet("https://player.kuwo.cn/webmusic/st/getNewMuiseByRid?rid=MUSIC_22804772"))
	return nil
}

func (k *Kuwo) GetPlaylist(meta Meta) ([]*player.Media, error) {
	return nil, ErrorExternalApi
}
