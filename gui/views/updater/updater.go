package updater

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateUpdaterPopUp() {
	global.EventBus.Subscribe(gctx.EventChannel,
		events.CheckUpdateResultUpdate, "gui.updater.check_update", gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			data := event.Data.(events.CheckUpdateResultUpdateEvent)
			msg := data.Info.Version.String() + "\n\n\n" + data.Info.Info
			if data.HasUpdate {
				dialog.ShowCustom(
					i18n.T("gui.update.new_version"),
					"OK",
					widget.NewRichTextFromMarkdown(msg),
					gctx.Context.Window)
			} else {
				dialog.ShowCustom(
					i18n.T("gui.update.already_latest_version"),
					"OK",
					widget.NewRichTextFromMarkdown(""),
					gctx.Context.Window)
			}
		}))
}
