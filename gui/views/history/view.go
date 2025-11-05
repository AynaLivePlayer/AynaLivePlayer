package history

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/gui/gctx"
	"AynaLivePlayer/gui/gutil"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/i18n"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

var medias []model.Media
var listWidget *widget.List
var mux sync.RWMutex

func CreateView() fyne.CanvasObject {
	view := createHistoryList()
	global.EventBus.Subscribe(gctx.EventChannel,
		events.PlayerPlayingUpdate,
		"gui.history.playing_update",
		gutil.ThreadSafeHandler(func(event *eventbus.Event) {
			if event.Data.(events.PlayerPlayingUpdateEvent).Removed {
				return
			}
			mux.Lock()
			medias = append(medias, event.Data.(events.PlayerPlayingUpdateEvent).Media)
			if len(medias) > 1000 {
				medias = medias[len(medias)-1000:]
			}
			listWidget.Refresh()
			mux.Unlock()
		}))
	return view
}

func createHistoryList() fyne.CanvasObject {
	listWidget = widget.NewList(
		func() int {
			return len(medias)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabel("index"),
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
					widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil),
				),
				container.NewGridWithColumns(3,
					component.NewLabelWithOpts("title", component.LabelTruncation(fyne.TextTruncateClip)),
					component.NewLabelWithOpts("artist", component.LabelTruncation(fyne.TextTruncateClip)),
					component.NewLabelWithOpts("user", component.LabelTruncation(fyne.TextTruncateClip))))
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			m := medias[len(medias)-id-1]
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(
				m.Info.Title)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(
				m.Info.Artist)
			object.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(
				m.ToUser().Name)
			object.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", id))
			btns := object.(*fyne.Container).Objects[2].(*fyne.Container).Objects
			m.User = model.SystemUser
			btns[0].(*widget.Button).OnTapped = func() {
				_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlayerPlayCmd, events.PlayerPlayCmdEvent{
					Media: m,
				})
			}
			btns[1].(*widget.Button).OnTapped = func() {
				_ = global.EventBus.PublishToChannel(gctx.EventChannel, events.PlaylistInsertCmd(model.PlaylistIDPlayer), events.PlaylistInsertCmdEvent{
					Media:    m,
					Position: -1,
				})
			}
		})
	return container.NewBorder(
		container.NewBorder(nil, nil,
			widget.NewLabel("#"), widget.NewLabel(i18n.T("gui.history.operation")),
			container.NewGridWithColumns(3,
				widget.NewLabel(i18n.T("gui.history.title")),
				widget.NewLabel(i18n.T("gui.history.artist")),
				widget.NewLabel(i18n.T("gui.history.user")))),
		nil, nil, nil,
		listWidget,
	)
}
