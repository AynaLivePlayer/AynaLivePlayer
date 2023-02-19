package gui

import (
	"AynaLivePlayer/common/i18n"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/internal"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
		l().Debugf("Search keyword: %s, provider: %s", keyword, pr)
		items, err := API.Provider().SearchWithProvider(keyword, pr)
		if err != nil {
			dialog.ShowError(err, MainWindow)
		}
		model.ApplyUser(items, internal.SystemUser)
		SearchResult.Items = items
		SearchResult.List.Refresh()
	})
	s := make([]string, len(API.Provider().GetPriority()))
	copy(s, API.Provider().GetPriority())

	SearchBar.UseSource = widget.NewSelect(s, func(s string) {})
	if len(s) > 0 {
		SearchBar.UseSource.SetSelected(s[0])
	}
	searchInput := container.NewBorder(
		nil, nil, widget.NewLabel(i18n.T("gui.search.search")), SearchBar.Button,
		container.NewBorder(nil, nil, SearchBar.UseSource, nil, SearchBar.Input))
	return container.NewVBox(
		searchInput,
		widget.NewSeparator())
}
