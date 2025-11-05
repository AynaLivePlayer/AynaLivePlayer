package source

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"github.com/AynaLivePlayer/miaosic"
)

func handleInfo() {
	err := global.EventBus.Subscribe("",
		events.CmdMiaosicGetMediaInfo, "internal.media_provider.getMediaInfo", func(event *eventbus.Event) {
			info, err := miaosic.GetMediaInfo(event.Data.(events.CmdMiaosicGetMediaInfoData).Meta)
			_ = global.EventBus.Reply(
				event, events.ReplyMiaosicGetMediaInfo,
				events.ReplyMiaosicGetMediaInfoData{
					Info:  info,
					Error: err,
				},
			)
		})
	if err != nil {
		log.ErrorW("Subscribe search event failed", "error", err)
	}
	err = global.EventBus.Subscribe("",
		events.CmdMiaosicGetMediaUrl, "internal.media_provider.getMediaUrl", func(event *eventbus.Event) {
			urls, err := miaosic.GetMediaUrl(event.Data.(events.CmdMiaosicGetMediaUrlData).Meta, event.Data.(events.CmdMiaosicGetMediaUrlData).Quality)
			_ = global.EventBus.Reply(
				event, events.ReplyMiaosicGetMediaUrl,
				events.ReplyMiaosicGetMediaUrlData{
					Urls:  urls,
					Error: err,
				},
			)
		})
	if err != nil {
		log.ErrorW("Subscribe search event failed", "error", err)
	}
}
