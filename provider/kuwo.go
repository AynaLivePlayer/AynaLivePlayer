package provider

import (
	"AynaLivePlayer/player"
	"fmt"
	"github.com/tidwall/gjson"
	"html"
	"net/url"
	"regexp"
	"strings"
)

type Kuwo struct {
	InfoApi        string
	FileApi        string
	SearchCookie   string
	SearchApi      string
	LyricApi       string
	PlaylistApi    string
	PlaylistRegex0 *regexp.Regexp
	PlaylistRegex1 *regexp.Regexp
	IdRegex0       *regexp.Regexp
	IdRegex1       *regexp.Regexp
}

func _newKuwo() *Kuwo {
	return &Kuwo{
		InfoApi: "http://www.kuwo.cn/api/www/music/musicInfo?mid=%s&httpsStatus=1",
		//FileApi:      "http://www.kuwo.cn/api/v1/www/music/playUrl?mid=%d&type=music&httpsStatus=1",
		FileApi:        "http://antiserver.kuwo.cn/anti.s?type=convert_url&format=mp3&response=url&rid=MUSIC_%s",
		SearchCookie:   "http://kuwo.cn/search/list?key=%s",
		SearchApi:      "http://www.kuwo.cn/api/www/search/searchMusicBykeyWord?key=%s&pn=%d&rn=%d",
		LyricApi:       "http://m.kuwo.cn/newh5/singles/songinfoandlrc?musicId=%s",
		PlaylistApi:    "http://www.kuwo.cn/api/www/playlist/playListInfo?pid=%s&pn=%d&rn=%d&httpsStatus=1",
		PlaylistRegex0: regexp.MustCompile("[0-9]+"),
		PlaylistRegex1: regexp.MustCompile("playlist/[0-9]+"),
		IdRegex0:       regexp.MustCompile("^[0-9]+"),
		IdRegex1:       regexp.MustCompile("^kw[0-9]+"),
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

func (k *Kuwo) MatchMedia(keyword string) *player.Media {
	if id := k.IdRegex0.FindString(keyword); id != "" {
		return &player.Media{
			Meta: Meta{
				Name: k.GetName(),
				Id:   id,
			},
		}
	}
	if id := k.IdRegex1.FindString(keyword); id != "" {
		return &player.Media{
			Meta: Meta{
				Name: k.GetName(),
				Id:   id[2:],
			},
		}
	}
	return nil
}

func (k *Kuwo) FormatPlaylistUrl(uri string) string {
	var id string
	id = k.PlaylistRegex0.FindString(uri)
	if id != "" {
		return id
	}
	id = k.PlaylistRegex1.FindString(uri)
	if id != "" {
		return id[9:]
	}
	return ""
}

//func (k *Kuwo) _kuwoGet(url string) string {
//	searchCookie, err := httpHead(fmt.Sprintf(k.SearchCookie, "any"), nil)
//	if err != nil {
//		return ""
//	}
//	kwToken, ok := util.SliceString(regexp.MustCompile("kw_token=([^;])*;").FindString(searchCookie.Header().Get("set-cookie")), 9, -1)
//	if !ok {
//		return ""
//	}
//	return httpGetString(url, map[string]string{
//		"cookie":  "kw_token=" + kwToken,
//		"csrf":    kwToken,
//		"referer": "http://www.kuwo.cn/",
//	})
//}

func (k *Kuwo) _kuwoGet(url string) string {
	return httpGetString(url, map[string]string{
		"cookie":  "kw_token=" + "95MWTYC4FP",
		"csrf":    "95MWTYC4FP",
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
			Cover:  player.Picture{Url: value.Get("pic").String()},
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
	media.Cover.Url = jresp.Get("data.pic").String()
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
	result := httpGetString(fmt.Sprintf(k.LyricApi, media.Meta.(Meta).Id), nil)
	if result == "" {
		return ErrorExternalApi
	}
	lrcs := make([]string, 0)
	gjson.Parse(result).Get("data.lrclist").ForEach(func(key, value gjson.Result) bool {
		lrcs = append(lrcs, fmt.Sprintf("[00:%s]%s", value.Get("time").String(), value.Get("lineLyric").String()))

		return true
	})
	media.Lyric = strings.Join(lrcs, "\n")
	return nil
}

func (k *Kuwo) GetPlaylist(meta Meta) ([]*player.Media, error) {
	medias := make([]*player.Media, 0)
	var resp string
	var jresp gjson.Result
	for i := 1; i <= 20; i++ {
		resp = k._kuwoGet(fmt.Sprintf(k.PlaylistApi, meta.Id, i, 128))
		if resp == "" {
			break
		}
		//fmt.Println(resp[:100])
		jresp = gjson.Parse(resp)
		//fmt.Println(jresp.Get("code").String())
		if jresp.Get("code").String() != "200" {
			break
		}
		cnt := int(jresp.Get("data.total").Int())
		//fmt.Println(cnt)
		//fmt.Println(len(jresp.Get("data.musicList").Array()))
		jresp.Get("data.musicList").ForEach(func(key, value gjson.Result) bool {
			medias = append(
				medias,
				&player.Media{
					Title:  html.UnescapeString(value.Get("name").String()),
					Artist: value.Get("artist").String(),
					Cover:  player.Picture{Url: value.Get("pic").String()},
					Album:  value.Get("album").String(),
					Meta: Meta{
						Name: k.GetName(),
						Id:   value.Get("rid").String(),
					},
				})
			return true
		})
		if cnt <= i*100 {
			break
		}
	}
	if len(medias) == 0 {
		return nil, ErrorExternalApi
	}
	return medias, nil
}
