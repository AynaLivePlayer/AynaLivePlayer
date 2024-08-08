package xfyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"reflect"
)

func EntryDisableUndoRedo(entry *widget.Entry) *widget.Entry {
	val := reflect.ValueOf(entry).Elem().FieldByName("shortcut").Addr().UnsafePointer()
	(*fyne.ShortcutHandler)(val).RemoveShortcut(&fyne.ShortcutRedo{})
	(*fyne.ShortcutHandler)(val).RemoveShortcut(&fyne.ShortcutUndo{})
	return entry
}
