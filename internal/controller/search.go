package controller

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/event"
	"github.com/AynaLivePlayer/miaosic"
)

func handleSearch() {
	log := global.Logger.WithPrefix("Search")
	global.EventManager.RegisterA(
		events.SearchCmd, "internal.controller.search.handleSearchCmd", func(event *event.Event) {
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
			global.EventManager.CallA(
				events.SearchResultUpdate, events.SearchResultUpdateEvent{
					Medias: medias,
				})
		})
	global.EventManager.CallA(
		events.SearchProviderUpdate, events.SearchProviderUpdateEvent{
			Providers: miaosic.ListAvailableProviders(),
		})
}
