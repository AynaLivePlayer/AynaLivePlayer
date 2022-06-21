package gui

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/player"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SearchBarContainer struct {
	Input     *widget.Entry
	Button    *widget.Button
	UseSource *widget.CheckGroup
	Items     []*player.Media
}

var SearchBar = &SearchBarContainer{}

func createSearchBar() fyne.CanvasObject {
	SearchBar.Input = widget.NewEntry()
	SearchBar.Input.SetPlaceHolder("Keyword")
	SearchBar.Button = widget.NewButton("Search", nil)
	SearchBar.Button.OnTapped = createAsyncOnTapped(SearchBar.Button, func() {
		keyword := SearchBar.Input.Text
		s := make([]string, len(SearchBar.UseSource.Selected))

		copy(s, SearchBar.UseSource.Selected)
		items := make([]*player.Media, 0)
		for _, p := range s {
			if r, err := controller.SearchWithProvider(keyword, p); err == nil {
				items = append(items, r...)
			}
		}
		controller.ApplyUser(items, player.SystemUser)
		SearchResult.Items = items
		SearchResult.List.Refresh()
	})
	s := make([]string, len(config.Provider.Priority))
	copy(s, config.Provider.Priority)

	SearchBar.UseSource = widget.NewCheckGroup(s, nil)
	SearchBar.UseSource.Horizontal = true
	SearchBar.UseSource.SetSelected(s)
	searchInput := container.NewBorder(
		nil, nil, widget.NewLabel("Search"), SearchBar.Button,
		SearchBar.Input)
	return container.NewVBox(
		searchInput,
		container.NewHBox(widget.NewLabel("Source filter: "), SearchBar.UseSource),
		widget.NewSeparator())
}
