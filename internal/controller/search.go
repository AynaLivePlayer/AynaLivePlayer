package controller

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"github.com/AynaLivePlayer/miaosic"
)

func handleSearch() {
	log := global.Logger.WithPrefix("Search")
	global.EventBus.Subscribe("",
		events.SearchCmd, "internal.controller.search.handleSearchCmd", func(event *eventbus.Event) {
			data := event.Data.(events.SearchCmdEvent)
			log.Infof("Search %s using %s", data.Keyword, data.Provider)
			searchResult, err := miaosic.SearchByProvider(data.Provider, data.Keyword, 1, 10)
			if err != nil {
				log.Warnf("Search %s using %s failed: %s", data.Keyword, data.Provider, err)
				return
			}
			medias := make([]model.Media, len(searchResult))
			for i, v := range searchResult {
				medias[i] = model.Media{
					Info: v,
					User: model.SystemUser,
				}
			}
			_ = global.EventBus.Publish(
				events.SearchResultUpdate, events.SearchResultUpdateEvent{
					Medias: medias,
				})
		})
}
