package xfyne

import (
	"fyne.io/fyne/v2/widget"
	"testing"
)

func TestEntryDisableUndoRedo(t *testing.T) {
	entry := widget.NewEntry()
	EntryDisableUndoRedo(entry)
}
