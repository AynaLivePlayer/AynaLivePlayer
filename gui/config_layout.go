package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TestConfig struct {
}

func (t *TestConfig) Title() string {
	return "Test Title"
}

func (T *TestConfig) Description() string {
	return "Test Description"
}

func (t *TestConfig) CreatePanel() fyne.CanvasObject {
	return widget.NewLabel("asdf")
}

func createConfigLayout() fyne.CanvasObject {
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
			return widget.NewLabel("AAAAAAAAAAAAAAAA")
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
		a := container.NewVScroll(ConfigList[id].CreatePanel())
		content.Objects = []fyne.CanvasObject{
			container.NewBorder(container.NewVBox(desc, widget.NewSeparator()), nil, nil, nil,
				a),
		}

		content.Refresh()
	}
	return container.NewBorder(
		nil, nil,
		container.NewHBox(entryList, widget.NewSeparator()), nil,
		content)
}
