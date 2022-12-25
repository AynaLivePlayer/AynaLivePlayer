package provider

import (
	"AynaLivePlayer/common/util"
	"AynaLivePlayer/model"
	"sort"
)

func MediaSort(keyword string, medias []*model.Media) {
	mediaDist := make([]struct {
		media *model.Media
		dist  int
	}, len(medias))
	for i, media := range medias {
		mediaDist[i].media = media
		mediaDist[i].dist = util.StrLen(util.LongestCommonString(keyword, media.Title)) +
			util.StrLen(util.LongestCommonString(keyword, media.Artist))
	}
	sort.Slice(mediaDist, func(i, j int) bool {
		return mediaDist[i].dist > mediaDist[j].dist
	})
	for i, media := range mediaDist {
		medias[i] = media.media
	}
	return
}
