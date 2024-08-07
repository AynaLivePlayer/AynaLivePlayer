package diange

import (
	"AynaLivePlayer/gui/xfyne"
	"AynaLivePlayer/pkg/i18n"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type blacklistItem struct {
	Exact bool
	Value string
}

type blacklist struct {
	panel fyne.CanvasObject
}

func (b *blacklist) Title() string {
	return i18n.T("plugin.diange.blacklist.title")
}

func (b *blacklist) Description() string {
	return i18n.T("plugin.diange.blacklist.description")
}

func (b *blacklist) CreatePanel() fyne.CanvasObject {
	if b.panel != nil {
		return b.panel
	}
	// UI组件
	input := xfyne.EntryDisableUndoRedo(widget.NewEntry())
	input.SetPlaceHolder(i18n.T("plugin.diange.blacklist.input.placeholder"))

	exactText := i18n.T("plugin.diange.blacklist.option.exact")
	containsText := i18n.T("plugin.diange.blacklist.option.contains")

	options := widget.NewRadioGroup([]string{exactText, containsText}, nil)
	options.SetSelected(containsText)

	var blackListGui *widget.List

	blackListGui = widget.NewList(
		func() int {
			return len(diange.blacklist)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabel(""),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
				widget.NewLabel(""))
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			if diange.blacklist[lii].Exact {
				co.(*fyne.Container).Objects[1].(*widget.Label).SetText(exactText)
			} else {
				co.(*fyne.Container).Objects[1].(*widget.Label).SetText(containsText)
			}
			co.(*fyne.Container).Objects[0].(*widget.Label).SetText(diange.blacklist[lii].Value)
			co.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				diange.blacklist = append(diange.blacklist[:lii], diange.blacklist[lii+1:]...)
				blackListGui.Refresh()
			}
		},
	)

	addButton := widget.NewButton(i18n.T("plugin.diange.blacklist.btn.add"), func() {
		if input.Text != "" {
			diange.blacklist = append(diange.blacklist, blacklistItem{
				Exact: options.Selected == exactText,
				Value: input.Text,
			})
			input.SetText("")
			blackListGui.Refresh()
		}
	})
	form := container.NewBorder(
		nil,
		nil,
		options, addButton, input,
	)
	b.panel = container.NewBorder(form, nil, nil, nil, blackListGui)
	return b.panel
}
