package source

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"github.com/AynaLivePlayer/miaosic"
)

func handleSourceLogin() {
	err := global.EventBus.Subscribe("",
		events.CmdMiaosicQrLogin, "internal.media_provider.qrlogin_handler", func(event *eventbus.Event) {
			data := event.Data.(events.CmdMiaosicQrLoginData)
			log.Infof("trying login %s", data.Provider)
			pvdr, ok := miaosic.GetProvider(data.Provider)
			if !ok {
				_ = global.EventBus.Reply(
					event, events.ReplyMiaosicQrLogin,
					events.ReplyMiaosicQrLoginData{
						Session: miaosic.QrLoginSession{},
						Error:   miaosic.ErrorNoSuchProvider,
					})
				return
			}
			result, ok := pvdr.(miaosic.Loginable)
			if !ok {
				_ = global.EventBus.Reply(
					event, events.ReplyMiaosicQrLogin,
					events.ReplyMiaosicQrLoginData{
						Session: miaosic.QrLoginSession{},
						Error:   miaosic.ErrNotImplemented,
					})
				return
			}
			var session miaosic.QrLoginSession
			sess, err := result.QrLogin()
			if err == nil && sess != nil {
				session = *sess
			}
			_ = global.EventBus.Reply(
				event, events.ReplyMiaosicQrLogin,
				events.ReplyMiaosicQrLoginData{
					Session: session,
					Error:   err,
				})
		})
	if err != nil {
		log.ErrorW("Subscribe search event failed", "error", err)
	}
}
