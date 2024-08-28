package xfyne

import (
	"fyne.io/fyne/v2/widget"
)

func EntryDisableUndoRedo(entry *widget.Entry) *widget.Entry {
	// do nothing because the bug has been fixed in fyne@v2.5.1
	return entry
	//val := reflect.ValueOf(entry).Elem().FieldByName("shortcut").Addr().UnsafePointer()
	//(*fyne.ShortcutHandler)(val).RemoveShortcut(&fyne.ShortcutRedo{})
	//(*fyne.ShortcutHandler)(val).RemoveShortcut(&fyne.ShortcutUndo{})
	//return entry
}
