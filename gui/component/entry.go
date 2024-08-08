package component

import (
	"AynaLivePlayer/gui/xfyne"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Entry struct {
	widget.Entry
	OnKeyUp   func(key *fyne.KeyEvent)
	OnKeyDown func(key *fyne.KeyEvent)
}

func NewEntry() *Entry {
	e := &Entry{}
	e.ExtendBaseWidget(e)
	xfyne.EntryDisableUndoRedo(&e.Entry)
	return e
}

func (m *Entry) KeyUp(key *fyne.KeyEvent) {
	m.Entry.KeyUp(key)
	if m.OnKeyUp != nil {
		m.OnKeyUp(key)
	}
}

func (m *Entry) KeyDown(key *fyne.KeyEvent) {
	m.Entry.KeyDown(key)
	if m.OnKeyDown != nil {
		m.OnKeyDown(key)
	}
}
