package gutil

import (
	"AynaLivePlayer/pkg/event"
	"fyne.io/fyne/v2"
)

// since 2.6.1, calls to fyne API from other go routine must be wrapped in fyne.Do
func ThreadSafeHandler(fn func(e *event.Event)) func(e *event.Event) {
	return func(e *event.Event) {
		fyne.Do(func() {
			fn(e)
		})
	}
}
