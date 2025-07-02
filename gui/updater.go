package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func checkUpdate() {
	global.EventManager.RegisterA(
		events.CheckUpdateResultUpdate, "gui.updater.check_update", gutil.ThreadSafeHandler(func(event *event.Event) {
			data := event.Data.(events.CheckUpdateResultUpdateEvent)
			msg := data.Info.Version.String() + "\n\n\n" + data.Info.Info
			if data.HasUpdate {
				dialog.ShowCustom(
					i18n.T("gui.update.new_version"),
					"OK",
					widget.NewRichTextFromMarkdown(msg),
					MainWindow)
			} else {
				dialog.ShowCustom(
					i18n.T("gui.update.already_latest_version"),
					"OK",
					widget.NewRichTextFromMarkdown(""),
					MainWindow)
			}
		}))
}
