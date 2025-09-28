package gutil

import (
	"AynaLivePlayer/pkg/eventbus"
	"fyne.io/fyne/v2"
)

// since 2.6.1, calls to fyne API from other go routine must be wrapped in fyne.Do
func ThreadSafeHandler(fn func(e *eventbus.Event)) func(e *eventbus.Event) {
	//return fn
	return func(e *eventbus.Event) {
		fyne.Do(func() {
			fn(e)
		})
	}
}

func RunInFyneThread(fn func()) {
	//fn()
	fyne.Do(fn)
}
