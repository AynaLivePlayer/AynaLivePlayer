package provider

import (
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/model"
	neteaseApi "github.com/XiaoMengXinX/Music163Api-Go/api"
	neteaseTypes "github.com/XiaoMengXinX/Music163Api-Go/types"
	neteaseUtil "github.com/XiaoMengXinX/Music163Api-Go/utils"
	"regexp"
	"strconv"
	"strings"
)

type Netease struct {
	PlaylistRegex0 *regexp.Regexp
	PlaylistRegex1 *regexp.Regexp
	ReqData        neteaseUtil.RequestData
	IdRegex0       *regexp.Regexp
	IdRegex1       *regexp.Regexp
	loginStatus    neteaseTypes.LoginStatusData
}

func NewNetease(config MediaProviderConfig) MediaProvider {
	return &Netease{
		PlaylistRegex0: regexp.MustCompile("^[0-9]+$"),
		// https://music.163.com/playlist?id=2382819181&userid=95906480
		PlaylistRegex1: regexp.MustCompile("playlist\\?id=[0-9]+"),
		ReqData: neteaseUtil.RequestData{
			Headers: neteaseUtil.Headers{
				{
					"X-Real-IP",
					"118.88.88.88",
				},
			},
		},
		IdRegex0: regexp.MustCompile("^[0-9]+"),
		IdRegex1: regexp.MustCompile("^wy[0-9]+"),
	}
}

func _newNetease() *Netease {
	return &Netease{
		PlaylistRegex0: regexp.MustCompile("^[0-9]+$"),
		// https://music.163.com/playlist?id=2382819181&userid=95906480
		PlaylistRegex1: regexp.MustCompile("playlist\\?id=[0-9]+"),
		ReqData: neteaseUtil.RequestData{
			Headers: neteaseUtil.Headers{
				{
					"X-Real-IP",
					"118.88.88.88",
				},
			},
		},
		IdRegex0: regexp.MustCompile("^[0-9]+"),
		IdRegex1: regexp.MustCompile("^wy[0-9]+"),
	}
}

var NeteaseAPI *Netease

func init() {
	NeteaseAPI = _newNetease()
	Providers[NeteaseAPI.GetName()] = NeteaseAPI
}

// Netease private helper method

func _neteaseGetArtistNames(data neteaseTypes.SongDetailData) string {
	artists := make([]string, 0)
	for _, a := range data.Ar {
		artists = append(artists, a.Name)
	}
	return strings.Join(artists, ",")
}

// MediaProvider implementation

func (n *Netease) GetName() string {
	return "netease"
}

func (n *Netease) MatchMedia(keyword string) *model.Media {
	if id := n.IdRegex0.FindString(keyword); id != "" {
		return &model.Media{
			Meta: model.Meta{
				Name: n.GetName(),
				Id:   id,
			},
		}
	}
	if id := n.IdRegex1.FindString(keyword); id != "" {
		return &model.Media{
			Meta: model.Meta{
				Name: n.GetName(),
				Id:   id[2:],
			},
		}
	}
	return nil
}

func (n *Netease) FormatPlaylistUrl(uri string) string {
	var id string
	id = n.PlaylistRegex0.FindString(uri)
	if id != "" {
		return id
	}
	id = n.PlaylistRegex1.FindString(uri)
	if id != "" {
		return id[12:]
	}
	return ""
}

func (n *Netease) GetPlaylist(playlist *model.Meta) ([]*model.Media, error) {
	result, err := neteaseApi.GetPlaylistDetail(
		n.ReqData, util.Atoi(playlist.Id))
	if err != nil || result.Code != 200 {
		return nil, ErrorExternalApi
	}
	cnt := len(result.Playlist.TrackIds)
	if cnt == 0 {
		return nil, ErrorExternalApi
	}

	ids := make([]int, len(result.Playlist.TrackIds))
	for i := 0; i < cnt; i++ {
		ids[i] = result.Playlist.TrackIds[i].Id
	}
	medias := make([]*model.Media, 0, cnt)
	for index := 0; index < len(ids); index += 1000 {
		result2, err := neteaseApi.GetSongDetail(
			n.ReqData,
			ids[index:util.IntMin(index+1000, len(ids))])
		if err != nil || result2.Code != 200 {
			break
		}
		cnt = len(result2.Songs)
		if cnt == 0 {
			break
		}
		for i := 0; i < cnt; i++ {
			medias = append(medias, &model.Media{
				Title:  result2.Songs[i].Name,
				Artist: _neteaseGetArtistNames(result2.Songs[i]),
				Cover:  model.Picture{Url: result2.Songs[i].Al.PicUrl},
				Album:  result2.Songs[i].Al.Name,
				Url:    "",
				Header: nil,
				User:   nil,
				Meta: model.Meta{
					Name: n.GetName(),
					Id:   strconv.Itoa(result2.Songs[i].Id),
				},
			})
		}
	}
	if len(medias) == 0 {
		return nil, ErrorExternalApi
	}
	return medias, nil
}

func (n *Netease) Search(keyword string) ([]*model.Media, error) {
	rawResult, err := neteaseApi.SearchSong(
		n.ReqData,
		neteaseApi.SearchSongConfig{
			Keyword: keyword,
			Limit:   30,
			Offset:  0,
		})
	if err != nil || rawResult.Code != 200 {
		return nil, ErrorExternalApi
	}
	medias := make([]*model.Media, 0)
	for _, song := range rawResult.Result.Songs {
		artists := make([]string, 0)
		for _, a := range song.Artists {
			artists = append(artists, a.Name)
		}
		medias = append(medias, &model.Media{
			Title:  song.Name,
			Artist: strings.Join(artists, ","),
			Cover:  model.Picture{},
			Album:  song.Album.Name,
			Url:    "",
			Header: nil,
			Meta: model.Meta{
				Name: n.GetName(),
				Id:   strconv.Itoa(song.Id),
			},
		})
	}
	return medias, nil
}

func (n *Netease) UpdateMedia(media *model.Media) error {
	result, err := neteaseApi.GetSongDetail(
		n.ReqData,
		[]int{util.Atoi(media.Meta.(model.Meta).Id)})
	if err != nil || result.Code != 200 {
		return ErrorExternalApi
	}
	if len(result.Songs) == 0 {
		return ErrorExternalApi
	}
	media.Title = result.Songs[0].Name
	media.Cover.Url = result.Songs[0].Al.PicUrl
	media.Album = result.Songs[0].Al.Name
	media.Artist = _neteaseGetArtistNames(result.Songs[0])
	return nil
}

func (n *Netease) UpdateMediaUrl(media *model.Media) error {
	result, err := neteaseApi.GetSongURL(
		n.ReqData,
		neteaseApi.SongURLConfig{Ids: []int{util.Atoi(media.Meta.(model.Meta).Id)}})
	if err != nil || result.Code != 200 {
		return ErrorExternalApi
	}
	if len(result.Data) == 0 {
		return ErrorExternalApi
	}
	if result.Data[0].Code != 200 {
		return ErrorExternalApi
	}
	media.Url = result.Data[0].Url
	return nil
}

func (n *Netease) UpdateMediaLyric(media *model.Media) error {
	result, err := neteaseApi.GetSongLyric(n.ReqData, util.Atoi(media.Meta.(model.Meta).Id))
	if err != nil || result.Code != 200 {
		return ErrorExternalApi
	}
	media.Lyric = result.Lrc.Lyric
	return nil
}
