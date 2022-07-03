package provider

import (
	"AynaLivePlayer/player"
	"AynaLivePlayer/util"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/tidwall/gjson"
	"regexp"
)

type BilibiliVideo struct {
	InfoApi   string
	FileApi   string
	SearchApi string
	BVRegex   *regexp.Regexp
	IdRegex   *regexp.Regexp
	PageRegex *regexp.Regexp
	header    map[string]string
}

func _newBilibiliVideo() *BilibiliVideo {
	return &BilibiliVideo{
		InfoApi:   "https://api.bilibili.com/x/web-interface/view/detail?bvid=%s&aid=&jsonp=jsonp",
		FileApi:   "https://api.bilibili.com/x/player/playurl?type=&otype=json&fourk=1&qn=32&avid=&bvid=%s&cid=%s",
		SearchApi: "",
		BVRegex:   regexp.MustCompile("^BV[0-9A-Za-z]+"),
		IdRegex:   regexp.MustCompile("^BV[0-9A-Za-z]+(\\?p=[0-9]+)?"),
		PageRegex: regexp.MustCompile("p=[0-9]+"),
		header: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:51.0) Gecko/20100101 Firefox/51.0",
			"Referer":    "https://www.bilibili.com/",
			"Origin":     "https://www.bilibili.com",
		},
	}
}

var BilibiliVideoAPI *BilibiliVideo

func init() {
	BilibiliVideoAPI = _newBilibiliVideo()
	Providers[BilibiliVideoAPI.GetName()] = BilibiliVideoAPI
}

func (b *BilibiliVideo) getPage(bv string) int {
	if page := b.PageRegex.FindString(bv); page != "" {
		return util.StringToInt(page[2:])
	}
	return 0
}

func (b *BilibiliVideo) getBv(bv string) string {
	return b.BVRegex.FindString(bv)
}

func (b *BilibiliVideo) GetName() string {
	return "bilibili-video"
}

func (b *BilibiliVideo) MatchMedia(keyword string) *player.Media {
	if id := b.IdRegex.FindString(keyword); id != "" {
		return &player.Media{
			Meta: Meta{
				Name: b.GetName(),
				Id:   id,
			},
		}
	}
	return nil
}

func (b *BilibiliVideo) GetPlaylist(playlist Meta) ([]*player.Media, error) {
	return nil, ErrorExternalApi
}

func (b *BilibiliVideo) FormatPlaylistUrl(uri string) string {
	return ""
}

func (b *BilibiliVideo) Search(keyword string) ([]*player.Media, error) {
	return nil, ErrorExternalApi
}

func (b *BilibiliVideo) UpdateMedia(media *player.Media) error {
	resp := httpGetString(fmt.Sprintf(b.InfoApi, media.Meta.(Meta).Id), nil)
	if resp == "" {
		return ErrorExternalApi
	}
	jresp := gjson.Parse(resp)
	if jresp.Get("data.View.title").String() == "" {
		return ErrorExternalApi
	}
	media.Title = jresp.Get("data.View.title").String()
	media.Artist = jresp.Get("data.View.owner.name").String()
	media.Cover = jresp.Get("data.View.pic").String()
	media.Album = media.Title
	return nil
}

func (b *BilibiliVideo) UpdateMediaUrl(media *player.Media) error {
	resp := httpGetString(fmt.Sprintf(b.InfoApi, media.Meta.(Meta).Id), nil)
	if resp == "" {
		return ErrorExternalApi
	}
	jresp := gjson.Parse(resp)
	page := b.getPage(media.Meta.(Meta).Id)
	cid := jresp.Get(fmt.Sprintf("data.View.pages.%d.cid", page)).String()
	if cid == "" {
		cid = jresp.Get("data.View.cid").String()
	}
	if cid == "" {
		return ErrorExternalApi
	}
	resp = httpGetString(fmt.Sprintf(b.FileApi, b.getBv(media.Meta.(Meta).Id), cid), b.header)
	if resp == "" {
		return ErrorExternalApi
	}
	jresp = gjson.Parse(resp)
	url := jresp.Get("data.durl.0.url").String()
	if url == "" {
		return ErrorExternalApi
	}
	media.Url = url
	header := make(map[string]string)
	_ = copier.Copy(&header, &b.header)
	header["Referer"] = fmt.Sprintf("https://www.bilibili.com/video/%s", b.getBv(media.Meta.(Meta).Id))
	media.Header = b.header
	return nil
}

func (b *BilibiliVideo) UpdateMediaLyric(media *player.Media) error {
	return nil
}
