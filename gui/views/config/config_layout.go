package config

import (
	"AynaLivePlayer/gui/component"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var ConfigList = []ConfigLayout{
	&bascicConfig{},
}

type ConfigLayout interface {
	Title() string
	Description() string
	CreatePanel() fyne.CanvasObject
}

func AddConfigLayout(cfgs ...ConfigLayout) {
	ConfigList = append(ConfigList, cfgs...)
}

func CreateView() fyne.CanvasObject {
	// initialize config panels
	for _, c := range ConfigList {
		c.CreatePanel()
	}
	content := container.NewMax()
	entryList := widget.NewList(
		func() int {
			return len(ConfigList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(ConfigList[id].Title())
		})
	entryList.OnSelected = func(id widget.ListItemID) {
		desc := widget.NewRichTextFromMarkdown("## " + ConfigList[id].Title() + " \n\n" + ConfigList[id].Description())
		for i := range desc.Segments {
			if seg, ok := desc.Segments[i].(*widget.TextSegment); ok {
				seg.Style.Alignment = fyne.TextAlignCenter
			}
		}
		content.Objects = []fyne.CanvasObject{
			container.NewVScroll(container.NewBorder(container.NewVBox(desc, widget.NewSeparator()), nil, nil, nil, ConfigList[id].CreatePanel())),
		}
		content.Refresh()
	}

	return component.NewFixedSplitContainer(entryList, content, true, 0.23)
}
