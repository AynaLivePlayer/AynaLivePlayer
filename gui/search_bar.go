package gui

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/pkg/event"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var SearchBar = &struct {
	Input     *component.Entry
	Button    *component.AsyncButton
	UseSource *widget.Select
}{}

func createSearchBar() fyne.CanvasObject {
	SearchBar.Input = component.NewEntry()
	SearchBar.Input.SetPlaceHolder(i18n.T("gui.search.placeholder"))
	SearchBar.Input.OnKeyUp = func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyReturn {
			SearchBar.Button.OnTapped()
		}
	}
	SearchBar.Button = component.NewAsyncButton(i18n.T("gui.search.search"), func() {
		keyword := SearchBar.Input.Text
		pr := SearchBar.UseSource.Selected
		logger.Debugf("Search keyword: %s, provider: %s", keyword, pr)
		SearchResult.mux.Lock()
		SearchResult.Items = make([]model.Media, 0)
		SearchResult.List.Refresh()
		SearchResult.mux.Unlock()
		global.EventManager.CallA(events.SearchCmd, events.SearchCmdEvent{
			Keyword:  keyword,
			Provider: pr,
		})
	})

	global.EventManager.RegisterA(events.MediaProviderUpdate,
		"gui.search.provider.update", func(event *event.Event) {
			providers := event.Data.(events.MediaProviderUpdateEvent)
			s := make([]string, len(providers.Providers))
			copy(s, providers.Providers)
			SearchBar.UseSource.Options = s
			if len(s) > 0 {
				SearchBar.UseSource.SetSelected(s[0])
			}
		})

	SearchBar.UseSource = widget.NewSelect([]string{}, func(s string) {
	})

	searchInput := container.NewBorder(
		nil, nil, widget.NewLabel(i18n.T("gui.search.search")), SearchBar.Button,
		container.NewBorder(nil, nil, SearchBar.UseSource, nil, SearchBar.Input))
	return container.NewVBox(
		searchInput,
		widget.NewSeparator())
}
